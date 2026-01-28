import { useState, useEffect } from 'react';
import { useWalletConnection } from '@solana/react-hooks';
import { DollarSign, TrendingUp, Clock, RefreshCw } from 'lucide-react';
import api from '../services/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';

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
      <Card className="border-zinc-800 bg-zinc-950">
        <CardContent className="flex items-center justify-center py-10">
          <p className="text-zinc-400">Connect your wallet to view merchant earnings</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="border-zinc-800 bg-zinc-950">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <TrendingUp className="h-5 w-5 text-purple-500" />
            <CardTitle>Merchant Earnings</CardTitle>
          </div>
          <Button
            onClick={loadEarnings}
            disabled={loading}
            variant="ghost"
            size="sm"
          >
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          </Button>
        </div>
        <CardDescription>
          Track your earnings from ShadowPay transactions
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {error && (
          <div className="rounded-lg border border-red-500/20 bg-red-500/10 p-3 text-sm text-red-400">
            {error}
          </div>
        )}

        {earnings && (
          <>
            <div className="grid gap-4 sm:grid-cols-3">
              <div className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-4">
                <div className="flex items-center gap-2 text-sm text-zinc-400">
                  <DollarSign className="h-4 w-4" />
                  Total Earnings
                </div>
                <div className="mt-2">
                  <div className="text-2xl font-bold text-purple-400">
                    {(earnings.total_earnings / 1e9).toFixed(4)} SOL
                  </div>
                  <div className="text-xs text-zinc-500">${earnings.total_usd_value}</div>
                </div>
              </div>

              <div className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-4">
                <div className="flex items-center gap-2 text-sm text-zinc-400">
                  <TrendingUp className="h-4 w-4" />
                  Withdrawable
                </div>
                <div className="mt-2">
                  <div className="text-2xl font-bold text-green-400">
                    {(earnings.withdrawable_sol / 1e9).toFixed(4)} SOL
                  </div>
                </div>
              </div>

              <div className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-4">
                <div className="flex items-center gap-2 text-sm text-zinc-400">
                  <Clock className="h-4 w-4" />
                  Pending
                </div>
                <div className="mt-2">
                  <div className="text-2xl font-bold text-yellow-400">
                    {(earnings.pending_settlement / 1e9).toFixed(4)} SOL
                  </div>
                </div>
              </div>
            </div>

            {earnings.token_breakdown.length > 0 && (
              <div className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-4">
                <div className="mb-4 text-sm font-medium text-zinc-200">Token Breakdown</div>
                <div className="space-y-3">
                  {earnings.token_breakdown.map((token) => (
                    <div
                      key={token.symbol}
                      className="flex items-center justify-between rounded-md border border-zinc-800 bg-zinc-900/30 p-3"
                    >
                      <span className="font-medium text-purple-400">{token.symbol}</span>
                      <span className="font-mono text-sm text-zinc-300">
                        {(token.amount / 1e9).toFixed(4)} SOL
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </>
        )}
      </CardContent>
    </Card>
  );
}
