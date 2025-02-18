<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<svelte:head>
	<title>Fernet Token</title>
</svelte:head>

<main class="mx-auto w-full p-4 lg:max-w-7xl lg:px-0">
	{@render children()}
</main>
