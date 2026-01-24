'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { createWallet } from '@/lib/api';

export default function CreateWallet() {
  const router = useRouter();
  const [step, setStep] = useState<'form' | 'mnemonic'>('form');
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  
  // Wallet creation result
  const [wallet, setWallet] = useState<{
    id: string;
    address: string;
    mnemonic: string;
  } | null>(null);
  
  const [savedConfirmed, setSavedConfirmed] = useState(false);

  const handleCreate = async (e: React.FormEvent) => {
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

    setLoading(true);
    try {
      const result = await createWallet(name.trim(), password);
      setWallet(result);
      setStep('mnemonic');
    } catch (err: any) {
      setError(err.message || 'Failed to create wallet');
    } finally {
      setLoading(false);
    }
  };

  const handleFinish = () => {
    if (!savedConfirmed) {
      alert('Please confirm you have saved your recovery phrase!');
      return;
    }
    // Store wallet ID in localStorage (simple approach for now)
    if (wallet) {
      localStorage.setItem('current_wallet', wallet.id);
    }
    router.push('/wallet/dashboard');
  };

  if (step === 'mnemonic' && wallet) {
    return (
      <div className="container mx-auto px-4 py-12">
        <div className="max-w-2xl mx-auto">
          {/* Critical Warning */}
          <div className="cyber-card border-red-500/50 mb-8">
            <h2 className="text-2xl font-bold mb-4 text-red-400">
              🔴 SAVE YOUR RECOVERY PHRASE
            </h2>
            <p className="text-zinc-300 mb-4">
              Write down these 12 words <strong>on paper</strong> and store them in a secure location.
              <strong className="block mt-2 text-red-400">
                This is the ONLY way to recover your wallet. This phrase is shown ONCE and never again.
              </strong>
            </p>
          </div>

          {/* Mnemonic Display */}
          <div className="cyber-card mb-6">
            <div className="grid grid-cols-3 gap-4">
              {wallet.mnemonic.split(' ').map((word, i) => (
                <div key={i} className="flex items-center gap-2 bg-black/50 p-3 clip-angle-sm border border-cyan-500/30">
                  <span className="text-cyan-400 font-mono text-sm">{i + 1}.</span>
                  <span className="font-mono font-bold">{word}</span>
                </div>
              ))}
            </div>
          </div>

          {/* Wallet Info */}
          <div className="cyber-card mb-6">
            <div className="space-y-3">
              <div>
                <div className="text-sm text-zinc-400 mb-1">Wallet Name</div>
                <div className="font-mono">{name}</div>
              </div>
              <div>
                <div className="text-sm text-zinc-400 mb-1">Address</div>
                <div className="font-mono text-sm break-all text-cyan-400">{wallet.address}</div>
              </div>
            </div>
          </div>

          {/* Confirmation */}
          <div className="cyber-card mb-6">
            <label className="flex items-start gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={savedConfirmed}
                onChange={(e) => setSavedConfirmed(e.target.checked)}
                className="mt-1 w-5 h-5 accent-cyan-500"
              />
              <span className="text-sm">
                I have written down my 12-word recovery phrase on paper and stored it securely.
                I understand this cannot be recovered if lost.
              </span>
            </label>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-4">
            <button
              onClick={handleFinish}
              disabled={!savedConfirmed}
              className="aggressive-btn flex-1"
            >
              Continue to Dashboard
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-12">
      <div className="max-w-xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2 text-cyan-400">Create New Wallet</h1>
          <p className="text-zinc-400">Generate a new Ethereum wallet with a recovery phrase</p>
        </div>

        {/* Form */}
        <form onSubmit={handleCreate}>
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
                  placeholder="My Wallet"
                  className="w-full bg-black/70 border border-cyan-500/50 px-4 py-3 text-white font-mono
                    focus:outline-none focus:border-cyan-400 focus:shadow-[0_0_15px_rgba(34,211,238,0.3)]
                    clip-angle-sm"
                  disabled={loading}
                />
              </div>

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
                <div className="mt-1 text-xs text-zinc-500">
                  Used to encrypt your wallet on the server
                </div>
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
              {loading ? 'Creating...' : 'Create Wallet'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}