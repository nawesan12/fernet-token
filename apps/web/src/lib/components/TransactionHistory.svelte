<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type TxResult } from '$lib/api';
	import TransactionRow from './TransactionRow.svelte';

	let { address }: { address: string } = $props();
	let results = $state<TxResult[]>([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			const res = await api.getAddressTransactions(address);
			results = res.transactions || [];
		} catch (err: any) {
			error = err.message;
		} finally {
			loading = false;
		}
	});
</script>

<div class="animate-slide-up">
	<h2 class="mb-4 text-xl font-bold text-zinc-100">Historial de Transacciones</h2>
	{#if loading}
		<div class="space-y-2">
			{#each Array(3) as _, i}
				<div class="glass-card h-16 animate-shimmer" style="animation-delay: {i * 100}ms"></div>
			{/each}
		</div>
	{:else if error}
		<p class="text-sm text-zinc-500">{error}</p>
	{:else if results.length === 0}
		<div class="glass-card p-8 text-center">
			<p class="text-zinc-500">No hay transacciones aun.</p>
		</div>
	{:else}
		<div class="space-y-2">
			{#each results as result (result.transaction.id)}
				<TransactionRow tx={result.transaction} />
			{/each}
		</div>
	{/if}
</div>
