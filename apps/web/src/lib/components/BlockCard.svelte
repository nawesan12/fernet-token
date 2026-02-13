<script lang="ts">
	import type { Block } from '$lib/api';
	import { timeAgo } from '$lib';

	let { block }: { block: Block } = $props();
</script>

<a href="/explorer/{block.index}" class="glass-card glass-card-interactive block p-5 group">
	<div class="flex items-start justify-between">
		<div class="flex items-center gap-3">
			<div
				class="flex h-10 w-10 items-center justify-center rounded-xl bg-amber-500/10 text-sm font-bold text-amber-400 transition-colors group-hover:bg-amber-500/20"
			>
				#{block.index}
			</div>
			<div>
				<div
					class="text-sm font-semibold text-zinc-200 transition-colors group-hover:text-zinc-50"
				>
					Bloque #{block.index}
				</div>
				<div class="mt-0.5 text-xs text-zinc-500">{timeAgo(block.timestamp)}</div>
			</div>
		</div>
		<span class="rounded-full bg-white/[0.05] px-2.5 py-0.5 text-xs text-zinc-400">
			{block.transactions.length} tx{block.transactions.length !== 1 ? 's' : ''}
		</span>
	</div>

	<div class="mt-3 truncate font-mono text-xs text-zinc-600">
		{block.hash}
	</div>

	{#if block.miner}
		<div class="mt-2 flex items-center gap-1.5 text-xs text-zinc-500">
			<svg
				class="h-3 w-3"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<path
					d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"
				/>
			</svg>
			<span class="font-mono">{block.miner.slice(0, 10)}...</span>
		</div>
	{/if}
</a>
