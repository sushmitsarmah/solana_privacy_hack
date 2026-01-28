# Umbra Privacy API

A TypeScript Node.js Express server integrating the **Umbra SDK** for privacy-focused payments on Solana using zero-knowledge proofs.

## Features

- **Zero-Knowledge Proofs**: Privacy-preserving transactions using ZK-SNARKs
- **Encrypted Balances**: Confidential balance encryption using Arcium MXE
- **Privacy Mixer**: Anonymous transaction mixing
- **Account Registration**: Register accounts for anonymity/confidentiality
- **Stealth Addresses**: One-time addresses for enhanced privacy

## Technology Stack

- **TypeScript**: Type-safe development
- **Express.js**: Web server framework
- **Umbra SDK**: Zero-knowledge privacy protocol
- **@solana/web3.js**: Solana blockchain interaction
- **@arcium-hq/client**: Encrypted computation
- **snarkjs**: Zero-knowledge proof generation

## Installation

```bash
npm install
```

## Configuration

Copy `.env.example` to `.env` and configure:

```env
PORT=3000
SOLANA_RPC_URL=https://api.devnet.solana.com
NODE_ENV=development
```

## Running the Server

### Development mode (with hot reload)
```bash
npm run dev
```

### Build
```bash
npm run build
```

### Production
```bash
npm start
```

The server will start on `http://localhost:3000`

## API Endpoints

### Root
- `GET /` - API information and health check

### Umbra SDK Endpoints

#### Create Umbra Wallet
```bash
POST /api/umbra/wallet/create
```

**Request Body:**
```json
{
  "privateKey": "base58_private_key"  // optional, generates new if not provided
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "publicKey": "SolanaPublicKey",
    "privateKey": "base58EncodedPrivateKey",
    "walletCreated": true
  },
  "message": "Umbra wallet created successfully!",
  "security": {
    "warning": "Keep your private key secure and never share it!"
  }
}
```

#### Register Account for Anonymity
```bash
POST /api/umbra/register
```

**Request Body:**
```json
{
  "privateKey": "base58_private_key",
  "rpcUrl": "https://api.devnet.solana.com"  // optional
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "signature": "transaction_signature",
    "publicKey": "SolanaPublicKey",
    "explorerUrl": "https://explorer.solana.com/tx/..."
  },
  "message": "Account registered for anonymity successfully!"
}
```

#### Deposit to Privacy Pool
```bash
POST /api/umbra/deposit
```

**Request Body:**
```json
{
  "privateKey": "base58_private_key",
  "amount": 0.1,  // in SOL
  "rpcUrl": "https://api.devnet.solana.com"  // optional
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "amount": 0.1,
    "amountLamports": 100000000,
    "publicKey": "SolanaPublicKey"
  },
  "message": "Deposit prepared.",
  "note": "Deposits involve complex zero-knowledge proofs and privacy pool interactions."
}
```

#### Get API Info
```bash
GET /api/umbra/info
```

Returns detailed information about the Umbra integration, endpoints, and features.

## Usage Example

### 1. Create a New Wallet

```bash
curl -X POST http://localhost:3000/api/umbra/wallet/create \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Response includes:**
- Public Key
- Private Key (keep secure!)

### 2. Register Account for Privacy

```bash
curl -X POST http://localhost:3000/api/umbra/register \
  -H "Content-Type: application/json" \
  -d '{
    "privateKey": "YOUR_PRIVATE_KEY"
  }'
```

**This:**
- Registers your account with Umbra protocol
- Enables zero-knowledge proof generation
- Sets up encrypted balance tracking

### 3. Deposit to Privacy Pool

```bash
curl -X POST http://localhost:3000/api/umbra/deposit \
  -H "Content-Type: application/json" \
  -d '{
    "privateKey": "YOUR_PRIVATE_KEY",
    "amount": 0.1
  }'
```

## How Umbra Works

### 1. **Zero-Knowledge Proofs**
- Uses Groth16 ZK-SNARKs for transaction privacy
- Proves transaction validity without revealing amounts or parties

### 2. **Encrypted Balances**
- Balances encrypted using Arcium MXE (Multi-Party Execution)
- Only the wallet owner can decrypt their balance

### 3. **Privacy Mixer**
- Deposits funds into a shared pool
- Withdrawals unlinked from deposits
- Breaks on-chain transaction graph analysis

### 4. **Stealth Addresses**
- One-time addresses for each transaction
- Recipient's main address remains private

### 5. **Master Viewing Keys**
- Allows scanning for received payments
- Does not expose spending keys

## Architecture

```
express_server/
├── src/
│   ├── index.ts           # Main Express server
│   ├── routes/
│   │   └── umbra.ts       # Umbra SDK routes
│   └── types/
│       └── umbra-sdk.d.ts # Type declarations
├── umbra-sdk/             # Compiled Umbra SDK
│   ├── index.cjs          # CommonJS bundle
│   ├── index.d.ts         # Type definitions
│   └── package.json
├── package.json
├── tsconfig.json
└── README.md
```

## Umbra SDK Integration

### What is Umbra?

Umbra is a privacy protocol for Solana that enables:
- **Confidential Transactions**: Hide transaction amounts
- **Anonymous Payments**: Conceal sender and receiver identities
- **Zero-Knowledge Proofs**: Cryptographically prove transaction validity
- **Encrypted State**: Keep account balances private

### SDK Installation

The Umbra SDK is included as a local module built from source:

```bash
# SDK was built from: https://github.com/umbra-defi/sdk
# Located in: ./umbra-sdk/
```

### Key Classes

1. **UmbraClient** - Main interface for protocol operations
   - Account registration
   - Deposits and withdrawals
   - Transaction management

2. **UmbraWallet** - Key management and encryption
   - Keypair handling
   - X25519 encryption keys
   - Master viewing keys

3. **WasmZkProver** - Zero-knowledge proof generation
   - Groth16 proofs
   - Circuit evaluation
   - WASM-based computation

## Security Considerations

⚠️ **Important Security Notes:**

1. **Private Keys**: Never expose private keys over insecure connections
2. **HTTPS Required**: Always use HTTPS in production
3. **Key Storage**: Store keys securely (hardware wallets, key management systems)
4. **Rate Limiting**: Implement rate limiting for production APIs
5. **Authentication**: Add authentication for production deployments
6. **Zero-Knowledge**: Even with privacy, follow best security practices
7. **Testnet First**: Test thoroughly on devnet before mainnet

## Development

### Project Structure

The server uses the locally built Umbra SDK with Express.js routes.

### Adding New Routes

1. Create route file in `src/routes/`
2. Import Umbra SDK classes
3. Register route in `src/index.ts`
4. Update documentation

### Building the SDK

If you need to rebuild the Umbra SDK:

```bash
cd umbra-sdk-temp
pnpm install
pnpm build
cp -r dist ../umbra-sdk/
```

## Testing

Test endpoints using curl, Postman, or any HTTP client:

```bash
# Health check
curl http://localhost:3000/

# Get API info
curl http://localhost:3000/api/umbra/info

# Create wallet
curl -X POST http://localhost:3000/api/umbra/wallet/create

# Register (requires private key)
curl -X POST http://localhost:3000/api/umbra/register \
  -H "Content-Type: application/json" \
  -d '{"privateKey": "YOUR_KEY"}'
```

## Troubleshooting

### Port Already in Use
```bash
PORT=3001 npm run dev
```

### Solana Connection Issues
- Verify `SOLANA_RPC_URL` in `.env`
- Check if RPC endpoint is operational
- Try alternative RPC endpoints

### Build Errors
```bash
npm run build
```

Check TypeScript errors and SDK imports.

### SDK Import Issues
The SDK is imported using CommonJS require with `@ts-ignore` for type compatibility.

## Resources

- **Umbra SDK**: https://github.com/umbra-defi/sdk
- **Solana Docs**: https://docs.solana.com
- **Arcium MXE**: https://arcium.com
- **snarkjs**: https://github.com/iden3/snarkjs

## License

ISC

## Contributing

Contributions welcome! Please feel free to submit a Pull Request.
