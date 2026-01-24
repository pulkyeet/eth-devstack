'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { importWallet } from '@/lib/api';

export default function ImportWallet() {
  const router = useRouter();
  const [method, setMethod] = useState<'private_key' | 'mnemonic'>('private_key');
  const [name, setName] = useState('');
  const [privateKey, setPrivateKey] = useState('');
  const [mnemonic, setMnemonic] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleImport = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    // Validation
    if (!name.trim()) {
      setError('Wallet name is required');
      return;
    }
    if (password.length < 12) {
      setError('Password must be at least 12 characters');
      return;
    }
    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    if (method === 'private_key') {
      if (!privateKey.trim()) {
        setError('Private key is required');
        return;
      }
      if (!privateKey.startsWith('0x') || privateKey.length !== 66) {
        setError('Invalid private key format (must be 0x followed by 64 hex chars)');
        return;
      }
    } else {
      const words = mnemonic.trim().split(/\s+/);
      if (words.length !== 12) {
        setError('Mnemonic must be exactly 12 words');
        return;
      }
    }

    setLoading(true);
    try {
      const result = await importWallet({
        method,
        private_key: method === 'private_key' ? privateKey.trim() : undefined,
        mnemonic: method === 'mnemonic' ? mnemonic.trim() : undefined,
        password,
        name: name.trim(),
      });

      // Store wallet ID
      localStorage.setItem('current_wallet', result.id);
      
      // Redirect to dashboard
      router.push('/wallet/dashboard');
    } catch (err: any) {
      setError(err.message || 'Failed to import wallet');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-12">
      <div className="max-w-xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2 text-cyan-400">Import Wallet</h1>
          <p className="text-zinc-400">Import an existing wallet using private key or recovery phrase</p>
        </div>

        {/* Method Selector */}
        <div className="flex gap-4 mb-6">
          <button
            onClick={() => setMethod('private_key')}
            className={`aggressive-btn flex-1 ${method === 'private_key' ? 'bg-cyan-500 text-black' : ''}`}
          >
            Private Key
          </button>
          <button
            onClick={() => setMethod('mnemonic')}
            className={`aggressive-btn flex-1 ${method === 'mnemonic' ? 'bg-cyan-500 text-black' : ''}`}
          >
            Recovery Phrase
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleImport}>
          <div className="cyber-card mb-6">
            <div className="space-y-4">
              {/* Wallet Name */}
              <div>
                <label className="block text-sm font-semibold mb-2 text-cyan-400">
                  Wallet Name
                </label>
                <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Imported Wallet"
                  className="w-full bg-black/70 border border-cyan-500/50 px-4 py-3 text-white font-mono
                    focus:outline-none focus:border-cyan-400 focus:shadow-[0_0_15px_rgba(34,211,238,0.3)]
                    clip-angle-sm"
                  disabled={loading}
                />
              </div>

              {/* Private Key Input */}
              {method === 'private_key' && (
                <div>
                  <label className="block text-sm font-semibold mb-2 text-cyan-400">
                    Private Key
                  </label>
                  <textarea
                    value={privateKey}
                    onChange={(e) => setPrivateKey(e.target.value)}
                    placeholder="0x..."
                    rows={3}
                    className="w-full bg-black/70 border border-cyan-500/50 px-4 py-3 text-white font-mono text-sm
                      focus:outline-none focus:border-cyan-400 focus:shadow-[0_0_15px_rgba(34,211,238,0.3)]
                      clip-angle-sm resize-none"
                    disabled={loading}
                  />
                  <div className="mt-1 text-xs text-zinc-500">
                    64 hex characters starting with 0x
                  </div>
                </div>
              )}

              {/* Mnemonic Input */}
              {method === 'mnemonic' && (
                <div>
                  <label className="block text-sm font-semibold mb-2 text-cyan-400">
                    Recovery Phrase <span className="text-zinc-500">(12 words)</span>
                  </label>
                  <textarea
                    value={mnemonic}
                    onChange={(e) => setMnemonic(e.target.value)}
                    placeholder="word1 word2 word3 ..."
                    rows={4}
                    className="w-full bg-black/70 border border-cyan-500/50 px-4 py-3 text-white font-mono text-sm
                      focus:outline-none focus:border-cyan-400 focus:shadow-[0_0_15px_rgba(34,211,238,0.3)]
                      clip-angle-sm resize-none"
                    disabled={loading}
                  />
                  <div className="mt-1 text-xs text-zinc-500">
                    Enter all 12 words separated by spaces
                  </div>
                </div>
              )}

              {/* Password */}
              <div>
                <label className="block text-sm font-semibold mb-2 text-cyan-400">
                  Password <span className="text-zinc-500">(min 12 characters)</span>
                </label>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••••••"
                  className="w-full bg-black/70 border border-cyan-500/50 px-4 py-3 text-white font-mono
                    focus:outline-none focus:border-cyan-400 focus:shadow-[0_0_15px_rgba(34,211,238,0.3)]
                    clip-angle-sm"
                  disabled={loading}
                />
              </div>

              {/* Confirm Password */}
              <div>
                <label className="block text-sm font-semibold mb-2 text-cyan-400">
                  Confirm Password
                </label>
                <input
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  placeholder="••••••••••••"
                  className="w-full bg-black/70 border border-cyan-500/50 px-4 py-3 text-white font-mono
                    focus:outline-none focus:border-cyan-400 focus:shadow-[0_0_15px_rgba(34,211,238,0.3)]
                    clip-angle-sm"
                  disabled={loading}
                />
              </div>
            </div>
          </div>

          {/* Error */}
          {error && (
            <div className="mb-4 p-4 bg-red-900/20 border border-red-500/50 clip-angle-sm">
              <p className="text-red-400 text-sm">{error}</p>
            </div>
          )}

          {/* Security Warning */}
          <div className="mb-6 cyber-card border-yellow-500/50">
            <div className="text-sm text-yellow-400">
              ⚠️ Never share your private key or recovery phrase with anyone. 
              This wallet stores keys on the server - only use for testing.
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-4">
            <button
              type="button"
              onClick={() => router.push('/wallet')}
              className="aggressive-btn"
              disabled={loading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="aggressive-btn flex-1"
              disabled={loading}
            >
              {loading ? 'Importing...' : 'Import Wallet'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}