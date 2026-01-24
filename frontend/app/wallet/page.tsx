'use client';

import Link from 'next/link';

export default function WalletHome() {
  return (
    <div className="container mx-auto px-4 py-12">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-12 text-center">
          <h1 className="text-5xl font-bold mb-4 bg-gradient-to-r from-cyan-400 to-purple-400 bg-clip-text text-transparent">
            ETHEREUM WALLET
          </h1>
          <p className="text-zinc-400 text-lg">
            Create or import a wallet to get started
          </p>
        </div>

        {/* Action Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Create Wallet */}
          <Link href="/wallet/create">
            <div className="cyber-card hover:scale-105 transition-transform cursor-pointer">
              <div className="text-center">
                <div className="text-4xl mb-4">🔐</div>
                <h2 className="text-2xl font-bold mb-2">Create New Wallet</h2>
                <p className="text-zinc-400 mb-4">
                  Generate a new wallet with a 12-word recovery phrase
                </p>
                <button className="aggressive-btn w-full">
                  Create Wallet
                </button>
              </div>
            </div>
          </Link>

          {/* Import Wallet */}
          <Link href="/wallet/import">
            <div className="cyber-card hover:scale-105 transition-transform cursor-pointer">
              <div className="text-center">
                <div className="text-4xl mb-4">📥</div>
                <h2 className="text-2xl font-bold mb-2">Import Wallet</h2>
                <p className="text-zinc-400 mb-4">
                  Import using private key or 12-word recovery phrase
                </p>
                <button className="aggressive-btn w-full">
                  Import Wallet
                </button>
              </div>
            </div>
          </Link>
        </div>

        {/* Security Notice */}
        <div className="mt-12 cyber-card border-red-500/50">
          <h3 className="text-xl font-bold mb-3 text-red-400">🔴 SECURITY WARNING</h3>
          <ul className="space-y-2 text-sm text-zinc-300">
            <li>• This wallet stores encrypted keys <strong>on the server</strong></li>
            <li>• <strong>NOT SAFE</strong> for production use or real funds</li>
            <li>• Built for <strong>learning and development</strong> only</li>
            <li>• For production, use MetaMask, Ledger, or hardware wallets</li>
          </ul>
        </div>
      </div>
    </div>
  );
}