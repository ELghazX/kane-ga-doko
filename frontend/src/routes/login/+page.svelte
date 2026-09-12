<script lang="ts">
	import { goto } from '$app/navigation';
	import { login, API_URL } from '$lib/auth.svelte';
	import { onMount } from 'svelte';

	let username = $state('');
	let password = $state('');
	let errorMsg = $state('');
	let successMsg = $state('');
	let loading = $state(false);
	
	let isRegister = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		errorMsg = '';
		successMsg = '';
		loading = true;

		try {
			if (isRegister) {
				const res = await fetch(`${API_URL}/register`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ username, password })
				});

				if (!res.ok) {
					const text = await res.text();
					throw new Error(text || 'Daftar gagal');
				}

				successMsg = 'Pendaftaran berhasil. Silakan tunggu Admin menyetujui akun Anda sebelum bisa masuk.';
				isRegister = false;
				password = '';
			} else {
				const res = await fetch(`${API_URL}/login`, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ username, password })
				});

				if (!res.ok) {
					const text = await res.text();
					throw new Error(text || 'Login gagal');
				}

				const data = await res.json();
				login(data);
				goto('/dashboard');
			}
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			loading = false;
		}
	}
</script>

<div class="pt-10 pb-4 space-y-8">
	<h1 class="px-4 text-4xl md:text-5xl font-bold tracking-tighter text-[#1d1d1f] text-center mb-6">
		{isRegister ? 'Daftar.' : 'Masuk.'}
	</h1>

	<!-- Tab Switcher -->
	<div class="mx-4 bg-gray-200 p-1 rounded-full flex text-[15px] font-medium">
		<button 
			class="flex-1 py-1.5 rounded-full transition-colors {!isRegister ? 'bg-white shadow-sm text-[#1d1d1f]' : 'text-gray-500'}"
			onclick={() => { isRegister = false; errorMsg = ''; successMsg = ''; }}
		>
			Masuk
		</button>
		<button 
			class="flex-1 py-1.5 rounded-full transition-colors {isRegister ? 'bg-white shadow-sm text-[#1d1d1f]' : 'text-gray-500'}"
			onclick={() => { isRegister = true; errorMsg = ''; successMsg = ''; }}
		>
			Daftar
		</button>
	</div>

	{#if errorMsg}
		<div class="mx-4 bg-red-50 text-red-600 p-3 rounded-ios-sm mb-4 text-[13px] text-center">{errorMsg}</div>
	{/if}

	{#if successMsg}
		<div class="mx-4 bg-green-50 text-green-600 p-3 rounded-ios-sm mb-4 text-[13px] text-center">{successMsg}</div>
	{/if}

	<form onsubmit={handleSubmit} class="space-y-6">
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
							placeholder="Masukkan username"
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
							placeholder="Masukkan password"
						/>
					</div>
				</div>
			</div>
		</div>

		<div class="mx-4 mt-8">
			<button 
				type="submit" 
				disabled={loading}
				class="w-full bg-[#1d1d1f] text-white font-medium text-[17px] py-3.5 rounded-full active:opacity-70 transition-opacity disabled:opacity-50 shadow-ios"
			>
				{loading ? 'Memproses...' : (isRegister ? 'Buat Akun' : 'Masuk')}
			</button>
		</div>
	</form>
</div>
