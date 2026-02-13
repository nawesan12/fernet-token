<script lang="ts">
	let { text }: { text: string } = $props();
	let copied = $state(false);

	async function copy() {
		await navigator.clipboard.writeText(text);
		copied = true;
		setTimeout(() => {
			copied = false;
		}, 2000);
	}
</script>

<button
	onclick={copy}
	class="inline-flex items-center text-zinc-600 hover:text-amber-500 transition-colors"
	title="Copiar"
>
	{#if copied}
		<svg
			class="w-3.5 h-3.5 text-emerald-400"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<polyline points="20 6 9 17 4 12" />
		</svg>
	{:else}
		<svg
			class="w-3.5 h-3.5"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
			<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
		</svg>
	{/if}
</button>
