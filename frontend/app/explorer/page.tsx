'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import Card from '@/components/ui/Card';

interface Block {
  block_number: number;
  hash: string;
  timestamp: string;
  miner: string;
  tx_count: number;
  gas_used: number;
  gas_limit: number;
  size: number;
}

function timeAgo(timestamp: string): string {
  const seconds = Math.floor((Date.now() - new Date(timestamp).getTime()) / 1000);
  if (seconds < 60) return `${seconds}s ago`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  return `${hours}h ago`;
}

export default function ExplorerHome() {
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('http://localhost:8080/api/v1/blocks?limit=10')
      .then((res) => res.json())
      .then((data) => {
        setBlocks(data.data.blocks || []);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="container mx-auto px-6 py-12">
        <div className="text-[var(--cyan)] text-center text-2xl animate-pulse">LOADING...</div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-6 py-8">
      {/* Removed RECENT BLOCKS heading */}
      
      <div className="space-y-4">
        {blocks.map((block) => {
          const gasPercent = ((block.gas_used / block.gas_limit) * 100).toFixed(1);
          
          return (
            <Link key={block.hash} href={`/explorer/blocks/${block.block_number}`}>
              <Card>
                <div className="grid grid-cols-[140px_1fr_100px_140px_140px_120px_200px] gap-8 items-center">
                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1 uppercase tracking-wider">Block</div>
                    <div className="text-2xl font-bold highlight-purple">
                      #{block.block_number}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1 uppercase tracking-wider">Hash</div>
                    <div className="mono text-base text-[var(--cyan)]">
                      {block.hash.slice(0, 24)}...{block.hash.slice(-10)}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1 uppercase tracking-wider">Age</div>
                    <div className="text-base font-semibold text-white">{timeAgo(block.timestamp)}</div>
                  </div>

                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1 uppercase tracking-wider">Transactions</div>
                    <div className="text-xl font-bold text-[var(--pink)]">{block.tx_count}</div>
                  </div>

                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1 uppercase tracking-wider">Gas Used</div>
                    <div className="text-base">
                      <span className="text-[var(--purple)] font-bold">{gasPercent}%</span>
                      <span className="text-[var(--text-dim)] text-sm ml-1">
                        ({(block.gas_used / 1000000).toFixed(2)}M)
                      </span>
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1 uppercase tracking-wider">Size</div>
                    <div className="mono text-base">{(block.size / 1024).toFixed(2)} KB</div>
                  </div>

                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1 uppercase tracking-wider">Miner</div>
                    <div className="mono text-sm text-[var(--cyan)]">
                      {block.miner.slice(0, 12)}...{block.miner.slice(-8)}
                    </div>
                  </div>
                </div>
              </Card>
            </Link>
          );
        })}
      </div>
    </div>
  );
}
