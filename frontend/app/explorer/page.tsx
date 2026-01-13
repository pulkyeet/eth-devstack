'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import Card from '@/components/ui/Card';
import LoadingSpinner from '@/components/ui/LoadingSpinner';
import { createBlockStream } from '@/lib/api';
import { formatTimestamp } from '@/lib/utils';
import type { Block } from '@/lib/types';

const CHAIN_ID = 1337;

export default function ExplorerHome() {
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [loading, setLoading] = useState(true);

  // Initial fetch
  useEffect(() => {
    fetch(`http://localhost:8080/api/v1/blocks?chain_id=${CHAIN_ID}&limit=10`)
      .then((res) => res.json())
      .then((data) => {
        setBlocks(data.data.blocks || []);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

  // Real-time updates via SSE
  useEffect(() => {
    const eventSource = createBlockStream(CHAIN_ID);

    eventSource.addEventListener('block', (event) => {
      const newBlock: Block = JSON.parse(event.data);
      setBlocks((prev) => [newBlock, ...prev.slice(0, 9)]);
    });

    eventSource.onerror = () => eventSource.close();
    return () => eventSource.close();
  }, []);

  if (loading) return <LoadingSpinner message="LOADING BLOCKS" />;

  return (
    <div className="container mx-auto px-6 py-8">
      <div className="space-y-10">
        {blocks.map((block) => {
          const gasPercent = block.gas_limit ? ((block.gas_used / block.gas_limit) * 100).toFixed(1) : '0';

          return (
            <Link key={block.hash} href={`/explorer/blocks/${block.block_number}`}>
              <Card>
                <div className="grid grid-cols-[120px_1fr_90px_100px_160px_180px_180px] gap-6 items-center">
                  <div>
                    <div className="text-xs text-zinc-500 mb-1 uppercase tracking-wider">Block</div>
                    <div className="text-2xl font-bold text-purple-400">
                      #{block.block_number}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1 uppercase tracking-wider">Hash</div>
                    <div className="font-mono text-sm text-cyan-400">
                      {block.hash.slice(0, 20)}...{block.hash.slice(-8)}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1 uppercase tracking-wider">Age</div>
                    <div className="text-sm font-semibold text-white">{formatTimestamp(block.timestamp)}</div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1 uppercase tracking-wider">Txs</div>
                    <div className="text-xl font-bold text-pink-400">{block.tx_count}</div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1 uppercase tracking-wider">Gas Used</div>
                    <div className="space-y-2">
                      <div className="text-sm">
                        <span className="text-purple-400 font-bold">{gasPercent}%</span>
                        <span className="text-zinc-500 text-xs ml-1">
                          ({block.gas_limit ? (block.gas_limit / 1000000).toFixed(1) : 'N/A'}M)
                        </span>
                      </div>
                      <div className="font-mono text-xs text-zinc-500">
                        {block.gas_used.toLocaleString()} wei
                      </div>
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1 uppercase tracking-wider">Gas Limit</div>
                    <div className="space-y-2">
                      <div className="text-sm font-semibold">
                        {(block.gas_limit / 1000000).toFixed(1)}M
                      </div>
                      <div className="font-mono text-xs text-zinc-500">
                        {block.gas_limit?.toLocaleString() || 'N/A'} wei
                      </div>
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1 uppercase tracking-wider">Miner</div>
                    <div className="font-mono text-xs text-cyan-400">
                      {block.miner ? `${block.miner.slice(0, 10)}...${block.miner.slice(-6)}` : '0x00...0000'}
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