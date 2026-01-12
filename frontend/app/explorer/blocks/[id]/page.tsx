'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import Card from '@/components/ui/Card';

interface BlockDetail {
  block_number: number;
  hash: string;
  parent_hash: string;
  timestamp: string;
  miner: string;
  gas_used: number;
  gas_limit: number;
  base_fee_per_gas: string;
  difficulty: string;
  nonce: string;
  size: number;
  tx_count: number;
}

interface Transaction {
  hash: string;
  from_address: string;
  to_address: string | null;
  value: string;
  gas_used: number;
  status: number;
}

export default function BlockDetailPage() {
  const params = useParams();
  const blockId = params.id as string;
  
  const [block, setBlock] = useState<BlockDetail | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      fetch(`http://localhost:8080/api/v1/blocks/${blockId}?chain_id=1337`).then(r => r.json()),
      fetch(`http://localhost:8080/api/v1/transactions?chain_id=1337&block_number=${blockId}&limit=100`).then(r => r.json())
    ])
    .then(([blockData, txData]) => {
      if (blockData.success) {
        setBlock(blockData.data);
      }
      if (txData.success && txData.data.transactions) {
        setTransactions(txData.data.transactions);
      }
      setLoading(false);
    })
    .catch(() => setLoading(false));
  }, [blockId]);

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-[var(--cyan)] text-center text-xl animate-pulse">LOADING...</div>
      </div>
    );
  }

  if (!block) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-[var(--pink)] text-center text-xl">BLOCK NOT FOUND</div>
      </div>
    );
  }

  const formatValue = (wei: string) => {
    const eth = BigInt(wei) / BigInt(10 ** 18);
    return eth.toString();
  };

  const details = [
    { label: 'Block Height', value: block.block_number, mono: false, color: 'purple' },
    { label: 'Timestamp', value: new Date(block.timestamp).toLocaleString(), mono: false },
    { label: 'Transactions', value: `${block.tx_count} transactions`, mono: false, color: 'pink' },
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
        <Link href="/explorer" className="text-[var(--cyan)] hover:text-[var(--pink)] transition-colors">
          ← BACK TO BLOCKS
        </Link>
      </div>

      <h1 className="text-4xl mb-2 text-[var(--cyan)]">Block #{block.block_number}</h1>

      <Card className="mb-8">
        <h2 className="text-2xl mb-6 text-white font-normal">Details</h2>
        
        <div className="space-y-4">
          {details.map((item, idx) => (
            <div key={idx} className="grid grid-cols-[200px_1fr] gap-4 pb-4 border-b border-[var(--cyan)] border-opacity-20">
              <div className="text-[var(--text-dim)] text-sm uppercase tracking-wider">
                {item.label}:
              </div>
              <div className={`${item.mono ? 'mono' : ''} ${
                item.color === 'purple' ? 'highlight-purple font-bold' :
                item.color === 'pink' ? 'text-[var(--pink)] font-bold' :
                item.color === 'cyan' ? 'text-[var(--cyan)]' :
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
          <div className="text-[var(--text-dim)] text-center py-8">
            No transactions in this block
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-[var(--cyan)] border-opacity-30">
                  <th className="text-left py-3 px-4 text-xs text-[var(--text-dim)] uppercase">TX HASH</th>
                  <th className="text-left py-3 px-4 text-xs text-[var(--text-dim)] uppercase">FROM</th>
                  <th className="text-left py-3 px-4 text-xs text-[var(--text-dim)] uppercase">TO</th>
                  <th className="text-right py-3 px-4 text-xs text-[var(--text-dim)] uppercase">VALUE</th>
                  <th className="text-right py-3 px-4 text-xs text-[var(--text-dim)] uppercase">GAS</th>
                  <th className="text-center py-3 px-4 text-xs text-[var(--text-dim)] uppercase">STATUS</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((tx) => (
                  <tr key={tx.hash} className="border-b border-[var(--cyan)] border-opacity-10 hover:bg-white/5">
                    <td className="py-3 px-4">
                      <Link 
                        href={`/explorer/tx/${tx.hash}`}
                        className="text-[var(--cyan)] hover:text-[var(--pink)] mono text-sm"
                      >
                        {tx.hash.slice(0, 10)}...{tx.hash.slice(-8)}
                      </Link>
                    </td>
                    <td className="py-3 px-4">
                      <Link 
                        href={`/explorer/address/${tx.from_address}`}
                        className="text-[var(--cyan)] hover:text-[var(--pink)] mono text-sm"
                      >
                        {tx.from_address.slice(0, 8)}...{tx.from_address.slice(-6)}
                      </Link>
                    </td>
                    <td className="py-3 px-4">
                      {tx.to_address ? (
                        <Link 
                          href={`/explorer/address/${tx.to_address}`}
                          className="text-[var(--cyan)] hover:text-[var(--pink)] mono text-sm"
                        >
                          {tx.to_address.slice(0, 8)}...{tx.to_address.slice(-6)}
                        </Link>
                      ) : (
                        <span className="text-[var(--text-dim)] text-sm">[Contract Creation]</span>
                      )}
                    </td>
                    <td className="py-3 px-4 text-right mono text-sm">
                      {formatValue(tx.value)} ETH
                    </td>
                    <td className="py-3 px-4 text-right mono text-sm text-[var(--text-dim)]">
                      {tx.gas_used.toLocaleString()}
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