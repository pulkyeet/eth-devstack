export default function LoadingSpinner({ message = "LOADING" }: { message?: string }) {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="text-center">
        <div className="inline-block h-12 w-12 animate-spin rounded-full border-4 border-cyan-400/30 border-t-cyan-400 mb-4" />
        <div className="text-cyan-400 text-xl tracking-wider font-bold">
          {message.toUpperCase()}...
        </div>
      </div>
    </div>
  );
}