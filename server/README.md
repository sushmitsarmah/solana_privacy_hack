# ShadowPay Go SDK

A clean, modular Go SDK for interacting with the ShadowPay API - a privacy-focused payment protocol for Solana.

## Features

- **Zero-Knowledge Payments**: Create and settle private payments using ZK proofs
- **Escrow Management**: Manage SOL and SPL token escrow accounts
- **Payment Intents**: Create and verify standard payment intents
- **API Key Management**: Generate, rotate, and manage API keys
- **X402 Verification**: Verify payment tokens and proofs

## Project Structure

```
.
├── shadowpay.go              # Main SDK entry point
├── cmd/
│   └── main.go              # Example usage
├── internal/
│   ├── client/              # HTTP client and core functionality
│   │   └── client.go
│   ├── errors/              # Error types and handling
│   │   └── errors.go
│   ├── types/               # Common types
│   │   └── transaction.go
│   ├── keys/                # API key management
│   │   └── keys.go
│   ├── escrow/              # Escrow operations
│   │   └── escrow.go
│   ├── payment/             # ZK payment operations
│   │   └── payment.go
│   ├── intent/              # Payment intent operations
│   │   └── intent.go
│   └── verify/              # X402 verification
│       └── verify.go
└── README.md
```

## Installation

```bash
go get sol_privacy
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    shadowpay "sol_privacy"
    "sol_privacy/internal/escrow"
)

func main() {
    // Initialize the SDK
    sdk := shadowpay.New("your-api-key")
    ctx := context.Background()

    // Check escrow balance
    balance, err := sdk.Escrow.GetBalance(ctx, "your-wallet-address")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Balance: %d lamports\n", balance.Balance)

    // Create a deposit transaction
    tx, err := sdk.Escrow.Deposit(ctx, escrow.TransactionRequest{
        WalletAddress: "your-wallet-address",
        Amount:        1000000, // 0.001 SOL
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Unsigned transaction: %s\n", tx.UnsignedTxBase64)
}
```

## Usage Examples

### API Key Management

```go
// Create a new API key
keyResp, err := sdk.Keys.Create(ctx, keys.GenerateRequest{
    WalletAddress: "your-wallet-address",
})

// Get API key by wallet
keyInfo, err := sdk.Keys.GetByWallet(ctx, "wallet-address")

// Rotate API key
newKey, err := sdk.Keys.Rotate(ctx)

// Check rate limits
limits, err := sdk.Keys.GetLimits(ctx)
```

### Escrow Operations

```go
// Get SOL balance
balance, err := sdk.Escrow.GetBalance(ctx, "wallet-address")

// Get token balance
tokenBalance, err := sdk.Escrow.GetTokenBalance(ctx, "wallet-address", "token-mint")

// Deposit SOL
depositTx, err := sdk.Escrow.Deposit(ctx, escrow.TransactionRequest{
    WalletAddress: "wallet-address",
    Amount:        1000000, // lamports
})

// Withdraw SOL
withdrawTx, err := sdk.Escrow.Withdraw(ctx, escrow.TransactionRequest{
    WalletAddress: "wallet-address",
    Amount:        1000000,
})

// Withdraw tokens
withdrawTokenTx, err := sdk.Escrow.WithdrawToken(ctx, escrow.TransactionRequest{
    WalletAddress: "wallet-address",
    Amount:        1000000,
    Mint:          "token-mint-address",
})
```

### ZK Payment Operations

```go
// Prepare a ZK payment
prepareResp, err := sdk.Payment.Prepare(ctx, payment.PrepareRequest{
    Amount:           1000000,
    Recipient:        "recipient-address",
    Resource:         "/api/resource",
    SenderCommitment: "commitment-hash",
})

// Settle a ZK payment
settleResp, err := sdk.Payment.Settle(ctx, payment.SettleRequest{
    X402Version:   1,
    PaymentHeader: "base64-encoded-header",
    Resource:      "/api/resource",
    PaymentRequirements: payment.Requirements{
        Scheme:            "zkproof",
        Network:           "solana-mainnet",
        MaxAmountRequired: "0.001",
        Resource:          "/api/resource",
        Description:       "Payment for resource access",
        MimeType:          "application/json",
        PayTo:             "recipient-address",
        MaxTimeoutSeconds: 300,
    },
})
```

### Payment Intents

```go
// Create a payment intent
intent, err := sdk.Intent.Create(ctx, intent.CreateRequest{
    Amount:    1000000,
    Recipient: "recipient-address",
    Reference: "order-12345",
})

// Verify a payment intent
verification, err := sdk.Intent.Verify(ctx, "intent-id")

// Get public key
pubKey, err := sdk.Intent.GetPublicKey(ctx)
```

### X402 Verification

```go
// Verify a payment token
result, err := sdk.Verify.X402(ctx, "payment-token")
if result.Valid {
    log.Printf("Token is valid, expires at: %d\n", result.ExpiresAt)
}
```

## Configuration

You can customize the SDK client with options:

```go
import "sol_privacy/internal/client"

sdk := shadowpay.New(
    "your-api-key",
    client.WithBaseURL("https://custom.api.url"),
    client.WithHTTPClient(&http.Client{
        Timeout: 60 * time.Second,
    }),
)
```

## Environment Variables

- `SHADOWPAY_API_KEY`: Your ShadowPay API key

## Running the Example

```bash
export SHADOWPAY_API_KEY=your-api-key
go run cmd/main.go
```

## API Documentation

For detailed API documentation, visit: https://registry.scalar.com/@radr/apis/shadowpay-api

## License

MIT
