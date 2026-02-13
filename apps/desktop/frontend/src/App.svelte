<script lang="ts">
  import {
    GetWalletInfo, CreateWallet, GetBalance, Mine,
    GetChainHeight, GetRecentBlocks, ConnectToPeer,
    GetPeerCount, GetPendingTxCount, SendTransaction
  } from '../wailsjs/go/main/App.js';

  let tab = 'wallet';

  // Wallet state
  let walletAddress = '';
  let walletPubKey = '';
  let balance = '0 FERNET';
  let walletStatus = '';

  // Chain state
  let chainHeight = 0;
  let recentBlocks: any[] = [];
  let pendingCount = 0;
  let peerCount = 0;
  let mineStatus = '';

  // Send form
  let sendReceiver = '';
  let sendAmount = '';
  let sendFee = '1000000';
  let sendStatus = '';

  // Peer form
  let peerAddr = '';
  let peerStatus = '';

  async function loadWallet() {
    try {
      const info = await GetWalletInfo();
      if (info.address) {
        walletAddress = info.address;
        walletPubKey = info.publicKey;
      }
      await refreshBalance();
    } catch (e) {}
  }

  async function createNewWallet() {
    try {
      const info = await CreateWallet();
      if (info.address) {
        walletAddress = info.address;
        walletPubKey = info.publicKey;
        walletStatus = 'Wallet created!';
        await refreshBalance();
      } else {
        walletStatus = info.error || 'Failed';
      }
    } catch (e: any) {
      walletStatus = e.message;
    }
  }

  async function refreshBalance() {
    try {
      const res = await GetBalance();
      balance = res.formatted || '0 FERNET';
    } catch (e) {}
  }

  async function refreshChain() {
    try {
      chainHeight = await GetChainHeight();
      recentBlocks = await GetRecentBlocks(10) || [];
      pendingCount = await GetPendingTxCount();
      peerCount = await GetPeerCount();
    } catch (e) {}
  }

  async function mineBlock() {
    mineStatus = 'Mining...';
    try {
      const res = await Mine();
      if (res.error) {
        mineStatus = res.error;
      } else {
        mineStatus = res.message;
        await refreshChain();
        await refreshBalance();
      }
    } catch (e: any) {
      mineStatus = e.message;
    }
  }

  async function sendTx() {
    sendStatus = 'Sending...';
    try {
      const amt = Math.floor(parseFloat(sendAmount) * 100_000_000);
      const fee = parseInt(sendFee);
      const res = await SendTransaction(sendReceiver, amt, fee);
      if (res.error) {
        sendStatus = res.error;
      } else {
        sendStatus = res.message;
        sendReceiver = '';
        sendAmount = '';
        await refreshBalance();
      }
    } catch (e: any) {
      sendStatus = e.message;
    }
  }

  async function connectPeer() {
    peerStatus = 'Connecting...';
    try {
      const res = await ConnectToPeer(peerAddr);
      peerStatus = res.error || res.message;
      peerAddr = '';
      await refreshChain();
    } catch (e: any) {
      peerStatus = e.message;
    }
  }

  // Load on mount
  loadWallet();
  refreshChain();
  setInterval(refreshChain, 10000);
</script>

<main>
  <header>
    <h1>Fernet Token Wallet</h1>
    <nav>
      <button class:active={tab === 'wallet'} on:click={() => tab = 'wallet'}>Wallet</button>
      <button class:active={tab === 'explorer'} on:click={() => { tab = 'explorer'; refreshChain(); }}>Explorer</button>
      <button class:active={tab === 'node'} on:click={() => { tab = 'node'; refreshChain(); }}>Node</button>
    </nav>
  </header>

  <div class="content">
    {#if tab === 'wallet'}
      <section>
        {#if walletAddress}
          <div class="card">
            <h2>Balance</h2>
            <div class="balance">{balance}</div>
            <div class="address">{walletAddress}</div>
            <button class="secondary" on:click={refreshBalance}>Refresh</button>
          </div>

          <div class="card">
            <h2>Send FERNET</h2>
            <label>
              Receiver
              <input type="text" bind:value={sendReceiver} placeholder="Address..." />
            </label>
            <div class="row">
              <label class="flex1">
                Amount (FERNET)
                <input type="number" step="0.00000001" bind:value={sendAmount} placeholder="0.00" />
              </label>
              <label>
                Fee
                <input type="number" bind:value={sendFee} />
              </label>
            </div>
            <button on:click={sendTx}>Send</button>
            {#if sendStatus}<p class="status">{sendStatus}</p>{/if}
          </div>
        {:else}
          <div class="card center">
            <h2>No Wallet Found</h2>
            <p>Create a wallet to get started</p>
            <button on:click={createNewWallet}>Create Wallet</button>
            {#if walletStatus}<p class="status">{walletStatus}</p>{/if}
          </div>
        {/if}
      </section>

    {:else if tab === 'explorer'}
      <section>
        <div class="stats">
          <div class="stat"><span class="label">Height</span><span class="value">{chainHeight}</span></div>
          <div class="stat"><span class="label">Pending</span><span class="value">{pendingCount}</span></div>
          <div class="stat"><span class="label">Peers</span><span class="value">{peerCount}</span></div>
        </div>

        <div class="card">
          <div class="row" style="justify-content: space-between; align-items: center;">
            <h2>Recent Blocks</h2>
            <button class="secondary" on:click={mineBlock}>Mine Block</button>
          </div>
          {#if mineStatus}<p class="status">{mineStatus}</p>{/if}

          {#if recentBlocks.length === 0}
            <p class="muted">No blocks mined yet</p>
          {:else}
            <div class="block-list">
              {#each recentBlocks as b}
                <div class="block-item">
                  <div class="row">
                    <span class="block-index">#{b.index}</span>
                    <span class="muted">{b.txCount} txn{b.txCount !== 1 ? 's' : ''}</span>
                  </div>
                  <div class="hash">{b.hash}</div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </section>

    {:else if tab === 'node'}
      <section>
        <div class="card">
          <h2>Connect to Peer</h2>
          <div class="row">
            <input type="text" bind:value={peerAddr} placeholder="host:port" class="flex1" />
            <button on:click={connectPeer}>Connect</button>
          </div>
          {#if peerStatus}<p class="status">{peerStatus}</p>{/if}
        </div>

        <div class="card">
          <h2>Node Info</h2>
          <p>P2P Port: 6100</p>
          <p>Connected Peers: {peerCount}</p>
          <p>Chain Height: {chainHeight}</p>
          <p>Pending Transactions: {pendingCount}</p>
        </div>
      </section>
    {/if}
  </div>
</main>

<style>
  :global(body) { margin: 0; background: #09090b; color: #e4e4e7; }
  main { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }

  header {
    display: flex; justify-content: space-between; align-items: center;
    padding: 16px 24px; border-bottom: 1px solid #27272a;
  }
  h1 { font-size: 18px; color: #f59e0b; margin: 0; }
  nav button {
    background: none; border: none; color: #a1a1aa; cursor: pointer;
    padding: 6px 12px; font-size: 14px; border-radius: 4px;
  }
  nav button:hover { color: #fbbf24; }
  nav button.active { color: #f59e0b; background: rgba(245,158,11,0.1); }

  .content { padding: 24px; max-width: 800px; margin: 0 auto; }
  section { display: flex; flex-direction: column; gap: 16px; }

  .card {
    background: rgba(39,39,42,0.5); border: 1px solid #27272a;
    border-radius: 8px; padding: 20px;
  }
  .card.center { text-align: center; }
  .card h2 { margin: 0 0 12px; font-size: 16px; color: #e4e4e7; }

  .balance { font-size: 28px; font-weight: bold; color: #f59e0b; }
  .address { font-family: monospace; font-size: 12px; color: #71717a; word-break: break-all; margin-top: 8px; }

  label { display: block; font-size: 12px; color: #71717a; margin-bottom: 8px; }
  input {
    width: 100%; box-sizing: border-box; padding: 8px 12px;
    background: #18181b; border: 1px solid #3f3f46; border-radius: 4px;
    color: #e4e4e7; font-size: 14px; margin-top: 4px;
  }
  input:focus { outline: none; border-color: #f59e0b; }

  button {
    background: #d97706; color: #09090b; border: none; border-radius: 4px;
    padding: 8px 16px; font-size: 14px; font-weight: 600; cursor: pointer;
    margin-top: 8px;
  }
  button:hover { background: #f59e0b; }
  button.secondary { background: none; border: 1px solid #3f3f46; color: #a1a1aa; }
  button.secondary:hover { border-color: #f59e0b; color: #f59e0b; }

  .row { display: flex; gap: 12px; }
  .flex1 { flex: 1; }
  .status { font-size: 12px; color: #71717a; margin-top: 8px; }
  .muted { color: #52525b; font-size: 14px; }

  .stats { display: flex; gap: 12px; }
  .stat {
    flex: 1; background: rgba(39,39,42,0.5); border: 1px solid #27272a;
    border-radius: 8px; padding: 12px; text-align: center;
  }
  .stat .label { display: block; font-size: 11px; color: #71717a; }
  .stat .value { display: block; font-size: 24px; font-weight: bold; color: #f59e0b; margin-top: 4px; }

  .block-list { display: flex; flex-direction: column; gap: 8px; }
  .block-item { padding: 8px; border: 1px solid #27272a; border-radius: 4px; }
  .block-index { font-weight: bold; color: #f59e0b; }
  .hash { font-family: monospace; font-size: 11px; color: #52525b; overflow: hidden; text-overflow: ellipsis; }
</style>
