<script lang="ts">
	import { page } from '$app/stores';
	import { walletStore } from '$lib/stores/wallet';

	let currentPath = $derived($page.url.pathname);

	const navLinks = [
		{ href: '/', label: 'Panel' },
		{ href: '/wallet', label: 'Billetera' },
		{ href: '/explorer', label: 'Explorador' }
	];

	function isActive(href: string, path: string): boolean {
		if (href === '/') return path === '/';
		return path.startsWith(href);
	}
</script>

<nav class="sticky top-0 z-50 border-b border-white/[0.06] bg-[#050507]/80 backdrop-blur-xl">
	<div class="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
		<a href="/" class="flex items-center gap-3 group">
			<svg
				width="32"
				height="32"
				viewBox="0 0 32 32"
				fill="none"
				class="transition-transform duration-300 group-hover:scale-110"
			>
				<path
					d="M16 2L28 9V23L16 30L4 23V9L16 2Z"
					fill="url(#logo-grad)"
					stroke="rgba(245,158,11,0.3)"
					stroke-width="0.5"
				/>
				<text
					x="16"
					y="20"
					text-anchor="middle"
					fill="#050507"
					font-weight="800"
					font-size="14"
					font-family="Inter, sans-serif">F</text
				>
				<defs>
					<linearGradient id="logo-grad" x1="4" y1="2" x2="28" y2="30">
						<stop stop-color="#fbbf24" />
						<stop offset="1" stop-color="#d97706" />
					</linearGradient>
				</defs>
			</svg>
			<span class="text-lg font-bold tracking-tight text-gradient-amber">Fernet Token</span>
		</a>

		<div class="flex items-center gap-1">
			{#each navLinks as { href, label }}
				<a
					{href}
					class="relative px-4 py-2 text-sm font-medium rounded-lg transition-all duration-200
						{isActive(href, currentPath)
						? 'text-amber-400 bg-amber-500/10'
						: 'text-zinc-400 hover:text-zinc-200 hover:bg-white/[0.04]'}"
				>
					{label}
					{#if isActive(href, currentPath)}
						<span
							class="absolute bottom-0 left-1/2 -translate-x-1/2 w-4 h-0.5 rounded-full bg-amber-400"
						></span>
					{/if}
				</a>
			{/each}

			{#if $walletStore.address}
				<div
					class="ml-4 flex items-center gap-2 rounded-full bg-white/[0.04] border border-white/[0.08] px-4 py-2 transition-all duration-200 hover:border-amber-500/30"
				>
					<span class="status-dot"></span>
					<span class="text-xs font-mono text-zinc-400">
						{$walletStore.address.slice(0, 6)}...{$walletStore.address.slice(-4)}
					</span>
				</div>
			{:else}
				<a
					href="/wallet"
					class="ml-4 flex items-center gap-2 rounded-full bg-amber-500/10 border border-amber-500/20 px-4 py-2 text-xs font-medium text-amber-400 hover:bg-amber-500/20 transition-all duration-200"
				>
					Conectar Billetera
				</a>
			{/if}
		</div>
	</div>
</nav>
