import { writable } from 'svelte/store';
import { api, type Block, type Transaction } from '$lib/api';

interface ChainState {
	height: number;
	recentBlocks: Block[];
	pendingTxns: Transaction[];
	pendingCount: number;
	peerCount: number;
}

const initial: ChainState = {
	height: 0,
	recentBlocks: [],
	pendingTxns: [],
	pendingCount: 0,
	peerCount: 0
};

export const chainStore = writable<ChainState>(initial);

export async function refreshChain() {
	try {
		const [chainRes, pendingRes, peersRes] = await Promise.all([
			api.getBlockchain(),
			api.getPending(),
			api.getPeers()
		]);

		const blocks = chainRes.chain || [];
		const recent = blocks.slice(-10).reverse();

		chainStore.set({
			height: chainRes.height,
			recentBlocks: recent,
			pendingTxns: pendingRes.transactions || [],
			pendingCount: pendingRes.count,
			peerCount: peersRes.count
		});
	} catch {
		// API might not be running yet
	}
}
