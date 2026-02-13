<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api, type TxResult } from '$lib/api';
	import { COINBASE_ADDRESS, formatAmount, formatDate } from '$lib';
	import CopyButton from '$lib/components/CopyButton.svelte';

	let result = $state<TxResult | null>(null);
	let error = $state('');
	let loading = $state(true);

	onMount(async () => {
		try {
			result = await api.getTransaction($page.params.id);
		} catch (err: any) {
			error = err.message;
		} finally {
			loading = false;
		}
	});
</script>

<div class="space-y-8">
	<div class="flex items-center gap-2 text-sm animate-fade-in">
		<a href="/explorer" class="text-zinc-500 transition-colors hover:text-amber-400">Explorador</a
		>
		<span class="text-zinc-700">/</span>
		<span class="text-zinc-300">Transaccion</span>
	</div>

	{#if loading}
		<div class="space-y-4">
			<div class="glass-card h-12 w-64 animate-shimmer"></div>
			<div class="glass-card h-64 animate-shimmer"></div>
		</div>
	{:else if error}
		<div class="glass-card animate-slide-up p-8 text-center">
			<div
				class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-red-500/10"
			>
				<svg
					class="h-8 w-8 text-red-400"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.5"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<circle cx="12" cy="12" r="10" />
					<line x1="15" y1="9" x2="9" y2="15" />
					<line x1="9" y1="9" x2="15" y2="15" />
				</svg>
			</div>
			<p class="text-lg font-semibold text-zinc-200">Transaccion no encontrada</p>
			<p class="mt-2 break-all font-mono text-xs text-zinc-600">{$page.params.id}</p>
			<p class="mt-3 text-sm text-zinc-500">
				Esta transaccion puede estar aun en el mempool (pendiente).
			</p>
			<a
				href="/explorer"
				class="mt-4 inline-block text-sm text-amber-500 transition-colors hover:text-amber-400"
			>
				&larr; Explorar la blockchain
			</a>
		</div>
	{:else if result}
		{@const tx = result.transaction}

		<h1 class="text-3xl font-bold tracking-tight text-zinc-100 animate-fade-in">Transaccion</h1>

		<!-- Visual Flow: From → Amount → To -->
		<div class="animate-slide-up flex flex-col items-center gap-4 md:flex-row md:gap-6">
			<!-- From -->
			<div class="glass-card flex-1 p-5 text-center w-full">
				<div class="section-label mb-2">De</div>
				{#if tx.sender === COINBASE_ADDRESS}
					<span
						class="inline-block rounded-full bg-emerald-500/10 px-3 py-1 text-sm font-semibold text-emerald-400"
					>
						Coinbase (Recompensa)
					</span>
				{:else}
					<p class="break-all font-mono text-xs text-zinc-300">{tx.sender}</p>
				{/if}
			</div>

			<!-- Arrow + Amount -->
			<div class="flex flex-col items-center gap-1 flex-shrink-0">
				<span class="text-lg font-bold text-amber-400">{formatAmount(tx.amount)}</span>
				<svg
					class="h-6 w-6 text-zinc-600 md:rotate-0 rotate-90"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<line x1="5" y1="12" x2="19" y2="12" />
					<polyline points="12 5 19 12 12 19" />
				</svg>
			</div>

			<!-- To -->
			<div class="glass-card flex-1 p-5 text-center w-full">
				<div class="section-label mb-2">Para</div>
				<p class="break-all font-mono text-xs text-zinc-300">{tx.receiver}</p>
			</div>
		</div>

		<!-- Details -->
		<div class="glass-card animate-slide-up stagger-1 p-6">
			<div class="grid gap-4 text-sm">
				<div>
					<span class="section-label">ID de Transaccion</span>
					<div class="mt-1 flex items-center gap-2">
						<span class="break-all font-mono text-xs text-zinc-300">{tx.id}</span>
						<CopyButton text={tx.id} />
					</div>
				</div>

				<div class="grid gap-4 sm:grid-cols-3">
					<div>
						<span class="section-label">Monto</span>
						<p class="mt-1 font-semibold text-amber-400">{formatAmount(tx.amount)}</p>
					</div>
					<div>
						<span class="section-label">Comision</span>
						<p class="mt-1 text-zinc-300">{formatAmount(tx.fee)}</p>
					</div>
					<div>
						<span class="section-label">Nonce</span>
						<p class="mt-1 text-zinc-300">{tx.nonce}</p>
					</div>
				</div>

				<div>
					<span class="section-label">Bloque</span>
					<div class="mt-1 flex items-center gap-2">
						<a
							href="/explorer/{result.blockIndex}"
							class="text-amber-500 transition-colors hover:text-amber-400"
						>
							#{result.blockIndex}
						</a>
						<span class="font-mono text-xs text-zinc-600">{result.blockHash}</span>
					</div>
				</div>

				{#if tx.pubKey}
					<div>
						<span class="section-label">Clave Publica</span>
						<p class="mt-1 break-all font-mono text-xs text-zinc-400">{tx.pubKey}</p>
					</div>
				{/if}

				{#if tx.signature}
					<div>
						<span class="section-label">Firma</span>
						<p class="mt-1 break-all font-mono text-xs text-zinc-400">{tx.signature}</p>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
