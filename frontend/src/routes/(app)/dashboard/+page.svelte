<script lang="ts">
	import { onMount } from 'svelte';
	import { authState, API_URL } from '$lib/auth.svelte';

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
	};

	let debts: Debt[] = $state([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			const res = await fetch(`${API_URL}/debts`, {
				headers: {
					'Authorization': `Bearer ${authState.token}`
				}
			});
			if (!res.ok) throw new Error('Gagal mengambil data');
			debts = await res.json();
		} catch (err: any) {
			error = err.message;
		} finally {
			loading = false;
		}
	});

	let piutang = $derived(debts.filter(d => d.creditor_id === authState.userId && d.remaining > 0));
	let hutang = $derived(debts.filter(d => d.debtor_id === authState.userId && d.remaining > 0));
	let lunas = $derived(debts.filter(d => d.remaining <= 0));

	let totalPiutang = $derived(piutang.reduce((sum, d) => sum + d.remaining, 0));
	let totalHutang = $derived(hutang.reduce((sum, d) => sum + d.remaining, 0));

	function formatRp(val: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(val);
	}
</script>

<div class="space-y-6">
	<div class="px-4 pt-10 pb-4">
		<h1 class="text-4xl md:text-5xl font-bold tracking-tighter text-[#1d1d1f] mb-8">Ringkasan.</h1>
		<div class="grid grid-cols-2 gap-4">
			<div class="bg-ios-card-light p-4 rounded-ios-md shadow-ios">
				<h3 class="text-[13px] text-gray-500 uppercase tracking-wide font-medium">Piutang</h3>
				<p class="text-xl font-bold text-green-600 mt-1">{formatRp(totalPiutang)}</p>
			</div>
			<div class="bg-ios-card-light p-4 rounded-ios-md shadow-ios">
				<h3 class="text-[13px] text-gray-500 uppercase tracking-wide font-medium">Hutang</h3>
				<p class="text-xl font-bold text-red-600 mt-1">{formatRp(totalHutang)}</p>
			</div>
		</div>
	</div>

	{#if error}
		<p class="text-red-500">{error}</p>
	{/if}

	{#if loading}
		<p class="text-center text-gray-500 py-4">Memuat data...</p>
	{:else}
		<div class="mt-8">
			<h2 class="px-4 pb-2 text-[13px] text-gray-500 uppercase tracking-wide font-medium">Daftar Tagihan</h2>
			{#if piutang.length === 0 && hutang.length === 0}
				<div class="mx-4 bg-ios-card-light p-6 text-center rounded-ios-md shadow-ios">
					<p class="text-gray-500">Belum ada catatan.</p>
				</div>
			{:else}
				<div class="mx-4 bg-ios-card-light rounded-ios-md shadow-ios overflow-hidden">
					<ul class="divide-y divide-ios-separator-light">
						{#each debts as debt (debt.id)}
							{#if debt.remaining > 0}
								<li transition:slide>
									<a href={`/detail/${debt.id}`} class="block px-4 py-3 active:bg-gray-100 transition-colors">
										<div class="flex justify-between items-center">
											<div class="flex flex-col">
												<span class="text-[17px] text-[#1d1d1f] font-medium">
													{#if debt.creditor_id === authState.userId}
														Menagih <span class="text-ios-blue-light">{debt.debtor_name}</span>
													{:else}
														Hutang ke <span class="text-red-500 font-bold">{debt.creditor_name}</span>
													{/if}
												</span>
												{#if debt.description}
													<span class="text-[13px] text-gray-500 line-clamp-1 mt-0.5">{debt.description}</span>
												{/if}
											</div>
											<div class="flex items-center space-x-2">
												<div class="text-right flex flex-col">
													<span class="text-[17px] font-medium {debt.creditor_id === authState.userId ? 'text-[#1d1d1f]' : 'text-red-500'}">{formatRp(debt.remaining)}</span>
													<span class="text-[13px] text-gray-400">{new Date(debt.created_at).toLocaleDateString('id-ID')}</span>
												</div>
												<!-- Chevron -->
												<svg class="w-4 h-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
											</div>
										</div>
									</a>
								</li>
							{/if}
						{/each}
					</ul>
				</div>
			{/if}
		</div>

		<div class="mt-8 mb-8">
			<h2 class="px-4 pb-2 text-[13px] text-gray-500 uppercase tracking-wide font-medium">Riwayat Lunas</h2>
			{#if lunas.length === 0}
				<div class="mx-4 bg-ios-card-light p-6 text-center rounded-ios-md shadow-ios">
					<p class="text-gray-400 text-sm">Belum ada tagihan yang lunas.</p>
				</div>
			{:else}
				<div class="mx-4 bg-ios-card-light rounded-ios-md shadow-ios overflow-hidden opacity-70">
					<ul class="divide-y divide-ios-separator-light">
						{#each lunas as debt (debt.id)}
							<li transition:slide>
								<a href={`/detail/${debt.id}`} class="block px-4 py-3 active:bg-gray-100 transition-colors">
									<div class="flex justify-between items-center">
										<div class="flex flex-col">
											<span class="text-[17px] text-gray-500 line-through">
												{#if debt.creditor_id === authState.userId}
													Menagih ke {debt.debtor_name}
												{:else}
													Hutang ke {debt.creditor_name}
												{/if}
											</span>
										</div>
										<div class="flex items-center space-x-2">
											<span class="text-[17px] text-gray-400">{formatRp(debt.amount)}</span>
											<svg class="w-4 h-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
										</div>
									</div>
								</a>
							</li>
						{/each}
					</ul>
				</div>
			{/if}
		</div>
	{/if}
</div>
