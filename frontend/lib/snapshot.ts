// frontend/src/lib/snapshot.ts
import fs from 'fs';
import path from 'path';

// For server-side rendering/static generation
export async function loadSnapshot<T>(filePath: string): Promise<T> {
  // During build, read from filesystem
  if (typeof window === 'undefined') {
    const fullPath = path.join(process.cwd(), 'public', 'snapshot', filePath);
    
    try {
      const fileContents = fs.readFileSync(fullPath, 'utf8');
      return JSON.parse(fileContents);
    } catch (error) {
      throw new Error(`Failed to load snapshot: ${filePath}`);
    }
  }
  
  // Client-side, fetch from public
  const res = await fetch(`/snapshot/${filePath}`);
  if (!res.ok) {
    throw new Error(`Failed to load snapshot: ${filePath}`);
  }
  return res.json();
}

export async function loadStats(): Promise<any> {
  return loadSnapshot('stats.json');
}

export async function loadBlocks(): Promise<any> {
  return loadSnapshot('blocks.json');
}

export async function loadTransactions(): Promise<any> {
  return loadSnapshot('transactions.json');
}

export async function loadBlockById(id: string): Promise<any> {
  // Try individual file first (has full data including tx hashes)
  try {
    return await loadSnapshot(`blocks/${id}.json`);
  } catch {
    // Fallback to blocks.json list (doesn't have tx hashes)
    const blocksData = await loadBlocks();
    const blocks = blocksData.data.blocks || [];
    const blockFromList = blocks.find((b: any) => b.block_number === Number(id));
    
    if (blockFromList) {
      return {
        success: true,
        data: blockFromList,
      };
    }
    
    return {
      success: false,
      error: { message: `Block ${id} not in snapshot` },
    };
  }
}

export async function loadTransactionByHash(hash: string): Promise<any> {
  // First try to find in transactions.json list
  const txData = await loadTransactions();
  const transactions = txData.data.transactions || [];
  const txFromList = transactions.find((t: any) => t.hash === hash);
  
  if (txFromList) {
    return {
      success: true,
      data: txFromList,
    };
  }
  
  // If not in list, try individual file (for featured txs)
  try {
    return await loadSnapshot(`transactions/${hash}.json`);
  } catch {
    return {
      success: false,
      error: { message: `Transaction ${hash} not in snapshot` },
    };
  }
}

export async function loadAddressByAddress(address: string): Promise<any> {
  try {
    return await loadSnapshot(`addresses/${address}.json`);
  } catch {
    return {
      success: false,
      error: { message: `Address ${address} not in snapshot` },
    };
  }
}