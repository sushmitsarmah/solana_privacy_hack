package cli

import (
	"context"
	"fmt"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"sol_privacy/internal/authorization"
	"sol_privacy/internal/merchant"
	"sol_privacy/internal/payment"
	"sol_privacy/internal/pool"
	"sol_privacy/internal/privacy"
	"sol_privacy/internal/shadowid"
	"sol_privacy/internal/token"
	"sol_privacy/internal/webhook"
)

// Message types for async operations
type operationSuccessMsg struct {
	message string
}

type operationErrorMsg struct {
	err error
}

type loadingMsg struct {
	message string
}

// Handle payment menu selections
func (m *Model) handlePaymentSelection() tea.Cmd {
	switch m.cursor {
	case 0: // Deposit Funds
		return m.showDepositForm()
	case 1: // Withdraw Funds
		return m.showWithdrawForm()
	case 2: // Prepare Payment
		return m.showPreparePaymentForm()
	case 3: // Authorize Payment
		return m.showAuthorizePaymentForm()
	case 4: // Verify Access
		return m.showVerifyAccessForm()
	case 5: // Settle Payment
		m.message = "Settle is complex - requires x402 payload. Use API directly."
		m.messageStyle = errorStyle
	case 6: // Back
		m.currentView = mainMenuView
		m.cursor = 0
	}
	return nil
}

func (m *Model) showDepositForm() tea.Cmd {
	m.inputForm = newInputForm(
		"💸 Deposit to Payment Account",
		[]string{"Wallet Address", "Amount (SOL)"},
		func(values []string) tea.Cmd {
			return m.performPaymentDeposit(values[0], values[1])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performPaymentDeposit(wallet, amountStr string) tea.Cmd {
	return func() tea.Msg {
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return operationErrorMsg{fmt.Errorf("invalid amount: %w", err)}
		}

		lamports := int64(amount * 1e9)
		req := payment.DepositRequest{
			WalletAddress: wallet,
			Amount:        lamports,
		}

		ctx := context.Background()
		resp, err := m.client.Payment.Deposit(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Deposit transaction created!\nBlockhash: %s\nSign and send the transaction to complete.", resp.RecentBlockhash),
		}
	}
}

func (m *Model) showWithdrawForm() tea.Cmd {
	m.inputForm = newInputForm(
		"📤 Withdraw from Payment Account",
		[]string{"Wallet Address", "Amount (SOL)"},
		func(values []string) tea.Cmd {
			return m.performPaymentWithdraw(values[0], values[1])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performPaymentWithdraw(wallet, amountStr string) tea.Cmd {
	return func() tea.Msg {
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return operationErrorMsg{fmt.Errorf("invalid amount: %w", err)}
		}

		lamports := int64(amount * 1e9)
		req := payment.WithdrawRequest{
			WalletAddress: wallet,
			Amount:        lamports,
		}

		ctx := context.Background()
		resp, err := m.client.Payment.Withdraw(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Withdraw transaction created!\nBlockhash: %s\n%s", resp.RecentBlockhash, resp.Message),
		}
	}
}

func (m *Model) showPreparePaymentForm() tea.Cmd {
	m.inputForm = newInputForm(
		"🔐 Prepare ZK Payment",
		[]string{"Receiver Commitment", "Amount (SOL)"},
		func(values []string) tea.Cmd {
			return m.performPreparePayment(values[0], values[1])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performPreparePayment(commitment, amountStr string) tea.Cmd {
	return func() tea.Msg {
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return operationErrorMsg{fmt.Errorf("invalid amount: %w", err)}
		}

		lamports := int64(amount * 1e9)
		req := payment.PrepareRequest{
			ReceiverCommitment: commitment,
			Amount:             lamports,
		}

		ctx := context.Background()
		resp, err := m.client.Payment.Prepare(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Payment prepared!\nPayment Hash: %s\nCommitment: %s\n%s", resp.PaymentHash, resp.Commitment[:20]+"...", resp.Message),
		}
	}
}

func (m *Model) showAuthorizePaymentForm() tea.Cmd {
	m.inputForm = newInputForm(
		"✅ Authorize Payment",
		[]string{"Commitment", "Nullifier", "Amount (SOL)", "Merchant Wallet"},
		func(values []string) tea.Cmd {
			return m.performAuthorizePayment(values[0], values[1], values[2], values[3])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performAuthorizePayment(commitment, nullifier, amountStr, merchant string) tea.Cmd {
	return func() tea.Msg {
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return operationErrorMsg{fmt.Errorf("invalid amount: %w", err)}
		}

		lamports := int64(amount * 1e9)
		req := payment.AuthorizeRequest{
			Commitment: commitment,
			Nullifier:  nullifier,
			Amount:     lamports,
			Merchant:   merchant,
		}

		ctx := context.Background()
		resp, err := m.client.Payment.Authorize(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Payment authorized!\nAccess Token: %s...\nExpires in: %d seconds\n%s",
				resp.AccessToken[:20], resp.ExpiresIn, resp.Message),
		}
	}
}

func (m *Model) showVerifyAccessForm() tea.Cmd {
	m.inputForm = newInputForm(
		"🔍 Verify Access Token",
		[]string{"Access Token"},
		func(values []string) tea.Cmd {
			return m.performVerifyAccess(values[0])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performVerifyAccess(token string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.client.Payment.VerifyAccess(ctx, token)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Invalid"
		if resp.Valid {
			status = "Valid ✓"
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Access verification: %s\nMerchant: %s\nAmount: %d\nExpires: %s\n%s",
				status, resp.Merchant, resp.Amount, resp.ExpiresAt, resp.Message),
		}
	}
}

// Pool operations
func (m *Model) handlePoolSelection() tea.Cmd {
	switch m.cursor {
	case 0: // Check Balance
		return m.showPoolBalanceForm()
	case 1: // Deposit
		return m.showPoolDepositForm()
	case 2: // Withdraw
		return m.showPoolWithdrawForm()
	case 3: // Get Deposit Address
		return m.performGetDepositAddress()
	case 4: // Back
		m.currentView = mainMenuView
		m.cursor = 0
	}
	return nil
}

func (m *Model) showPoolBalanceForm() tea.Cmd {
	m.inputForm = newInputForm(
		"💰 Check Pool Balance",
		[]string{"Wallet Address"},
		func(values []string) tea.Cmd {
			return m.performPoolBalance(values[0])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performPoolBalance(wallet string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		balance, err := m.client.Pool.GetBalance(ctx, wallet)
		if err != nil {
			return operationErrorMsg{err}
		}

		solBalance := float64(balance.Balance) / 1e9
		minDeposit := float64(balance.MinDeposit) / 1e9
		return operationSuccessMsg{
			message: fmt.Sprintf("Pool Balance: %.4f SOL (%d lamports)\nMin Deposit: %.4f SOL",
				solBalance, balance.Balance, minDeposit),
		}
	}
}

func (m *Model) showPoolDepositForm() tea.Cmd {
	m.inputForm = newInputForm(
		"💰 Deposit to Pool",
		[]string{"Wallet Address", "Amount (SOL)"},
		func(values []string) tea.Cmd {
			return m.performPoolDeposit(values[0], values[1])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performPoolDeposit(wallet, amountStr string) tea.Cmd {
	return func() tea.Msg {
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return operationErrorMsg{fmt.Errorf("invalid amount: %w", err)}
		}

		lamports := int64(amount * 1e9)
		req := pool.DepositRequest{
			WalletAddress: wallet,
			Amount:        lamports,
		}

		ctx := context.Background()
		resp, err := m.client.Pool.Deposit(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Pool deposit transaction created!\n%s", resp.Message),
		}
	}
}

func (m *Model) showPoolWithdrawForm() tea.Cmd {
	m.inputForm = newInputForm(
		"📤 Withdraw from Pool",
		[]string{"Wallet Address", "Amount (SOL)"},
		func(values []string) tea.Cmd {
			return m.performPoolWithdraw(values[0], values[1])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performPoolWithdraw(wallet, amountStr string) tea.Cmd {
	return func() tea.Msg {
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return operationErrorMsg{fmt.Errorf("invalid amount: %w", err)}
		}

		lamports := int64(amount * 1e9)
		req := pool.WithdrawRequest{
			WalletAddress: wallet,
			Amount:        lamports,
		}

		ctx := context.Background()
		resp, err := m.client.Pool.Withdraw(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		netSol := float64(resp.NetAmount) / 1e9
		feeSol := float64(resp.Fee) / 1e9
		return operationSuccessMsg{
			message: fmt.Sprintf("Pool withdrawal created!\nNet Amount: %.4f SOL\nFee: %.4f SOL\n%s",
				netSol, feeSol, resp.Message),
		}
	}
}

func (m *Model) performGetDepositAddress() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.client.Pool.GetDepositAddress(ctx)
		if err != nil {
			return operationErrorMsg{err}
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Pool Deposit Address:\n%s\nNetwork: %s", resp.DepositAddress, resp.Network),
		}
	}
}

// Token operations
func (m *Model) handleTokenSelection() tea.Cmd {
	switch m.cursor {
	case 0: // List Tokens
		return m.performListTokens()
	case 1: // Add Token
		return m.showAddTokenForm()
	case 2: // Update Token
		return m.showUpdateTokenForm()
	case 3: // Remove Token
		return m.showRemoveTokenForm()
	case 4: // Back
		m.currentView = mainMenuView
		m.cursor = 0
	}
	return nil
}

func (m *Model) performListTokens() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.client.Token.ListSupported(ctx)
		if err != nil {
			return operationErrorMsg{err}
		}

		var tokenList string
		if len(resp.Tokens) == 0 {
			tokenList = "No tokens configured"
		} else {
			for _, t := range resp.Tokens {
				status := "❌ Disabled"
				if t.Enabled {
					status = "✓ Enabled"
				}
				tokenList += fmt.Sprintf("\n• %s (%s) %s\n  Mint: %s\n  Decimals: %d",
					t.Symbol, status, t.Mint, t.Mint, t.Decimals)
			}
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Supported Tokens:%s", tokenList),
		}
	}
}

func (m *Model) showAddTokenForm() tea.Cmd {
	m.inputForm = newInputForm(
		"➕ Add Token",
		[]string{"Mint Address", "Symbol", "Decimals"},
		func(values []string) tea.Cmd {
			return m.performAddToken(values[0], values[1], values[2])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performAddToken(mint, symbol, decimalsStr string) tea.Cmd {
	return func() tea.Msg {
		decimals, err := strconv.Atoi(decimalsStr)
		if err != nil {
			return operationErrorMsg{fmt.Errorf("invalid decimals: %w", err)}
		}

		req := token.AddRequest{
			Mint:     mint,
			Symbol:   symbol,
			Decimals: decimals,
			Enabled:  true,
		}

		ctx := context.Background()
		resp, err := m.client.Token.Add(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Failed"
		if resp.Success {
			status = "Success ✓"
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Add Token: %s\n%s", status, resp.Message),
		}
	}
}

func (m *Model) showUpdateTokenForm() tea.Cmd {
	m.inputForm = newInputForm(
		"✏️ Update Token",
		[]string{"Mint Address", "New Symbol (optional)", "Enabled (true/false)"},
		func(values []string) tea.Cmd {
			return m.performUpdateToken(values[0], values[1], values[2])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performUpdateToken(mint, symbol, enabledStr string) tea.Cmd {
	return func() tea.Msg {
		req := token.UpdateRequest{}

		if symbol != "" {
			req.Symbol = &symbol
		}

		if enabledStr != "" {
			enabled := enabledStr == "true"
			req.Enabled = &enabled
		}

		ctx := context.Background()
		resp, err := m.client.Token.Update(ctx, mint, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Failed"
		if resp.Success {
			status = "Success ✓"
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Update Token: %s\n%s", status, resp.Message),
		}
	}
}

func (m *Model) showRemoveTokenForm() tea.Cmd {
	m.inputForm = newInputForm(
		"🗑️ Remove Token",
		[]string{"Mint Address"},
		func(values []string) tea.Cmd {
			return m.performRemoveToken(values[0])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performRemoveToken(mint string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.client.Token.Remove(ctx, mint)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Failed"
		if resp.Success {
			status = "Success ✓"
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Remove Token: %s\n%s", status, resp.Message),
		}
	}
}

// Merchant operations
func (m *Model) handleMerchantSelection() tea.Cmd {
	switch m.cursor {
	case 0: // View Earnings
		return m.performViewEarnings()
	case 1: // Get Analytics
		return m.showGetAnalyticsForm()
	case 2: // Withdraw Earnings
		return m.showWithdrawEarningsForm()
	case 3: // Decrypt Amount
		return m.showDecryptAmountForm()
	case 4: // Back
		m.currentView = mainMenuView
		m.cursor = 0
	}
	return nil
}

func (m *Model) performViewEarnings() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.client.Merchant.GetEarnings(ctx)
		if err != nil {
			return operationErrorMsg{err}
		}

		totalSol := float64(resp.TotalEarnings) / 1e9
		withdrawableSol := float64(resp.WithdrawableSOL) / 1e9
		pendingSol := float64(resp.PendingSettlement) / 1e9

		var tokenBreakdown string
		for _, t := range resp.TokenBreakdown {
			tokenSol := float64(t.Amount) / 1e9
			tokenBreakdown += fmt.Sprintf("\n  • %s: %.4f SOL", t.Symbol, tokenSol)
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Merchant Earnings:\nTotal: %.4f SOL ($%s)\nWithdrawable: %.4f SOL\nPending: %.4f SOL\n\nToken Breakdown:%s",
				totalSol, resp.TotalUsdValue, withdrawableSol, pendingSol, tokenBreakdown),
		}
	}
}

func (m *Model) showGetAnalyticsForm() tea.Cmd {
	m.inputForm = newInputForm(
		"📊 Get Analytics",
		[]string{"Start Date (YYYY-MM-DD, optional)", "End Date (YYYY-MM-DD, optional)"},
		func(values []string) tea.Cmd {
			return m.performGetAnalytics(values[0], values[1])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performGetAnalytics(startDate, endDate string) tea.Cmd {
	return func() tea.Msg {
		req := merchant.AnalyticsRequest{
			StartDate: startDate,
			EndDate:   endDate,
		}

		ctx := context.Background()
		resp, err := m.client.Merchant.GetAnalytics(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		totalVolSol := float64(resp.TotalVolume) / 1e9
		avgSol := float64(resp.AveragePayment) / 1e9

		var topResources string
		for i, r := range resp.TopResources {
			if i >= 5 {
				break
			}
			amtSol := float64(r.TotalAmount) / 1e9
			topResources += fmt.Sprintf("\n  %d. %s: %d payments, %.4f SOL", i+1, r.Resource, r.PaymentCount, amtSol)
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Analytics:\nTotal Payments: %d\nTotal Volume: %.4f SOL\nAvg Payment: %.4f SOL\nUnique Customers: %d\nSuccess Rate: %.1f%%\nPending: %d\n\nTop Resources:%s",
				resp.TotalPayments, totalVolSol, avgSol, resp.UniqueCustomers, resp.SuccessRate, resp.PendingPayments, topResources),
		}
	}
}

func (m *Model) showWithdrawEarningsForm() tea.Cmd {
	m.inputForm = newInputForm(
		"📤 Withdraw Earnings",
		[]string{"Amount (SOL)", "Destination Wallet"},
		func(values []string) tea.Cmd {
			return m.performWithdrawEarnings(values[0], values[1])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performWithdrawEarnings(amountStr, destination string) tea.Cmd {
	return func() tea.Msg {
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return operationErrorMsg{fmt.Errorf("invalid amount: %w", err)}
		}

		lamports := int64(amount * 1e9)
		req := merchant.WithdrawRequest{
			Amount:      lamports,
			Destination: destination,
		}

		ctx := context.Background()
		resp, err := m.client.Merchant.Withdraw(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Failed"
		if resp.Success {
			status = "Success ✓"
		}

		netSol := float64(resp.NetAmount) / 1e9
		feeSol := float64(resp.Fee) / 1e9

		return operationSuccessMsg{
			message: fmt.Sprintf("Withdraw Earnings: %s\nWithdrawal ID: %s\nAmount: %.4f SOL\nFee: %.4f SOL\nNet: %.4f SOL\n%s",
				status, resp.WithdrawalID, amount, feeSol, netSol, resp.Message),
		}
	}
}

func (m *Model) showDecryptAmountForm() tea.Cmd {
	m.inputForm = newInputForm(
		"🔓 Decrypt Amount",
		[]string{"Encrypted Ciphertext (hex)", "Private Key (hex)"},
		func(values []string) tea.Cmd {
			return m.performDecryptAmount(values[0], values[1])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performDecryptAmount(ciphertext, privKey string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		req := privacy.DecryptRequest{
			Ciphertext: ciphertext,
			PrivateKey: privKey,
		}
		resp, err := m.client.Privacy.Decrypt(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		solAmount := float64(resp.Amount) / 1e9
		return operationSuccessMsg{
			message: fmt.Sprintf("Decrypted Amount: %.4f SOL (%d lamports)", solAmount, resp.Amount),
		}
	}
}

// Webhook operations
func (m *Model) handleWebhookSelection() tea.Cmd {
	switch m.cursor {
	case 0: // Register Webhook
		return m.showRegisterWebhookForm()
	case 1: // Get Configuration
		return m.performGetWebhookConfig()
	case 2: // Test Webhook
		return m.showTestWebhookForm()
	case 3: // View Logs
		return m.showViewLogsForm()
	case 4: // Get Stats
		return m.performGetWebhookStats()
	case 5: // Deactivate Webhook
		return m.showDeactivateWebhookForm()
	case 6: // Back
		m.currentView = mainMenuView
		m.cursor = 0
	}
	return nil
}

func (m *Model) showRegisterWebhookForm() tea.Cmd {
	m.inputForm = newInputForm(
		"➕ Register Webhook",
		[]string{"Webhook URL (https://...)", "Events (comma-separated)", "Secret (optional)"},
		func(values []string) tea.Cmd {
			return m.performRegisterWebhook(values[0], values[1], values[2])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performRegisterWebhook(url, eventsStr, secret string) tea.Cmd {
	return func() tea.Msg {
		// Parse comma-separated events
		events := []string{}
		if eventsStr != "" {
			for _, event := range splitAndTrim(eventsStr, ",") {
				events = append(events, event)
			}
		}

		req := webhook.RegisterRequest{
			URL:    url,
			Events: events,
			Secret: secret,
		}

		ctx := context.Background()
		resp, err := m.client.Webhook.Register(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Failed"
		if resp.Success {
			status = "Success ✓"
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Register Webhook: %s\nWebhook ID: %s\nURL: %s\nEvents: %v\nCreated: %s\n%s",
				status, resp.WebhookID, resp.URL, resp.Events, resp.CreatedAt, resp.Message),
		}
	}
}

func (m *Model) performGetWebhookConfig() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.client.Webhook.GetConfig(ctx)
		if err != nil {
			return operationErrorMsg{err}
		}

		activeStatus := "Inactive ❌"
		if resp.Active {
			activeStatus = "Active ✓"
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Webhook Configuration:\nWebhook ID: %s\nURL: %s\nEvents: %v\nStatus: %s\nCreated: %s\nUpdated: %s",
				resp.WebhookID, resp.URL, resp.Events, activeStatus, resp.CreatedAt, resp.UpdatedAt),
		}
	}
}

func (m *Model) showTestWebhookForm() tea.Cmd {
	m.inputForm = newInputForm(
		"🧪 Test Webhook",
		[]string{"Webhook ID (optional)", "Event Type (optional)"},
		func(values []string) tea.Cmd {
			return m.performTestWebhook(values[0], values[1])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performTestWebhook(webhookID, event string) tea.Cmd {
	return func() tea.Msg {
		req := webhook.TestRequest{
			WebhookID: webhookID,
			Event:     event,
		}

		ctx := context.Background()
		resp, err := m.client.Webhook.Test(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Failed ❌"
		if resp.Success {
			status = "Success ✓"
		}

		errorInfo := ""
		if resp.Error != "" {
			errorInfo = fmt.Sprintf("\nError: %s", resp.Error)
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Test Webhook: %s\nStatus Code: %d\nResponse Time: %d ms\n%s%s",
				status, resp.StatusCode, resp.ResponseTime, resp.Message, errorInfo),
		}
	}
}

func (m *Model) showViewLogsForm() tea.Cmd {
	m.inputForm = newInputForm(
		"📜 View Webhook Logs",
		[]string{"Webhook ID (optional)", "Limit (default 50)"},
		func(values []string) tea.Cmd {
			return m.performViewLogs(values[0], values[1])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performViewLogs(webhookID, limitStr string) tea.Cmd {
	return func() tea.Msg {
		limit := 50
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}

		req := webhook.LogsRequest{
			WebhookID: webhookID,
			Limit:     limit,
		}

		ctx := context.Background()
		resp, err := m.client.Webhook.GetLogs(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		var logsStr string
		if len(resp.Logs) == 0 {
			logsStr = "No logs found"
		} else {
			for i, log := range resp.Logs {
				if i >= 10 { // Show max 10 logs
					logsStr += fmt.Sprintf("\n... and %d more", len(resp.Logs)-10)
					break
				}
				status := "❌"
				if log.Success {
					status = "✓"
				}
				logsStr += fmt.Sprintf("\n%s %s | Event: %s | Code: %d | Time: %dms | Attempt: %d\n   %s",
					status, log.Timestamp, log.Event, log.StatusCode, log.ResponseTime, log.Attempt, log.ID)
			}
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Webhook Logs (Total: %d):%s", resp.TotalCount, logsStr),
		}
	}
}

func (m *Model) performGetWebhookStats() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.client.Webhook.GetStats(ctx)
		if err != nil {
			return operationErrorMsg{err}
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Webhook Statistics:\nTotal Deliveries: %d\nSuccessful: %d\nFailed: %d\nSuccess Rate: %.1f%%\nAvg Response Time: %d ms\nLast Delivery: %s\nLast Success: %s\nLast Failure: %s",
				resp.TotalDeliveries, resp.SuccessfulDeliveries, resp.FailedDeliveries,
				resp.SuccessRate, resp.AverageResponseTime, resp.LastDelivery, resp.LastSuccess, resp.LastFailure),
		}
	}
}

func (m *Model) showDeactivateWebhookForm() tea.Cmd {
	m.inputForm = newInputForm(
		"🚫 Deactivate Webhook",
		[]string{"Webhook ID"},
		func(values []string) tea.Cmd {
			return m.performDeactivateWebhook(values[0])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performDeactivateWebhook(webhookID string) tea.Cmd {
	return func() tea.Msg {
		req := webhook.DeactivateRequest{
			WebhookID: webhookID,
		}

		ctx := context.Background()
		resp, err := m.client.Webhook.Deactivate(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Failed"
		if resp.Success {
			status = "Success ✓"
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Deactivate Webhook: %s\nWebhook ID: %s\n%s",
				status, resp.WebhookID, resp.Message),
		}
	}
}

// ShadowID operations
func (m *Model) handleShadowIDSelection() tea.Cmd {
	switch m.cursor {
	case 0: // Auto Register
		return m.showAutoRegisterForm()
	case 1: // Register Commitment
		return m.showRegisterCommitmentForm()
	case 2: // Get Proof
		return m.showGetProofForm()
	case 3: // Get Tree Root
		return m.performGetTreeRoot()
	case 4: // Check Status
		return m.showCheckStatusForm()
	case 5: // Back
		m.currentView = mainMenuView
		m.cursor = 0
	}
	return nil
}

func (m *Model) showAutoRegisterForm() tea.Cmd {
	m.inputForm = newInputForm(
		"✅ Auto Register ShadowID",
		[]string{"Wallet Address", "Signature (base58)", "Message"},
		func(values []string) tea.Cmd {
			return m.performAutoRegister(values[0], values[1], values[2])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performAutoRegister(wallet, signature, message string) tea.Cmd {
	return func() tea.Msg {
		req := shadowid.AutoRegisterRequest{
			WalletAddress: wallet,
			Signature:     signature,
			Message:       message,
		}

		ctx := context.Background()
		resp, err := m.client.ShadowID.AutoRegister(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Failed"
		if resp.Success {
			status = "Success ✓"
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Auto Register: %s\nCommitment: %s\nLeaf Index: %d\n%s",
				status, resp.Commitment, resp.LeafIndex, resp.Message),
		}
	}
}

func (m *Model) showRegisterCommitmentForm() tea.Cmd {
	m.inputForm = newInputForm(
		"📝 Register Commitment",
		[]string{"Poseidon Hash Commitment"},
		func(values []string) tea.Cmd {
			return m.performRegisterCommitment(values[0])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performRegisterCommitment(commitment string) tea.Cmd {
	return func() tea.Msg {
		req := shadowid.RegisterRequest{
			Commitment: commitment,
		}

		ctx := context.Background()
		resp, err := m.client.ShadowID.Register(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Failed"
		if resp.Success {
			status = "Success ✓"
		}

		txInfo := ""
		if resp.TxHash != "" {
			txInfo = fmt.Sprintf("\nTx Hash: %s", resp.TxHash)
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Register Commitment: %s\nLeaf Index: %d%s\n%s",
				status, resp.LeafIndex, txInfo, resp.Message),
		}
	}
}

func (m *Model) showGetProofForm() tea.Cmd {
	m.inputForm = newInputForm(
		"🔍 Get Merkle Proof",
		[]string{"Commitment"},
		func(values []string) tea.Cmd {
			return m.performGetProof(values[0])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performGetProof(commitment string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.client.ShadowID.GetProof(ctx, commitment)
		if err != nil {
			return operationErrorMsg{err}
		}

		proofStr := ""
		maxProofShow := 3
		for i, p := range resp.Proof {
			if i >= maxProofShow {
				proofStr += fmt.Sprintf("\n  ... and %d more hashes", len(resp.Proof)-maxProofShow)
				break
			}
			displayHash := p
			if len(p) > 20 {
				displayHash = p[:20] + "..."
			}
			proofStr += fmt.Sprintf("\n  [%d] %s", i, displayHash)
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Merkle Proof:\nCommitment: %s\nLeaf Index: %d\nRoot: %s\nProof (%d hashes):%s",
				resp.Commitment[:20]+"...", resp.LeafIndex, resp.Root[:20]+"...", len(resp.Proof), proofStr),
		}
	}
}

func (m *Model) performGetTreeRoot() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.client.ShadowID.GetRoot(ctx)
		if err != nil {
			return operationErrorMsg{err}
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Merkle Tree Root:\nRoot: %s\nTree Depth: %d\nLeaf Count: %d",
				resp.Root, resp.TreeDepth, resp.LeafCount),
		}
	}
}

func (m *Model) showCheckStatusForm() tea.Cmd {
	m.inputForm = newInputForm(
		"📋 Check Registration Status",
		[]string{"Commitment"},
		func(values []string) tea.Cmd {
			return m.performCheckStatus(values[0])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performCheckStatus(commitment string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.client.ShadowID.GetStatus(ctx, commitment)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Not Registered ❌"
		leafInfo := ""
		if resp.Registered {
			status = "Registered ✓"
			leafInfo = fmt.Sprintf("\nLeaf Index: %d", resp.LeafIndex)
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Registration Status: %s\nCommitment: %s%s",
				status, resp.Commitment[:20]+"...", leafInfo),
		}
	}
}

// Authorization operations
func (m *Model) handleAuthorizationSelection() tea.Cmd {
	switch m.cursor {
	case 0: // Authorize Bot Spending
		return m.showAuthorizeSpendingForm()
	case 1: // List Authorizations
		return m.showListAuthorizationsForm()
	case 2: // Revoke Authorization
		return m.showRevokeAuthorizationForm()
	case 3: // Back
		m.currentView = mainMenuView
		m.cursor = 0
	}
	return nil
}

func (m *Model) showAuthorizeSpendingForm() tea.Cmd {
	m.inputForm = newInputForm(
		"✅ Authorize Bot Spending",
		[]string{"User Wallet", "Authorized Service", "Max Per Tx (SOL)", "Max Daily (SOL)", "Valid Until (days from now)", "User Signature (base58)"},
		func(values []string) tea.Cmd {
			return m.performAuthorizeSpending(values[0], values[1], values[2], values[3], values[4], values[5])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performAuthorizeSpending(wallet, service, maxPerTx, maxDaily, validDays, signature string) tea.Cmd {
	return func() tea.Msg {
		// Calculate valid until timestamp
		days, err := strconv.Atoi(validDays)
		if err != nil {
			return operationErrorMsg{fmt.Errorf("invalid valid days: %w", err)}
		}
		validUntil := time.Now().Add(time.Duration(days) * 24 * time.Hour).Unix()

		req := authorization.AuthorizeSpendingRequest{
			UserWallet:        wallet,
			AuthorizedService: service,
			MaxAmountPerTx:    maxPerTx,
			MaxDailySpend:     maxDaily,
			ValidUntil:        validUntil,
			UserSignature:     signature,
		}

		ctx := context.Background()
		resp, err := m.client.Authorization.AuthorizeSpending(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Failed"
		if resp.Success {
			status = "Success ✓"
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Authorize Spending: %s\nAuthorization ID: %d\n%s",
				status, resp.AuthorizationID, resp.Message),
		}
	}
}

func (m *Model) showListAuthorizationsForm() tea.Cmd {
	m.inputForm = newInputForm(
		"📋 List Authorizations",
		[]string{"Wallet Address"},
		func(values []string) tea.Cmd {
			return m.performListAuthorizations(values[0])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performListAuthorizations(wallet string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.client.Authorization.ListAuthorizations(ctx, wallet)
		if err != nil {
			return operationErrorMsg{err}
		}

		var authList string
		if len(resp.Authorizations) == 0 {
			authList = "No authorizations found"
		} else {
			for i, auth := range resp.Authorizations {
				status := "Active ✓"
				if auth.Revoked {
					status = "Revoked ❌"
				}

				maxPerTxSol := float64(auth.MaxAmountPerTx) / 1e9
				maxDailySol := float64(auth.MaxDailySpend) / 1e9
				spentTodaySol := float64(auth.SpentToday) / 1e9

				validUntilTime := time.Unix(auth.ValidUntil, 0)
				createdTime := time.Unix(auth.CreatedAt, 0)

				authList += fmt.Sprintf("\n\n[%d] %s\nService: %s\nMax Per Tx: %.4f SOL\nMax Daily: %.4f SOL\nSpent Today: %.4f SOL\nValid Until: %s\nCreated: %s\nLast Reset: %s",
					i+1, status, auth.AuthorizedService, maxPerTxSol, maxDailySol, spentTodaySol,
					validUntilTime.Format("2006-01-02 15:04:05"), createdTime.Format("2006-01-02 15:04:05"), auth.LastResetDate)
			}
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Authorizations for %s:%s", wallet, authList),
		}
	}
}

func (m *Model) showRevokeAuthorizationForm() tea.Cmd {
	m.inputForm = newInputForm(
		"🚫 Revoke Authorization",
		[]string{"User Wallet", "Authorized Service", "User Signature (base58)"},
		func(values []string) tea.Cmd {
			return m.performRevokeAuthorization(values[0], values[1], values[2])
		},
	)
	m.showingInput = true
	return nil
}

func (m *Model) performRevokeAuthorization(wallet, service, signature string) tea.Cmd {
	return func() tea.Msg {
		req := authorization.RevokeAuthorizationRequest{
			UserWallet:        wallet,
			AuthorizedService: service,
			UserSignature:     signature,
		}

		ctx := context.Background()
		resp, err := m.client.Authorization.RevokeAuthorization(ctx, req)
		if err != nil {
			return operationErrorMsg{err}
		}

		status := "Failed"
		if resp.Success {
			status = "Success ✓"
		}

		return operationSuccessMsg{
			message: fmt.Sprintf("Revoke Authorization: %s\nAuthorization ID: %d\n%s",
				status, resp.AuthorizationID, resp.Message),
		}
	}
}

// Helper function for splitting and trimming strings
func splitAndTrim(s, sep string) []string {
	var result []string
	parts := []string{}
	current := ""
	for _, char := range s {
		if string(char) == sep {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}

	for _, part := range parts {
		trimmed := ""
		// Trim leading/trailing spaces
		start := 0
		end := len(part)
		for start < len(part) && part[start] == ' ' {
			start++
		}
		for end > start && part[end-1] == ' ' {
			end--
		}
		if start < end {
			trimmed = part[start:end]
		}
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
