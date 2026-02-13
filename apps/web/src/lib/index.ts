export const COINBASE_ADDRESS = '0000000000000000000000000000000000000000';
export const ONE_FERNET = 100_000_000;

export function formatAmount(fernetoshi: number): string {
	const whole = Math.floor(fernetoshi / ONE_FERNET);
	const frac = fernetoshi % ONE_FERNET;
	if (frac === 0) return `${whole} FERNET`;
	return `${whole}.${frac.toString().padStart(8, '0')} FERNET`;
}

export function shortAddr(address: string, chars = 10): string {
	if (!address || address.length <= chars * 2) return address;
	return `${address.slice(0, chars)}...${address.slice(-6)}`;
}

export function shortHash(hash: string, chars = 16): string {
	if (!hash || hash.length <= chars) return hash;
	return `${hash.slice(0, chars)}...`;
}

export function timeAgo(timestamp: number): string {
	const seconds = Math.floor(Date.now() / 1000 - timestamp);
	if (seconds < 5) return 'ahora mismo';
	if (seconds < 60) return `hace ${seconds}s`;
	if (seconds < 3600) return `hace ${Math.floor(seconds / 60)}m`;
	if (seconds < 86400) return `hace ${Math.floor(seconds / 3600)}h`;
	return `hace ${Math.floor(seconds / 86400)}d`;
}

export function formatDate(timestamp: number): string {
	return new Date(timestamp * 1000).toLocaleString('es-AR');
}

export function animateValue(
	start: number,
	end: number,
	duration: number,
	callback: (value: number) => void
) {
	const startTime = performance.now();
	function update(currentTime: number) {
		const elapsed = currentTime - startTime;
		const progress = Math.min(elapsed / duration, 1);
		const eased = 1 - Math.pow(1 - progress, 3);
		callback(Math.floor(start + (end - start) * eased));
		if (progress < 1) requestAnimationFrame(update);
	}
	requestAnimationFrame(update);
}
