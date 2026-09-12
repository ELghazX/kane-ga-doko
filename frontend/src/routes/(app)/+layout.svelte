<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { authState, logout } from '$lib/auth.svelte';
	import { Home, PlusCircle, Users, LogOut, Contact } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';

	let { children } = $props();

	onMount(() => {
		if (!authState.token) {
			goto('/login');
		}
	});

	function handleLogout() {
		logout();
		goto('/login');
	}
</script>

{#if authState.token}
	{#key page.url.pathname}
		<div in:fly={{ y: 20, duration: 400, delay: 150 }} out:fade={{ duration: 150 }}>
			{@render children()}
		</div>
	{/key}

	<nav class="fixed bottom-0 left-0 right-0 bg-white/70 backdrop-blur-md border-t border-ios-separator-light flex justify-around p-3 pb-safe shadow-[0_-4px_24px_rgba(0,0,0,0.02)] z-20">
		<a href="/dashboard" class="flex flex-col items-center transition-colors duration-300 {page.url.pathname === '/dashboard' ? 'text-[#1d1d1f]' : 'text-gray-400'}">
			<Home size={24} strokeWidth={page.url.pathname === '/dashboard' ? 2.5 : 2} class="transition-all duration-300" />
			<span class="text-[11px] font-medium mt-1 transition-all duration-300">Beranda</span>
		</a>
		<a href="/kontak" class="flex flex-col items-center transition-colors duration-300 {page.url.pathname === '/kontak' ? 'text-[#1d1d1f]' : 'text-gray-400'}">
			<Contact size={24} strokeWidth={page.url.pathname === '/kontak' ? 2.5 : 2} class="transition-all duration-300" />
			<span class="text-[11px] font-medium mt-1 transition-all duration-300">Kontak</span>
		</a>
		<a href="/tambah" class="flex flex-col items-center transition-colors duration-300 {page.url.pathname === '/tambah' ? 'text-[#1d1d1f]' : 'text-gray-400'}">
			<PlusCircle size={24} strokeWidth={page.url.pathname === '/tambah' ? 2.5 : 2} class="transition-all duration-300" />
			<span class="text-[11px] font-medium mt-1 transition-all duration-300">Catat</span>
		</a>
		{#if authState.isAdmin}
		<a href="/admin" class="flex flex-col items-center transition-colors duration-300 {page.url.pathname === '/admin' ? 'text-[#1d1d1f]' : 'text-gray-400'}">
			<Users size={24} strokeWidth={page.url.pathname === '/admin' ? 2.5 : 2} class="transition-all duration-300" />
			<span class="text-[11px] font-medium mt-1 transition-all duration-300">Admin</span>
		</a>
		{/if}
		<button onclick={handleLogout} class="flex flex-col items-center text-red-500 transition-colors duration-300 active:opacity-50">
			<LogOut size={24} />
			<span class="text-[11px] font-medium mt-1">Keluar</span>
		</button>
	</nav>
{/if}
