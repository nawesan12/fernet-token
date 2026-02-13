<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Block } from '$lib/api';
	import BlockCard from '$lib/components/BlockCard.svelte';

	let allBlocks = $state<Block[]>([]);
	let loading = $state(true);
	let currentPage = $state(1);
	const perPage = 12;

	$effect(() => {
		if (allBlocks.length > 0 && currentPage > totalPages) {
			currentPage = totalPages;
		}
	});

	let totalPages = $derived(Math.max(1, Math.ceil(allBlocks.length / perPage)));
	let displayedBlocks = $derived(
		allBlocks.slice((currentPage - 1) * perPage, currentPage * perPage)
	);

	let lookupAddr = $state('');
	let lookupResult = $state<{ balance: number; formatted: string } | null>(null);
	let lookupError = $state('');

	onMount(async () => {
		try {
			const res = await api.getBlockchain();
			allBlocks = (res.chain || []).slice().reverse();
		} catch (err) {
			console.error('Failed to load blockchain:', err);
		} finally {
			loading = false;
		}
	});

	async function lookupAddress() {
		if (!lookupAddr.trim()) return;
		lookupError = '';
		lookupResult = null;
		try {
			lookupResult = await api.getBalance(lookupAddr.trim());
		} catch (err: any) {
			lookupError = err.message;
		}
	}
</script>

<div class="space-y-8">
	<div class="animate-fade-in">
		<h1 class="text-3xl font-bold tracking-tight text-zinc-100">Explorador de Bloques</h1>
		<p class="mt-1 text-sm text-zinc-500">
			Explora todos los bloques y busca balances por direccion.
		</p>
	</div>

	<!-- Search -->
	<div class="glass-card animate-slide-up p-5">
		<div class="flex gap-3">
			<div class="relative flex-1">
				<svg
					class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-zinc-600"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<circle cx="11" cy="11" r="8" />
					<line x1="21" y1="21" x2="16.65" y2="16.65" />
				</svg>
				<input
					type="text"
					bind:value={lookupAddr}
					placeholder="Buscar balance por direccion..."
					class="input-field pl-10 font-mono text-xs"
					onkeydown={(e) => {
						if (e.key === 'Enter') lookupAddress();
					}}
				/>
			</div>
			<button onclick={lookupAddress} class="btn-primary">Buscar</button>
		</div>
		{#if lookupResult}
			<div class="mt-3 animate-fade-in text-sm">
				<span class="text-zinc-500">Balance:</span>
				<span class="ml-2 font-semibold text-amber-400">{lookupResult.formatted}</span>
			</div>
		{/if}
		{#if lookupError}
			<p class="mt-2 text-xs text-red-400 animate-fade-in">{lookupError}</p>
		{/if}
	</div>

	{#if loading}
		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
			{#each Array(6) as _, i}
				<div
					class="glass-card h-32 animate-shimmer"
					style="animation-delay: {i * 100}ms"
				></div>
			{/each}
		</div>
	{:else if allBlocks.length === 0}
		<div class="glass-card p-12 text-center">
			<p class="text-zinc-500">No se encontraron bloques.</p>
		</div>
	{:else}
		<!-- Pagination bar -->
		<div class="flex items-center justify-between">
			<p class="text-sm text-zinc-500">{allBlocks.length} bloques en total</p>
			<div class="flex items-center gap-2">
				<button
					onclick={() => {
						if (currentPage > 1) currentPage--;
					}}
					disabled={currentPage <= 1}
					class="btn-secondary !px-3 !py-1 text-xs"
				>
					Anterior
				</button>
				<span class="text-xs text-zinc-500">Pagina {currentPage} de {totalPages}</span>
				<button
					onclick={() => {
						if (currentPage < totalPages) currentPage++;
					}}
					disabled={currentPage >= totalPages}
					class="btn-secondary !px-3 !py-1 text-xs"
				>
					Siguiente
				</button>
			</div>
		</div>

		<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
			{#each displayedBlocks as block, i (block.index)}
				<div class="animate-slide-up" style="animation-delay: {Math.min(i, 5) * 50}ms">
					<BlockCard {block} />
				</div>
			{/each}
		</div>
	{/if}
</div>
