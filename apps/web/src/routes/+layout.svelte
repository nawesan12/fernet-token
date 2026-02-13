<script lang="ts">
	import '../app.css';
	import Navbar from '$lib/components/Navbar.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import { initWallet } from '$lib/stores/wallet';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';

	let { children } = $props();

	onMount(() => {
		initWallet();
	});

	let currentPath = $derived($page.url.pathname);
</script>

<svelte:head>
	<title>Fernet Token</title>
</svelte:head>

<div class="flex min-h-screen flex-col">
	<Navbar />
	<main class="mx-auto w-full max-w-7xl flex-1 px-6 py-10">
		{#key currentPath}
			<div class="animate-fade-in">
				{@render children()}
			</div>
		{/key}
	</main>
	<Footer />
</div>
