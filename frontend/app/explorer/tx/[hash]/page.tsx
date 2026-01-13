'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import Card from '@/components/ui/Card';
import LoadingSpinner from '@/components/ui/LoadingSpinner';
import ErrorMessage from '@/components/ui/ErrorMessage';
import { formatWei } from '@/lib/utils';
import type { Transaction, Log } from '@/lib/types';

const CHAIN_ID = 1337;
const API_BASE = 'http://localhost:8080/api/v1';

export default function TransactionDetailPage() {
  const params = useParams();
  const hash = params.hash as string;

  const [tx, setTx] = useState<Transaction | null>(null);
  const [logs, setLogs] = useState<Log[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetch(`${API_BASE}/transactions/${hash}?chain_id=${CHAIN_ID}`)
      .then(res => res.json())
      .then(data => {
        if (data.success) {
          setTx(data.data);
          setLogs(data.data.logs || []);
        } else {
          setError('Transaction not found');
        }
        setLoading(false);
      })
      .catch(() => {
        setError('Failed to load transaction');
        setLoading(false);
      });
  }, [hash]);

  if (loading) return <LoadingSpinner message="LOADING TRANSACTION" />;
  if (error || !tx) return <ErrorMessage message={error || 'Transaction not found'} />;

  const formatGas = (gas: string) => (Number(gas) / 1e9).toFixed(2);

  const getStatusBadge = (status: number) => {
    if (status === 1) {
      return <span className="px-3 py-1 bg-green-500/20 text-green-400 border border-green-500/50 rounded font-mono text-sm">✓ SUCCESS</span>;
    } else if (status === 0) {
      return <span className="px-3 py-1 bg-red-500/20 text-red-400 border border-red-500/50 rounded font-mono text-sm">✗ FAILED</span>;
    }
    return <span className="px-3 py-1 bg-yellow-500/20 text-yellow-400 border border-yellow-500/50 rounded font-mono text-sm">⏳ PENDING</span>;
  };

  const formatTimestamp = (ts: string) => new Date(ts).toLocaleString();

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-4xl font-bold mb-2 bg-gradient-to-r from-cyan-400 to-purple-400 bg-clip-text text-transparent">
          Transaction Details
        </h1>
        <div className="flex items-center gap-3">
          <span className="text-zinc-500 font-mono text-sm">{hash}</span>
          {getStatusBadge(tx.status || 2)}
        </div>
      </div>

      {/* Overview Card */}
      <Card className="mb-6">
        <h2 className="text-2xl font-semibold mb-6 text-cyan-400">Overview</h2>

        <div className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pb-4 border-b border-zinc-800">
            <div>
              <div className="text-xs text-zinc-500 mb-1">BLOCK</div>
              <Link
                href={`/explorer/blocks/${tx.block_number}`}
                className="text-purple-400 hover:text-purple-300 font-mono"
              >
                #{tx.block_number}
              </Link>
            </div>
            <div>
              <div className="text-xs text-zinc-500 mb-1">TIMESTAMP</div>
              <div className="font-mono text-sm">{formatTimestamp(tx.timestamp)}</div>
            </div>
            <div>
              <div className="text-xs text-zinc-500 mb-1">INDEX</div>
              <div className="font-mono text-sm">{tx.transaction_index}</div>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4 pb-4 border-b border-zinc-800">
            <div>
              <div className="text-xs text-zinc-500 mb-1">FROM</div>
              <Link
                href={`/explorer/address/${tx.from_address}`}
                className="text-cyan-400 hover:text-cyan-300 font-mono text-sm break-all"
              >
                {tx.from_address}
              </Link>
            </div>
            <div>
              <div className="text-xs text-zinc-500 mb-1">TO</div>
              {tx.to_address ? (
                <Link
                  href={`/explorer/address/${tx.to_address}`}
                  className="text-cyan-400 hover:text-cyan-300 font-mono text-sm break-all"
                >
                  {tx.to_address}
                </Link>
              ) : (
                <span className="text-yellow-400 font-mono text-sm">
                  [Contract Creation]
                </span>
              )}
            </div>
            {tx.contract_address && (
              <div>
                <div className="text-xs text-zinc-500 mb-1">CONTRACT CREATED</div>
                <Link
                  href={`/explorer/address/${tx.contract_address}`}
                  className="text-pink-400 hover:text-pink-300 font-mono text-sm break-all"
                >
                  {tx.contract_address}
                </Link>
              </div>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <div className="text-xs text-zinc-500 mb-1">VALUE</div>
              <div className="font-mono text-lg text-cyan-400">{formatWei(tx.value)} ETH</div>
            </div>
            <div>
              <div className="text-xs text-zinc-500 mb-1">GAS USED</div>
              <div className="font-mono text-sm">{tx.gas_used?.toLocaleString() || 'N/A'} / {tx.gas.toLocaleString()}</div>
            </div>
            <div>
              <div className="text-xs text-zinc-500 mb-1">GAS PRICE</div>
              <div className="font-mono text-sm">{tx.gas_price ? formatGas(tx.gas_price) : 'N/A'} Gwei</div>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4 border-t border-zinc-800">
            <div>
              <div className="text-xs text-zinc-500 mb-1">NONCE</div>
              <div className="font-mono text-sm">{tx.nonce}</div>
            </div>
            <div>
              <div>
                <div className="text-xs text-zinc-500 mb-1">TYPE</div>
                <div className="font-mono text-sm">
                  {tx.transaction_type === 0 ? 'Legacy' :
                    tx.transaction_type === 1 ? 'EIP-2930' :
                      tx.transaction_type === 2 ? 'EIP-1559' :
                        'Unknown'}
                </div>
              </div>
            </div>
          </div>
        </div>
      </Card>

      {/* Input Data Card */}
      {tx.input && tx.input !== '0x' && (
        <Card className="mb-6">
          <h2 className="text-2xl font-semibold mb-4 text-cyan-400">Input Data</h2>
          <div className="bg-black/50 p-4 rounded border border-zinc-800 overflow-x-auto">
            <pre className="font-mono text-xs text-zinc-400 whitespace-pre-wrap break-all">
              {tx.input}
            </pre>
          </div>
        </Card>
      )}

      {/* Logs Card */}
      {logs.length > 0 && (
        <Card>
          <h2 className="text-2xl font-semibold mb-4 text-cyan-400">
            Event Logs ({logs.length})
          </h2>
          <div className="space-y-4">
            {logs.map((log, idx) => (
              <div key={idx} className="bg-black/30 p-4 rounded border border-zinc-800">
                <div className="grid grid-cols-1 gap-2 mb-3">
                  <div>
                    <span className="text-xs text-zinc-500">Log Index: </span>
                    <span className="font-mono text-sm text-pink-400">{log.log_index}</span>
                  </div>
                  <div>
                    <span className="text-xs text-zinc-500">Address: </span>
                    <Link
                      href={`/explorer/address/${log.address}`}
                      className="font-mono text-sm text-cyan-400 hover:text-cyan-300"
                    >
                      {log.address}
                    </Link>
                  </div>
                </div>

                {log.topics.length > 0 && (
                  <div className="mb-2">
                    <div className="text-xs text-zinc-500 mb-2">Topics:</div>
                    {log.topics.map((topic, i) => (
                      <div key={i} className="font-mono text-xs text-zinc-400 mb-1 pl-4">
                        [{i}] {topic}
                      </div>
                    ))}
                  </div>
                )}

                {log.data && log.data !== '0x' && (
                  <div>
                    <div className="text-xs text-zinc-500 mb-1">Data:</div>
                    <div className="font-mono text-xs text-zinc-400 break-all bg-black/50 p-2 rounded">
                      {log.data}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </Card>
      )}

      {/* Back link */}
      <div className="mt-6 text-center">
        <Link
          href={`/explorer/blocks/${tx.block_number}`}
          className="text-cyan-400 hover:text-cyan-300"
        >
          ← Back to Block #{tx.block_number}
        </Link>
      </div>
    </div>
  );
}