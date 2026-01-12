'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

export default function Header() {
  const pathname = usePathname();

  const tabs = [
    { name: 'BLOCKS', path: '/explorer' },
    { name: 'SEARCH', path: '/explorer/search' },
    { name: 'STATS', path: '/explorer/stats' },
  ];

  return (
    <header className="py-6 mb-12">
      <div className="container mx-auto px-6 flex items-center justify-between">
        <Link href="/explorer" className="text-3xl font-black text-[var(--cyan)] mono tracking-tight hover:text-[var(--pink)] transition-colors">
          ETH_EXPLORER
        </Link>

        <nav className="flex gap-4">
          {tabs.map((tab) => {
            const isActive = pathname === tab.path;
            return (
              <Link
                key={tab.path}
                href={tab.path}
                className={`aggressive-btn ${isActive ? 'active' : ''}`}
              >
                {tab.name}
              </Link>
            );
          })}
        </nav>
      </div>
    </header>
  );
}