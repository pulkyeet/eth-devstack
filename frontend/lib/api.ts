import type { ApiResponse, Block, Transaction, Address, Stats, Chain, Pagination } from './types';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

async function fetchAPI<T>(endpoint: string): Promise<T> {
  const res = await fetch(`${API_BASE}${endpoint}`);
  if (!res.ok) {
    throw new Error(`API error: ${res.status}`);
  }
  const json: ApiResponse<T> = await res.json();
  if (!json.success) {
    throw new Error(json.error?.message || 'API request failed');
  }
  return json.data;
}

export async function getChains(): Promise<{ chains: Chain[] }> {
  return fetchAPI('/chains');
}

export async function getBlocks(chainId: number, page = 1, limit = 20): Promise<{
  blocks: Block[];
  pagination: Pagination;
}> {
  return fetchAPI(`/blocks?chain_id=${chainId}&page=${page}&limit=${limit}`);
}

export async function getBlock(chainId: number, blockId: string): Promise<Block> {
  return fetchAPI(`/blocks/${blockId}?chain_id=${chainId}`);
}

export async function getTransactions(chainId: number, page = 1, limit = 20): Promise<{
  transactions: Transaction[];
  pagination: Pagination;
}> {
  return fetchAPI(`/transactions?chain_id=${chainId}&page=${page}&limit=${limit}`);
}

export async function getTransaction(chainId: number, hash: string): Promise<Transaction> {
  return fetchAPI(`/transactions/${hash}?chain_id=${chainId}`);
}

export async function getAddress(chainId: number, address: string): Promise<Address> {
  return fetchAPI(`/addresses/${address}?chain_id=${chainId}`);
}

export async function getAddressTransactions(
  chainId: number,
  address: string,
  page = 1,
  limit = 20
): Promise<{
  address: string;
  transactions: Transaction[];
  pagination: Pagination;
}> {
  return fetchAPI(`/addresses/${address}/transactions?chain_id=${chainId}&page=${page}&limit=${limit}`);
}

export async function search(chainId: number, query: string): Promise<{
  type: 'block' | 'transaction' | 'address';
  result: Block | Transaction | Address;
}> {
  return fetchAPI(`/search?chain_id=${chainId}&q=${encodeURIComponent(query)}`);
}

export async function getStats(chainId: number): Promise<Stats> {
  return fetchAPI(`/stats?chain_id=${chainId}`);
}

export function createBlockStream(chainId: number): EventSource {
  return new EventSource(`${API_BASE}/stream/blocks?chain_id=${chainId}`);
}