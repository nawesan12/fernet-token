<script lang="ts">
	let {
		mining,
		mineSuccess,
		onclick
	}: { mining: boolean; mineSuccess: boolean; onclick: () => void } = $props();
</script>

<button
	{onclick}
	disabled={mining}
	class="btn-primary w-full py-4 text-base
		{mineSuccess ? '!bg-gradient-to-r !from-emerald-500 !to-emerald-600 !shadow-emerald-500/20' : ''}"
>
	{#if mining}
		<svg
			class="h-5 w-5 animate-spin-slow"
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
		Minando...
	{:else if mineSuccess}
		<svg
			class="h-5 w-5"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2.5"
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<polyline points="20 6 9 17 4 12" />
		</svg>
		Bloque Minado!
	{:else}
		<svg
			class="h-5 w-5"
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
		Minar Bloque
	{/if}
</button>
