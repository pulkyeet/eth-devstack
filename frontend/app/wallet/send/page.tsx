'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { sendTransaction, getWalletBalance } from '@/lib/api';
import { ethToWei, weiToEth } from '@/lib/types';

export default function SendTransaction() {
  const router = useRouter();
  const [walletId, setWalletId] = useState<string | null>(null);
  const [balance, setBalance] = useState('0');
  const [chainId, setChainId] = useState(1337);
  
  const [to, setTo] = useState('');
  const [amount, setAmount] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [txHash, setTxHash] = useState('');

  useEffect(() => {
    const id = localStorage.getItem('current_wallet');
    if (!id) {
      router.push('/wallet');
      return;
    }
    setWalletId(id);
    loadBalance(id);
  }, [router]);

  const loadBalance = async (id: string) => {
    try {
      const data = await getWalletBalance(id);
      if (data.balances.length > 0) {
        setBalance(data.balances[0].balance);
        setChainId(data.balances[0].chain_id);
      }
    } catch (err) {
      console.error('Failed to load balance:', err);
    }
  };

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setTxHash('');

    // Validation
    if (!to.trim()) {
      setError('Recipient address is required');
      return;
    }
    if (!to.startsWith('0x') || to.length !== 42) {
      setError('Invalid address format');
      return;
    }
    if (!amount || parseFloat(amount) <= 0) {
      setError('Amount must be greater than 0');
      return;
    }
    if (parseFloat(amount) > parseFloat(weiToEth(balance))) {
      setError('Insufficient balance');
      return;
    }
    if (!password) {
      setError('Password is required');
      return;
    }

    setLoading(true);
    try {
      const result = await sendTransaction({
        wallet_id: walletId!,
        chain_id: chainId,
        to: to.trim(),
        value: ethToWei(amount),
        password,
      });
      
      setTxHash(result.tx_hash);
      
      // Clear form
      setTo('');
      setAmount('');
      setPassword('');
      
      // Refresh balance after 3 seconds
      setTimeout(() => {
        if (walletId) loadBalance(walletId);
      }, 3000);
    } catch (err: any) {
      setError(err.message || 'Failed to send transaction');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-12">
      <div className="max-w-2xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2 text-cyan-400">Send Transaction</h1>
          <p className="text-zinc-400">Transfer ETH to another address</p>
        </div>

        {/* Balance Card */}
        <div className="cyber-card mb-6">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm text-zinc-400 mb-1">Available Balance</div>
              <div className="text-3xl font-bold font-mono text-cyan-400">
                {weiToEth(balance)} ETH
              </div>
            </div>
            <div className="text-right">
              <div className="text-sm text-zinc-400">Chain ID</div>
              <div className="font-mono">{chainId}</div>
            </div>
          </div>
        </div>

        {/* Success Message */}
        {txHash && (
          <div className="mb-6 cyber-card border-green-500/50">
            <h3 className="text-lg font-bold mb-2 text-green-400">✅ Transaction Sent</h3>
            <div className="text-sm text-zinc-300 mb-2">
              Your transaction has been broadcast to the network
            </div>
            <div className="font-mono text-xs break-all text-cyan-400 bg-black/50 p-3 clip-angle-sm">
              {txHash}
            </div>
            <Link href={`/explorer/tx/${txHash}`} className="mt-3 block">
              <button className="aggressive-btn w-full">
                View in Explorer →
              </button>
            </Link>
          </div>
        )}

        {/* Send Form */}
        {!txHash && (
        <form onSubmit={handleSend}>
          <div className="cyber-card mb-6">
            <div className="space-y-4">
              {/* Recipient Address */}
              <div>
                <label className="block text-sm font-semibold mb-2 text-cyan-400">
                  Recipient Address
                </label>
                <input
                  type="text"
                  value={to}
                  onChange={(e) => setTo(e.target.value)}
                  placeholder="0x..."
                  className="w-full bg-black/70 border border-cyan-500/50 px-4 py-3 text-white font-mono text-sm
                    focus:outline-none focus:border-cyan-400 focus:shadow-[0_0_15px_rgba(34,211,238,0.3)]
                    clip-angle-sm"
                  disabled={loading}
                />
              </div>

              {/* Amount */}
              <div>
                <label className="block text-sm font-semibold mb-2 text-cyan-400">
                  Amount (ETH)
                </label>
                <div className="relative">
                  <input
                    type="number"
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
                    placeholder="0.0"
                    step="0.0001"
                    className="w-full bg-black/70 border border-cyan-500/50 px-4 py-3 text-white font-mono
                      focus:outline-none focus:border-cyan-400 focus:shadow-[0_0_15px_rgba(34,211,238,0.3)]
                      clip-angle-sm pr-20"
                    disabled={loading}
                  />
                  <button
                    type="button"
                    onClick={() => setAmount(weiToEth(balance))}
                    className="absolute right-2 top-1/2 -translate-y-1/2 text-xs text-cyan-400 hover:text-cyan-300 font-bold"
                    disabled={loading}
                  >
                    MAX
                  </button>
                </div>
                <div className="mt-1 text-xs text-zinc-500">
                  Available: {weiToEth(balance)} ETH
                </div>
              </div>

              {/* Password */}
              <div>
                <label className="block text-sm font-semibold mb-2 text-cyan-400">
                  Wallet Password
                </label>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••••••"
                  className="w-full bg-black/70 border border-cyan-500/50 px-4 py-3 text-white font-mono
                    focus:outline-none focus:border-cyan-400 focus:shadow-[0_0_15px_rgba(34,211,238,0.3)]
                    clip-angle-sm"
                  disabled={loading}
                />
                <div className="mt-1 text-xs text-zinc-500">
                  Required to sign the transaction
                </div>
              </div>
            </div>
          </div>

          {/* Error */}
          {error && (
            <div className="mb-4 p-4 bg-red-900/20 border border-red-500/50 clip-angle-sm">
              <p className="text-red-400 text-sm">{error}</p>
            </div>
          )}

          {/* Transaction Summary */}
          {to && amount && (
            <div className="cyber-card border-purple-500/50 mb-6">
              <h3 className="text-sm font-semibold mb-3 text-purple-400">Transaction Summary</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-zinc-400">To:</span>
                  <span className="font-mono text-xs">{to.slice(0, 10)}...{to.slice(-8)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-400">Amount:</span>
                  <span className="font-mono font-bold">{amount} ETH</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-400">Gas:</span>
                  <span className="font-mono text-xs">~21,000 gas</span>
                </div>
              </div>
            </div>
          )}

          {/* Actions */}
          <div className="flex gap-4">
            <Link href="/wallet/dashboard" className="flex-1">
              <button
                type="button"
                className="aggressive-btn w-full"
                disabled={loading}
              >
                Cancel
              </button>
            </Link>
            <button
              type="submit"
              className="aggressive-btn flex-1"
              disabled={loading || !walletId}
            >
              {loading ? 'Sending...' : 'Send Transaction'}
            </button>
          </div>
        </form>
        )}
      </div>
    </div>
  );
}