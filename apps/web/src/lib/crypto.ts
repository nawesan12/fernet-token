// Client-side wallet crypto using Web Crypto API (ECDSA P-256)
// Private keys never leave the browser.

export interface WalletKeys {
	privateKey: CryptoKey;
	publicKey: CryptoKey;
	publicKeyHex: string;
	address: string;
}

const ALGORITHM = { name: 'ECDSA', namedCurve: 'P-256' };

export async function generateWallet(): Promise<WalletKeys> {
	const keyPair = await crypto.subtle.generateKey(ALGORITHM, true, ['sign', 'verify']);

	const publicKeyHex = await exportPublicKeyHex(keyPair.publicKey);
	const address = await deriveAddress(publicKeyHex);

	return {
		privateKey: keyPair.privateKey,
		publicKey: keyPair.publicKey,
		publicKeyHex,
		address
	};
}

export async function exportPublicKeyHex(key: CryptoKey): Promise<string> {
	// Export as raw - gives 65 bytes: 0x04 || X (32) || Y (32)
	const raw = await crypto.subtle.exportKey('raw', key);
	const bytes = new Uint8Array(raw);
	// Skip the 0x04 prefix byte, take X||Y (64 bytes)
	const xy = bytes.slice(1);
	return bufToHex(xy);
}

async function deriveAddress(publicKeyHex: string): Promise<string> {
	const pubKeyBytes = hexToBuf(publicKeyHex);
	const hashBuf = await crypto.subtle.digest('SHA-256', pubKeyBytes);
	const hashBytes = new Uint8Array(hashBuf);
	// First 20 bytes as address (matches Go implementation)
	return bufToHex(hashBytes.slice(0, 20));
}

export async function signTransaction(
	privateKey: CryptoKey,
	signableData: string
): Promise<string> {
	// SHA-256 hash of the signable data string
	const encoder = new TextEncoder();
	const dataBytes = encoder.encode(signableData);
	const hashBuf = await crypto.subtle.digest('SHA-256', dataBytes);

	// Sign the hash with ECDSA
	const sigBuf = await crypto.subtle.sign(
		{ name: 'ECDSA', hash: 'SHA-256' },
		privateKey,
		hashBuf
	);

	// P1363 format: R (32 bytes) || S (32 bytes) - matches Go's R||S format
	return bufToHex(new Uint8Array(sigBuf));
}

export function buildSignableData(
	sender: string,
	receiver: string,
	amount: number,
	fee: number,
	nonce: number,
	timestamp: number
): string {
	return `${sender}:${receiver}:${amount}:${fee}:${nonce}:${timestamp}`;
}

export async function exportPrivateKeyJWK(key: CryptoKey): Promise<JsonWebKey> {
	return crypto.subtle.exportKey('jwk', key);
}

export async function importPrivateKeyJWK(jwk: JsonWebKey): Promise<CryptoKey> {
	return crypto.subtle.importKey('jwk', jwk, ALGORITHM, true, ['sign']);
}

export async function importPublicKeyFromJWK(jwk: JsonWebKey): Promise<CryptoKey> {
	// Remove private key component for public key import
	const publicJwk = { ...jwk };
	delete publicJwk.d;
	return crypto.subtle.importKey('jwk', publicJwk, ALGORITHM, true, ['verify']);
}

export async function restoreWallet(jwk: JsonWebKey): Promise<WalletKeys> {
	const privateKey = await importPrivateKeyJWK(jwk);

	// Derive public key from JWK
	const publicJwk = { ...jwk };
	delete publicJwk.d;
	const publicKey = await importPublicKeyFromJWK(publicJwk);

	const publicKeyHex = await exportPublicKeyHex(publicKey);
	const address = await deriveAddress(publicKeyHex);

	return { privateKey, publicKey, publicKeyHex, address };
}

// Utility functions
function bufToHex(buf: Uint8Array): string {
	return Array.from(buf)
		.map((b) => b.toString(16).padStart(2, '0'))
		.join('');
}

function hexToBuf(hex: string): Uint8Array {
	const bytes = new Uint8Array(hex.length / 2);
	for (let i = 0; i < hex.length; i += 2) {
		bytes[i / 2] = parseInt(hex.substring(i, i + 2), 16);
	}
	return bytes;
}

export async function computeTxHash(
	sender: string,
	receiver: string,
	amount: number,
	fee: number,
	nonce: number,
	timestamp: number
): Promise<string> {
	const data = `${sender}:${receiver}:${amount}:${fee}:${nonce}:${timestamp}`;
	const encoder = new TextEncoder();
	const hashBuf = await crypto.subtle.digest('SHA-256', encoder.encode(data));
	return bufToHex(new Uint8Array(hashBuf));
}
