<script lang="ts">
	import { animateValue } from '$lib';

	let {
		label,
		value,
		format = 'number',
		icon,
		accent = false
	}: {
		label: string;
		value: number | string;
		format?: 'number' | 'text';
		icon: 'blocks' | 'pending' | 'network' | 'wallet';
		accent?: boolean;
	} = $props();

	let displayValue = $state(0);
	let prevValue = $state(0);

	$effect(() => {
		if (format === 'number' && typeof value === 'number') {
			animateValue(prevValue, value, 600, (v) => {
				displayValue = v;
			});
			prevValue = value;
		}
	});
</script>

<div class="glass-card p-5 transition-all duration-300 hover:border-amber-500/20">
	<div class="mb-3 flex items-center gap-3">
		<div
			class="flex h-8 w-8 items-center justify-center rounded-lg {accent
				? 'bg-amber-500/15'
				: 'bg-white/[0.05]'}"
		>
			{#if icon === 'blocks'}
				<svg
					class="h-4 w-4 {accent ? 'text-amber-400' : 'text-zinc-400'}"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<rect x="3" y="3" width="7" height="7" />
					<rect x="14" y="3" width="7" height="7" />
					<rect x="14" y="14" width="7" height="7" />
					<rect x="3" y="14" width="7" height="7" />
				</svg>
			{:else if icon === 'pending'}
				<svg
					class="h-4 w-4 text-zinc-400"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<circle cx="12" cy="12" r="10" />
					<polyline points="12 6 12 12 16 14" />
				</svg>
			{:else if icon === 'network'}
				<svg
					class="h-4 w-4 text-zinc-400"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<circle cx="12" cy="5" r="3" />
					<circle cx="5" cy="19" r="3" />
					<circle cx="19" cy="19" r="3" />
					<line x1="12" y1="8" x2="5" y2="16" />
					<line x1="12" y1="8" x2="19" y2="16" />
				</svg>
			{:else if icon === 'wallet'}
				<svg
					class="h-4 w-4 text-zinc-400"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M21 12V7H5a2 2 0 0 1 0-4h14v4" />
					<path d="M3 5v14a2 2 0 0 0 2 2h16v-5" />
					<path d="M18 12a2 2 0 0 0 0 4h4v-4z" />
				</svg>
			{/if}
		</div>
		<span class="section-label">{label}</span>
	</div>
	<div
		class="text-2xl font-bold tracking-tight {accent ? 'text-gradient-amber' : 'text-zinc-100'}"
	>
		{#if format === 'number'}
			{displayValue.toLocaleString('es-AR')}
		{:else}
			{value}
		{/if}
	</div>
</div>
