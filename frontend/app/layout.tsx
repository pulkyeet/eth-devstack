import type { Metadata } from 'next';
import './globals.css';
import Header from '@/components/layout/Header';
import CursorGlow from '../components/CursorGlow';

export const metadata: Metadata = {
  title: 'ETH DevStack',
  description: 'Complete Ethereum development environment',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <CursorGlow />
        <Header />
        <main className="min-h-screen relative z-10">
          {children}
        </main>
      </body>
    </html>
  );
}
