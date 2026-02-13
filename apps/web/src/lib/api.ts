const API_BASE = 'http://localhost:8080/api';

export interface Transaction {
	id: string;
	sender: string;
	receiver: string;
	amount: number;
	fee: number;
	nonce: number;
	timestamp: number;
	pubKey: string;
	signature: string;
}

export interface Block {
	index: number;
	timestamp: number;
	transactions: Transaction[];
	prevHash: string;
	hash: string;
	nonce: number;
	miner: string;
}

export interface BalanceResponse {
	address: string;
	balance: number;
	formatted: string;
}

export interface HeightResponse {
	height: number;
}

export interface BlockchainResponse {
	chain: Block[];
	height: number;
}

export interface PendingResponse {
	transactions: Transaction[];
	count: number;
}

export interface PeersResponse {
	peers: string[] | null;
	count: number;
}

export interface WalletResponse {
	address: string;
	publicKey: string;
}

export interface NonceResponse {
	address: string;
	nonce: number;
}

export interface TxResult {
	transaction: Transaction;
	blockIndex: number;
	blockHash: string;
}

export interface AddressTransactionsResponse {
	address: string;
	transactions: TxResult[];
	count: number;
}

async function fetchJSON<T>(url: string, options?: RequestInit): Promise<T> {
	const res = await fetch(url, options);
	const data = await res.json();
	if (!res.ok) {
		throw new Error(data.error || `HTTP ${res.status}`);
	}
	return data as T;
}

export const api = {
	getBlockchain: () => fetchJSON<BlockchainResponse>(`${API_BASE}/blockchain`),

	getHeight: () => fetchJSON<HeightResponse>(`${API_BASE}/blockchain/height`),

	getBlock: (index: number) => fetchJSON<Block>(`${API_BASE}/block/${index}`),

	getBalance: (address: string) => fetchJSON<BalanceResponse>(`${API_BASE}/balance/${address}`),

	getNonce: (address: string) => fetchJSON<NonceResponse>(`${API_BASE}/nonce/${address}`),

	getPending: () => fetchJSON<PendingResponse>(`${API_BASE}/tx/pending`),

	getPeers: () => fetchJSON<PeersResponse>(`${API_BASE}/peers`),

	createWallet: () =>
		fetchJSON<WalletResponse>(`${API_BASE}/wallet/create`, { method: 'POST' }),

	submitTransaction: (tx: Transaction) =>
		fetchJSON<{ message: string; txId: string }>(`${API_BASE}/transaction`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(tx)
		}),

	mine: (miner: string) =>
		fetchJSON<{ message: string; block: Block }>(`${API_BASE}/mine`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ miner })
		}),

	connectPeer: (address: string) =>
		fetchJSON<{ message: string }>(`${API_BASE}/peers/connect`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ address })
		}),

	faucet: (address: string) =>
		fetchJSON<{ message: string; amount: number; formatted: string }>(`${API_BASE}/faucet`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ address })
		}),

	getTransaction: (id: string) =>
		fetchJSON<TxResult>(`${API_BASE}/tx/${id}`),

	getAddressTransactions: (address: string) =>
		fetchJSON<AddressTransactionsResponse>(`${API_BASE}/address/${address}/transactions`)
};
