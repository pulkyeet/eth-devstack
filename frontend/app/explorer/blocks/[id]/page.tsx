'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import Card from '@/components/ui/Card';
import LoadingSpinner from '@/components/ui/LoadingSpinner';
import ErrorMessage from '@/components/ui/ErrorMessage';
import { formatWei } from '@/lib/utils';
import type { Block, Transaction } from '@/lib/types';

const CHAIN_ID = 1337;
const API_BASE = 'http://localhost:8080/api/v1';

export default function BlockDetailPage() {
  const params = useParams();
  const blockId = params.id as string;
  
  const [block, setBlock] = useState<Block | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    Promise.all([
      fetch(`${API_BASE}/blocks/${blockId}?chain_id=${CHAIN_ID}`).then(r => r.json()),
      fetch(`${API_BASE}/transactions?chain_id=${CHAIN_ID}&block_number=${blockId}&limit=100`).then(r => r.json())
    ])
    .then(([blockData, txData]) => {
      if (blockData.success) {
        setBlock(blockData.data);
      } else {
        setError('Block not found');
      }
      if (txData.success && txData.data.transactions) {
        setTransactions(txData.data.transactions);
      }
      setLoading(false);
    })
    .catch(() => {
      setError('Failed to load block');
      setLoading(false);
    });
  }, [blockId]);

  if (loading) return <LoadingSpinner message="LOADING BLOCK" />;
  if (error || !block) return <ErrorMessage message={error || 'Block not found'} />;

  const details = [
    { label: 'Block Height', value: block.block_number, color: 'purple' },
    { label: 'Timestamp', value: new Date(block.timestamp).toLocaleString() },
    { label: 'Transactions', value: `${block.tx_count} transactions`, color: 'pink' },
    { label: 'Miner', value: block.miner, mono: true, color: 'cyan' },
    { label: 'Gas Used', value: `${block.gas_used.toLocaleString()} (${((block.gas_used / block.gas_limit) * 100).toFixed(2)}%)`, mono: true },
    { label: 'Gas Limit', value: block.gas_limit.toLocaleString(), mono: true },
    { label: 'Base Fee', value: `${block.base_fee_per_gas || '0'} wei`, mono: true },
    { label: 'Hash', value: block.hash, mono: true, color: 'cyan' },
    { label: 'Parent Hash', value: block.parent_hash, mono: true, color: 'cyan' },
  ];

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      <div className="mb-8">
        <Link href="/explorer" className="text-cyan-400 hover:text-pink-400 transition-colors">
          ← BACK TO BLOCKS
        </Link>
      </div>

      <h1 className="text-4xl mb-2 text-cyan-400">Block #{block.block_number}</h1>

      <Card className="mb-8">
        <h2 className="text-2xl mb-6 text-white font-normal">Details</h2>
        
        <div className="space-y-4">
          {details.map((item, idx) => (
            <div key={idx} className="grid grid-cols-[200px_1fr] gap-4 pb-4 border-b border-cyan-400/20">
              <div className="text-zinc-500 text-sm uppercase tracking-wider">
                {item.label}:
              </div>
              <div className={`${item.mono ? 'font-mono' : ''} ${
                item.color === 'purple' ? 'text-purple-400 font-bold' :
                item.color === 'pink' ? 'text-pink-400 font-bold' :
                item.color === 'cyan' ? 'text-cyan-400' :
                'text-white'
              } break-all`}>
                {item.value}
              </div>
            </div>
          ))}
        </div>
      </Card>

      <Card>
        <h2 className="text-2xl mb-6 text-white font-normal">
          Transactions ({block.tx_count})
        </h2>
        
        {transactions.length === 0 ? (
          <div className="text-zinc-500 text-center py-8">
            No transactions in this block
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-cyan-400/30">
                  <th className="text-left py-3 px-4 text-xs text-zinc-500 uppercase">TX HASH</th>
                  <th className="text-left py-3 px-4 text-xs text-zinc-500 uppercase">FROM</th>
                  <th className="text-left py-3 px-4 text-xs text-zinc-500 uppercase">TO</th>
                  <th className="text-right py-3 px-4 text-xs text-zinc-500 uppercase">VALUE</th>
                  <th className="text-right py-3 px-4 text-xs text-zinc-500 uppercase">GAS</th>
                  <th className="text-center py-3 px-4 text-xs text-zinc-500 uppercase">STATUS</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((tx) => (
                  <tr key={tx.hash} className="border-b border-cyan-400/10 hover:bg-white/5">
                    <td className="py-3 px-4">
                      <Link 
                        href={`/explorer/tx/${tx.hash}`}
                        className="text-cyan-400 hover:text-pink-400 font-mono text-sm"
                      >
                        {tx.hash.slice(0, 10)}...{tx.hash.slice(-8)}
                      </Link>
                    </td>
                    <td className="py-3 px-4">
                      <Link 
                        href={`/explorer/address/${tx.from_address}`}
                        className="text-cyan-400 hover:text-pink-400 font-mono text-sm"
                      >
                        {tx.from_address.slice(0, 8)}...{tx.from_address.slice(-6)}
                      </Link>
                    </td>
                    <td className="py-3 px-4">
                      {tx.to_address ? (
                        <Link 
                          href={`/explorer/address/${tx.to_address}`}
                          className="text-cyan-400 hover:text-pink-400 font-mono text-sm"
                        >
                          {tx.to_address.slice(0, 8)}...{tx.to_address.slice(-6)}
                        </Link>
                      ) : (
                        <span className="text-zinc-500 text-sm">[Contract Creation]</span>
                      )}
                    </td>
                    <td className="py-3 px-4 text-right font-mono text-sm">
                      {formatWei(tx.value)} ETH
                    </td>
                    <td className="py-3 px-4 text-right font-mono text-sm text-zinc-500">
                      {tx.gas_used?.toLocaleString() || 'N/A'}
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
        )}
      </Card>
    </div>
  );
}