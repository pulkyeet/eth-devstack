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
          <form onSubmit={handleSearch} className="flex-1 max-w-3xl mx-8">
            <div className="flex items-stretch gap-3">
              <div className="relative flex-1">
                {/* Corner accents */}
                <div className="absolute top-0 left-0 w-2 h-2 border-l border-t border-[var(--cyan)]"
                  style={{ clipPath: 'polygon(0 0, 100% 0, 0 100%)' }}></div>
                <div className="absolute bottom-0 right-0 w-2 h-2 border-r border-b border-[var(--cyan)]"
                  style={{ clipPath: 'polygon(100% 0, 100% 100%, 0 100%)' }}></div>

                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Search by Address / Tx Hash / Block"
                  className="w-full h-full bg-black/70 border border-[var(--cyan)] px-4 py-3
                 text-sm font-mono text-white placeholder-zinc-500
                 focus:outline-none focus:border-[var(--purple)] focus:shadow-[0_0_15px_rgba(180,0,255,0.4)]
                 transition-all"
                  style={{ clipPath: 'polygon(10px 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%, 0 10px)' }}
                />
              </div>
              <span>
                <button type="submit" className="aggressive-btn" style={{ padding: '0 24px', height: '100%' }}>
                  SEARCH
                </button>
              </span>
            </div>
          </form>

          {/* Nav Buttons */}
          <div className="flex items-center gap-3">
            <Link href="/explorer">
              <button className="aggressive-btn">Explorer</button>
            </Link>
            <Link href="/explorer/blocks">
              <button className="aggressive-btn">Blocks</button>
            </Link>
            <Link href="/explorer/transactions">
              <button className="aggressive-btn">Transactions</button>
            </Link>
          </div>
        </div>
      </div>

      {/* Pink accent bar */}
      <div className="h-0.5 bg-gradient-to-r from-transparent via-pink-500 to-transparent"></div>
    </header>
  );
}