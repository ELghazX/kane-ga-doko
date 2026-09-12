	<script lang="ts">
	import { authState, API_URL } from '$lib/auth.svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { slide, fade } from 'svelte/transition';

	let username = $state('');
	let password = $state('');
	let loading = $state(false);
	let error = $state('');
	let successMsg = $state('');
	let users: any[] = $state([]);
	let activeUsers = $derived(users.filter(u => u.is_approved));
	let pendingUsers = $derived(users.filter(u => !u.is_approved));

	let editingUserId = $state<number | null>(null);
	let editUsername = $state('');
	let editPassword = $state('');

	onMount(() => {
		if (!authState.isAdmin) {
			goto('/dashboard');
		} else {
			fetchUsers();
		}
	});

	async function fetchUsers() {
		try {
			const res = await fetch(`${API_URL}/users`, {
				headers: { 'Authorization': `Bearer ${authState.token}` }
			});
			if (res.ok) {
				users = await res.json();
			}
		} catch (e) {
			console.error(e);
		}
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		loading = true;
		error = '';
		successMsg = '';

		try {
			const res = await fetch(`${API_URL}/users`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${authState.token}`
				},
				body: JSON.stringify({ username, password })
			});

			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Gagal membuat pengguna');
			}

			successMsg = `Pengguna ${username} berhasil dibuat!`;
			username = '';
			password = '';
			await fetchUsers();
		} catch (err: any) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	function startEdit(u: any) {
		editingUserId = u.id;
		editUsername = u.username;
		editPassword = '';
	}

	function cancelEdit() {
		editingUserId = null;
		editUsername = '';
		editPassword = '';
	}

	async function saveEdit(id: number) {
		try {
			const res = await fetch(`${API_URL}/users/${id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${authState.token}`
				},
				body: JSON.stringify({ username: editUsername, password: editPassword })
			});

			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Gagal mengubah pengguna');
			}
			
			cancelEdit();
			await fetchUsers();
		} catch (err: any) {
			alert(err.message);
		}
	}

	async function handleDelete(id: number, uname: string) {
		if (!confirm(`Hapus pengguna ${uname}? Pastikan ia tidak memiliki piutang/hutang.`)) return;
		try {
			const res = await fetch(`${API_URL}/users/${id}`, {
				method: 'DELETE',
				headers: { 'Authorization': `Bearer ${authState.token}` }
			});

			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Gagal menghapus pengguna');
			}

			await fetchUsers();
		} catch (err: any) {
			alert(err.message);
		}
	}

	async function handleApprove(id: number, uname: string) {
		if (!confirm(`Setujui pengguna ${uname} untuk masuk ke sistem?`)) return;
		try {
			const res = await fetch(`${API_URL}/users/${id}/approve`, {
				method: 'PUT',
				headers: { 'Authorization': `Bearer ${authState.token}` }
			});

			if (!res.ok) {
				const text = await res.text();
				throw new Error(text || 'Gagal menyetujui pengguna');
			}

			await fetchUsers();
			successMsg = `Pengguna ${uname} berhasil disetujui!`;
		} catch (err: any) {
			alert(err.message);
		}
	}
</script>

<div class="pt-10 pb-4 space-y-8">
	<h1 class="px-4 text-4xl md:text-5xl font-bold tracking-tighter text-[#1d1d1f]">Manajemen Pengguna.</h1>
	
	<p class="px-4 text-[15px] text-gray-500">Hanya Admin yang dapat mengelola pengguna di sistem ini.</p>

	{#if error}
		<div transition:slide class="mx-4 bg-red-50 text-red-600 p-3 rounded-ios-sm mb-4 text-[13px]">{error}</div>
	{/if}

	{#if successMsg}
		<div transition:slide class="mx-4 bg-green-50 text-green-600 p-3 rounded-ios-sm mb-4 text-[13px]">{successMsg}</div>
	{/if}

	<div class="mx-4">
		{#if pendingUsers.length > 0}
			<h2 class="pb-2 text-[13px] text-orange-500 uppercase tracking-wide font-medium">Menunggu Persetujuan</h2>
			<div class="bg-ios-card-light rounded-ios-md shadow-ios overflow-hidden mb-8 border border-orange-100">
				<ul class="divide-y divide-ios-separator-light">
					{#each pendingUsers as u (u.id)}
						<li transition:slide={{ duration: 300 }} class="px-4 py-3 flex justify-between items-center">
							<span class="text-[17px] text-[#1d1d1f] font-medium">{u.username}</span>
							<div class="flex gap-3">
								<button onclick={() => handleDelete(u.id, u.username)} class="text-red-500 text-[15px] font-medium active:opacity-50 transition">Tolak (Hapus)</button>
								<button onclick={() => handleApprove(u.id, u.username)} class="bg-green-500 text-white px-3 py-1 rounded-full text-[13px] font-bold active:opacity-50 transition shadow-sm">Setujui</button>
							</div>
						</li>
					{/each}
				</ul>
			</div>
		{/if}

		<h2 class="pb-2 text-[13px] text-gray-500 uppercase tracking-wide font-medium">Pengguna Aktif</h2>
		<div class="bg-ios-card-light rounded-ios-md shadow-ios overflow-hidden mb-8">
			<ul class="divide-y divide-ios-separator-light">
				{#each activeUsers as u (u.id)}
					<li transition:slide={{ duration: 300 }} class="px-4 py-3 flex flex-col gap-2">
						{#if editingUserId === u.id}
							<div class="flex flex-col gap-2">
								<input type="text" bind:value={editUsername} class="w-full bg-gray-50 p-2 rounded border focus:outline-none" placeholder="Username" />
								<input type="password" bind:value={editPassword} class="w-full bg-gray-50 p-2 rounded border focus:outline-none" placeholder="Password baru (opsional)" />
								<div class="flex gap-2 justify-end mt-1">
									<button onclick={cancelEdit} class="text-gray-500 text-sm px-3 py-1">Batal</button>
									<button onclick={() => saveEdit(u.id)} class="bg-[#1d1d1f] text-white rounded px-3 py-1 text-sm font-medium shadow-sm">Simpan</button>
								</div>
							</div>
						{:else}
							<div class="flex justify-between items-center">
								<div>
									<span class="text-[17px] text-[#1d1d1f] font-medium">{u.username}</span>
									{#if u.is_admin}
										<span class="ml-2 text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full font-bold">Admin</span>
									{/if}
								</div>
								<div class="flex gap-3">
									<button onclick={() => startEdit(u)} class="text-ios-blue-light text-[15px] font-medium active:opacity-50 transition">Ubah</button>
									{#if !u.is_admin}
										<button onclick={() => handleDelete(u.id, u.username)} class="text-red-500 text-[15px] font-medium active:opacity-50 transition">Hapus</button>
									{/if}
								</div>
							</div>
						{/if}
					</li>
				{/each}
			</ul>
		</div>
	</div>

	<form onsubmit={handleSubmit} class="space-y-6">
		<div>
			<h2 class="px-4 pb-2 text-[13px] text-gray-500 uppercase tracking-wide font-medium">Tambah Pengguna Baru</h2>
			<div class="mx-4 bg-ios-card-light rounded-ios-md shadow-ios overflow-hidden">
				<div class="divide-y divide-ios-separator-light">
					<div class="flex items-center px-4 py-3 bg-white">
						<label for="username" class="w-1/3 text-[17px] text-[#1d1d1f] font-medium">Username</label>
						<div class="w-2/3 flex items-center justify-end">
							<input 
								type="text" 
								id="username" 
								bind:value={username} 
								required 
								class="w-full bg-transparent text-[17px] text-right focus:outline-none text-[#1d1d1f]"
								placeholder="Username"
							/>
						</div>
					</div>
					
					<div class="flex items-center px-4 py-3 bg-white">
						<label for="password" class="w-1/3 text-[17px] text-[#1d1d1f] font-medium">Password</label>
						<div class="w-2/3 flex items-center justify-end">
							<input 
								type="password" 
								id="password" 
								bind:value={password} 
								required 
								class="w-full bg-transparent text-[17px] text-right focus:outline-none text-[#1d1d1f]"
								placeholder="Password awal"
							/>
						</div>
					</div>
				</div>
			</div>
		</div>

		<div class="mx-4 mt-8 pb-8">
			<button 
				type="submit" 
				disabled={loading}
				class="w-full bg-[#1d1d1f] text-white font-medium text-[17px] py-3.5 rounded-full active:opacity-70 transition-opacity disabled:opacity-50 shadow-ios"
			>
				{loading ? 'Menyimpan...' : 'Buat Pengguna'}
			</button>
		</div>
	</form>
</div>
