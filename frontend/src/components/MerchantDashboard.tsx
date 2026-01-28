import { useState, useEffect } from 'react';
import { useWalletConnection } from '@solana/react-hooks';
import api from '../services/api';

interface Earnings {
  total_earnings: number;
  withdrawable_sol: number;
  pending_settlement: number;
  total_usd_value: string;
  token_breakdown: Array<{
    symbol: string;
    amount: number;
  }>;
}

export function MerchantDashboard() {
  const { wallet } = useWalletConnection();
  const [earnings, setEarnings] = useState<Earnings | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const walletAddress = wallet?.account.address.toString();

  useEffect(() => {
    if (walletAddress) {
      loadEarnings();
    }
  }, [walletAddress]);

  const loadEarnings = async () => {
    try {
      setLoading(true);
      const data = await api.getMerchantEarnings();
      setEarnings(data);
      setError('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load earnings');
    } finally {
      setLoading(false);
    }
  };

  if (!walletAddress) {
    return (
      <div className="rounded-2xl border border-border-low bg-card p-6">
        <p className="text-muted">Connect your wallet to view merchant earnings</p>
      </div>
    );
  }

  return (
    <div className="space-y-4 rounded-2xl border border-border-low bg-card p-6 shadow-[0_20px_80px_-50px_rgba(0,0,0,0.35)]">
      <div className="flex items-start justify-between">
        <div className="space-y-1">
          <h2 className="text-xl font-semibold">Merchant Earnings</h2>
          <p className="text-sm text-muted">
            Track your earnings from ShadowPay transactions
          </p>
        </div>
        <button
          onClick={loadEarnings}
          disabled={loading}
          className="rounded-lg border border-border-low bg-card px-3 py-1.5 text-sm font-medium transition hover:-translate-y-0.5 hover:shadow-sm disabled:cursor-not-allowed disabled:opacity-60"
        >
          {loading ? 'Loading...' : 'Refresh'}
        </button>
      </div>

      {error && (
        <div className="rounded-lg border border-red-500/20 bg-red-500/10 p-3 text-sm text-red-600">
          {error}
        </div>
      )}

      {earnings && (
        <div className="space-y-3">
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="rounded-lg border border-border-low bg-cream p-4">
              <div className="text-sm text-muted">Total Earnings</div>
              <div className="mt-1 text-2xl font-semibold">
                {(earnings.total_earnings / 1e9).toFixed(4)} SOL
              </div>
              <div className="text-xs text-muted">${earnings.total_usd_value}</div>
            </div>

            <div className="rounded-lg border border-border-low bg-cream p-4">
              <div className="text-sm text-muted">Withdrawable</div>
              <div className="mt-1 text-2xl font-semibold">
                {(earnings.withdrawable_sol / 1e9).toFixed(4)} SOL
              </div>
            </div>

            <div className="rounded-lg border border-border-low bg-cream p-4">
              <div className="text-sm text-muted">Pending Settlement</div>
              <div className="mt-1 text-2xl font-semibold">
                {(earnings.pending_settlement / 1e9).toFixed(4)} SOL
              </div>
            </div>
          </div>

          {earnings.token_breakdown.length > 0 && (
            <div className="rounded-lg border border-border-low p-4">
              <div className="mb-3 text-sm font-medium">Token Breakdown</div>
              <div className="space-y-2">
                {earnings.token_breakdown.map((token) => (
                  <div
                    key={token.symbol}
                    className="flex items-center justify-between text-sm"
                  >
                    <span className="font-medium">{token.symbol}</span>
                    <span className="font-mono">
                      {(token.amount / 1e9).toFixed(4)} SOL
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
