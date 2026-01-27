package main

import (
	"context"
	"fmt"
	"log"
	"os"

	shadowpay "sol_privacy"
	"sol_privacy/internal/escrow"
	"sol_privacy/internal/intent"
)

func main() {
	apiKey := os.Getenv("SHADOWPAY_API_KEY")
	if apiKey == "" {
		log.Println("Warning: SHADOWPAY_API_KEY not set. Some operations may fail.")
	}

	walletAddr := "AVSSWPbWRYDF7w8GZcrP6yVWsmRWPshMnziHqFQ5RaDR"

	// Initialize ShadowPay SDK
	sdk := shadowpay.New(apiKey)
	ctx := context.Background()

	// Example 1: API Key Management
	fmt.Println("=== API Key Management ===")
	if err := demonstrateKeyManagement(ctx, sdk, walletAddr); err != nil {
		log.Printf("Key management error: %v\n", err)
	}

	// Example 2: Escrow Operations
	fmt.Println("\n=== Escrow Operations ===")
	if err := demonstrateEscrow(ctx, sdk, walletAddr); err != nil {
		log.Printf("Escrow error: %v\n", err)
	}

	// Example 3: Payment Intent
	fmt.Println("\n=== Payment Intent ===")
	if err := demonstratePaymentIntent(ctx, sdk, walletAddr); err != nil {
		log.Printf("Payment intent error: %v\n", err)
	}
}

func demonstrateKeyManagement(ctx context.Context, sdk *shadowpay.ShadowPay, wallet string) error {
	// Get API key by wallet
	keyInfo, err := sdk.Keys.GetByWallet(ctx, wallet)
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}
	fmt.Printf("API Key for wallet %s: %s\n", keyInfo.Wallet, keyInfo.APIKey)

	// Get rate limits
	limits, err := sdk.Keys.GetLimits(ctx)
	if err != nil {
		return fmt.Errorf("failed to get limits: %w", err)
	}
	fmt.Printf("Rate Limits - Limit: %d, Remaining: %d, Reset: %d\n",
		limits.Limit, limits.Remaining, limits.Reset)

	return nil
}

func demonstrateEscrow(ctx context.Context, sdk *shadowpay.ShadowPay, wallet string) error {
	// Check escrow balance
	balance, err := sdk.Escrow.GetBalance(ctx, wallet)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}
	fmt.Printf("Escrow Balance: %d lamports (%.4f SOL)\n",
		balance.Balance, float64(balance.Balance)/1e9)

	// Generate deposit transaction
	depositReq := escrow.TransactionRequest{
		WalletAddress: wallet,
		Amount:        50000000, // 0.05 SOL
	}

	txResponse, err := sdk.Escrow.Deposit(ctx, depositReq)
	if err != nil {
		return fmt.Errorf("failed to create deposit tx: %w", err)
	}

	fmt.Printf("Deposit Transaction Generated:\n")
	fmt.Printf("  Blockhash: %s\n", txResponse.RecentBlockhash)
	fmt.Printf("  Last Valid Block Height: %d\n", txResponse.LastValidBlockHeight)
	fmt.Printf("  Unsigned Tx (base64): %s...\n", txResponse.UnsignedTxBase64[:50])

	return nil
}

func demonstratePaymentIntent(ctx context.Context, sdk *shadowpay.ShadowPay, wallet string) error {
	// Create a payment intent
	intentReq := intent.CreateRequest{
		Amount:    1000000, // 0.001 SOL
		Recipient: wallet,
		Reference: "order-12345",
	}

	intentResp, err := sdk.Intent.Create(ctx, intentReq)
	if err != nil {
		return fmt.Errorf("failed to create payment intent: %w", err)
	}

	fmt.Printf("Payment Intent Created:\n")
	fmt.Printf("  Intent ID: %s\n", intentResp.IntentID)
	fmt.Printf("  Status: %s\n", intentResp.Status)
	fmt.Printf("  Client Secret: %s\n", intentResp.ClientSecret)

	// Verify the payment intent
	verifyResp, err := sdk.Intent.Verify(ctx, intentResp.IntentID)
	if err != nil {
		return fmt.Errorf("failed to verify payment intent: %w", err)
	}

	fmt.Printf("\nPayment Intent Verification:\n")
	fmt.Printf("  Status: %s\n", verifyResp.Status)
	fmt.Printf("  Verified: %v\n", verifyResp.Verified)

	return nil
}
