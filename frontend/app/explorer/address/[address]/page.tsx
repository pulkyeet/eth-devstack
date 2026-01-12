'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import Card from '@/components/ui/Card';

interface AddressData {
  address: string;
  balance: string;
  nonce: number;
  is_contract: boolean;
  tx_count: number;
  first_seen_block: number;
  last_seen_block: number;
  first_seen_at: string;
  last_seen_at: string;
  token_balances?: TokenBalance[];
}

interface TokenBalance {
  token_address: string;
  token_name: string;
  token_symbol: string;
  token_type: string;
  balance: string;
  decimals: number;
}

interface Transaction {
  hash: string;
  block_number: number;
  timestamp: string;
  from_address: string;
  to_address: string | null;
  value: string;
  gas_used: number;
  status: number;
}

const API_BASE = 'http://localhost:8080/api/v1';

export default function AddressDetailPage() {
  const params = useParams();
  const address = params.address as string;
  
  const [addressData, setAddressData] = useState<AddressData | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [txPage, setTxPage] = useState(1);
  const [totalTxPages, setTotalTxPages] = useState(1);

  useEffect(() => {
    Promise.all([
      fetch(`${API_BASE}/addresses/${address}?chain_id=1337`).then(r => r.json()),
      fetch(`${API_BASE}/addresses/${address}/transactions?chain_id=1337&page=${txPage}&limit=20`).then(r => r.json())
    ])
    .then(([addrData, txData]) => {
      if (addrData.success) {
        setAddressData(addrData.data);
      } else {
        setError('Address not found');
      }
      
      if (txData.success) {
        setTransactions(txData.data.transactions || []);
        setTotalTxPages(txData.data.pagination?.total_pages || 1);
      }
      
      setLoading(false);
    })
    .catch(err => {
      setError('Failed to load address data');
      setLoading(false);
    });
  }, [address, txPage]);

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center text-zinc-400">Loading address...</div>
      </div>
    );
  }

  if (error || !addressData) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Card>
          <div className="text-center py-8">
            <div className="text-red-500 text-xl mb-2">❌ {error}</div>
            <Link href="/explorer" className="text-cyan-400 hover:text-cyan-300">
              ← Back to Explorer
            </Link>
          </div>
        </Card>
      </div>
    );
  }

  const formatValue = (wei: string) => {
    return (BigInt(wei) / BigInt(10 ** 18)).toString();
  };

  const formatTimestamp = (ts: string) => {
    const date = new Date(ts);
    const now = new Date();
    const diff = Math.floor((now.getTime() - date.getTime()) / 1000);
    
    if (diff < 60) return `${diff}s ago`;
    if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
    if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
    return `${Math.floor(diff / 86400)}d ago`;
  };

  const formatTokenBalance = (balance: string, decimals: number) => {
    const value = BigInt(balance) / BigInt(10 ** decimals);
    return value.toString();
  };

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-4xl font-bold mb-2 bg-gradient-to-r from-cyan-400 to-purple-400 bg-clip-text text-transparent">
          {addressData.is_contract ? 'Contract' : 'Address'}
        </h1>
        <div className="flex items-center gap-3">
          <span className="text-zinc-400 font-mono text-sm break-all">{address}</span>
          {addressData.is_contract && (
            <span className="px-3 py-1 bg-pink-500/20 text-pink-400 border border-pink-500/50 rounded font-mono text-xs">
              CONTRACT
            </span>
          )}
        </div>
      </div>

      {/* Overview Card */}
      <Card className="mb-6">
        <h2 className="text-2xl font-semibold mb-6 text-cyan-400">Overview</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div>
            <div className="text-xs text-zinc-500 mb-2">BALANCE</div>
            <div className="font-mono text-2xl text-cyan-400">
              {formatValue(addressData.balance)} ETH
            </div>
          </div>
          
          <div>
            <div className="text-xs text-zinc-500 mb-2">TRANSACTIONS</div>
            <div className="font-mono text-2xl text-purple-400">
              {addressData.tx_count.toLocaleString()}
            </div>
          </div>
          
          {!addressData.is_contract && (
            <div>
              <div className="text-xs text-zinc-500 mb-2">NONCE</div>
              <div className="font-mono text-2xl">
                {addressData.nonce}
              </div>
            </div>
          )}
          
          <div>
            <div className="text-xs text-zinc-500 mb-2">FIRST SEEN</div>
            <div className="font-mono text-sm">
              Block <Link 
                href={`/explorer/blocks/${addressData.first_seen_block}`}
                className="text-purple-400 hover:text-purple-300"
              >
                #{addressData.first_seen_block}
              </Link>
              <div className="text-zinc-500 text-xs mt-1">
                {formatTimestamp(addressData.first_seen_at)}
              </div>
            </div>
          </div>
          
          <div>
            <div className="text-xs text-zinc-500 mb-2">LAST SEEN</div>
            <div className="font-mono text-sm">
              Block <Link 
                href={`/explorer/blocks/${addressData.last_seen_block}`}
                className="text-purple-400 hover:text-purple-300"
              >
                #{addressData.last_seen_block}
              </Link>
              <div className="text-zinc-500 text-xs mt-1">
                {formatTimestamp(addressData.last_seen_at)}
              </div>
            </div>
          </div>
        </div>
      </Card>

      {/* Token Balances Card (if any) */}
      {addressData.token_balances && addressData.token_balances.length > 0 && (
        <Card className="mb-6">
          <h2 className="text-2xl font-semibold mb-4 text-cyan-400">
            Token Balances ({addressData.token_balances.length})
          </h2>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-zinc-800">
                  <th className="text-left py-3 px-4 text-xs text-zinc-500">TOKEN</th>
                  <th className="text-left py-3 px-4 text-xs text-zinc-500">TYPE</th>
                  <th className="text-right py-3 px-4 text-xs text-zinc-500">BALANCE</th>
                  <th className="text-left py-3 px-4 text-xs text-zinc-500">CONTRACT</th>
                </tr>
              </thead>
              <tbody>
                {addressData.token_balances.map((token, idx) => (
                  <tr key={idx} className="border-b border-zinc-800/50 hover:bg-zinc-900/50">
                    <td className="py-3 px-4">
                      <div className="font-semibold">{token.token_name}</div>
                      <div className="text-xs text-zinc-500">{token.token_symbol}</div>
                    </td>
                    <td className="py-3 px-4">
                      <span className="px-2 py-1 bg-purple-500/20 text-purple-400 border border-purple-500/50 rounded text-xs">
                        {token.token_type}
                      </span>
                    </td>
                    <td className="py-3 px-4 text-right font-mono">
                      {formatTokenBalance(token.balance, token.decimals)}
                    </td>
                    <td className="py-3 px-4">
                      <Link 
                        href={`/explorer/address/${token.token_address}`}
                        className="font-mono text-xs text-cyan-400 hover:text-cyan-300"
                      >
                        {token.token_address.slice(0, 10)}...{token.token_address.slice(-8)}
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>
      )}

      {/* Transactions Card */}
      <Card>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-semibold text-cyan-400">
            Transactions
          </h2>
          <div className="text-sm text-zinc-500">
            Page {txPage} of {totalTxPages}
          </div>
        </div>

        {transactions.length > 0 ? (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-zinc-800">
                    <th className="text-left py-3 px-4 text-xs text-zinc-500">TX HASH</th>
                    <th className="text-left py-3 px-4 text-xs text-zinc-500">BLOCK</th>
                    <th className="text-left py-3 px-4 text-xs text-zinc-500">AGE</th>
                    <th className="text-left py-3 px-4 text-xs text-zinc-500">FROM/TO</th>
                    <th className="text-right py-3 px-4 text-xs text-zinc-500">VALUE</th>
                    <th className="text-center py-3 px-4 text-xs text-zinc-500">STATUS</th>
                  </tr>
                </thead>
                <tbody>
                  {transactions.map((tx) => (
                    <tr key={tx.hash} className="border-b border-zinc-800/50 hover:bg-zinc-900/50">
                      <td className="py-3 px-4">
                        <Link 
                          href={`/explorer/tx/${tx.hash}`}
                          className="font-mono text-cyan-400 hover:text-cyan-300 text-sm"
                        >
                          {tx.hash.slice(0, 10)}...{tx.hash.slice(-8)}
                        </Link>
                      </td>
                      <td className="py-3 px-4">
                        <Link 
                          href={`/explorer/blocks/${tx.block_number}`}
                          className="font-mono text-purple-400 hover:text-purple-300 text-sm"
                        >
                          {tx.block_number}
                        </Link>
                      </td>
                      <td className="py-3 px-4 text-sm text-zinc-400">
                        {formatTimestamp(tx.timestamp)}
                      </td>
                      <td className="py-3 px-4">
                        <div className="text-xs space-y-1">
                          <div className="flex items-center gap-2">
                            <span className="text-zinc-500">From:</span>
                            {tx.from_address === address ? (
                              <span className="text-yellow-400 font-mono">this address</span>
                            ) : (
                              <Link 
                                href={`/explorer/address/${tx.from_address}`}
                                className="font-mono text-cyan-400 hover:text-cyan-300"
                              >
                                {tx.from_address.slice(0, 6)}...{tx.from_address.slice(-4)}
                              </Link>
                            )}
                          </div>
                          <div className="flex items-center gap-2">
                            <span className="text-zinc-500">To:</span>
                            {tx.to_address ? (
                              tx.to_address === address ? (
                                <span className="text-yellow-400 font-mono">this address</span>
                              ) : (
                                <Link 
                                  href={`/explorer/address/${tx.to_address}`}
                                  className="font-mono text-cyan-400 hover:text-cyan-300"
                                >
                                  {tx.to_address.slice(0, 6)}...{tx.to_address.slice(-4)}
                                </Link>
                              )
                            ) : (
                              <span className="text-zinc-500">[Contract Creation]</span>
                            )}
                          </div>
                        </div>
                      </td>
                      <td className="py-3 px-4 text-right font-mono text-sm">
                        {formatValue(tx.value)} ETH
                      </td>
                      <td className="py-3 px-4 text-center">
                        {tx.status === 1 ? (
                          <span className="text-green-400">✓</span>
                        ) : (
                          <span className="text-red-400">✗</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {totalTxPages > 1 && (
              <div className="flex justify-center gap-2 mt-6">
                <button
                  onClick={() => setTxPage(p => Math.max(1, p - 1))}
                  disabled={txPage === 1}
                  className="cyber-button px-4 py-2 disabled:opacity-50"
                >
                  Previous
                </button>
                <button
                  onClick={() => setTxPage(p => Math.min(totalTxPages, p + 1))}
                  disabled={txPage === totalTxPages}
                  className="cyber-button px-4 py-2 disabled:opacity-50"
                >
                  Next
                </button>
              </div>
            )}
          </>
        ) : (
          <div className="text-center py-8 text-zinc-500">
            No transactions found for this address
          </div>
        )}
      </Card>

      {/* Back link */}
      <div className="mt-6 text-center">
        <Link 
          href="/explorer"
          className="text-cyan-400 hover:text-cyan-300"
        >
          ← Back to Explorer
        </Link>
      </div>
    </div>
  );
}