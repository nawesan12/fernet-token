<script lang="ts">
	import { onMount } from 'svelte';
	import { chainStore, refreshChain } from '$lib/stores/chain';
	import { walletStore } from '$lib/stores/wallet';
	import BlockCard from '$lib/components/BlockCard.svelte';
	import StatCard from '$lib/components/StatCard.svelte';
	import MineButton from '$lib/components/MineButton.svelte';
	import PendingTicker from '$lib/components/PendingTicker.svelte';
	import { api } from '$lib/api';

	let mining = $state(false);
	let mineStatus = $state('');
	let mineSuccess = $state(false);

	onMount(() => {
		refreshChain();
		const interval = setInterval(refreshChain, 10000);
		return () => clearInterval(interval);
	});

	async function mineBlock() {
		if (!$walletStore.address) {
			mineStatus = 'Primero crea una billetera';
			return;
		}
		mining = true;
		mineStatus = 'Minando...';
		mineSuccess = false;
		try {
			const res = await api.mine($walletStore.address);
			mineStatus = `Bloque #${res.block.index} minado`;
			mineSuccess = true;
			await refreshChain();
			setTimeout(() => {
				mineSuccess = false;
			}, 3000);
		} catch (err: any) {
			mineStatus = `Error: ${err.message}`;
		} finally {
			mining = false;
		}
	}
</script>

<div class="space-y-12">
	<!-- Hero -->
	<section class="animate-slide-up">
		<h1 class="text-5xl font-extrabold tracking-tight text-gradient-amber sm:text-6xl">
			Fernet Token
		</h1>
		<p class="mt-3 max-w-xl text-lg text-zinc-500">
			Criptomoneda proof-of-work descentralizada. Mina bloques, envia transacciones y explora la
			cadena.
		</p>
	</section>

	<!-- Stats -->
	<section class="grid grid-cols-2 gap-4 lg:grid-cols-4 animate-slide-up stagger-1">
		<StatCard label="Altura de Cadena" value={$chainStore.height} icon="blocks" accent />
		<StatCard label="Txns Pendientes" value={$chainStore.pendingCount} icon="pending" />
		<StatCard label="Pares Conectados" value={$chainStore.peerCount} icon="network" />
		<StatCard
			label="Tu Balance"
			value={$walletStore.address ? $walletStore.balanceFormatted : 'Sin billetera'}
			format="text"
			icon="wallet"
		/>
	</section>

	<!-- Mining & Activity -->
	<section class="grid gap-6 lg:grid-cols-3 animate-slide-up stagger-2">
		<!-- Mine Panel -->
		<div class="glass-card p-6 lg:col-span-1">
			<h2 class="section-label mb-4">Mineria</h2>
			<MineButton {mining} {mineSuccess} onclick={mineBlock} />
			{#if mineStatus}
				<p
					class="mt-3 text-xs animate-fade-in {mineSuccess
						? 'text-emerald-400'
						: 'text-zinc-500'}"
				>
					{mineStatus}
				</p>
			{/if}
			{#if !$walletStore.address}
				<a
					href="/wallet"
					class="mt-4 inline-block text-sm text-amber-500 transition-colors hover:text-amber-400"
				>
					Crear billetera para minar &rarr;
				</a>
			{/if}
		</div>

		<!-- Pending Transactions -->
		<div class="glass-card p-6 lg:col-span-2">
			<div class="mb-4 flex items-center justify-between">
				<h2 class="section-label">Actividad Pendiente</h2>
				{#if $chainStore.pendingCount > 0}
					<span class="flex items-center gap-2 text-xs text-amber-400">
						<span class="status-dot-amber"></span>
						{$chainStore.pendingCount} pendiente{$chainStore.pendingCount !== 1 ? 's' : ''}
					</span>
				{/if}
			</div>
			<PendingTicker transactions={$chainStore.pendingTxns} />
		</div>
	</section>

	<!-- Recent Blocks -->
	<section class="animate-slide-up stagger-3">
		<div class="mb-6 flex items-center justify-between">
			<h2 class="text-xl font-bold text-zinc-100">Bloques Recientes</h2>
			<a
				href="/explorer"
				class="text-sm text-amber-500 transition-colors hover:text-amber-400"
			>
				Ver todos &rarr;
			</a>
		</div>
		{#if $chainStore.recentBlocks.length === 0}
			<div class="glass-card p-12 text-center">
				<div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-amber-500/10">
					<svg
						class="h-8 w-8 text-amber-400"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.5"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<rect x="3" y="3" width="7" height="7" />
						<rect x="14" y="3" width="7" height="7" />
						<rect x="14" y="14" width="7" height="7" />
						<rect x="3" y="14" width="7" height="7" />
					</svg>
				</div>
				<p class="text-zinc-500">No hay bloques aun. Comienza a minar!</p>
			</div>
		{:else}
			<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each $chainStore.recentBlocks.slice(0, 6) as block, i (block.index)}
					<div class="animate-slide-up" style="animation-delay: {(i + 1) * 50}ms">
						<BlockCard {block} />
					</div>
				{/each}
			</div>
		{/if}
	</section>
</div>
