<script lang="ts">
	import { walletStore, refreshBalance } from '$lib/stores/wallet';
	import { shortAddr } from '$lib';
	import { onMount } from 'svelte';

	onMount(() => {
		if ($walletStore.address) refreshBalance();
	});
</script>

<div class="glass-card relative overflow-hidden p-8">
	<!-- Ambient glow -->
	<div class="absolute -right-10 -top-10 h-40 w-40 rounded-full bg-amber-500/10 blur-3xl"></div>

	<div class="relative">
		<div class="section-label mb-2">Balance</div>
		<div class="text-4xl font-extrabold tracking-tight text-gradient-amber">
			{$walletStore.balanceFormatted}
		</div>
		<div class="mt-4 flex items-center gap-3">
			<span class="max-w-[220px] truncate font-mono text-xs text-zinc-600"
				>{$walletStore.address}</span
			>
			<button
				onclick={() => refreshBalance()}
				class="text-xs text-amber-600 transition-colors hover:text-amber-400"
			>
				Actualizar
			</button>
		</div>
	</div>
</div>
