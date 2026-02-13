<script lang="ts">
	import type { Transaction } from '$lib/api';
	import { COINBASE_ADDRESS, formatAmount, shortAddr } from '$lib';

	let { tx }: { tx: Transaction } = $props();
</script>

<a
	href="/tx/{tx.id}"
	class="glass-card glass-card-interactive flex items-center justify-between px-5 py-3.5 group"
>
	<div class="flex flex-col gap-1.5 min-w-0">
		<span class="truncate font-mono text-xs text-zinc-500 max-w-[200px]"
			>{tx.id.slice(0, 16)}...</span
		>
		<div class="flex items-center gap-2 text-xs">
			{#if tx.sender === COINBASE_ADDRESS}
				<span
					class="rounded-full bg-emerald-500/10 px-2 py-0.5 text-[10px] font-semibold text-emerald-400"
					>Coinbase</span
				>
			{:else}
				<span class="font-mono text-zinc-400">{shortAddr(tx.sender, 8)}</span>
			{/if}
			<svg
				class="h-3 w-3 flex-shrink-0 text-zinc-600"
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
			<span class="font-mono text-zinc-400">{shortAddr(tx.receiver, 8)}</span>
		</div>
	</div>
	<span class="text-sm font-semibold text-amber-400 whitespace-nowrap ml-4"
		>{formatAmount(tx.amount)}</span
	>
</a>
