import Link from 'next/link';
import Card from '@/components/ui/Card';
import { formatWei } from '@/lib/utils';
import { loadBlockById, loadTransactions, loadTransactionByHash } from '@/lib/snapshot';

export default async function BlockDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const resolvedParams = await params;
  const blockId = resolvedParams.id;

  const blockData = await loadBlockById(blockId);

  if (!blockData.success) {
    return (
      <div className="container mx-auto px-4 py-12 text-center">
        <div className="bg-zinc-900 border border-zinc-800 rounded-lg p-8 max-w-2xl mx-auto">
          <h2 className="text-2xl font-bold mb-4 text-yellow-400">Block Not in Snapshot</h2>
          <p className="text-zinc-400 mb-6">Block #{blockId} is not part of the frozen snapshot.</p>
          <Link href="/explorer" className="text-cyan-400 hover:text-cyan-300">← Back to Explorer</Link>
        </div>
      </div>
    );
  }

  const block = blockData.data;

  // Get all transactions from transactions.json
  // Load individual transaction files
  const txHashes = block.transactions || [];
  const blockTxs = await Promise.all(
    txHashes.map(async (hash: string) => {
      try {
        const txData = await loadTransactionByHash(hash);
        return txData.success ? txData.data : null;
      } catch {
        return null;
      }
    })
  );
  const validTxs = blockTxs.filter(tx => tx !== null);

  const details = [
    { label: 'Block Height', value: block.block_number, color: 'purple' },
    { label: 'Timestamp', value: new Date(block.timestamp).toLocaleString() },
    { label: 'Transactions', value: `${block.tx_count} transactions`, color: 'pink' },
    { label: 'Miner', value: block.miner, mono: true, color: 'cyan' },
    { label: 'Gas Used', value: `${block.gas_used?.toLocaleString() || 0} (${((block.gas_used / block.gas_limit) * 100).toFixed(2)}%)`, mono: true },
    { label: 'Gas Limit', value: block.gas_limit?.toLocaleString() || 0, mono: true },
    { label: 'Base Fee', value: `${block.base_fee_per_gas || '0'} wei`, mono: true },
    { label: 'Hash', value: block.hash, mono: true, color: 'cyan' },
    { label: 'Parent Hash', value: block.parent_hash, mono: true, color: 'cyan' },
  ];

  // Show tr
  const hasTxHashes = Array.isArray(txHashes) && txHashes.length > 0;

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
              <div className="text-zinc-500 text-sm uppercase tracking-wider">{item.label}:</div>
              <div className={`${item.mono ? 'font-mono' : ''} ${item.color === 'purple' ? 'text-purple-400 font-bold' :
                  item.color === 'pink' ? 'text-pink-400 font-bold' :
                    item.color === 'cyan' ? 'text-cyan-400' : 'text-white'
                } break-all`}>
                {item.value}
              </div>
            </div>
          ))}
        </div>
      </Card>

      <Card>
        <h2 className="text-2xl mb-6 text-white font-normal">
          Transactions ({block.tx_count || 0})
        </h2>

        {validTxs.length > 0 ? (
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
                {validTxs.map((tx: any) => (
                  <tr key={tx.hash} className="border-b border-cyan-400/10 hover:bg-white/5">
                    <td className="py-3 px-4">
                      <Link href={`/explorer/tx/${tx.hash}`} className="text-cyan-400 hover:text-pink-400 font-mono text-sm">
                        {tx.hash.slice(0, 10)}...{tx.hash.slice(-8)}
                      </Link>
                    </td>
                    <td className="py-3 px-4">
                      <Link href={`/explorer/address/${tx.from_address}`} className="text-cyan-400 hover:text-pink-400 font-mono text-sm">
                        {tx.from_address.slice(0, 8)}...{tx.from_address.slice(-6)}
                      </Link>
                    </td>
                    <td className="py-3 px-4">
                      {tx.to_address ? (
                        <Link href={`/explorer/address/${tx.to_address}`} className="text-cyan-400 hover:text-pink-400 font-mono text-sm">
                          {tx.to_address.slice(0, 8)}...{tx.to_address.slice(-6)}
                        </Link>
                      ) : (
                        <span className="text-zinc-500 text-sm">[Contract Creation]</span>
                      )}
                    </td>
                    <td className="py-3 px-4 text-right font-mono text-sm">{formatWei(tx.value)} ETH</td>
                    <td className="py-3 px-4 text-right font-mono text-sm text-zinc-500">
                      {tx.gas_used?.toLocaleString() || 'N/A'}
                    </td>
                    <td className="py-3 px-4 text-center">
                      {tx.status === 1 ? <span className="text-green-400">✓</span> : <span className="text-red-400">✗</span>}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : hasTxHashes ? (
          <div>
            <div className="bg-yellow-500/10 border-l-4 border-yellow-500 p-3 mb-4 text-sm text-yellow-300">
              This block has {txHashes.length} transactions, but detailed data is only in snapshot for the most recent 100 transactions across all blocks.
            </div>
            <div className="space-y-2">
              <div className="text-sm text-zinc-500 mb-2">Transaction Hashes:</div>
              {txHashes.slice(0, 20).map((hash: string, idx: number) => (
                <div key={idx} className="font-mono text-sm text-cyan-400 py-1">
                  {hash}
                </div>
              ))}
              {txHashes.length > 20 && (
                <div className="text-sm text-zinc-500 mt-2">... and {txHashes.length - 20} more</div>
              )}
            </div>
          </div>
        ) : (
          <div className="text-zinc-500 text-center py-8">
            {block.tx_count > 0 ? 'Transaction details not in snapshot' : 'No transactions in this block'}
          </div>
        )}
      </Card>
    </div>
  );
}