<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api, type Block } from '$lib/api';
	import { formatDate } from '$lib';
	import TransactionRow from '$lib/components/TransactionRow.svelte';
	import CopyButton from '$lib/components/CopyButton.svelte';

	let block = $state<Block | null>(null);
	let error = $state('');

	onMount(async () => {
		const index = parseInt($page.params.index);
		try {
			block = await api.getBlock(index);
		} catch (err: any) {
			error = err.message;
		}
	});
</script>

<div class="space-y-8">
	{#if error}
		<div class="glass-card p-8 text-center">
			<p class="text-red-400">Error: {error}</p>
			<a
				href="/explorer"
				class="mt-3 inline-block text-sm text-amber-500 transition-colors hover:text-amber-400"
			>
				&larr; Volver al explorador
			</a>
		</div>
	{:else if !block}
		<div class="space-y-4">
			<div class="glass-card h-12 w-64 animate-shimmer"></div>
			<div class="glass-card h-64 animate-shimmer"></div>
		</div>
	{:else}
		<!-- Breadcrumb + Navigation -->
		<div class="flex items-center justify-between animate-fade-in">
			<div class="flex items-center gap-2 text-sm">
				<a
					href="/explorer"
					class="text-zinc-500 transition-colors hover:text-amber-400">Explorador</a
				>
				<span class="text-zinc-700">/</span>
				<span class="text-zinc-300">Bloque #{block.index}</span>
			</div>
			<div class="flex items-center gap-3">
				{#if block.index > 0}
					<a
						href="/explorer/{block.index - 1}"
						class="btn-secondary !px-3 !py-1 text-xs">&larr; Anterior</a
					>
				{/if}
				<a href="/explorer/{block.index + 1}" class="btn-secondary !px-3 !py-1 text-xs"
					>Siguiente &rarr;</a
				>
			</div>
		</div>

		<!-- Block Number Hero -->
		<div class="animate-slide-up">
			<h1 class="text-5xl font-extrabold tracking-tight text-gradient-amber">
				Bloque #{block.index}
			</h1>
		</div>

		<!-- Metadata -->
		<div class="glass-card animate-slide-up p-6 stagger-1">
			<div class="grid gap-4 text-sm">
				<div>
					<span class="section-label">Hash</span>
					<div class="mt-1 flex items-center gap-2">
						<span class="break-all font-mono text-xs text-zinc-300">{block.hash}</span>
						<CopyButton text={block.hash} />
					</div>
				</div>
				<div>
					<span class="section-label">Hash Anterior</span>
					<div class="mt-1 flex items-center gap-2">
						<span class="break-all font-mono text-xs text-zinc-400">{block.prevHash}</span>
						<CopyButton text={block.prevHash} />
					</div>
				</div>
				<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
					<div>
						<span class="section-label">Fecha</span>
						<p class="mt-1 text-zinc-300">{formatDate(block.timestamp)}</p>
					</div>
					<div>
						<span class="section-label">Minero</span>
						<p class="mt-1 font-mono text-xs text-zinc-300">
							{block.miner || 'Genesis'}
						</p>
					</div>
					<div>
						<span class="section-label">Nonce</span>
						<p class="mt-1 text-zinc-300">{block.nonce.toLocaleString('es-AR')}</p>
					</div>
					<div>
						<span class="section-label">Transacciones</span>
						<p class="mt-1 text-zinc-300">{block.transactions.length}</p>
					</div>
				</div>
			</div>
		</div>

		<!-- Transactions -->
		{#if block.transactions.length > 0}
			<div class="animate-slide-up stagger-2">
				<h2 class="mb-4 text-xl font-bold text-zinc-100">Transacciones</h2>
				<div class="space-y-2">
					{#each block.transactions as tx (tx.id)}
						<TransactionRow {tx} />
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>
