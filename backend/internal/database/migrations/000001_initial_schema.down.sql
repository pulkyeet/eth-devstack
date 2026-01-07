DROP TRIGGER IF EXISTS update_addresses_updated_at ON addresses;
DROP TRIGGER IF EXISTS update_chains_updated_at ON chains;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS reorgs;
DROP TABLE IF EXISTS sync_status;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS transaction_logs;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS blocks;
DROP TABLE IF EXISTS chains;

DROP EXTENSION IF EXISTS "uuid-ossp";