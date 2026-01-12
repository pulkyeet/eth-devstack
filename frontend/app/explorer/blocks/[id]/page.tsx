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
  transactions: string[];
}

export default function BlockDetailPage() {
  const params = useParams();
  const [block, setBlock] = useState<BlockDetail | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`http://localhost:8080/api/v1/blocks/${params.id}`)
      .then((res) => res.json())
      .then((data) => {
        setBlock(data.data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [params.id]);

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

  const details = [
    { label: 'Block Height', value: block.block_number, mono: false, color: 'purple' },
    { label: 'Timestamp', value: new Date(block.timestamp).toLocaleString(), mono: false },
    { label: 'Transactions', value: `${block.tx_count} transactions`, mono: false, color: 'pink' },
    { label: 'Miner', value: block.miner, mono: true, color: 'cyan' },
    { label: 'Gas Used', value: `${block.gas_used.toLocaleString()} (${((block.gas_used / block.gas_limit) * 100).toFixed(2)}%)`, mono: true },
    { label: 'Gas Limit', value: block.gas_limit.toLocaleString(), mono: true },
    { label: 'Base Fee', value: `${block.base_fee_per_gas || '0'} ETH`, mono: true },
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
        <h2 className="text-2xl mb-6 text-white font-normal">Transactions</h2>
        
        {block.tx_count === 0 ? (
          <div className="text-[var(--text-dim)] text-center py-8">
            No transactions in this block
          </div>
        ) : (
          <div className="text-[var(--cyan)]">
            {block.tx_count} transaction(s) - Transaction detail page coming soon
          </div>
        )}
      </Card>
    </div>
  );
}
