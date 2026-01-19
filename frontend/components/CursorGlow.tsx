'use client';

import { useEffect } from 'react';

export default function CursorGlow() {
  useEffect(() => {
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

  return null;
}