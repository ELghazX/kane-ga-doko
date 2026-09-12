<script lang="ts">
	import { onMount } from 'svelte';
	import { authState, API_URL } from '$lib/auth.svelte';
	import { slide } from 'svelte/transition';

	let contacts: any[] = $state([]);
	let loading = $state(true);
	let errorMsg = $state('');
	
	let newUsername = $state('');
	let addLoading = $state(false);
	let addError = $state('');
	let addSuccess = $state('');

	onMount(() => {
		fetchContacts();
	});

	async function fetchContacts() {
		try {
			const res = await fetch(`${API_URL}/contacts`, {
				headers: { 'Authorization': `Bearer ${authState.token}` }
			});
			if (!res.ok) throw new Error('Gagal memuat kontak');
			contacts = await res.json();
		} catch (e: any) {
			errorMsg = e.message;
		} finally {
			loading = false;
		}
	}

	async function addContact(e: Event) {
		e.preventDefault();
		addLoading = true;
		addError = '';
		addSuccess = '';

		try {
			const res = await fetch(`${API_URL}/contacts`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${authState.token}`
				},
				body: JSON.stringify({ username: newUsername })
			});

			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Gagal menambahkan kontak');
			}

			addSuccess = `Kontak ${newUsername} berhasil ditambahkan!`;
			newUsername = '';
			await fetchContacts();
		} catch (err: any) {
			addError = err.message;
		} finally {
			addLoading = false;
		}
	}

	async function removeContact(id: number, name: string) {
		if (!confirm(`Hapus ${name} dari kontak?`)) return;
		try {
			const res = await fetch(`${API_URL}/contacts/${id}`, {
				method: 'DELETE',
				headers: { 'Authorization': `Bearer ${authState.token}` }
			});

			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Gagal menghapus kontak');
			}

			await fetchContacts();
		} catch (err: any) {
			alert(err.message);
		}
	}
</script>

<div class="pt-10 pb-4 space-y-8">
	<h1 class="px-4 text-4xl md:text-5xl font-bold tracking-tighter text-[#1d1d1f] mb-8">Kontak.</h1>

	<!-- Form Add Contact -->
	<div class="mx-4">
		<h2 class="pb-2 text-[13px] text-gray-500 uppercase tracking-wide font-medium">Tambah Teman</h2>
		
		{#if addError}
			<div transition:slide class="bg-red-50 text-red-600 p-3 rounded-ios-sm mb-4 text-[13px]">{addError}</div>
		{/if}
		{#if addSuccess}
			<div transition:slide class="bg-green-50 text-green-600 p-3 rounded-ios-sm mb-4 text-[13px]">{addSuccess}</div>
		{/if}

		<form onsubmit={addContact} class="bg-ios-card-light rounded-ios-md shadow-ios overflow-hidden p-4">
			<div class="flex gap-2">
				<input 
					type="text" 
					bind:value={newUsername} 
					required 
					placeholder="Masukkan username teman..."
					class="flex-1 bg-gray-100 px-3 py-2.5 rounded-lg focus:outline-none text-[#1d1d1f] text-[15px]"
				/>
				<button 
					type="submit" 
					disabled={addLoading}
					class="bg-[#1d1d1f] text-white px-5 rounded-lg font-medium text-[15px] active:opacity-70 transition-opacity disabled:opacity-50"
				>
					{addLoading ? '...' : 'Tambah'}
				</button>
			</div>
		</form>
	</div>

	<!-- Contact List -->
	<div class="mx-4">
		<h2 class="pb-2 text-[13px] text-gray-500 uppercase tracking-wide font-medium">Daftar Kontak</h2>
		
		{#if errorMsg}
			<div class="bg-red-50 text-red-600 p-3 rounded-ios-sm text-[13px]">{errorMsg}</div>
		{/if}

		{#if loading}
			<p class="text-center text-gray-500 py-4 text-sm">Memuat kontak...</p>
		{:else if contacts.length === 0}
			<p class="text-gray-400 text-sm text-center py-8 bg-ios-card-light rounded-ios-md shadow-ios">Belum ada kontak. Tambahkan teman di atas.</p>
		{:else}
			<div class="bg-ios-card-light rounded-ios-md shadow-ios overflow-hidden mb-8">
				<ul class="divide-y divide-ios-separator-light">
					{#each contacts as contact (contact.id)}
						<li transition:slide={{ duration: 300 }}>
							<div class="px-4 py-3.5 flex justify-between items-center">
								<span class="text-[17px] text-[#1d1d1f] font-medium">{contact.username}</span>
								<button onclick={() => removeContact(contact.id, contact.username)} class="text-red-500 text-[15px] font-medium active:opacity-50 transition-opacity">
									Hapus
								</button>
							</div>
						</li>
					{/each}
				</ul>
			</div>
		{/if}
	</div>
</div>
