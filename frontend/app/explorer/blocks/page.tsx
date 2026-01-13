'use client';

import { useState, useEffect } from 'react';
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
}

function timeAgo(timestamp: string): string {
  const seconds = Math.floor((Date.now() - new Date(timestamp).getTime()) / 1000);
  if (seconds < 60) return `${seconds}s ago`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  return `${hours}h ago`;
}

export default function BlocksPage() {
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [loading, setLoading] = useState(true);
  const limit = 20;

  useEffect(() => {
    setLoading(true);
    fetch(`http://localhost:8080/api/v1/blocks?chain_id=1337&page=${page}&limit=${limit}`)
      .then((res) => res.json())
      .then((data) => {
        setBlocks(data.data.blocks || []);
        setTotalPages(data.data.pagination?.total_pages || 0);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [page]);

  if (loading) {
    return (
      <div className="container mx-auto px-6 py-12">
        <div className="text-[var(--cyan)] text-center text-2xl animate-pulse">LOADING...</div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-6 py-8">
      <div className="mb-8 flex justify-between items-center">
        <h1 className="text-4xl font-black uppercase tracking-tight text-white">
          ALL BLOCKS
        </h1>
        <div className="text-[var(--text-dim)]">
          Page {page} of {totalPages}
        </div>
      </div>

      <div className="space-y-3">
        {blocks.map((block) => {
          const gasPercent = ((block.gas_used / block.gas_limit) * 100).toFixed(1);
          
          return (
            <Link key={block.hash} href={`/explorer/blocks/${block.block_number}`}>
              <Card>
                <div className="grid grid-cols-[120px_1fr_100px_120px_120px] gap-6 items-center">
                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1">BLOCK</div>
                    <div className="text-xl font-bold highlight-purple">
                      #{block.block_number}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1">HASH</div>
                    <div className="mono text-sm text-[var(--cyan)]">
                      {block.hash.slice(0, 20)}...{block.hash.slice(-8)}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1">AGE</div>
                    <div className="text-sm font-semibold">{timeAgo(block.timestamp)}</div>
                  </div>

                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1">TXS</div>
                    <div className="text-lg font-bold text-[var(--pink)]">{block.tx_count}</div>
                  </div>

                  <div>
                    <div className="text-xs text-[var(--text-dim)] mb-1">GAS</div>
                    <div className="text-sm">
                      <span className="text-[var(--purple)] font-bold">{gasPercent}%</span>
                    </div>
                  </div>
                </div>
              </Card>
            </Link>
          );
        })}
      </div>

      {/* Pagination */}
      <div className="mt-8 flex justify-center gap-2">
        <button
          onClick={() => setPage(p => Math.max(1, p - 1))}
          disabled={page === 1}
          className="px-6 py-2 bg-[var(--card-bg)] border border-[var(--border)] hover:border-[var(--cyan)] disabled:opacity-30 disabled:cursor-not-allowed transition-colors font-bold uppercase"
        >
          PREV
        </button>
        
        <div className="px-6 py-2 bg-[var(--card-bg)] border border-[var(--cyan)] font-bold">
          {page}
        </div>
        
        <button
          onClick={() => setPage(p => Math.min(totalPages, p + 1))}
          disabled={page === totalPages}
          className="px-6 py-2 bg-[var(--card-bg)] border border-[var(--border)] hover:border-[var(--cyan)] disabled:opacity-30 disabled:cursor-not-allowed transition-colors font-bold uppercase"
        >
          NEXT
        </button>
      </div>
    </div>
  );
}