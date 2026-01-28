import { useState, useEffect } from 'react';
import { useWalletConnection } from '@solana/react-hooks';
import { Shield, Send, ArrowDown, ArrowUp, Wallet, CheckCircle, AlertCircle } from 'lucide-react';
import api from '../services/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Badge } from './ui/badge';

interface PaymentHistoryItem {
  id: string;
  type: 'deposit' | 'withdraw' | 'payment';
  amount: number;
  status: 'pending' | 'completed' | 'failed';
  timestamp: number;
  hash?: string;
}

export function ZKPaymentManager() {
  const { wallet } = useWalletConnection();
  const [activeTab, setActiveTab] = useState('account');

  // Account management state
  const [accountBalance, setAccountBalance] = useState<number | null>(null);
  const [amount, setAmount] = useState('');

  // Payment preparation state
  const [receiverCommitment, setReceiverCommitment] = useState('');
  const [paymentAmount, setPaymentAmount] = useState('');
  const [tokenMint, setTokenMint] = useState('');
  const [paymentHash, setPaymentHash] = useState('');
  const [preparedCommitment, setPreparedCommitment] = useState('');

  // Authorization state
  const [nullifier, setNullifier] = useState('');
  const [authorizingCommitment, setAuthorizingCommitment] = useState('');
  const [authAmount, setAuthAmount] = useState('');
  const [merchant, setMerchant] = useState('');
  const [accessToken, setAccessToken] = useState('');

  // Settlement state
  const [paymentHeader, setPaymentHeader] = useState('');
  const [resource, setResource] = useState('');
  const [settlementResult, setSettlementResult] = useState<{ success: boolean; tx_sig?: string; message?: string } | null>(null);

  // Verification state
  const [verifyToken, setVerifyToken] = useState('');
  const [verifyResult, setVerifyResult] = useState<any>(null);

  // UI state
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [paymentHistory, setPaymentHistory] = useState<PaymentHistoryItem[]>([]);

  const walletAddress = wallet?.account.address.toString();

  useEffect(() => {
    if (walletAddress) {
      loadAccountBalance();
      loadPaymentHistory();
    }
  }, [walletAddress]);

  const loadAccountBalance = async () => {
    if (!walletAddress) return;

    try {
      setLoading(true);
      // For demo purposes, we'll use the pool balance endpoint
      // In a real implementation, there should be a payment account balance endpoint
      const data = await api.getPoolBalance(walletAddress);
      setAccountBalance(data.balance);
      setError('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load account balance');
    } finally {
      setLoading(false);
    }
  };

  const loadPaymentHistory = async () => {
    if (!walletAddress) return;

    try {
      setLoading(true);
      // Load from localStorage for demo
      const history = localStorage.getItem(`payment_history_${walletAddress}`);
      if (history) {
        setPaymentHistory(JSON.parse(history));
      }
    } catch (err) {
      console.error('Failed to load payment history');
    } finally {
      setLoading(false);
    }
  };

  const saveToHistory = (item: PaymentHistoryItem) => {
    if (!walletAddress) return;

    const updatedHistory = [item, ...paymentHistory].slice(0, 10);
    setPaymentHistory(updatedHistory);
    localStorage.setItem(`payment_history_${walletAddress}`, JSON.stringify(updatedHistory));
  };

  const handleDeposit = async () => {
    if (!walletAddress || !amount) return;

    try {
      setLoading(true);
      setMessage('');
      setError('');

      const lamports = Math.floor(parseFloat(amount) * 1e9);
      const response = await api.depositToPayment(walletAddress, lamports);

      setMessage('Deposit transaction prepared! Please sign and submit the transaction. ' + response.message);

      const historyItem: PaymentHistoryItem = {
        id: `deposit_${Date.now()}`,
        type: 'deposit',
        amount: lamports,
        status: 'pending',
        timestamp: Date.now(),
        hash: response.recent_blockhash
      };
      saveToHistory(historyItem);

      setAmount('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Deposit failed');
    } finally {
      setLoading(false);
    }
  };

  const handleWithdraw = async () => {
    if (!walletAddress || !amount) return;

    try {
      setLoading(true);
      setMessage('');
      setError('');

      const lamports = Math.floor(parseFloat(amount) * 1e9);
      const response = await api.withdrawFromPayment(walletAddress, lamports);

      setMessage('Withdrawal transaction prepared! Please sign and submit the transaction. ' + response.message);

      const historyItem: PaymentHistoryItem = {
        id: `withdraw_${Date.now()}`,
        type: 'withdraw',
        amount: lamports,
        status: 'pending',
        timestamp: Date.now(),
        hash: response.recent_blockhash
      };
      saveToHistory(historyItem);

      setAmount('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Withdrawal failed');
    } finally {
      setLoading(false);
    }
  };

  const handlePreparePayment = async () => {
    if (!receiverCommitment || !paymentAmount) return;

    try {
      setLoading(true);
      setMessage('');
      setError('');

      const lamports = Math.floor(parseFloat(paymentAmount) * 1e9);
      const response = await api.preparePayment(receiverCommitment, lamports, tokenMint || undefined);

      setPaymentHash(response.payment_hash);
      setPreparedCommitment(response.commitment);
      setMessage('Payment prepared successfully! Commitment: ' + response.commitment.substring(0, 16) + '...');

      const historyItem: PaymentHistoryItem = {
        id: `prepare_${Date.now()}`,
        type: 'payment',
        amount: lamports,
        status: 'pending',
        timestamp: Date.now(),
        hash: response.payment_hash
      };
      saveToHistory(historyItem);

      // Clear form
      setReceiverCommitment('');
      setPaymentAmount('');
      setTokenMint('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Payment preparation failed');
    } finally {
      setLoading(false);
    }
  };

  const handleAuthorizePayment = async () => {
    if (!authorizingCommitment || !nullifier || !authAmount || !merchant) return;

    try {
      setLoading(true);
      setMessage('');
      setError('');

      const lamports = Math.floor(parseFloat(authAmount) * 1e9);
      const response = await api.authorizePayment(
        authorizingCommitment,
        nullifier,
        lamports,
        merchant
      );

      setAccessToken(response.access_token);
      setMessage('Payment authorized! Access token received. Expires in ' + response.expires_in + ' seconds.');

      // Clear form
      setAuthorizingCommitment('');
      setNullifier('');
      setAuthAmount('');
      setMerchant('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Payment authorization failed');
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyAccess = async () => {
    if (!verifyToken) return;

    try {
      setLoading(true);
      setMessage('');
      setError('');

      const response = await api.verifyPaymentAccess(verifyToken);
      setVerifyResult(response);

      if (response.valid) {
        setMessage('Access token is valid!');
      } else {
        setMessage('Access token is invalid or expired.');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Access verification failed');
    } finally {
      setLoading(false);
    }
  };

  const handleSettlePayment = async () => {
    if (!paymentHeader || !resource) return;

    try {
      setLoading(true);
      setMessage('');
      setError('');

      const response = await api.settlePayment(
        1, // x402Version
        paymentHeader,
        resource,
        {
          scheme: 'zkproof',
          network: 'solana-mainnet',
          maxAmountRequired: '1.0',
          resource: resource,
          description: 'ZK Payment Settlement',
          mimeType: 'application/json',
          payTo: merchant || 'Merchant',
          maxTimeoutSeconds: 300
        }
      );

      setSettlementResult(response);

      if (response.success) {
        setMessage('Payment settled successfully! Transaction: ' + response.tx_sig);
      } else {
        setMessage('Payment settlement failed: ' + response.message);
      }

      // Clear form
      setPaymentHeader('');
      setResource('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Payment settlement failed');
    } finally {
      setLoading(false);
    }
  };

  if (!walletAddress) {
    return (
      <Card className="border-zinc-800 bg-zinc-950">
        <CardContent className="flex items-center justify-center py-10">
          <p className="text-zinc-400">Connect your wallet to access ZK payments</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="border-zinc-800 bg-zinc-950">
      <CardHeader>
        <div className="flex items-center gap-2">
          <Shield className="h-5 w-5 text-purple-500" />
          <CardTitle>ZK Payment Manager</CardTitle>
        </div>
        <CardDescription>
          Send and receive private payments using zero-knowledge proofs
        </CardDescription>
      </CardHeader>
      <CardContent>
        {/* Account Balance */}
        <div className="mb-6 rounded-lg border border-zinc-800 bg-zinc-900/50 p-4">
          <div className="flex items-center justify-between">
            <span className="text-sm text-zinc-400">Payment Account Balance</span>
            <div className="text-right">
              {accountBalance !== null ? (
                <>
                  <div className="text-xl font-bold text-purple-400">
                    {(accountBalance / 1e9).toFixed(4)} SOL
                  </div>
                  <div className="text-xs text-zinc-500">
                    {accountBalance.toLocaleString()} lamports
                  </div>
                </>
              ) : (
                <span className="text-sm text-zinc-500">Loading...</span>
              )}
            </div>
          </div>
        </div>

        {/* Main Tabs */}
        <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
          <TabsList className="bg-zinc-900">
            <TabsTrigger value="account">Account</TabsTrigger>
            <TabsTrigger value="prepare">Prepare Payment</TabsTrigger>
            <TabsTrigger value="authorize">Authorize</TabsTrigger>
            <TabsTrigger value="settle">Settle</TabsTrigger>
            <TabsTrigger value="verify">Verify</TabsTrigger>
          </TabsList>

          {/* Account Management Tab */}
          <TabsContent value="account" className="space-y-4">
            <div className="space-y-4">
              <div>
                <label className="mb-2 block text-sm font-medium text-zinc-200">
                  Amount (SOL)
                </label>
                <Input
                  type="number"
                  step="0.001"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  placeholder="0.001"
                  className="font-mono"
                />
              </div>

              <div className="grid gap-3 sm:grid-cols-2">
                <Button
                  onClick={handleDeposit}
                  disabled={loading || !amount}
                  className="bg-purple-600 hover:bg-purple-700"
                >
                  <ArrowDown className="mr-2 h-4 w-4" />
                  {loading ? 'Processing...' : 'Deposit'}
                </Button>
                <Button
                  onClick={handleWithdraw}
                  disabled={loading || !amount}
                  variant="outline"
                  className="border-zinc-700"
                >
                  <ArrowUp className="mr-2 h-4 w-4" />
                  {loading ? 'Processing...' : 'Withdraw'}
                </Button>
              </div>
            </div>
          </TabsContent>

          {/* Prepare Payment Tab */}
          <TabsContent value="prepare" className="space-y-4">
            <div className="space-y-4">
              <div>
                <label className="mb-2 block text-sm font-medium text-zinc-200">
                  Receiver Commitment
                </label>
                <Input
                  value={receiverCommitment}
                  onChange={(e) => setReceiverCommitment(e.target.value)}
                  placeholder="Enter receiver's zk commitment"
                  className="font-mono"
                />
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-zinc-200">
                  Amount (SOL)
                </label>
                <Input
                  type="number"
                  step="0.001"
                  value={paymentAmount}
                  onChange={(e) => setPaymentAmount(e.target.value)}
                  placeholder="0.001"
                  className="font-mono"
                />
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-zinc-200">
                  Token Mint (Optional)
                </label>
                <Input
                  value={tokenMint}
                  onChange={(e) => setTokenMint(e.target.value)}
                  placeholder="Leave empty for SOL"
                  className="font-mono"
                />
              </div>

              <Button
                onClick={handlePreparePayment}
                disabled={loading || !receiverCommitment || !paymentAmount}
                className="w-full bg-purple-600 hover:bg-purple-700"
              >
                <Send className="mr-2 h-4 w-4" />
                {loading ? 'Preparing...' : 'Prepare Payment'}
              </Button>

              {paymentHash && (
                <div className="rounded-lg border border-green-500/20 bg-green-500/10 p-3">
                  <p className="text-sm text-green-400">
                    <strong>Payment Hash:</strong> {paymentHash}
                  </p>
                  <p className="text-sm text-green-400">
                    <strong>Commitment:</strong> {preparedCommitment?.substring(0, 32)}...
                  </p>
                </div>
              )}
            </div>
          </TabsContent>

          {/* Authorize Payment Tab */}
          <TabsContent value="authorize" className="space-y-4">
            <div className="space-y-4">
              <div>
                <label className="mb-2 block text-sm font-medium text-zinc-200">
                  Your Commitment
                </label>
                <Input
                  value={authorizingCommitment}
                  onChange={(e) => setAuthorizingCommitment(e.target.value)}
                  placeholder="Enter your commitment"
                  className="font-mono"
                />
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-zinc-200">
                  Nullifier
                </label>
                <Input
                  value={nullifier}
                  onChange={(e) => setNullifier(e.target.value)}
                  placeholder="Enter nullifier"
                  className="font-mono"
                />
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-zinc-200">
                  Amount (SOL)
                </label>
                <Input
                  type="number"
                  step="0.001"
                  value={authAmount}
                  onChange={(e) => setAuthAmount(e.target.value)}
                  placeholder="0.001"
                  className="font-mono"
                />
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-zinc-200">
                  Merchant Address
                </label>
                <Input
                  value={merchant}
                  onChange={(e) => setMerchant(e.target.value)}
                  placeholder="Enter merchant wallet address"
                  className="font-mono"
                />
              </div>

              <Button
                onClick={handleAuthorizePayment}
                disabled={loading || !authorizingCommitment || !nullifier || !authAmount || !merchant}
                className="w-full bg-purple-600 hover:bg-purple-700"
              >
                <CheckCircle className="mr-2 h-4 w-4" />
                {loading ? 'Authorizing...' : 'Authorize Payment'}
              </Button>

              {accessToken && (
                <div className="rounded-lg border border-green-500/20 bg-green-500/10 p-3">
                  <p className="text-sm text-green-400">
                    <strong>Access Token:</strong> {accessToken.substring(0, 32)}...
                  </p>
                  <p className="text-xs text-green-400 mt-1">
                    Save this token securely!
                  </p>
                </div>
              )}
            </div>
          </TabsContent>

          {/* Settle Payment Tab */}
          <TabsContent value="settle" className="space-y-4">
            <div className="space-y-4">
              <div>
                <label className="mb-2 block text-sm font-medium text-zinc-200">
                  Payment Header (Base64)
                </label>
                <Input
                  value={paymentHeader}
                  onChange={(e) => setPaymentHeader(e.target.value)}
                  placeholder="Enter payment header"
                  className="font-mono"
                />
              </div>

              <div>
                <label className="mb-2 block text-sm font-medium text-zinc-200">
                  Resource
                </label>
                <Input
                  value={resource}
                  onChange={(e) => setResource(e.target.value)}
                  placeholder="Enter resource URL or identifier"
                  className="font-mono"
                />
              </div>

              <Button
                onClick={handleSettlePayment}
                disabled={loading || !paymentHeader || !resource}
                className="w-full bg-purple-600 hover:bg-purple-700"
              >
                <Wallet className="mr-2 h-4 w-4" />
                {loading ? 'Settling...' : 'Settle Payment'}
              </Button>

              {settlementResult && (
                <div className={`rounded-lg border p-3 ${
                  settlementResult.success
                    ? 'border-green-500/20 bg-green-500/10'
                    : 'border-red-500/20 bg-red-500/10'
                }`}>
                  <p className={`text-sm ${
                    settlementResult.success ? 'text-green-400' : 'text-red-400'
                  }`}>
                    {settlementResult.message}
                  </p>
                  {settlementResult.tx_sig && (
                    <p className="text-sm text-green-400 mt-1">
                      <strong>Transaction:</strong> {settlementResult.tx_sig}
                    </p>
                  )}
                </div>
              )}
            </div>
          </TabsContent>

          {/* Verify Access Tab */}
          <TabsContent value="verify" className="space-y-4">
            <div className="space-y-4">
              <div>
                <label className="mb-2 block text-sm font-medium text-zinc-200">
                  Access Token
                </label>
                <Input
                  value={verifyToken}
                  onChange={(e) => setVerifyToken(e.target.value)}
                  placeholder="Enter access token to verify"
                  className="font-mono"
                />
              </div>

              <Button
                onClick={handleVerifyAccess}
                disabled={loading || !verifyToken}
                className="w-full bg-purple-600 hover:bg-purple-700"
              >
                <AlertCircle className="mr-2 h-4 w-4" />
                {loading ? 'Verifying...' : 'Verify Access Token'}
              </Button>

              {verifyResult && (
                <div className={`rounded-lg border p-3 ${
                  verifyResult.valid
                    ? 'border-green-500/20 bg-green-500/10'
                    : 'border-red-500/20 bg-red-500/10'
                }`}>
                  <p className={`text-sm font-medium ${
                    verifyResult.valid ? 'text-green-400' : 'text-red-400'
                  }`}>
                    Status: {verifyResult.valid ? 'VALID' : 'INVALID'}
                  </p>
                  {verifyResult.valid && (
                    <div className="mt-2 space-y-1 text-sm text-green-400">
                      <p><strong>Commitment:</strong> {verifyResult.commitment?.substring(0, 16)}...</p>
                      <p><strong>Merchant:</strong> {verifyResult.merchant?.substring(0, 16)}...</p>
                      <p><strong>Amount:</strong> {verifyResult.amount ? (verifyResult.amount / 1e9).toFixed(4) : '0'} SOL</p>
                      <p><strong>Expires:</strong> {verifyResult.expires_at}</p>
                    </div>
                  )}
                </div>
              )}
            </div>
          </TabsContent>
        </Tabs>

        {/* Messages */}
        {message && (
          <div className="mt-4 rounded-lg border border-green-500/20 bg-green-500/10 p-3 text-sm text-green-400">
            {message}
          </div>
        )}

        {error && (
          <div className="mt-4 rounded-lg border border-red-500/20 bg-red-500/10 p-3 text-sm text-red-400">
            {error}
          </div>
        )}

        {/* Payment History */}
        {paymentHistory.length > 0 && (
          <div className="mt-6">
            <h3 className="mb-3 text-sm font-medium text-zinc-200">Recent Activity</h3>
            <div className="space-y-2">
              {paymentHistory.map((item) => (
                <div key={item.id} className="flex items-center justify-between rounded-lg border border-zinc-800 bg-zinc-900/30 p-3">
                  <div className="flex items-center gap-3">
                    {item.type === 'deposit' && <ArrowDown className="h-4 w-4 text-green-400" />}
                    {item.type === 'withdraw' && <ArrowUp className="h-4 w-4 text-red-400" />}
                    {item.type === 'payment' && <Send className="h-4 w-4 text-purple-400" />}
                    <div>
                      <p className="text-sm font-medium text-zinc-200">
                        {item.type.charAt(0).toUpperCase() + item.type.slice(1)}
                      </p>
                      <p className="text-xs text-zinc-500">
                        {(item.amount / 1e9).toFixed(4)} SOL
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <Badge
                      variant={item.status === 'completed' ? 'default' : 'secondary'}
                      className={`text-xs ${
                        item.status === 'completed'
                          ? 'bg-green-600 hover:bg-green-700'
                          : 'bg-zinc-700'
                      }`}
                    >
                      {item.status}
                    </Badge>
                    <p className="mt-1 text-xs text-zinc-500">
                      {new Date(item.timestamp).toLocaleTimeString()}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
