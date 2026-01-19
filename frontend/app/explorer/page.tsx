import Link from 'next/link';
import Card from '@/components/ui/Card';
import { formatTimestamp, formatWei, shortenHash } from '@/lib/utils';
import { loadStats, loadBlocks, loadTransactions } from '@/lib/snapshot';

export default async function ExplorerDashboard() {
  const statsData = await loadStats();
  const blocksData = await loadBlocks();
  const txsData = await loadTransactions();

  const stats = statsData.data;
  const recentBlocks = blocksData.data.blocks || [];
  const recentTxs = txsData.data.transactions || [];

  return (
    <div className="container mx-auto px-6 py-8 space-y-10">
      {/* Snapshot Banner */}
      <div className="bg-yellow-500/10 border-l-4 border-yellow-500 p-4 rounded">
        <p className="text-sm text-yellow-300">
          📸 <strong>Frozen Snapshot</strong> - Showing blocks {stats.start_block} to {stats.latest_block}
          <br />
          <span className="text-zinc-400">
            Backend runs locally. <a href="https://youtube.com/YOUR_VIDEO" className="text-cyan-400 hover:text-cyan-300 underline">Watch demo video</a> for full features.
          </span>
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-5 gap-6">
        <Card>
          <div className="text-xs text-zinc-500 uppercase tracking-wider mb-2">Latest Block</div>
          <div className="text-3xl font-bold text-cyan-400">
            #{stats.latest_block || 0}
          </div>
        </Card>

        <Card>
          <div className="text-xs text-zinc-500 uppercase tracking-wider mb-2">Total Transactions</div>
          <div className="text-3xl font-bold text-purple-400">
            {stats.total_transactions?.toLocaleString() || 0}
          </div>
        </Card>

        <Card>
          <div className="text-xs text-zinc-500 uppercase tracking-wider mb-2">Avg Block Time</div>
          <div className="text-3xl font-bold text-pink-400">
            {stats.avg_block_time || '0.0'}s
          </div>
        </Card>

        <Card>
          <div className="text-xs text-zinc-500 uppercase tracking-wider mb-2">TPS</div>
          <div className="text-3xl font-bold text-cyan-400">
            {stats.tps || 0}
          </div>
        </Card>
      </div>

      {/* Recent Blocks */}
      <div>
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-black uppercase tracking-tight text-white">
            Recent Blocks
          </h2>
          <Link href="/explorer/blocks">
            <button className="aggressive-btn text-sm">VIEW ALL</button>
          </Link>
        </div>

        <div className="space-y-3">
          {recentBlocks.slice(0, 5).map((block: any) => (
            <Link key={block.hash} href={`/explorer/blocks/${block.block_number}`}>
              <Card>
                <div className="grid grid-cols-[100px_1fr_120px_80px] gap-6 items-center">
                  <div>
                    <div className="text-xs text-zinc-500 mb-1">BLOCK</div>
                    <div className="text-xl font-bold text-purple-400">
                      #{block.block_number}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1">HASH</div>
                    <div className="font-mono text-sm text-cyan-400">
                      {shortenHash(block.hash)}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1">AGE</div>
                    <div className="text-sm font-semibold">{formatTimestamp(block.timestamp)}</div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1">TXS</div>
                    <div className="text-lg font-bold text-pink-400">{block.tx_count}</div>
                  </div>
                </div>
              </Card>
            </Link>
          ))}
        </div>
      </div>

      {/* Recent Transactions */}
      <div>
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-black uppercase tracking-tight text-white">
            Recent Transactions
          </h2>
          <Link href="/explorer/transactions">
            <button className="aggressive-btn text-sm">VIEW ALL</button>
          </Link>
        </div>

        <div className="space-y-3">
          {recentTxs.slice(0, 10).map((tx: any) => (
            <Link key={tx.hash} href={`/explorer/tx/${tx.hash}`}>
              <Card>
                <div className="grid grid-cols-[120px_1fr_1fr_120px_100px] gap-6 items-center">
                  <div>
                    <div className="text-xs text-zinc-500 mb-1">HASH</div>
                    <div className="font-mono text-sm text-cyan-400">
                      {shortenHash(tx.hash)}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1">FROM</div>
                    <div className="font-mono text-xs text-zinc-400">
                      {shortenHash(tx.from_address)}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1">TO</div>
                    <div className="font-mono text-xs text-zinc-400">
                      {tx.to_address ? shortenHash(tx.to_address) : 'Contract Creation'}
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1">VALUE</div>
                    <div className="text-sm font-semibold text-purple-400">
                      {formatWei(tx.value)} ETH
                    </div>
                  </div>

                  <div>
                    <div className="text-xs text-zinc-500 mb-1">AGE</div>
                    <div className="text-sm">{formatTimestamp(tx.timestamp)}</div>
                  </div>
                </div>
              </Card>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}