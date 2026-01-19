# Wallet Service

## 🚨 CRITICAL SECURITY WARNING

**This wallet implementation stores encrypted private keys on the server.**

**THIS IS FUNDAMENTALLY INSECURE AND SHOULD NEVER BE USED IN PRODUCTION.**

---

## Why This is Unsafe

### The Core Problem

Private keys stored server-side = **single point of failure**

Even with encryption:
1. **Server compromise** → All keys exposed
2. **Database breach** → Keys vulnerable to offline attacks  
3. **Memory dumps** → Keys exist in plaintext during signing
4. **Insider threats** → Server admins can access keys
5. **Legal seizure** → Government can seize server
6. **Cloud provider access** → AWS/GCP admins theoretically have access

### Attack Scenarios

**Scenario 1: Server Hack**
```
Attacker gains root access
→ Reads memory during signing operations
→ Obtains plaintext private keys
→ Steals all funds
```

**Scenario 2: Database Breach**
```
Attacker dumps database
→ Has encrypted keys + IVs + tags
→ Runs offline brute-force
→ If weak password: keys recovered
→ Steals all funds
```

**Scenario 3: Insider Threat**
```
Malicious employee
→ Has database access
→ Has encryption key from env
→ Decrypts keys
→ Steals all funds
```

**Scenario 4: Social Engineering**
```
Attacker calls support
→ Tricks admin into server access
→ Reads encryption key from environment
→ Decrypts all keys
→ Steals all funds
```

---

## Production Alternatives

### Client-Side Signing (Standard for dApps)

**How it works:**
```
User's browser                    Your server
     │                                 │
     │  Request: "Sign transaction"   │
     │ ◄─────────────────────────────┼
     │                                 │
     │  [User approves in MetaMask]   │
     │                                 │
     │  Signed transaction             │
     │ ────────────────────────────►  │
     │                                 │
     │         Transaction hash        │
     │ ◄─────────────────────────────┼
```

**Benefits:**
- Keys never leave user's device
- User has full control
- Standard UX (MetaMask, WalletConnect)
- No custody liability

**Implementations:**
- [wagmi](https://wagmi.sh/) + [viem](https://viem.sh/)
- [Web3.js](https://web3js.readthedocs.io/)
- [Ethers.js](https://docs.ethers.org/)

### Hardware Wallets

**How it works:**
```
User's computer          Hardware device         Your server
     │                        │                       │
     │  Request: Sign tx      │                       │
     │ ──────────────────►   │                       │
     │                        │                       │
     │  [User presses button] │                       │
     │                        │                       │
     │  Signed transaction    │                       │
     │ ◄──────────────────   │                       │
     │                        │                       │
     │           Broadcast to server                  │
     │ ────────────────────────────────────────────► │
```

**Benefits:**
- Private keys stored in secure hardware
- Immune to malware
- Physical confirmation required
- Industry gold standard

**Devices:**
- Ledger Nano S/X
- Trezor Model T
- GridPlus Lattice1

### Smart Contract Wallets

**How it works:**
```
No private keys at all!
Wallet is a smart contract on-chain
Recovery via social recovery, time locks, etc.
```

**Benefits:**
- No keys to steal
- Social recovery (friends help recover)
- Multi-sig for security
- Programmable security rules

**Examples:**
- Gnosis Safe (multi-sig)
- Argent (social recovery)
- Braavos (StarkNet)

### MPC Wallets (Multi-Party Computation)

**How it works:**
```
Key split into shares
No single party has full key
Signing requires coordination
```

**Benefits:**
- No single point of failure
- Distributed trust
- High security

**Providers:**
- Fireblocks
- ZenGo
- Qredo

---

## Why We Built This Anyway

### Educational Value

This implementation teaches:

1. **HD Wallets**
   - BIP39 (mnemonic generation)
   - BIP44 (key derivation paths)
   - Hierarchical deterministic keys

2. **Cryptography**
   - Argon2id password hashing
   - AES-256-GCM authenticated encryption
   - Initialization vectors (IVs)
   - Authentication tags

3. **Transaction Signing**
   - RLP encoding
   - ECDSA signatures
   - EIP-155 (replay protection)
   - Nonce management

4. **Ethereum APIs**
   - go-ethereum client library
   - Account management
   - Transaction building
   - Gas estimation

### Appropriate Use Cases

✅ **OK to use for:**
- Local development
- Automated testing
- CI/CD test wallets
- Demo applications (testnet)
- Learning blockchain programming

❌ **NEVER use for:**
- Production applications
- Real funds (mainnet)
- User-facing wallets
- Custody services
- Any scenario with financial risk

---

## Security Measures Implemented

Despite fundamental architectural flaws, we implement defense-in-depth:

### 1. Password Hashing (Argon2id)
```go
// High memory cost (64MB)
// High time cost (1 iteration minimum)
// Salt unique per user
// Output: 32-byte encryption key
argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
```

**Protects against:**
- Rainbow tables (salt)
- Brute force (high cost)
- GPU attacks (memory-hard)

**Does NOT protect against:**
- Server compromise (attacker gets parameters)
- Offline attacks (given enough time/resources)

### 2. AES-256-GCM Encryption
```go
// Cipher: AES-256
// Mode: GCM (Galois/Counter Mode)
// Key size: 32 bytes (256 bits)
// IV size: 12 bytes (unique per encryption)
// Tag size: 16 bytes (authentication)
```

**Features:**
- Authenticated encryption (prevents tampering)
- Unique IV per key (prevents pattern analysis)
- Standard, battle-tested cipher

**Does NOT protect against:**
- Weak passwords (entropy is key)
- Encryption key theft (stored in environment)

### 3. Secure Random Generation
```go
// Cryptographically secure random bytes
import "crypto/rand"
rand.Read(iv)
```

**Ensures:**
- Unpredictable IVs
- Unpredictable salts
- Secure mnemonic generation

### 4. Constant-Time Operations
```go
// Prevents timing attacks
import "crypto/subtle"
subtle.ConstantTimeCompare(a, b)
```

**Protects against:**
- Timing-based password guessing
- Side-channel attacks

### 5. Memory Safety
```go
// Clear sensitive data after use
copy(privateKey, make([]byte, len(privateKey)))
```

**Reduces (but doesn't eliminate):**
- Memory dump attacks
- Core dump exposure

### 6. No Plaintext Logging
```
❌ log.Printf("Private key: %s", key)  // NEVER
✅ log.Printf("Wallet created: %s", address)  // OK
```

**Ensures:**
- Keys never in log files
- Debugging doesn't expose secrets

### 7. Session Management
```
- JWT tokens expire after 30 minutes
- Automatic logout on timeout
- No persistent sessions
```

**Limits:**
- Window of vulnerability
- Replay attack timeframe

---

## Architecture

### Database Schema
```sql
CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    address VARCHAR(42) UNIQUE NOT NULL,
    encrypted_private_key TEXT NOT NULL,  -- AES-256-GCM ciphertext
    encryption_iv VARCHAR(32) NOT NULL,   -- Unique IV per key
    encryption_tag VARCHAR(32) NOT NULL,  -- GCM auth tag
    name VARCHAR(100),
    derivation_path VARCHAR(100),         -- m/44'/60'/0'/0/N
    is_imported BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Encryption Flow
```
User Password
     │
     ▼
Argon2id (1 iter, 64MB, 4 threads)
     │
     ▼
32-byte Encryption Key
     │
     ├──► Generate 12-byte random IV
     │
     ▼
AES-256-GCM Encrypt (plaintext private key)
     │
     ▼
Ciphertext + 16-byte Auth Tag
     │
     ▼
Store in Database
```

### Decryption Flow
```
User Password
     │
     ▼
Argon2id (same params as encryption)
     │
     ▼
32-byte Decryption Key
     │
     ├──► Retrieve IV from database
     ├──► Retrieve Tag from database
     ├──► Retrieve Ciphertext from database
     │
     ▼
AES-256-GCM Decrypt + Verify Tag
     │
     ▼
Plaintext Private Key (in memory only)
     │
     ├──► Sign Transaction
     │
     ▼
Immediately Zero Out Memory
```

### Transaction Signing Flow
```
1. User provides password
2. Decrypt private key to memory
3. Get current nonce from blockchain
4. Build transaction (to, value, gas, etc.)
5. Sign transaction with private key
6. Zero out private key from memory
7. Broadcast signed transaction
8. Return transaction hash
```

---

## API Endpoints

### POST /wallet/create
Create new HD wallet from generated mnemonic.

**Request:**
```json
{
  "name": "My Wallet",
  "password": "secure_password_min_12_chars"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "address": "0x123...",
    "mnemonic": "word1 word2 ... word12",
    "warning": "⚠️ SAVE YOUR MNEMONIC SECURELY - Cannot be recovered!"
  }
}
```

**Security notes:**
- Mnemonic shown ONCE, never stored
- Password must be >= 12 characters
- Wallet locked after creation (requires password to use)

### POST /wallet/import
Import wallet from mnemonic or private key.

**Request:**
```json
{
  "method": "mnemonic",
  "mnemonic": "word1 word2 ...",
  "password": "secure_password",
  "name": "Imported Wallet"
}
```

**Security notes:**
- Never store mnemonic
- Derive key from mnemonic, encrypt, store only encrypted key

### POST /wallet/send
Send transaction.

**Request:**
```json
{
  "wallet_id": "uuid",
  "chain_id": 1337,
  "to": "0xto...",
  "value": "1000000000000000000",
  "password": "secure_password"
}
```

**Security flow:**
1. Verify password → decrypt key (memory)
2. Sign transaction → send to blockchain
3. Zero key from memory
4. Never log key or password

---

## Development Notes

### Environment Variables
```bash
# Encryption key for wallet service
# Generate: openssl rand -hex 32
WALLET_ENCRYPTION_KEY=your-32-byte-hex-key-here

# Password hashing parameters
ARGON2_TIME=1
ARGON2_MEMORY=64
ARGON2_THREADS=4
ARGON2_KEY_LENGTH=32

# Session management
JWT_SECRET=your-jwt-secret-here
JWT_EXPIRY=30m
```

### Testing
```bash
# Run wallet tests
go test ./internal/wallet/... -v -cover

# Generate coverage report
go test ./internal/wallet/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Debugging

**NEVER log sensitive data:**
```go
// ❌ NEVER
log.Printf("Private key: %x", privateKey)
log.Printf("Password: %s", password)
log.Printf("Mnemonic: %s", mnemonic)

// ✅ OK
log.Printf("Wallet created: %s", address)
log.Printf("Transaction sent: %s", txHash)
```

---

## Future Improvements (Still Not Production-Ready)

Even with these improvements, server-side key storage remains fundamentally insecure:

- [ ] Hardware Security Module (HSM) integration
- [ ] Multi-party computation (MPC)
- [ ] Key sharding across multiple servers
- [ ] Biometric authentication
- [ ] Yubikey second factor
- [ ] Audit logging (all key access logged)
- [ ] Key rotation policies
- [ ] Automated security testing

**But still:**  
**Client-side signing >> Any server-side solution**

---

## Responsible Disclosure

If you find a security vulnerability:

1. **DO NOT** open a public GitHub issue
2. **DO** email: [your-email]
3. Include:
   - Description of vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We'll respond within 48 hours.

---

## License & Liability

**MIT License** - Use at your own risk.

**NO WARRANTY** - This code is provided "as-is" without any warranty.

**NO LIABILITY** - Not responsible for loss of funds.

**BY USING THIS CODE, YOU ACKNOWLEDGE:**
- You understand the security risks
- You will NOT use this for real funds
- You will NOT deploy this to production
- You will NOT hold the authors liable for losses

---

## Learn More

**HD Wallets:**
- [BIP39 Specification](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)
- [BIP44 Specification](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)

**Cryptography:**
- [Argon2 RFC](https://www.rfc-editor.org/rfc/rfc9106.html)
- [AES-GCM NIST](https://csrc.nist.gov/publications/detail/sp/800-38d/final)

**Ethereum:**
- [go-ethereum Documentation](https://geth.ethereum.org/docs)
- [Ethereum Yellow Paper](https://ethereum.github.io/yellowpaper/paper.pdf)

---

**Remember:** This is a learning exercise. Production wallets use client-side signing.

**Last Updated:** January 19, 2026
