export default function WalletLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen">
      {/* ⚠️ DEV ONLY WARNING */}
      <div className="bg-red-900/20 border-b border-red-500/50 py-3">
        <div className="container mx-auto px-4">
          <div className="flex items-center gap-3 text-red-400">
            <span className="text-2xl">⚠️</span>
            <div className="text-sm">
              <strong>DEVELOPMENT WALLET - NOT FOR PRODUCTION</strong>
              <span className="ml-2 opacity-75">
                This wallet stores private keys on the server (encrypted). 
                NEVER use for real funds. For learning purposes only.
              </span>
            </div>
          </div>
        </div>
      </div>
      
      {children}
    </div>
  );
}