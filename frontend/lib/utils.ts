import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { formatEther } from "viem";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function shortenAddress(address: string, chars = 4): string {
  if (!address) return 'N/A';
  return `${address.slice(0, chars + 2)}...${address.slice(-chars)}`;
}

export function shortenHash(hash: string, chars = 6): string {
  if (!hash) return 'N/A';
  return `${hash.slice(0, chars + 2)}...${hash.slice(-chars)}`;
}

export function formatTimestamp(timestamp: string): string {
  const date = new Date(timestamp);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  
  if (seconds < 60) return `${seconds}s ago`;
  if (minutes < 60) return `${minutes}m ago`;
  if (hours < 24) return `${hours}h ago`;
  return `${days}d ago`;
}

export function formatWei(wei: string | bigint, decimals = 4): string {
  const eth = formatEther(BigInt(wei));
  const num = parseFloat(eth);
  return num.toFixed(decimals);
}

export function formatGas(gas: number): string {
  return gas.toLocaleString();
}

export function formatNumber(num: number): string {
  if (num >= 1_000_000) {
    return `${(num / 1_000_000).toFixed(2)}M`;
  }
  if (num >= 1_000) {
    return `${(num / 1_000).toFixed(2)}K`;
  }
  return num.toString();
}