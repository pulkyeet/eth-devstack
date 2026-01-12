import Link from 'next/link';

export default function Home() {
  return (
    <div className="container mx-auto px-4 py-16">
      <div className="max-w-3xl">
        <h1 className="text-5xl font-bold mb-6">ETH DevStack</h1>
        <p className="text-xl text-zinc-600 dark:text-zinc-400 mb-8">
          A complete Ethereum development environment with block explorer, wallet, and dApps.
        </p>
        
        <div className="grid gap-6 md:grid-cols-2">
          <Link
            href="/explorer"
            className="group rounded-lg border p-6 hover:border-blue-600 transition-colors"
          >
            <h2 className="text-2xl font-semibold mb-2 group-hover:text-blue-600">
              Block Explorer →
            </h2>
            <p className="text-zinc-600 dark:text-zinc-400">
              View blocks, transactions, and addresses on the local testnet.
            </p>
          </Link>
          
          <div className="rounded-lg border p-6 opacity-50">
            <h2 className="text-2xl font-semibold mb-2">Wallet (Coming Soon)</h2>
            <p className="text-zinc-600 dark:text-zinc-400">
              Manage accounts and send transactions.
            </p>
          </div>
          
          <div className="rounded-lg border p-6 opacity-50">
            <h2 className="text-2xl font-semibold mb-2">dApps (Coming Soon)</h2>
            <p className="text-zinc-600 dark:text-zinc-400">
              Interact with deployed smart contracts.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}