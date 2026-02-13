<script lang="ts">
	import { walletStore, createWallet, refreshBalance, logout } from '$lib/stores/wallet';
	import WalletBalance from '$lib/components/WalletBalance.svelte';
	import SendForm from '$lib/components/SendForm.svelte';
	import TransactionHistory from '$lib/components/TransactionHistory.svelte';
	import { api } from '$lib/api';

	let faucetStatus = $state('');
	let faucetLoading = $state(false);
	let showDetails = $state(false);

	async function requestFaucet() {
		if (!$walletStore.address) return;
		faucetLoading = true;
		faucetStatus = '';
		try {
			const res = await api.faucet($walletStore.address);
			faucetStatus = `Se acreditaron ${res.formatted}`;
			await refreshBalance();
		} catch (err: any) {
			faucetStatus = `Error: ${err.message}`;
		} finally {
			faucetLoading = false;
		}
	}
</script>

<div class="space-y-10">
	<h1 class="text-3xl font-bold tracking-tight text-zinc-100 animate-fade-in">Billetera</h1>

	{#if !$walletStore.address}
		<!-- Onboarding -->
		<div class="flex min-h-[50vh] items-center justify-center">
			<div class="glass-card max-w-md animate-slide-up p-12 text-center">
				<div
					class="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-amber-500/10"
				>
					<svg
						width="32"
						height="32"
						viewBox="0 0 24 24"
						fill="none"
						stroke="#f59e0b"
						stroke-width="1.5"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
					</svg>
				</div>
				<h2 class="mb-3 text-xl font-bold text-zinc-100">Crea tu Billetera</h2>
				<p class="mb-8 text-sm leading-relaxed text-zinc-500">
					Genera un par de claves ECDSA P-256 directamente en tu navegador. Tu clave privada nunca
					sale de este dispositivo.
				</p>
				<button onclick={createWallet} class="btn-primary w-full">Crear Billetera</button>
			</div>
		</div>
	{:else}
		<div class="grid animate-slide-up gap-8 lg:grid-cols-5">
			<!-- Left: Balance + Actions -->
			<div class="space-y-6 lg:col-span-3">
				<WalletBalance />

				<div class="flex flex-wrap gap-3">
					<button onclick={requestFaucet} disabled={faucetLoading} class="btn-secondary">
						{faucetLoading ? 'Solicitando...' : 'Grifo (100 FERNET)'}
					</button>
					<button onclick={() => (showDetails = !showDetails)} class="btn-secondary">
						{showDetails ? 'Ocultar Detalles' : 'Detalles'}
					</button>
					<button
						onclick={logout}
						class="btn-secondary !border-red-500/20 !text-red-400 hover:!bg-red-500/10"
					>
						Cerrar Sesion
					</button>
				</div>

				{#if faucetStatus}
					<p class="text-xs text-zinc-400 animate-fade-in">{faucetStatus}</p>
				{/if}

				{#if showDetails}
					<div class="glass-card animate-fade-in p-6">
						<h3 class="section-label mb-4">Detalles de Billetera</h3>
						<div class="space-y-3 text-sm">
							<div>
								<span class="text-xs text-zinc-500">Direccion</span>
								<p class="mt-1 break-all font-mono text-xs text-zinc-300">
									{$walletStore.address}
								</p>
							</div>
							<div>
								<span class="text-xs text-zinc-500">Clave Publica</span>
								<p class="mt-1 break-all font-mono text-xs text-zinc-400">
									{$walletStore.publicKey}
								</p>
							</div>
							<div>
								<span class="text-xs text-zinc-500">Nonce</span>
								<p class="mt-1 text-zinc-300">{$walletStore.nonce}</p>
							</div>
						</div>
					</div>
				{/if}
			</div>

			<!-- Right: Send Form -->
			<div class="lg:col-span-2">
				<SendForm />
			</div>
		</div>

		<!-- Transaction History -->
		<TransactionHistory address={$walletStore.address} />
	{/if}
</div>
