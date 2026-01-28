import { useState, useEffect } from 'react';
import { useWalletConnection } from '@solana/react-hooks';
import api from '../services/api';

export function PoolManager() {
  const { wallet } = useWalletConnection();
  const [balance, setBalance] = useState<number | null>(null);
  const [minDeposit, setMinDeposit] = useState<number | null>(null);
  const [amount, setAmount] = useState('');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const walletAddress = wallet?.account.address.toString();

  useEffect(() => {
    if (walletAddress) {
      loadBalance();
    }
  }, [walletAddress]);

  const loadBalance = async () => {
    if (!walletAddress) return;

    try {
      setLoading(true);
      const data = await api.getPoolBalance(walletAddress);
      setBalance(data.balance);
      setMinDeposit(data.min_deposit);
      setError('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load balance');
    } finally {
      setLoading(false);
    }
  };

  const handleDeposit = async () => {
    if (!walletAddress || !amount) return;

    try {
      setLoading(true);
      setMessage('');
      setError('');

      const lamports = Math.floor(parseFloat(amount) * 1e9);
      const response = await api.depositToPool(walletAddress, lamports);

      setMessage('Deposit transaction created! ' + response.message);
      setAmount('');

      // Reload balance after deposit
      setTimeout(loadBalance, 2000);
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
      const response = await api.withdrawFromPool(walletAddress, lamports);

      const netSol = response.net_amount / 1e9;
      const feeSol = response.fee / 1e9;
      setMessage('Withdrawal created! Net: ' + netSol.toFixed(4) + ' SOL (Fee: ' + feeSol.toFixed(4) + ' SOL)');
      setAmount('');

      // Reload balance after withdrawal
      setTimeout(loadBalance, 2000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Withdrawal failed');
    } finally {
      setLoading(false);
    }
  };

  if (!walletAddress) {
    return (
      <div className="rounded-2xl border border-border-low bg-card p-6">
        <p className="text-muted">Connect your wallet to manage privacy pool</p>
      </div>
    );
  }

  return (
    <div className="space-y-4 rounded-2xl border border-border-low bg-card p-6 shadow-[0_20px_80px_-50px_rgba(0,0,0,0.35)]">
      <div className="space-y-1">
        <h2 className="text-xl font-semibold">Privacy Pool</h2>
        <p className="text-sm text-muted">
          Mix your funds with others for maximum anonymity on-chain
        </p>
      </div>

      <div className="rounded-lg border border-border-low bg-cream p-4">
        <div className="flex items-baseline justify-between">
          <span className="text-sm text-muted">Pool Balance</span>
          <div className="text-right">
            {balance !== null ? (
              <>
                <div className="text-lg font-semibold">{(balance / 1e9).toFixed(4)} SOL</div>
                <div className="text-xs text-muted">{balance} lamports</div>
              </>
            ) : (
              <span className="text-sm text-muted">Loading...</span>
            )}
          </div>
        </div>
        {minDeposit !== null && (
          <div className="mt-2 text-xs text-muted">
            Min deposit: {(minDeposit / 1e9).toFixed(4)} SOL
          </div>
        )}
      </div>

      <div className="space-y-3">
        <div>
          <label className="mb-2 block text-sm font-medium">
            Amount (SOL)
          </label>
          <input
            type="number"
            step="0.001"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            placeholder="0.001"
            className="w-full rounded-lg border border-border-low bg-card px-4 py-2 font-mono text-sm focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
          />
        </div>

        <div className="grid gap-2 sm:grid-cols-2">
          <button
            onClick={handleDeposit}
            disabled={loading || !amount}
            className="rounded-lg border border-border-low bg-primary px-4 py-2 font-medium transition hover:-translate-y-0.5 hover:shadow-sm disabled:cursor-not-allowed disabled:opacity-60"
          >
            {loading ? 'Processing...' : 'Deposit'}
          </button>
          <button
            onClick={handleWithdraw}
            disabled={loading || !amount}
            className="rounded-lg border border-border-low bg-card px-4 py-2 font-medium transition hover:-translate-y-0.5 hover:shadow-sm disabled:cursor-not-allowed disabled:opacity-60"
          >
            {loading ? 'Processing...' : 'Withdraw'}
          </button>
        </div>

        <button
          onClick={loadBalance}
          disabled={loading}
          className="w-full rounded-lg border border-border-low bg-card px-4 py-2 text-sm font-medium transition hover:-translate-y-0.5 hover:shadow-sm disabled:cursor-not-allowed disabled:opacity-60"
        >
          Refresh Balance
        </button>
      </div>

      {message && (
        <div className="rounded-lg border border-green-500/20 bg-green-500/10 p-3 text-sm text-green-600">
          {message}
        </div>
      )}

      {error && (
        <div className="rounded-lg border border-red-500/20 bg-red-500/10 p-3 text-sm text-red-600">
          {error}
        </div>
      )}
    </div>
  );
}
