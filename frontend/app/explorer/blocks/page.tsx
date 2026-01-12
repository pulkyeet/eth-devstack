'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { getBlocks } from '@/lib/api';
import { formatTimestamp, shortenHash, formatGas } from '@/lib/utils';
import type { Block, Pagination } from '@/lib/types';
import LoadingSpinner from '@/components/ui/LoadingSpinner';
import ErrorMessage from '@/components/ui/ErrorMessage';
import Card from '@/components/ui/Card';

export default function BlockListPage() {
  const searchParams = useSearchParams();
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [pagination, setPagination] = useState<Pagination | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  const chainId = 1337;
  const page = Number(searchParams.get('page')) || 1;

  useEffect(() => {
    setLoading(true);
    getBlocks(chainId, page, 20)
      .then(data => {
        setBlocks(data.blocks);
        setPagination(data.pagination);
        setLoading(false);
      })
      .catch(err => {
        setError(err.message);
        setLoading(false);
      });
  }, [page]);

  if (loading) return <LoadingSpinner />;
  if (error) return (
    <div className="container mx-auto px-4 py-8">
      <ErrorMessage message={error} />
    </div>
  );

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold mb-8">Blocks</h1>

      <Card>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="border-b text-left text-sm text-zinc-500">
              <tr>
                <th className="pb-3">Block</th>
                <th className="pb-3">Time</th>
                <th className="pb-3">Miner</th>
                <th className="pb-3">Transactions</th>
                <th className="pb-3">Gas Used</th>
                <th className="pb-3">Gas Limit</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {blocks.map(block => (
                <tr key={`${block.block_number}-${block.hash}`} className="hover:bg-zinc-50 dark:hover:bg-zinc-900">
                  <td className="py-4">
                    <Link href={`/explorer/blocks/${block.block_number}`} className="font-mono text-blue-600 hover:underline">
                      {block.block_number}
                    </Link>
                  </td>
                  <td className="py-4 text-sm text-zinc-500">
                    {formatTimestamp(block.timestamp)}
                  </td>
                  <td className="py-4">
                    <Link href={`/explorer/address/${block.miner}`} className="font-mono text-sm text-blue-600 hover:underline">
                      {shortenHash(block.miner, 6)}
                    </Link>
                  </td>
                  <td className="py-4">{block.tx_count}</td>
                  <td className="py-4 font-mono text-sm">{formatGas(block.gas_used)}</td>
                  <td className="py-4 font-mono text-sm">{formatGas(block.gas_limit)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {pagination && pagination.total_pages > 1 && (
          <div className="mt-6 flex items-center justify-between border-t pt-4">
            <div className="text-sm text-zinc-500">
              Page {pagination.page} of {pagination.total_pages} ({pagination.total.toLocaleString()} total blocks)
            </div>
            <div className="flex gap-2">
              <Link
                href={`/explorer/blocks?page=${page - 1}`}
                className={`rounded-lg border px-4 py-2 text-sm ${
                  page <= 1 ? 'pointer-events-none opacity-50' : 'hover:bg-zinc-50 dark:hover:bg-zinc-900'
                }`}
              >
                Previous
              </Link>
              <Link
                href={`/explorer/blocks?page=${page + 1}`}
                className={`rounded-lg border px-4 py-2 text-sm ${
                  page >= pagination.total_pages ? 'pointer-events-none opacity-50' : 'hover:bg-zinc-50 dark:hover:bg-zinc-900'
                }`}
              >
                Next
              </Link>
            </div>
          </div>
        )}
      </Card>
    </div>
  );
}