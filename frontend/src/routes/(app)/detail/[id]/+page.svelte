<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authState, API_URL } from '$lib/auth.svelte';
	import { slide } from 'svelte/transition';

	let id = $derived(page.params.id);
	
	type Debt = {
		id: number;
		creditor_id: number;
		debtor_id: number;
		amount: number;
		description: string;
		created_at: string;
		creditor_name: string;
		debtor_name: string;
		total_paid: number;
		remaining: number;
		payments: {
			id: number;
			amount: number;
			created_at: string;
		}[];
	};

	let debt: Debt | null = $state(null);
	let loading = $state(true);
	let error = $state('');

	let paymentAmount = $state('');
	let paymentLoading = $state(false);
	let paymentError = $state('');

	async function fetchDebt() {
		try {
			const res = await fetch(`${API_URL}/debts/${id}`, {
				headers: { 'Authorization': `Bearer ${authState.token}` }
			});
			if (!res.ok) throw new Error('Gagal memuat data tagihan');
			debt = await res.json();
		} catch (err: any) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchDebt();
	});

	function playApplePaySound() {
		try {
			const AudioContext = window.AudioContext || (window as any).webkitAudioContext;
			if (!AudioContext) return;
			const ctx = new AudioContext();
			
			const playNote = (freq: number, startTime: number, duration: number) => {
				const osc = ctx.createOscillator();
				const gain = ctx.createGain();
				
				osc.type = 'sine';
				osc.frequency.setValueAtTime(freq, startTime);
				
				gain.gain.setValueAtTime(0, startTime);
				gain.gain.linearRampToValueAtTime(0.5, startTime + 0.03);
				gain.gain.exponentialRampToValueAtTime(0.01, startTime + duration);
				
				osc.connect(gain);
				gain.connect(ctx.destination);
				
				osc.start(startTime);
				osc.stop(startTime + duration);
			};

			const now = ctx.currentTime;
			// Dua nada cepat mirip Apple Pay success chime
			playNote(987.77, now, 0.15);       // Nada pertama (B5)
			playNote(1318.51, now + 0.12, 0.4); // Nada kedua (E6)
		} catch (e) {
			console.error('Audio playback failed', e);
		}
	}

	async function handlePayment(e: Event) {
		e.preventDefault();
		paymentLoading = true;
		paymentError = '';

		try {
			const res = await fetch(`${API_URL}/debts/${id}/payments`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${authState.token}`
				},
				body: JSON.stringify({ amount: paymentAmount })
			});

			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Gagal mencatat pembayaran');
			}

			// Mainkan suara jika berhasil (terutama jika otomatis diterima / pengguna adalah kreditur)
			if (debt && debt.creditor_id === authState.userId) {
				playApplePaySound();
			}

			paymentAmount = 0;
			await fetchDebt();
		} catch (err: any) {
			paymentError = err.message;
		} finally {
			paymentLoading = false;
		}
	}

	async function updatePaymentStatus(paymentId: number, status: string) {
		if (!confirm(`Apakah Anda yakin ingin ${status === 'confirmed' ? 'menerima' : 'menolak'} pembayaran ini?`)) return;
		
		try {
			const res = await fetch(`${API_URL}/debts/${id}/payments/${paymentId}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${authState.token}`
				},
				body: JSON.stringify({ status })
			});

			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Gagal mengubah status pembayaran');
			}

			if (status === 'confirmed') {
				playApplePaySound();
			}

			await fetchDebt();
		} catch (err: any) {
			alert(err.message);
		}
	}

	let deleteLoading = $state(false);

	async function handleDelete() {
		if (!confirm('Apakah Anda yakin ingin menghapus tagihan ini? (Tindakan ini tidak dapat dibatalkan)')) return;
		
		deleteLoading = true;
		try {
			const res = await fetch(`${API_URL}/debts/${id}`, {
				method: 'DELETE',
				headers: {
					'Authorization': `Bearer ${authState.token}`
				}
			});

			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Gagal menghapus piutang');
			}

			goto('/dashboard');
		} catch (err: any) {
			alert(err.message);
			deleteLoading = false;
		}
	}

	function formatRp(val: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(val);
	}
</script>

{#if loading}
	<p class="text-center text-gray-500 py-4">Memuat data...</p>
{:else if error}
	<p class="text-red-500 text-center">{error}</p>
{:else if debt}
	<div class="pt-10 pb-4 space-y-8">
		<h1 class="px-4 text-4xl md:text-5xl font-bold tracking-tighter text-[#1d1d1f]">Detail Tagihan.</h1>

		<div class="mx-4 bg-ios-card-light rounded-ios-md shadow-ios overflow-hidden">
			<ul class="divide-y divide-ios-separator-light">
				<li class="flex justify-between items-center px-4 py-3">
					<span class="text-[17px] text-[#1d1d1f] font-medium">Kreditur</span>
					<span class="text-[17px] text-gray-500">{debt.creditor_name}</span>
				</li>
				<li class="flex justify-between items-center px-4 py-3">
					<span class="text-[17px] text-[#1d1d1f] font-medium">Debitur</span>
					<span class="text-[17px] text-gray-500">{debt.debtor_name}</span>
				</li>
				<li class="flex justify-between items-center px-4 py-3">
					<span class="text-[17px] text-[#1d1d1f] font-medium">Total Nominal</span>
					<span class="text-[17px] text-gray-500">{formatRp(debt.amount)}</span>
				</li>
				<li class="flex justify-between items-center px-4 py-3 bg-green-50/50">
					<span class="text-[17px] text-[#1d1d1f] font-medium">Telah Dibayar</span>
					<span class="text-[17px] text-green-600 font-medium">{formatRp(debt.total_paid)}</span>
				</li>
				<li class="flex justify-between items-center px-4 py-3 bg-red-50/50">
					<span class="text-[17px] text-[#1d1d1f] font-medium">Sisa Tagihan</span>
					<span class="text-[17px] text-red-600 font-bold">{formatRp(debt.remaining)}</span>
				</li>
			</ul>
		</div>

		{#if debt.description}
			<div class="mx-4">
				<h2 class="px-4 pb-2 text-[13px] text-gray-500 uppercase tracking-wide font-medium">Keterangan</h2>
				<div class="bg-ios-card-light rounded-ios-md shadow-ios p-4">
					<p class="text-[17px] text-[#1d1d1f]">{debt.description}</p>
				</div>
			</div>
		{/if}

		{#if debt.payments.length === 0 && debt.creditor_id === authState.userId}
			<div class="mx-4 mt-8">
				<button 
					onclick={handleDelete}
					disabled={deleteLoading}
					class="w-full bg-ios-card-light text-red-500 font-medium text-[17px] py-3.5 rounded-full shadow-ios active:opacity-70 transition-opacity disabled:opacity-50 border border-red-100"
				>
					{deleteLoading ? 'Menghapus...' : 'Hapus Piutang'}
				</button>
			</div>
		{/if}

		{#if debt.remaining > 0}
			<div class="mx-4 mt-8">
				<h2 class="px-4 pb-2 text-[13px] text-gray-500 uppercase tracking-wide font-medium">Pembayaran</h2>
				<div class="bg-ios-card-light rounded-ios-md shadow-ios overflow-hidden">
					{#if paymentError}
						<div class="bg-red-50 text-red-600 p-3 text-[13px] border-b border-red-100">{paymentError}</div>
					{/if}
					<form onsubmit={handlePayment}>
						<div class="flex items-center px-4 py-3 bg-white border-b border-ios-separator-light">
							<label for="payAmount" class="w-1/3 text-[17px] text-[#1d1d1f] font-medium">Nominal</label>
							<div class="w-2/3 flex items-center justify-end">
								<span class="text-gray-400 mr-1">Rp</span>
								<input 
									type="number" 
									id="payAmount" 
									bind:value={paymentAmount} 
									required 
									min="1"
									max={debt.remaining}
									class="w-full bg-transparent text-[17px] text-right focus:outline-none text-ios-blue-dark font-medium"
									placeholder={debt.remaining.toString()}
								/>
							</div>
						</div>
						<div class="bg-ios-bg-light p-4">
							<button 
								type="submit" 
								disabled={paymentLoading || !paymentAmount}
								class="w-full bg-[#1d1d1f] text-white font-medium text-[17px] py-3.5 rounded-full active:opacity-70 transition-opacity disabled:opacity-50 shadow-sm"
							>
								{#if paymentLoading}
									Memproses...
								{:else if debt.creditor_id === authState.userId}
									Catat Pembayaran
								{:else}
									Ajukan Pembayaran
								{/if}
							</button>
						</div>
					</form>
				</div>
			</div>
		{:else}
			<div class="mx-4">
				<div class="bg-green-100 border border-green-200 p-4 rounded-ios-md shadow-sm text-center">
					<p class="text-green-800 font-medium text-[17px]">Tagihan ini sudah lunas!</p>
				</div>
			</div>
		{/if}

		{#if debt.payments && debt.payments.length > 0}
			<div class="mx-4 pb-8 mt-8">
				<h2 class="px-4 pb-2 text-[13px] text-gray-500 uppercase tracking-wide font-medium">Riwayat Pembayaran</h2>
				<div class="bg-ios-card-light rounded-ios-md shadow-ios overflow-hidden">
					<ul class="divide-y divide-ios-separator-light">
						{#each debt.payments as payment (payment.id)}
							<li transition:slide class="px-4 py-3 flex flex-col">
								<div class="flex justify-between items-start mb-1">
									<div class="flex flex-col">
										<span class="text-[17px] {payment.status === 'pending' ? 'text-gray-500' : 'text-green-600'} font-medium">
											+{formatRp(payment.amount)}
										</span>
										<span class="text-[13px] text-gray-400 mt-0.5">
											{new Date(payment.created_at).toLocaleDateString('id-ID')} {new Date(payment.created_at).toLocaleTimeString('id-ID', {hour: '2-digit', minute:'2-digit'})}
										</span>
									</div>
									<div class="text-right">
										{#if payment.status === 'pending'}
											<span class="text-[13px] text-ios-blue-dark bg-ios-blue-light/10 px-2.5 py-1 rounded-full font-medium">Menunggu</span>
										{:else}
											<span class="text-[13px] text-green-700 bg-green-100 px-2.5 py-1 rounded-full font-medium">Selesai</span>
										{/if}
									</div>
								</div>
								
								{#if payment.status === 'pending' && debt.creditor_id === authState.userId}
									<div class="mt-3 flex gap-2">
										<button onclick={() => updatePaymentStatus(payment.id, 'rejected')} class="flex-1 bg-red-100 text-red-600 font-medium text-[15px] py-2 rounded-full active:opacity-70 transition-opacity">Tolak</button>
										<button onclick={() => updatePaymentStatus(payment.id, 'confirmed')} class="flex-1 bg-[#1d1d1f] text-white font-medium text-[15px] py-2 rounded-full active:opacity-70 transition-opacity shadow-sm">Terima</button>
									</div>
								{/if}
							</li>
						{/each}
					</ul>
				</div>
			</div>
		{/if}
	</div>
{/if}
