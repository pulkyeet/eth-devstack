import Link from 'next/link';
import Card from '@/components/ui/Card';
import { formatWei, formatTimestamp, formatTokenBalance } from '@/lib/utils';
import { loadAddressByAddress, loadTransactions } from '@/lib/snapshot';

export default async function AddressDetailPage({ params }: { params: { address: string } }) {
  const resolvedParams = await params;
  const result = await loadAddressByAddress(resolvedParams.address);
  
  if (!result.success) {
    return (
      <div className="container mx-auto px-4 py-12 text-center">
        <div className="bg-zinc-900 border border-zinc-800 rounded-lg p-8 max-w-2xl mx-auto">
          <h2 className="text-2xl font-bold mb-4 text-yellow-400">Address Not in Snapshot</h2>
          <p className="text-zinc-400 mb-6">
            Address {resolvedParams.address.slice(0, 10)}...{resolvedParams.address.slice(-8)} is not part of the frozen snapshot.
          </p>
          <div className="space-y-3 text-sm text-zinc-500">
            <p>
              <a href="https://github.com/YOUR_REPO" className="text-cyan-400 hover:text-cyan-300 underline">
                Clone the repo
              </a> and run locally to explore all addresses.
            </p>
            <p>Or <a href="https://youtube.com/YOUR_VIDEO" className="text-cyan-400 hover:text-cyan-300 underline">watch the demo video</a></p>
          </div>
          <div className="mt-6">
            <Link href="/explorer" className="text-cyan-400 hover:text-cyan-300">
              ← Back to Explorer
            </Link>
          </div>
        </div>
      </div>
    );
  }

  const addressData = result.data;
  
  // Load transactions involving this address
  const txData = await loadTransactions();
  const allTxs = txData.data.transactions || [];
  const transactions = allTxs.filter((tx: any) => 
    tx.from_address === resolvedParams.address || tx.to_address === resolvedParams.address
  ).slice(0, 20); // Show max 20 in snapshot

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-4xl font-bold mb-2 bg-gradient-to-r from-cyan-400 to-purple-400 bg-clip-text text-transparent">
          {addressData.is_contract ? 'Contract' : 'Address'}
        </h1>
        <div className="flex items-center gap-3">
          <span className="text-zinc-400 font-mono text-sm break-all">{resolvedParams.address}</span>
          {addressData.is_contract && (
            <span className="px-3 py-1 bg-pink-500/20 text-pink-400 border border-pink-500/50 rounded font-mono text-xs">
              CONTRACT
            </span>
          )}
        </div>
      </div>

      {/* Overview Card */}
      <Card className="mb-6">
        <h2 className="text-2xl font-semibold mb-6 text-cyan-400">Overview</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div>
            <div className="text-xs text-zinc-500 mb-2">BALANCE</div>
            <div className="font-mono text-2xl text-cyan-400">
              {formatWei(addressData.balance)} ETH
            </div>
          </div>
          
          <div>
            <div className="text-xs text-zinc-500 mb-2">TRANSACTIONS</div>
            <div className="font-mono text-2xl text-purple-400">
              {addressData.tx_count.toLocaleString()}
            </div>
          </div>
          
          {!addressData.is_contract && (
            <div>
              <div className="text-xs text-zinc-500 mb-2">NONCE</div>
              <div className="font-mono text-2xl">
                {addressData.nonce}
              </div>
            </div>
          )}
          
          {addressData.first_seen_block && (
            <div>
              <div className="text-xs text-zinc-500 mb-2">FIRST SEEN</div>
              <div className="font-mono text-sm">
                Block <Link 
                  href={`/explorer/blocks/${addressData.first_seen_block}`}
                  className="text-purple-400 hover:text-purple-300"
                >
                  #{addressData.first_seen_block}
                </Link>
                {addressData.first_seen_at && (
                  <div className="text-zinc-500 text-xs mt-1">
                    {formatTimestamp(addressData.first_seen_at)}
                  </div>
                )}
              </div>
            </div>
          )}
          
          {addressData.last_seen_block && (
            <div>
              <div className="text-xs text-zinc-500 mb-2">LAST SEEN</div>
              <div className="font-mono text-sm">
                Block <Link 
                  href={`/explorer/blocks/${addressData.last_seen_block}`}
                  className="text-purple-400 hover:text-purple-300"
                >
                  #{addressData.last_seen_block}
                </Link>
                {addressData.last_seen_at && (
                  <div className="text-zinc-500 text-xs mt-1">
                    {formatTimestamp(addressData.last_seen_at)}
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </Card>

      {/* Token Balances Card */}
      {addressData.token_balances && addressData.token_balances.length > 0 && (
        <Card className="mb-6">
          <h2 className="text-2xl font-semibold mb-4 text-cyan-400">
            Token Balances ({addressData.token_balances.length})
          </h2>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-zinc-800">
                  <th className="text-left py-3 px-4 text-xs text-zinc-500">TOKEN</th>
                  <th className="text-left py-3 px-4 text-xs text-zinc-500">TYPE</th>
                  <th className="text-right py-3 px-4 text-xs text-zinc-500">BALANCE</th>
                  <th className="text-left py-3 px-4 text-xs text-zinc-500">CONTRACT</th>
                </tr>
              </thead>
              <tbody>
                {addressData.token_balances.map((token: any, idx: number) => (
                  <tr key={idx} className="border-b border-zinc-800/50 hover:bg-zinc-900/50">
                    <td className="py-3 px-4">
                      <div className="font-semibold">{token.token_name}</div>
                      <div className="text-xs text-zinc-500">{token.token_symbol}</div>
                    </td>
                    <td className="py-3 px-4">
                      <span className="px-2 py-1 bg-purple-500/20 text-purple-400 border border-purple-500/50 rounded text-xs">
                        {token.token_type}
                      </span>
                    </td>
                    <td className="py-3 px-4 text-right font-mono">
                      {formatTokenBalance(token.balance, token.decimals)}
                    </td>
                    <td className="py-3 px-4">
                      <Link 
                        href={`/explorer/address/${token.token_address}`}
                        className="font-mono text-xs text-cyan-400 hover:text-cyan-300"
                      >
                        {token.token_address.slice(0, 10)}...{token.token_address.slice(-8)}
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>
      )}

      {/* Transactions Card */}
      <Card>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-semibold text-cyan-400">
            Recent Transactions
          </h2>
          {transactions.length > 0 && addressData.tx_count > transactions.length && (
            <div className="text-sm text-zinc-500">
              Showing {transactions.length} of {addressData.tx_count}
            </div>
          )}
        </div>

        {transactions.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-zinc-800">
                  <th className="text-left py-3 px-4 text-xs text-zinc-500">TX HASH</th>
                  <th className="text-left py-3 px-4 text-xs text-zinc-500">BLOCK</th>
                  <th className="text-left py-3 px-4 text-xs text-zinc-500">AGE</th>
                  <th className="text-left py-3 px-4 text-xs text-zinc-500">FROM/TO</th>
                  <th className="text-right py-3 px-4 text-xs text-zinc-500">VALUE</th>
                  <th className="text-center py-3 px-4 text-xs text-zinc-500">STATUS</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((tx: any) => (
                  <tr key={tx.hash} className="border-b border-zinc-800/50 hover:bg-zinc-900/50">
                    <td className="py-3 px-4">
                      <Link 
                        href={`/explorer/tx/${tx.hash}`}
                        className="font-mono text-cyan-400 hover:text-cyan-300 text-sm"
                      >
                        {tx.hash.slice(0, 10)}...{tx.hash.slice(-8)}
                      </Link>
                    </td>
                    <td className="py-3 px-4">
                      <Link 
                        href={`/explorer/blocks/${tx.block_number}`}
                        className="font-mono text-purple-400 hover:text-purple-300 text-sm"
                      >
                        {tx.block_number}
                      </Link>
                    </td>
                    <td className="py-3 px-4 text-sm text-zinc-400">
                      {formatTimestamp(tx.timestamp)}
                    </td>
                    <td className="py-3 px-4">
                      <div className="text-xs space-y-1">
                        <div className="flex items-center gap-2">
                          <span className="text-zinc-500">From:</span>
                          {tx.from_address === resolvedParams.address ? (
                            <span className="text-yellow-400 font-mono">this address</span>
                          ) : (
                            <Link 
                              href={`/explorer/address/${tx.from_address}`}
                              className="font-mono text-cyan-400 hover:text-cyan-300"
                            >
                              {tx.from_address.slice(0, 6)}...{tx.from_address.slice(-4)}
                            </Link>
                          )}
                        </div>
                        <div className="flex items-center gap-2">
                          <span className="text-zinc-500">To:</span>
                          {tx.to_address ? (
                            tx.to_address === resolvedParams.address ? (
                              <span className="text-yellow-400 font-mono">this address</span>
                            ) : (
                              <Link 
                                href={`/explorer/address/${tx.to_address}`}
                                className="font-mono text-cyan-400 hover:text-cyan-300"
                              >
                                {tx.to_address.slice(0, 6)}...{tx.to_address.slice(-4)}
                              </Link>
                            )
                          ) : (
                            <span className="text-zinc-500">[Contract Creation]</span>
                          )}
                        </div>
                      </div>
                    </td>
                    <td className="py-3 px-4 text-right font-mono text-sm">
                      {formatWei(tx.value)} ETH
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
        ) : (
          <div className="text-center py-8 text-zinc-500">
            No transactions found in snapshot for this address
          </div>
        )}
      </Card>

      {/* Back link */}
      <div className="mt-6 text-center">
        <Link 
          href="/explorer"
          className="text-cyan-400 hover:text-cyan-300"
        >
          ← Back to Explorer
        </Link>
      </div>
    </div>
  );
}