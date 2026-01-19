#!/bin/bash
set -e

OUT="$HOME/eth-devstack/frontend/public/snapshot"
DB="postgresql://eth_user:eth_pass_dev_only@localhost:5432/ethereum_explorer"

echo "=== Exporting Snapshot from Database ==="
mkdir -p $OUT

# Find last 100 blocks WITH transactions (before empty blocks)
END_BLOCK=$(psql $DB -t -c "SELECT MAX(block_number) FROM blocks WHERE chain_id = 1337 AND tx_count > 0;")
START_BLOCK=$((END_BLOCK - 99))

# Get actual counts
TOTAL_TXS=$(psql $DB -t -c "SELECT COUNT(*) FROM transactions WHERE chain_id = 1337 AND block_number BETWEEN $START_BLOCK AND $END_BLOCK;")

# Calculate TPS from these 100 blocks
FIRST_TS=$(psql $DB -t -c "SELECT EXTRACT(EPOCH FROM timestamp) FROM blocks WHERE chain_id = 1337 AND block_number = $START_BLOCK;")
LAST_TS=$(psql $DB -t -c "SELECT EXTRACT(EPOCH FROM timestamp) FROM blocks WHERE chain_id = 1337 AND block_number = $END_BLOCK;")
DURATION=$(echo "$LAST_TS - $FIRST_TS" | bc)
TPS=$(echo "scale=1; $TOTAL_TXS / $DURATION" | bc)

echo "Block range: $START_BLOCK to $END_BLOCK (100 blocks)"
echo "Total transactions: $TOTAL_TXS"
echo "Duration: ${DURATION}s"
echo "TPS: $TPS"

# 1. Stats with TPS
cat > $OUT/stats.json << STATS
{
  "success": true,
  "data": {
    "chain_id": 1337,
    "latest_block": $END_BLOCK,
    "start_block": $START_BLOCK,
    "avg_block_time": 2.0,
    "tps": $TPS,
    "total_transactions": $TOTAL_TXS,
    "total_addresses": $(psql $DB -t -c "SELECT COUNT(*) FROM addresses WHERE chain_id = 1337;"),
    "note": "Frozen snapshot - last 100 blocks with data"
  }
}
STATS

# 2. Blocks list (last 100)
echo "Exporting blocks list..."
psql $DB -t -A -c "
SELECT json_build_object(
  'success', true,
  'data', json_build_object(
    'blocks', json_agg(
      json_build_object(
        'block_number', block_number,
        'hash', hash,
        'parent_hash', parent_hash,
        'timestamp', timestamp,
        'miner', miner,
        'tx_count', tx_count,
        'gas_used', gas_used,
        'gas_limit', gas_limit,
        'base_fee_per_gas', base_fee_per_gas
      ) ORDER BY block_number DESC
    ),
    'pagination', json_build_object(
      'page', 1,
      'limit', 100,
      'total', 100
    )
  )
)
FROM blocks 
WHERE chain_id = 1337 
AND block_number BETWEEN $START_BLOCK AND $END_BLOCK;
" > $OUT/blocks.json

# 3. Transactions list (last 100 txs)
echo "Exporting transactions list..."
psql $DB -t -A -c "
SELECT json_build_object(
  'success', true,
  'data', json_build_object(
    'transactions', json_agg(
      json_build_object(
        'hash', hash,
        'block_number', block_number,
        'timestamp', timestamp,
        'from_address', from_address,
        'to_address', to_address,
        'value', value,
        'gas_used', gas_used,
        'gas_price', gas_price,
        'status', status
      ) ORDER BY block_number DESC, transaction_index DESC
    ),
    'pagination', json_build_object(
      'page', 1,
      'limit', 100,
      'total', 100
    )
  )
)
FROM (
  SELECT * FROM transactions 
  WHERE chain_id = 1337 
  AND block_number BETWEEN $START_BLOCK AND $END_BLOCK
  ORDER BY block_number DESC, transaction_index DESC 
  LIMIT 100
) sub;
" > $OUT/transactions.json

# 4. ALL individual block files (100 blocks)
echo "Exporting individual blocks..."
mkdir -p $OUT/blocks
for ((BLOCK_NUM=$START_BLOCK; BLOCK_NUM<=END_BLOCK; BLOCK_NUM++)); do
  echo -ne "  Progress: $((BLOCK_NUM - START_BLOCK + 1))/100\r"
  psql $DB -t -A -c "
  WITH block_data AS (
    SELECT * FROM blocks WHERE chain_id = 1337 AND block_number = $BLOCK_NUM
  ),
  tx_data AS (
    SELECT json_agg(hash ORDER BY transaction_index) as tx_hashes
    FROM transactions 
    WHERE chain_id = 1337 AND block_number = $BLOCK_NUM
  )
  SELECT json_build_object(
    'success', true,
    'data', json_build_object(
      'block_number', b.block_number,
      'hash', b.hash,
      'parent_hash', b.parent_hash,
      'timestamp', b.timestamp,
      'miner', b.miner,
      'difficulty', b.difficulty,
      'total_difficulty', b.total_difficulty,
      'size', b.size,
      'gas_used', b.gas_used,
      'gas_limit', b.gas_limit,
      'base_fee_per_gas', b.base_fee_per_gas,
      'nonce', b.nonce,
      'extra_data', b.extra_data,
      'tx_count', b.tx_count,
      'transactions', COALESCE(t.tx_hashes, '[]'::json)
    )
  )
  FROM block_data b
  LEFT JOIN tx_data t ON true;
  " > $OUT/blocks/$BLOCK_NUM.json
done
echo ""

# 5. ALL transaction details from these blocks
echo "Exporting individual transactions..."
mkdir -p $OUT/transactions
psql $DB -t -c "SELECT hash FROM transactions WHERE chain_id = 1337 AND block_number BETWEEN $START_BLOCK AND $END_BLOCK ORDER BY block_number DESC;" | while read TX_HASH; do
  TX_HASH=$(echo $TX_HASH | xargs)
  if [ ! -z "$TX_HASH" ]; then
    psql $DB -t -A -c "
    SELECT json_build_object(
      'success', true,
      'data', json_build_object(
        'hash', hash,
        'block_number', block_number,
        'block_hash', block_hash,
        'timestamp', timestamp,
        'from_address', from_address,
        'to_address', to_address,
        'value', value,
        'gas', gas,
        'gas_price', gas_price,
        'gas_used', gas_used,
        'effective_gas_price', effective_gas_price,
        'nonce', nonce,
        'transaction_index', transaction_index,
        'transaction_type', transaction_type,
        'input', input,
        'status', status
      )
    )
    FROM transactions 
    WHERE chain_id = 1337 AND hash = '$TX_HASH';
    " > $OUT/transactions/$TX_HASH.json
  fi
done

# 6. Sample addresses (10 addresses)
echo "Exporting sample addresses..."
mkdir -p $OUT/addresses
psql $DB -t -c "SELECT address FROM addresses WHERE chain_id = 1337 ORDER BY tx_count DESC LIMIT 128;" | while read ADDR; do
  ADDR=$(echo $ADDR | xargs)
  if [ ! -z "$ADDR" ]; then
    psql $DB -t -A -c "
    SELECT json_build_object(
      'success', true,
      'data', json_build_object(
        'address', address,
        'balance', balance,
        'nonce', nonce,
        'is_contract', is_contract,
        'tx_count', tx_count,
        'first_seen_block', first_seen_block,
        'last_seen_block', last_seen_block
      )
    )
    FROM addresses 
    WHERE chain_id = 1337 AND address = '$ADDR';
    " > $OUT/addresses/$ADDR.json
  fi
done

echo ""
echo "=== Export Complete ==="
ls -lh $OUT/
echo "Total size: $(du -sh $OUT | cut -f1)"