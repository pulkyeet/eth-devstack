interface CardProps {
  children: React.ReactNode;
  className?: string;
  hover?: boolean;
}

export default function Card({ children, className = '', hover = true }: CardProps) {
  return (
    <div className={`relative ${hover ? 'group cursor-pointer' : ''} ${className}`}>
      <div className={`
        relative bg-[var(--bg-card)] 
        border border-[var(--cyan)]
        ${hover ? 'group-hover:border-[var(--purple)] group-hover:-translate-y-1 group-hover:shadow-[0_0_20px_rgba(180,0,255,0.4)]' : ''}
        transition-all duration-300
      `} style={{
        /* This clipPath cuts 12px off the corners. 
           With 32px left padding, your text starts 20px away from the cut. */
        clipPath: 'polygon(12px 0, 100% 0, 100% calc(100% - 12px), calc(100% - 12px) 100%, 0 100%, 0 12px)'
      }}>
        
        {/* Pink accent bar */}
        <div className="absolute left-0 top-0 w-1 h-full bg-[var(--pink)] z-20" 
             style={{ boxShadow: 'var(--glow-pink)' }}></div>
        
        {/* Content wrapper with EXPLICIT pixel padding */}
        <div 
          className="relative z-10"
          style={{ 
            paddingTop: '24px', 
            paddingBottom: '24px', 
            paddingRight: '24px', 
            paddingLeft: '32px' 
          }}
        >
          {children}
        </div>

        {/* Corner accents */}
        <div className="absolute -top-[1px] -left-[1px] w-3 h-3 border-l border-t border-[var(--cyan)]" 
             style={{ clipPath: 'polygon(10px 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%, 0 10px)' }}></div>
        <div className="absolute -bottom-[1px] -right-[1px] w-3 h-3 border-r border-b border-[var(--cyan)]" 
             style={{ clipPath: 'polygon(10px 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%, 0 10px)' }}></div>
      </div>
    </div>
  );
}