export interface Chain {
  chain_id: number;
  name: string;
  short_name: string;
  native_symbol: string;
  rpc_endpoint: string;
  ws_endpoint?: string;
  block_time_seconds: number;
  is_active: boolean;
  is_testnet: boolean;
  explorer_url?: string;
  icon_url?: string;
}

export interface Block {
  id: number;
  chain_id: number;
  block_number: number;
  hash: string;
  parent_hash: string;
  miner: string;
  timestamp: string;
  gas_used: number;
  gas_limit: number;
  tx_count: number;
  base_fee_per_gas?: string;
}

export interface Transaction {
  id: number;
  chain_id: number;
  hash: string;
  block_number: number;
  block_hash: string;
  transaction_index: number;
  from_address: string;
  to_address?: string;
  value: string;
  gas: number;
  gas_price?: string;
  max_fee_per_gas?: string;
  max_priority_fee_per_gas?: string;
  input?: string;
  nonce: number;
  status?: number;
  gas_used?: number;
  timestamp: string;
  contract_address?: string | null;
  logs?: Log[];
  effective_gas_price?: string;
  transaction_type?: number;
}

export interface Address {
  address: string;
  balance: string;
  nonce: number;
  is_contract: boolean;
  tx_count: number;
  first_seen_block?: number;
  last_seen_block?: number;
  first_seen_at?: string;
  last_seen_at?: string;
  token_balances?: TokenBalance[];
}

export interface TokenBalance {
  token_address: string;
  token_name: string;
  token_symbol: string;
  token_type: string;
  balance: string;
  decimals: number;
}

export interface Log {
  log_index: number;
  address: string;
  topics: string[];
  data: string;
}

export interface Stats {
  chain_id: number;
  latest_block: number;
  avg_block_time: number;
  total_transactions: number;
  total_addresses: number;
  active_addresses_24h: number;
  tps_24h: number;
  gas_price?: {
    low: string;
    medium: string;
    high: string;
  };
}

export interface Pagination {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: {
    code: string;
    message: string;
  };
}

export interface Wallet {
  id: string;
  address: string;
  name: string;
}

export interface CreateWalletRequest {
  name: string;
  password: string;
}

export interface CreateWalletResponse {
  id: string;
  address: string;
  mnemonic: string; // ⚠️ SHOWN ONCE
  name: string;
}

export interface ImportWalletRequest {
  method: 'private_key' | 'mnemonic';
  private_key?: string;
  mnemonic?: string;
  password: string;
  name: string;
}

export interface Balance {
  chain_id: number;
  chain_name: string;
  balance: string; // Wei as string
  symbol: string;
}

export interface BalanceResponse {
  wallet_id: string;
  balances: Balance[];
}

export interface SendTransactionRequest {
  wallet_id: string;
  chain_id: number;
  to: string;
  value: string; // Wei as string
  password: string;
  gas_limit?: number;
  gas_price?: string;
}

export interface SendTransactionResponse {
  tx_hash: string;
  status: 'pending';
}

// Helper: Convert wei to ETH
export function weiToEth(wei: string): string {
  return (parseFloat(wei) / 1e18).toFixed(4);
}

// Helper: Convert ETH to wei
export function ethToWei(eth: string): string {
  return (parseFloat(eth) * 1e18).toString();
}