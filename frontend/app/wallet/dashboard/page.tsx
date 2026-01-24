'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { getWalletBalance } from '@/lib/api';
import { weiToEth } from '@/lib/types';
import type { BalanceResponse } from '@/lib/types';

export default function WalletDashboard() {
  const router = useRouter();
  const [walletId, setWalletId] = useState<string | null>(null);
  const [balances, setBalances] = useState<BalanceResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    // Get wallet ID from localStorage
    const id = localStorage.getItem('current_wallet');
    if (!id) {
      router.push('/wallet');
      return;
    }
    setWalletId(id);
    loadBalance(id);
  }, [router]);

  const loadBalance = async (id: string) => {
    setLoading(true);
    setError('');
    try {
      const data = await getWalletBalance(id);
      setBalances(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load balance');
    } finally {
      setLoading(false);
    }
  };

  const handleRefresh = () => {
    if (walletId) {
      loadBalance(walletId);
    }
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-12">
        <div className="text-center">
          <div className="text-cyan-400 text-xl">Loading wallet...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-12">
        <div className="max-w-2xl mx-auto">
          <div className="cyber-card border-red-500/50">
            <h2 className="text-xl font-bold mb-2 text-red-400">Error</h2>
            <p className="text-zinc-300 mb-4">{error}</p>
            <button onClick={handleRefresh} className="aggressive-btn">
              Retry
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-12">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-4xl font-bold text-cyan-400 mb-2">Wallet Dashboard</h1>
            <p className="text-zinc-400 font-mono text-sm">
              {balances?.balances[0]?.chain_name || 'Local Testnet'}
            </p>
          </div>
          <button onClick={handleRefresh} className="aggressive-btn">
            Refresh Balance
          </button>
        </div>

        {/* Balance Cards */}
        <div className="grid grid-cols-1 gap-6 mb-8">
          {balances?.balances.map((balance) => (
            <div key={balance.chain_id} className="cyber-card">
              <div className="flex items-start justify-between mb-4">
                <div>
                  <div className="text-sm text-zinc-400 mb-1">Balance</div>
                  <div className="text-5xl font-bold font-mono text-cyan-400">
                    {weiToEth(balance.balance)} <span className="text-2xl">{balance.symbol}</span>
                  </div>
                </div>
                <div className="text-right">
                  <div className="text-sm text-zinc-400 mb-1">Network</div>
                  <div className="font-semibold">{balance.chain_name}</div>
                  <div className="text-xs text-zinc-500 mono">Chain ID: {balance.chain_id}</div>
                </div>
              </div>

              {/* Address */}
              <div className="mt-4 pt-4 border-t border-cyan-500/20">
                <div className="text-sm text-zinc-400 mb-1">Wallet Address</div>
                <div className="font-mono text-sm break-all text-zinc-300">
                  {balances.wallet_id}
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Quick Actions */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Link href="/wallet/send">
            <div className="cyber-card hover:scale-105 transition-transform cursor-pointer">
              <div className="text-center">
                <div className="text-4xl mb-3">📤</div>
                <h3 className="text-xl font-bold mb-2">Send Transaction</h3>
                <p className="text-sm text-zinc-400 mb-4">
                  Transfer ETH to another address
                </p>
                <button className="aggressive-btn w-full">
                  Send
                </button>
              </div>
            </div>
          </Link>

          <div className="cyber-card opacity-50">
            <div className="text-center">
              <div className="text-4xl mb-3">📜</div>
              <h3 className="text-xl font-bold mb-2">Transaction History</h3>
              <p className="text-sm text-zinc-400 mb-4">
                View your past transactions
              </p>
              <button className="aggressive-btn w-full" disabled>
                Coming Soon
              </button>
            </div>
          </div>
        </div>

        {/* Back to Wallet */}
        <div className="mt-8 text-center">
          <Link href="/wallet" className="text-cyan-400 hover:text-cyan-300">
            ← Back to Wallet Home
          </Link>
        </div>
      </div>
    </div>
  );
}