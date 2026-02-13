<script lang="ts">
	import type { Transaction } from '$lib/api';
	import { formatAmount, shortAddr } from '$lib';

	let { transactions }: { transactions: Transaction[] } = $props();
</script>

{#if transactions.length === 0}
	<div class="flex items-center justify-center py-8 text-sm text-zinc-600">
		No hay transacciones pendientes
	</div>
{:else}
	<div class="max-h-52 space-y-2 overflow-y-auto pr-1">
		{#each transactions.slice(0, 10) as tx, i (tx.id)}
			<a
				href="/tx/{tx.id}"
				class="flex items-center justify-between rounded-lg bg-white/[0.02] px-3 py-2.5 transition-colors hover:bg-white/[0.04] animate-fade-in"
				style="animation-delay: {i * 50}ms"
			>
				<div class="flex items-center gap-2 text-xs">
					<span class="status-dot-amber" style="width:6px;height:6px;"></span>
					<span class="font-mono text-zinc-500">{shortAddr(tx.sender, 6)}</span>
					<svg
						class="h-3 w-3 flex-shrink-0 text-zinc-700"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
					>
						<line x1="5" y1="12" x2="19" y2="12" />
						<polyline points="12 5 19 12 12 19" />
					</svg>
					<span class="font-mono text-zinc-500">{shortAddr(tx.receiver, 6)}</span>
				</div>
				<span class="ml-3 whitespace-nowrap text-xs font-semibold text-amber-400"
					>{formatAmount(tx.amount)}</span
				>
			</a>
		{/each}
	</div>
{/if}
