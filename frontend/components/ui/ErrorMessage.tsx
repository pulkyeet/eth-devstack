import Link from 'next/link';
import Card from './Card';

interface Props {
  message: string;
  showBackLink?: boolean;
}

export default function ErrorMessage({ message, showBackLink = true }: Props) {
  return (
    <div className="container mx-auto px-4 py-8">
      <Card hover={false}>
        <div className="text-center py-8">
          <div className="text-pink-400 text-2xl mb-4 font-bold">✗ ERROR</div>
          <div className="text-zinc-400 mb-6">{message}</div>
          {showBackLink && (
            <Link href="/explorer" className="text-cyan-400 hover:text-cyan-300 transition-colors">
              ← BACK TO EXPLORER
            </Link>
          )}
        </div>
      </Card>
    </div>
  );
}