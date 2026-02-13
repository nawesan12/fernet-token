<script lang="ts">
	import { walletStore, refreshBalance } from '$lib/stores/wallet';
	import { signTransaction, buildSignableData, computeTxHash } from '$lib/crypto';
	import { api } from '$lib/api';

	let receiver = $state('');
	let amount = $state('');
	let fee = $state('1000000');
	let status = $state('');
	let loading = $state(false);

	async function send() {
		if (!$walletStore.keys || !receiver || !amount) return;

		loading = true;
		status = '';

		try {
			const amountVal = Math.floor(parseFloat(amount) * 100_000_000);
			const feeVal = parseInt(fee);
			const nonce = $walletStore.nonce;
			const timestamp = Date.now() * 1_000_000;

			const signableData = buildSignableData(
				$walletStore.address,
				receiver,
				amountVal,
				feeVal,
				nonce,
				timestamp
			);

			const signature = await signTransaction($walletStore.keys.privateKey, signableData);
			const id = await computeTxHash(
				$walletStore.address,
				receiver,
				amountVal,
				feeVal,
				nonce,
				timestamp
			);

			const tx = {
				id,
				sender: $walletStore.address,
				receiver,
				amount: amountVal,
				fee: feeVal,
				nonce,
				timestamp,
				pubKey: $walletStore.publicKey,
				signature
			};

			const res = await api.submitTransaction(tx);
			status = `Transaccion enviada: ${res.txId.slice(0, 16)}...`;
			receiver = '';
			amount = '';
			await refreshBalance();
		} catch (err: any) {
			status = `Error: ${err.message}`;
		} finally {
			loading = false;
		}
	}
</script>

<div class="glass-card p-6">
	<h3 class="section-label mb-5">Enviar FERNET</h3>

	<div class="space-y-4">
		<div>
			<label for="receiver" class="mb-1.5 block text-xs text-zinc-500"
				>Direccion del Destinatario</label
			>
			<input
				id="receiver"
				type="text"
				bind:value={receiver}
				placeholder="Ingresa direccion..."
				class="input-field font-mono text-xs"
			/>
		</div>

		<div class="flex gap-3">
			<div class="flex-1">
				<label for="amount" class="mb-1.5 block text-xs text-zinc-500">Monto (FERNET)</label>
				<input
					id="amount"
					type="number"
					step="0.00000001"
					bind:value={amount}
					placeholder="0.00"
					class="input-field"
				/>
			</div>
			<div class="w-32">
				<label for="fee" class="mb-1.5 block text-xs text-zinc-500">Comision</label>
				<input id="fee" type="number" bind:value={fee} class="input-field" />
			</div>
		</div>

		<button
			onclick={send}
			disabled={loading || !receiver || !amount}
			class="btn-primary w-full"
		>
			{#if loading}
				<svg
					class="h-4 w-4 animate-spin-slow"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
				>
					<circle cx="12" cy="12" r="10" />
				</svg>
				Enviando...
			{:else}
				Enviar Transaccion
			{/if}
		</button>

		{#if status}
			<p class="mt-1 text-xs text-zinc-400 animate-fade-in">{status}</p>
		{/if}
	</div>
</div>
