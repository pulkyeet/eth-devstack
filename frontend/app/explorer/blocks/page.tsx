import Link from 'next/link';
import Card from '@/components/ui/Card';
import { loadBlocks } from '@/lib/snapshot';

function timeAgo(timestamp: string): string {
  const seconds = Math.floor((Date.now() - new Date(timestamp).getTime()) / 1000);
  if (seconds < 60) return `${seconds}s ago`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  return `${hours}h ago`;
}

export default async function BlocksPage() {
  const data = await loadBlocks();
  const blocks = data.data.blocks || [];
  const pagination = data.data.pagination || {};

  return (
    <div className="container mx-auto px-6 py-8">
      {/* Snapshot Notice */}
      <div className="mb-6 bg-yellow-500/10 border border-yellow-500/30 p-4 rounded">
        <p className="text-sm text-yellow-300">
          📸 Showing {blocks.length} most recent blocks from snapshot.
          <a href="https://github.com/YOUR_REPO" className="ml-2 text-cyan-400 hover:text-cyan-300 underline">
            Run locally
          </a> to browse all {pagination.total?.toLocaleString() || 0} blocks.
        </p>
      </div>

      <div className="mb-8 flex justify-between items-center">
        <h1 className="text-4xl font-black uppercase tracking-tight text-white">
          RECENT BLOCKS
        </h1>
      </div>

      <div className="space-y-3">
        {blocks.map((block: any) => {
          const gasPercent = ((block.gas_used / block.gas_limit) * 100).toFixed(1);

          return (
            <Link key={block.hash} href={`/explorer/blocks/${block.block_number}`}>
              <Card>
                <div className="grid grid-cols-[120px_1fr_100px_120px_120px] gap-6 items-center">
                  <div>
                    <div className="text-xs text-zinc-500 mb-1">BLOCK</div>
                    <div className="text-xl font-bold text-purple-400">
                      #{block.block_number}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1">HASH</div>
                    <div className="font-mono text-sm text-cyan-400">
                      {block.hash.slice(0, 20)}...{block.hash.slice(-8)}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1">AGE</div>
                    <div className="text-sm font-semibold">{timeAgo(block.timestamp)}</div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1">TXS</div>
                    <div className="text-lg font-bold text-pink-400">{block.tx_count}</div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1">GAS</div>
                    <div className="text-sm">
                      <span className="text-purple-400 font-bold">{gasPercent}%</span>
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