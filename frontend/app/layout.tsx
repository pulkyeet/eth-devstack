'use client';

import { useEffect } from 'react';
import type { Metadata } from 'next';
import './globals.css';
import Header from '@/components/layout/Header';

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  useEffect(() => {
    // Cursor glow effect
    const glow = document.createElement('div');
    glow.id = 'cursor-glow';
    document.body.appendChild(glow);

    const handleMouseMove = (e: MouseEvent) => {
      glow.style.left = `${e.clientX - 150}px`;
      glow.style.top = `${e.clientY - 150}px`;
    };

    window.addEventListener('mousemove', handleMouseMove);

    return () => {
      window.removeEventListener('mousemove', handleMouseMove);
      glow.remove();
    };
  }, []);

  return (
    <html lang="en">
      <body>
        <Header />
        <main className="min-h-screen relative z-10">
          {children}
        </main>
      </body>
    </html>
  );
}
