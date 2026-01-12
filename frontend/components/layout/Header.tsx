'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function Header() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchQuery.trim()) return;

    const query = searchQuery.trim();
    
    // Detect what type of query it is
    if (query.startsWith('0x')) {
      if (query.length === 66) {
        // 64 hex chars + 0x = transaction hash or block hash
        // Try transaction first (most common)
        router.push(`/explorer/tx/${query}`);
      } else if (query.length === 42) {
        // 40 hex chars + 0x = address
        router.push(`/explorer/address/${query}`);
      } else {
        // Invalid hex string
        alert('Invalid hash or address format');
      }
    } else if (/^\d+$/.test(query)) {
      // Numeric = block number
      router.push(`/explorer/blocks/${query}`);
    } else {
      alert('Invalid search query. Enter:\n- Block number (e.g., 123)\n- Block/tx hash (0x...)\n- Address (0x...)');
    }
    
    setSearchQuery('');
  };

  return (
    <header className="border-b border-pink-500/30 bg-black/50 backdrop-blur-sm sticky top-0 z-50">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link href="/explorer" className="flex items-center gap-3">
            <div className="text-2xl font-bold bg-gradient-to-r from-cyan-400 to-purple-400 bg-clip-text text-transparent tracking-wider">
              ETH_EXPLORER
            </div>
          </Link>

          {/* Search Bar */}
          <form onSubmit={handleSearch} className="flex-1 max-w-2xl mx-8">
            <div className="relative">
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search by Address / Tx Hash / Block #"
                className="w-full bg-black/50 border border-cyan-500/30 rounded px-4 py-2 pr-10 
                         text-sm font-mono text-white placeholder-zinc-500
                         focus:outline-none focus:border-cyan-400 focus:ring-1 focus:ring-cyan-400
                         transition-all"
              />
              <button
                type="submit"
                className="absolute right-2 top-1/2 -translate-y-1/2 text-cyan-400 hover:text-cyan-300"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </button>
            </div>
          </form>

          {/* Nav Buttons */}
          <div className="flex items-center gap-3">
            <Link href="/explorer" className="cyber-button px-4 py-2">
              Blocks
            </Link>
            <Link href="/explorer" className="cyber-button px-4 py-2">
              Transactions
            </Link>
          </div>
        </div>
      </div>

      {/* Pink accent bar */}
      <div className="h-0.5 bg-gradient-to-r from-transparent via-pink-500 to-transparent"></div>
    </header>
  );
}