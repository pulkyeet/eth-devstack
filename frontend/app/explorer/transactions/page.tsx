import Link from 'next/link';
import Card from '@/components/ui/Card';
import { formatTimestamp, formatWei, shortenHash } from '@/lib/utils';
import { loadTransactions } from '@/lib/snapshot';

export default async function TransactionsPage() {
    const data = await loadTransactions();
    const transactions = data.data.transactions || [];
    const pagination = data.data.pagination || {};

    return (
        <div className="container mx-auto px-6 py-8">
            {/* Snapshot Notice */}
            <div className="mb-6 bg-yellow-500/10 border border-yellow-500/30 p-4 rounded">
                <p className="text-sm text-yellow-300">
                    📸 Showing {transactions.length} most recent transactions from snapshot.
                    <a href="https://github.com/YOUR_REPO" className="ml-2 text-cyan-400 hover:text-cyan-300 underline">
                        Run locally
                    </a> to browse all {pagination.total?.toLocaleString() || 0} transactions.
                </p>
            </div>

            <div className="mb-8 flex justify-between items-center">
                <h1 className="text-4xl font-black uppercase tracking-tight text-white">
                    RECENT TRANSACTIONS
                </h1>
            </div>

            <div className="space-y-3">
                {transactions.map((tx: any) => {
                    const statusColor = tx.status === 1 ? 'text-green-400' : tx.status === 0 ? 'text-red-400' : 'text-yellow-400';
                    const statusText = tx.status === 1 ? 'SUCCESS' : tx.status === 0 ? 'FAILED' : 'PENDING';

                    return (
                        <Link key={tx.hash} href={`/explorer/tx/${tx.hash}`}>
                            <Card>
                                <div className="grid grid-cols-[140px_80px_1fr_1fr_120px_100px_80px] gap-6 items-center">
                                    <div>
                                        <div className="text-xs text-zinc-500 mb-1">HASH</div>
                                        <div className="font-mono text-sm text-cyan-400">
                                            {shortenHash(tx.hash)}
                                        </div>
                                    </div>

                                    <div>
                                        <div className="text-xs text-zinc-500 mb-1">BLOCK</div>
                                        <Link href={`/explorer/blocks/${tx.block_number}`}>
                                            <div className="text-sm font-bold text-purple-400 hover:text-purple-300">
                                                #{tx.block_number}
                                            </div>
                                        </Link>
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
                                            {tx.to_address ? shortenHash(tx.to_address) : (
                                                <span className="text-purple-400">Contract Creation</span>
                                            )}
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

                                    <div>
                                        <div className="text-xs text-zinc-500 mb-1">STATUS</div>
                                        <div className={`text-xs font-bold ${statusColor}`}>
                                            {statusText}
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