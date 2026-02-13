import { writable, get } from 'svelte/store';
import {
	generateWallet,
	exportPrivateKeyJWK,
	restoreWallet,
	type WalletKeys
} from '$lib/crypto';
import { api } from '$lib/api';

interface WalletState {
	loaded: boolean;
	keys: WalletKeys | null;
	address: string;
	publicKey: string;
	balance: number;
	balanceFormatted: string;
	nonce: number;
}

const initial: WalletState = {
	loaded: false,
	keys: null,
	address: '',
	publicKey: '',
	balance: 0,
	balanceFormatted: '0 FERNET',
	nonce: 0
};

export const walletStore = writable<WalletState>(initial);

const STORAGE_KEY = 'fernet-wallet';

export async function initWallet() {
	const saved = localStorage.getItem(STORAGE_KEY);
	if (saved) {
		try {
			const jwk = JSON.parse(saved);
			const keys = await restoreWallet(jwk);
			walletStore.set({
				loaded: true,
				keys,
				address: keys.address,
				publicKey: keys.publicKeyHex,
				balance: 0,
				balanceFormatted: '0 FERNET',
				nonce: 0
			});
			await refreshBalance();
		} catch {
			localStorage.removeItem(STORAGE_KEY);
			walletStore.set({ ...initial, loaded: true });
		}
	} else {
		walletStore.set({ ...initial, loaded: true });
	}
}

export async function createWallet() {
	const keys = await generateWallet();
	const jwk = await exportPrivateKeyJWK(keys.privateKey);
	localStorage.setItem(STORAGE_KEY, JSON.stringify(jwk));

	walletStore.set({
		loaded: true,
		keys,
		address: keys.address,
		publicKey: keys.publicKeyHex,
		balance: 0,
		balanceFormatted: '0 FERNET',
		nonce: 0
	});
}

export async function refreshBalance() {
	const state = get(walletStore);
	if (!state.address) return;

	try {
		const [balRes, nonceRes] = await Promise.all([
			api.getBalance(state.address),
			api.getNonce(state.address)
		]);
		walletStore.update((s) => ({
			...s,
			balance: balRes.balance,
			balanceFormatted: balRes.formatted,
			nonce: nonceRes.nonce
		}));
	} catch {
		// API might not be running yet
	}
}

export function logout() {
	localStorage.removeItem(STORAGE_KEY);
	walletStore.set({ ...initial, loaded: true });
}
