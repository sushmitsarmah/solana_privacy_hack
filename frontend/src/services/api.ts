// ShadowPay API Service
const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

class ShadowPayAPI {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({
        error: response.statusText,
      }));
      throw new Error(error.error || 'API request failed');
    }

    return response.json();
  }

  // Pool Operations
  async getPoolBalance(wallet: string) {
    return this.request<{
      balance: number;
      min_deposit: number;
    }>(`/pool/balance/${wallet}`);
  }

  async depositToPool(wallet: string, amount: number) {
    return this.request<{
      unsigned_tx_base64: string;
      recent_blockhash: string;
      message: string;
    }>(`/pool/deposit`, {
      method: 'POST',
      body: JSON.stringify({
        wallet_address: wallet,
        amount,
      }),
    });
  }

  async withdrawFromPool(wallet: string, amount: number) {
    return this.request<{
      net_amount: number;
      fee: number;
      unsigned_tx_base64: string;
      recent_blockhash: string;
      message: string;
    }>(`/pool/withdraw`, {
      method: 'POST',
      body: JSON.stringify({
        wallet_address: wallet,
        amount,
      }),
    });
  }

  async getDepositAddress() {
    return this.request<{
      deposit_address: string;
      network: string;
    }>(`/pool/deposit-address`);
  }

  // Payment Operations
  async depositToPayment(wallet: string, amount: number) {
    return this.request<{
      unsigned_tx_base64: string;
      recent_blockhash: string;
      message: string;
    }>(`/payment/deposit`, {
      method: 'POST',
      body: JSON.stringify({
        wallet_address: wallet,
        amount,
      }),
    });
  }

  async withdrawFromPayment(wallet: string, amount: number) {
    return this.request<{
      unsigned_tx_base64: string;
      recent_blockhash: string;
      message: string;
    }>(`/payment/withdraw`, {
      method: 'POST',
      body: JSON.stringify({
        wallet_address: wallet,
        amount,
      }),
    });
  }

  async preparePayment(receiverCommitment: string, amount: number) {
    return this.request<{
      payment_hash: string;
      commitment: string;
      message: string;
    }>(`/payment/prepare`, {
      method: 'POST',
      body: JSON.stringify({
        receiver_commitment: receiverCommitment,
        amount,
      }),
    });
  }

  async authorizePayment(
    commitment: string,
    nullifier: string,
    amount: number,
    merchant: string
  ) {
    return this.request<{
      access_token: string;
      expires_in: number;
      message: string;
    }>(`/payment/authorize`, {
      method: 'POST',
      body: JSON.stringify({
        commitment,
        nullifier,
        amount,
        merchant,
      }),
    });
  }

  // Token Operations
  async listTokens() {
    return this.request<{
      tokens: Array<{
        mint: string;
        symbol: string;
        decimals: number;
        enabled: boolean;
      }>;
    }>(`/token/list`);
  }

  async addToken(mint: string, symbol: string, decimals: number) {
    return this.request<{
      success: boolean;
      message: string;
    }>(`/token/add`, {
      method: 'POST',
      body: JSON.stringify({
        mint,
        symbol,
        decimals,
        enabled: true,
      }),
    });
  }

  // Merchant Operations
  async getMerchantEarnings() {
    return this.request<{
      total_earnings: number;
      withdrawable_sol: number;
      pending_settlement: number;
      total_usd_value: string;
      token_breakdown: Array<{
        symbol: string;
        amount: number;
      }>;
    }>(`/merchant/earnings`);
  }

  async withdrawEarnings(amount: number, destination: string) {
    return this.request<{
      success: boolean;
      withdrawal_id: string;
      amount: number;
      fee: number;
      net_amount: number;
      message: string;
    }>(`/merchant/withdraw`, {
      method: 'POST',
      body: JSON.stringify({
        amount,
        destination,
      }),
    });
  }

  // ShadowID Operations
  async getShadowIDRoot() {
    return this.request<{
      root: string;
      tree_depth: number;
      leaf_count: number;
    }>(`/shadowid/root`);
  }

  async checkShadowIDStatus(commitment: string) {
    return this.request<{
      commitment: string;
      registered: boolean;
      leaf_index?: number;
    }>(`/shadowid/status/${commitment}`);
  }

  // Authorization Operations
  async listAuthorizations(wallet: string) {
    return this.request<{
      authorizations: Array<{
        id: number;
        user_wallet: string;
        authorized_service: string;
        max_amount_per_tx: number;
        max_daily_spend: number;
        spent_today: number;
        valid_until: number;
        revoked: boolean;
        created_at: number;
      }>;
    }>(`/authorization/list/${wallet}`);
  }
}

export const api = new ShadowPayAPI();
export default api;
