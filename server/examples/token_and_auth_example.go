package main

import (
	"context"
	"fmt"
	"log"

	"sol_privacy"
	"sol_privacy/internal/authorization"
	"sol_privacy/internal/token"
)

// Example demonstrating Token Management and Bot Authorization features
func main() {
	ctx := context.Background()
	client := shadowpay.New("your_api_key_here")

	fmt.Println("=== Token Management Examples ===\n")

	// 1. List supported SPL tokens
	fmt.Println("1. Listing supported tokens...")
	tokens, err := client.Token.ListSupported(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   Found %d supported tokens\n", len(tokens.Tokens))
		for _, t := range tokens.Tokens {
			fmt.Printf("   - %s (%s): %d decimals, enabled=%v\n",
				t.Symbol, t.Mint, t.Decimals, t.Enabled)
		}
	}

	// 2. Add new SPL token (admin operation)
	fmt.Println("\n2. Adding new token (USDC)...")
	addResp, err := client.Token.Add(ctx, token.AddRequest{
		Mint:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Symbol:   "USDC",
		Decimals: 6,
		Enabled:  true,
	})
	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   ✅ %s\n", addResp.Message)
	}

	// 3. Update token configuration
	fmt.Println("\n3. Updating token configuration...")
	enabled := false
	_, err = client.Token.Update(ctx, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", token.UpdateRequest{
		Enabled: &enabled,
	})
	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   ✅ Token updated successfully\n")
	}

	// 4. Remove token
	fmt.Println("\n4. Removing token...")
	_, err = client.Token.Remove(ctx, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   ✅ Token removed successfully\n")
	}

	fmt.Println("\n\n=== Bot/Agent Authorization Examples ===\n")

	// 5. Authorize bot spending
	fmt.Println("5. Authorizing bot spending...")
	authResp, err := client.Authorization.AuthorizeSpending(ctx, authorization.AuthorizeSpendingRequest{
		UserWallet:        "AVSSWPbWRYDF7w8GZcrP6yVWsmRWPshMnziHqFQ5RaDR",
		AuthorizedService: "MyTradingBot_v1",
		MaxAmountPerTx:    "0.01",  // Max 0.01 SOL per transaction
		MaxDailySpend:     "1.0",   // Max 1 SOL per day
		ValidUntil:        1735689600, // Expires Dec 31, 2024
		UserSignature:     "base58_signature_here",
	})
	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   ✅ %s\n", authResp.Message)
		fmt.Printf("   Authorization ID: %d\n", authResp.AuthorizationID)
	}

	// 6. List user authorizations
	fmt.Println("\n6. Listing active bot authorizations...")
	auths, err := client.Authorization.ListAuthorizations(ctx, "AVSSWPbWRYDF7w8GZcrP6yVWsmRWPshMnziHqFQ5RaDR")
	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   Found %d active authorizations:\n", len(auths.Authorizations))
		for _, auth := range auths.Authorizations {
			fmt.Printf("\n   Authorization #%d:\n", auth.ID)
			fmt.Printf("   - Service: %s\n", auth.AuthorizedService)
			fmt.Printf("   - Max per tx: %d lamports\n", auth.MaxAmountPerTx)
			fmt.Printf("   - Daily limit: %d lamports\n", auth.MaxDailySpend)
			fmt.Printf("   - Spent today: %d lamports\n", auth.SpentToday)
			fmt.Printf("   - Valid until: %d\n", auth.ValidUntil)
			fmt.Printf("   - Revoked: %v\n", auth.Revoked)
		}
	}

	// 7. Revoke bot authorization
	fmt.Println("\n7. Revoking bot authorization...")
	revokeResp, err := client.Authorization.RevokeAuthorization(ctx, authorization.RevokeAuthorizationRequest{
		UserWallet:        "AVSSWPbWRYDF7w8GZcrP6yVWsmRWPshMnziHqFQ5RaDR",
		AuthorizedService: "MyTradingBot_v1",
		UserSignature:     "base58_signature_here",
	})
	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   ✅ %s\n", revokeResp.Message)
		fmt.Printf("   Authorization ID %d revoked\n", revokeResp.AuthorizationID)
	}

	fmt.Println("\n✅ All Token Management and Authorization features demonstrated!")
}
