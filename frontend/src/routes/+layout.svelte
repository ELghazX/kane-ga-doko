<script lang="ts">
	import '../app.css';
	import { page } from '$app/state';
	import { fade, fly } from 'svelte/transition';
	let { children } = $props();
</script>

<div class="min-h-screen bg-ios-bg-light text-[#1d1d1f] font-sans pb-20">
	<header class="bg-[#1d1d1f] text-white p-4 shadow-sm sticky top-0 z-10 flex justify-center items-center">
		<h1 class="text-xl font-bold tracking-tight">Kane ga doko?.</h1>
	</header>
	<main class="max-w-md mx-auto overflow-hidden">
		{@render children()}
	</main>
</div>
