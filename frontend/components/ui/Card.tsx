interface CardProps {
  children: React.ReactNode;
  className?: string;
  hover?: boolean;
  accent?: boolean;
}

export default function Card({ children, className = '', hover = true, accent = false }: CardProps) {
  return (
    <div
      className={`cyber-card clip-angle ${hover ? '' : 'hover:transform-none'} ${
        accent ? 'accent-pink' : ''
      } ${className}`}
    >
      {children}
    </div>
  );
}
