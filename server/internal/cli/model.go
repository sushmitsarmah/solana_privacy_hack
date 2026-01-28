package cli

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sol_privacy"
)

type view int

const (
	mainMenuView view = iota
	paymentView
	poolView
	tokenView
	authorizationView
	merchantView
	webhookView
	shadowIDView
	settingsView
)

type Model struct {
	client       *shadowpay.ShadowPay
	ctx          context.Context
	currentView  view
	cursor       int
	apiKey       string
	width        int
	height       int
	message      string
	messageStyle lipgloss.Style
	loading      bool
	loadingMsg   string
	showingInput bool
	inputForm    inputForm

	// Sub-models for different views
	paymentModel       *PaymentModel
	poolModel          *PoolModel
	tokenModel         *TokenModel
	authorizationModel *AuthorizationModel
}

func NewModel(apiKey string) Model {
	ctx := context.Background()
	var client *shadowpay.ShadowPay
	if apiKey != "" {
		client = shadowpay.New(apiKey)
	}
	return Model{
		ctx:          ctx,
		client:       client,
		currentView:  mainMenuView,
		cursor:       0,
		apiKey:       apiKey,
		messageStyle: successStyle,
		loading:      false,
		showingInput: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case operationSuccessMsg:
		m.loading = false
		m.showingInput = false
		m.message = msg.message
		m.messageStyle = successStyle
		return m, nil

	case operationErrorMsg:
		m.loading = false
		m.showingInput = false
		m.message = fmt.Sprintf("Error: %v", msg.err)
		m.messageStyle = errorStyle
		return m, nil

	case loadingMsg:
		m.loading = true
		m.loadingMsg = msg.message
		return m, nil

	case tea.KeyMsg:
		// Handle input form
		if m.showingInput {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.showingInput = false
				m.message = ""
				return m, nil
			default:
				cmd := m.inputForm.Update(msg)
				return m, cmd
			}
		}

		switch msg.String() {
		case "ctrl+c", "q":
			if m.currentView == mainMenuView {
				return m, tea.Quit
			}
			// Return to main menu from sub-views
			m.currentView = mainMenuView
			m.message = ""
			return m, nil

		case "esc":
			if m.currentView != mainMenuView {
				m.currentView = mainMenuView
				m.message = ""
				return m, nil
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			maxCursor := m.getMaxCursor()
			if m.cursor < maxCursor {
				m.cursor++
			}

		case "enter":
			return m.handleEnter()
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if m.loading {
		return m.renderLoading()
	}

	if m.showingInput {
		return m.inputForm.View(m.width, m.height)
	}

	switch m.currentView {
	case mainMenuView:
		return m.renderMainMenu()
	case paymentView:
		return m.renderPaymentView()
	case poolView:
		return m.renderPoolView()
	case tokenView:
		return m.renderTokenView()
	case authorizationView:
		return m.renderAuthorizationView()
	case merchantView:
		return m.renderMerchantView()
	case webhookView:
		return m.renderWebhookView()
	case shadowIDView:
		return m.renderShadowIDView()
	case settingsView:
		return m.renderSettingsView()
	default:
		return m.renderMainMenu()
	}
}

func (m Model) renderMainMenu() string {
	title := titleStyle.Render("🔒 ShadowPay CLI")

	var statusText string
	if m.client != nil {
		statusText = successStyle.Render("✓ Connected")
	} else {
		statusText = errorStyle.Render("✗ Not Connected (Set API Key)")
	}

	menu := []string{
		"💸 ZK Payments",
		"🏊 Privacy Pool",
		"🪙 Token Management",
		"🤖 Bot Authorization",
		"💰 Merchant Tools",
		"🔔 Webhooks",
		"👤 ShadowID",
		"⚙️  Settings",
		"🚪 Exit",
	}

	var menuStr string
	for i, item := range menu {
		cursor := "  "
		style := menuItemStyle
		if m.cursor == i {
			cursor = "❯ "
			style = selectedMenuItemStyle
		}
		menuStr += cursor + style.Render(item) + "\n"
	}

	help := helpStyle.Render("↑/↓: navigate • enter: select • q: quit")

	var messageBox string
	if m.message != "" {
		messageBox = "\n" + infoBoxStyle.Render(m.messageStyle.Render(m.message))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		statusText,
		"",
		menuStr,
		messageBox,
		"",
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) getMaxCursor() int {
	switch m.currentView {
	case mainMenuView:
		return 8 // 9 menu items (0-8)
	default:
		return 5
	}
}

func (m *Model) handleEnter() (tea.Model, tea.Cmd) {
	if m.currentView == mainMenuView {
		switch m.cursor {
		case 0: // ZK Payments
			if m.client == nil {
				m.message = "Please set API key in Settings first"
				m.messageStyle = errorStyle
				return *m, nil
			}
			m.currentView = paymentView
			m.cursor = 0
			m.message = ""

		case 1: // Privacy Pool
			if m.client == nil {
				m.message = "Please set API key in Settings first"
				m.messageStyle = errorStyle
				return *m, nil
			}
			m.currentView = poolView
			m.cursor = 0
			m.message = ""

		case 2: // Token Management
			if m.client == nil {
				m.message = "Please set API key in Settings first"
				m.messageStyle = errorStyle
				return *m, nil
			}
			m.currentView = tokenView
			m.cursor = 0
			m.message = ""

		case 3: // Bot Authorization
			if m.client == nil {
				m.message = "Please set API key in Settings first"
				m.messageStyle = errorStyle
				return *m, nil
			}
			m.currentView = authorizationView
			m.cursor = 0
			m.message = ""

		case 4: // Merchant Tools
			if m.client == nil {
				m.message = "Please set API key in Settings first"
				m.messageStyle = errorStyle
				return *m, nil
			}
			m.currentView = merchantView
			m.cursor = 0
			m.message = ""

		case 5: // Webhooks
			if m.client == nil {
				m.message = "Please set API key in Settings first"
				m.messageStyle = errorStyle
				return *m, nil
			}
			m.currentView = webhookView
			m.cursor = 0
			m.message = ""

		case 6: // ShadowID
			if m.client == nil {
				m.message = "Please set API key in Settings first"
				m.messageStyle = errorStyle
				return *m, nil
			}
			m.currentView = shadowIDView
			m.cursor = 0
			m.message = ""

		case 7: // Settings
			m.currentView = settingsView
			m.cursor = 0
			m.message = ""

		case 8: // Exit
			return *m, tea.Quit
		}
	} else {
		// Handle sub-menu selections
		switch m.currentView {
		case paymentView:
			return *m, m.handlePaymentSelection()
		case poolView:
			return *m, m.handlePoolSelection()
		case tokenView:
			return *m, m.handleTokenSelection()
		case authorizationView:
			return *m, m.handleAuthorizationSelection()
		case merchantView:
			return *m, m.handleMerchantSelection()
		case webhookView:
			return *m, m.handleWebhookSelection()
		case shadowIDView:
			return *m, m.handleShadowIDSelection()
		}
	}

	return *m, nil
}

func (m Model) renderPaymentView() string {
	title := titleStyle.Render("💸 ZK Payments")

	menu := []string{
		"📥 Deposit Funds",
		"📤 Withdraw Funds",
		"🔐 Prepare Payment",
		"✅ Authorize Payment",
		"🔍 Verify Access",
		"⚡ Settle Payment",
		"◀ Back",
	}

	var menuStr string
	for i, item := range menu {
		cursor := "  "
		style := menuItemStyle
		if m.cursor == i {
			cursor = "❯ "
			style = selectedMenuItemStyle
		}
		menuStr += cursor + style.Render(item) + "\n"
	}

	help := helpStyle.Render("↑/↓: navigate • enter: select • esc: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		headerStyle.Render("Select an operation:"),
		menuStr,
		"",
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) renderPoolView() string {
	title := titleStyle.Render("🏊 Privacy Pool")

	info := infoBoxStyle.Render(
		"Privacy pools mix your funds with other users\n" +
		"for maximum anonymity on-chain.",
	)

	menu := []string{
		"💰 Check Balance",
		"📥 Deposit to Pool",
		"📤 Withdraw from Pool",
		"📍 Get Deposit Address",
		"◀ Back",
	}

	var menuStr string
	for i, item := range menu {
		cursor := "  "
		style := menuItemStyle
		if m.cursor == i {
			cursor = "❯ "
			style = selectedMenuItemStyle
		}
		menuStr += cursor + style.Render(item) + "\n"
	}

	help := helpStyle.Render("↑/↓: navigate • enter: select • esc: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		info,
		headerStyle.Render("Select an operation:"),
		menuStr,
		"",
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) renderTokenView() string {
	title := titleStyle.Render("🪙 Token Management")

	menu := []string{
		"📋 List Supported Tokens",
		"➕ Add New Token",
		"✏️  Update Token",
		"🗑️  Remove Token",
		"◀ Back",
	}

	var menuStr string
	for i, item := range menu {
		cursor := "  "
		style := menuItemStyle
		if m.cursor == i {
			cursor = "❯ "
			style = selectedMenuItemStyle
		}
		menuStr += cursor + style.Render(item) + "\n"
	}

	help := helpStyle.Render("↑/↓: navigate • enter: select • esc: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		headerStyle.Render("Manage SPL tokens:"),
		menuStr,
		"",
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) renderAuthorizationView() string {
	title := titleStyle.Render("🤖 Bot Authorization")

	info := infoBoxStyle.Render(
		"Allow bots and services to spend from your\n" +
		"escrow with custom limits and expiration.",
	)

	menu := []string{
		"✅ Authorize Bot Spending",
		"📋 List Authorizations",
		"🚫 Revoke Authorization",
		"◀ Back",
	}

	var menuStr string
	for i, item := range menu {
		cursor := "  "
		style := menuItemStyle
		if m.cursor == i {
			cursor = "❯ "
			style = selectedMenuItemStyle
		}
		menuStr += cursor + style.Render(item) + "\n"
	}

	help := helpStyle.Render("↑/↓: navigate • enter: select • esc: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		info,
		headerStyle.Render("Manage bot permissions:"),
		menuStr,
		"",
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) renderMerchantView() string {
	title := titleStyle.Render("💰 Merchant Tools")

	menu := []string{
		"💵 View Earnings",
		"📊 Get Analytics",
		"📤 Withdraw Earnings",
		"🔓 Decrypt Amount",
		"◀ Back",
	}

	var menuStr string
	for i, item := range menu {
		cursor := "  "
		style := menuItemStyle
		if m.cursor == i {
			cursor = "❯ "
			style = selectedMenuItemStyle
		}
		menuStr += cursor + style.Render(item) + "\n"
	}

	help := helpStyle.Render("↑/↓: navigate • enter: select • esc: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		headerStyle.Render("Merchant operations:"),
		menuStr,
		"",
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) renderWebhookView() string {
	title := titleStyle.Render("🔔 Webhooks")

	menu := []string{
		"➕ Register Webhook",
		"📋 Get Configuration",
		"🧪 Test Webhook",
		"📜 View Logs",
		"📊 Get Stats",
		"🚫 Deactivate Webhook",
		"◀ Back",
	}

	var menuStr string
	for i, item := range menu {
		cursor := "  "
		style := menuItemStyle
		if m.cursor == i {
			cursor = "❯ "
			style = selectedMenuItemStyle
		}
		menuStr += cursor + style.Render(item) + "\n"
	}

	help := helpStyle.Render("↑/↓: navigate • enter: select • esc: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		headerStyle.Render("Webhook operations:"),
		menuStr,
		"",
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) renderShadowIDView() string {
	title := titleStyle.Render("👤 ShadowID")

	info := infoBoxStyle.Render(
		"Anonymous identity system using Merkle trees\n" +
		"for privacy-preserving authentication.",
	)

	menu := []string{
		"✅ Auto Register",
		"📝 Register Commitment",
		"🔍 Get Proof",
		"🌳 Get Tree Root",
		"📋 Check Status",
		"◀ Back",
	}

	var menuStr string
	for i, item := range menu {
		cursor := "  "
		style := menuItemStyle
		if m.cursor == i {
			cursor = "❯ "
			style = selectedMenuItemStyle
		}
		menuStr += cursor + style.Render(item) + "\n"
	}

	help := helpStyle.Render("↑/↓: navigate • enter: select • esc: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		info,
		headerStyle.Render("ShadowID operations:"),
		menuStr,
		"",
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) renderSettingsView() string {
	title := titleStyle.Render("⚙️  Settings")

	var statusBox string
	if m.apiKey == "" {
		statusBox = infoBoxStyle.Render(errorStyle.Render("⚠ API Key not set"))
	} else {
		maskedKey := m.apiKey
		if len(maskedKey) > 12 {
			maskedKey = maskedKey[:4] + "..." + maskedKey[len(maskedKey)-4:]
		}
		statusBox = infoBoxStyle.Render(
			successStyle.Render("✓ API Key: ") + maskedKey,
		)
	}

	instructions := lipgloss.NewStyle().
		Foreground(subtleColor).
		Render(
			"To set your API key, run:\n" +
			"export SHADOWPAY_API_KEY=your_key_here\n\n" +
			"Or create a .env file with:\n" +
			"SHADOWPAY_API_KEY=your_key_here",
		)

	help := helpStyle.Render("esc: back to main menu")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		statusBox,
		"",
		instructions,
		"",
		help,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m Model) renderLoading() string {
	spinner := "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
	frame := spinner[0:1] // Simple static spinner for now

	title := titleStyle.Render("Processing...")
	loadingText := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true).
		Render(frame + " " + m.loadingMsg)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		loadingText,
		"",
		helpStyle.Render("Please wait..."),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// Placeholder types for sub-models
type PaymentModel struct{}
type PoolModel struct{}
type TokenModel struct{}
type AuthorizationModel struct{}
