<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authState, API_URL } from '$lib/auth.svelte';

	type User = {
		id: number;
		username: string;
		is_admin: boolean;
	};

	let users: User[] = $state([]);
	let selectedUser = $state('');
	let amount = $state('');
	let description = $state('');
	
	let loading = $state(false);
	let error = $state('');

	onMount(async () => {
		try {
			const res = await fetch(`${API_URL}/contacts`, {
				headers: { 'Authorization': `Bearer ${authState.token}` }
			});
			if (!res.ok) throw new Error('Gagal memuat kontak');
			const data = await res.json();
			users = data;
		} catch (e: any) {
			error = e.message;
		}
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		loading = true;
		error = '';

		try {
			const res = await fetch(`${API_URL}/debts`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${authState.token}`
				},
				body: JSON.stringify({
					debtor_id: Number(selectedUser),
					amount: Number(amount),
					description
				})
			});

			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Gagal menyimpan data');
			}

			goto('/dashboard');
		} catch (err: any) {
			error = err.message;
			loading = false;
		}
	}
</script>

<div class="pt-10 pb-4 space-y-8">
	<h1 class="px-4 text-4xl md:text-5xl font-bold tracking-tighter text-[#1d1d1f]">Catat Piutang.</h1>

	{#if error}
		<div class="mx-4 bg-red-50 text-red-600 p-3 rounded-ios-sm mb-4 text-[13px]">{error}</div>
	{/if}

	<form onsubmit={handleSubmit} class="space-y-6">
		<div>
			<h2 class="px-4 pb-2 text-[13px] text-gray-500 uppercase tracking-wide font-medium">Detail Pinjaman</h2>
			<div class="mx-4 bg-ios-card-light rounded-ios-md shadow-ios overflow-hidden">
				<div class="divide-y divide-ios-separator-light">
					<div class="flex items-center px-4 py-3 bg-white">
						<label for="user" class="w-1/3 text-[17px] text-[#1d1d1f] font-medium">Peminjam</label>
						<select 
							id="user" 
							bind:value={selectedUser} 
							required 
							class="w-2/3 bg-transparent text-[17px] text-ios-blue-light text-right focus:outline-none appearance-none"
						>
							<option value="" disabled>Pilih pengguna</option>
							{#each users as u}
								<option value={u.id}>{u.username}</option>
							{/each}
						</select>
					</div>

					<div class="flex items-center px-4 py-3 bg-white">
						<label for="amount" class="w-1/3 text-[17px] text-[#1d1d1f] font-medium">Nominal</label>
						<div class="w-2/3 flex items-center justify-end">
							<span class="text-gray-400 mr-1">Rp</span>
							<input 
								type="number" 
								id="amount" 
								bind:value={amount} 
								required 
								min="1"
								class="w-full bg-transparent text-[17px] text-right focus:outline-none text-ios-blue-dark font-medium"
								placeholder="50000"
							/>
						</div>
					</div>

					<div class="flex flex-col px-4 py-3 bg-white">
						<label for="desc" class="text-[17px] text-[#1d1d1f] font-medium mb-2">Keterangan</label>
						<textarea 
							id="desc" 
							bind:value={description}
							required
							rows="2"
							class="w-full bg-transparent text-[17px] text-gray-500 focus:outline-none resize-none"
							placeholder="Makan siang kemarin..."
						></textarea>
					</div>
				</div>
			</div>
		</div>

		<div class="mx-4 mt-8">
			<button 
				type="submit" 
				disabled={loading || !selectedUser}
				class="w-full bg-[#1d1d1f] text-white font-medium text-[17px] py-3.5 rounded-full active:opacity-70 transition-opacity disabled:opacity-50 shadow-ios"
			>
				{loading ? 'Menyimpan...' : 'Simpan Piutang'}
			</button>
		</div>
	</form>
</div>
