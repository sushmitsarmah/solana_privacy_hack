import { useState, useEffect } from 'react';
import { useWalletConnection } from '@solana/react-hooks';
import { Wallet, ArrowDown, ArrowUp, RefreshCw } from 'lucide-react';
import api from '../services/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { Input } from './ui/input';

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

      setTimeout(loadBalance, 2000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Withdrawal failed');
    } finally {
      setLoading(false);
    }
  };

  if (!walletAddress) {
    return (
      <Card className="border-zinc-800 bg-zinc-950">
        <CardContent className="flex items-center justify-center py-10">
          <p className="text-zinc-400">Connect your wallet to manage privacy pool</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="border-zinc-800 bg-zinc-950">
      <CardHeader>
        <div className="flex items-center gap-2">
          <Wallet className="h-5 w-5 text-purple-500" />
          <CardTitle>Privacy Pool</CardTitle>
        </div>
        <CardDescription>
          Mix your funds with others for maximum anonymity on-chain
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-4">
          <div className="flex items-center justify-between">
            <span className="text-sm text-zinc-400">Pool Balance</span>
            <div className="text-right">
              {balance !== null ? (
                <>
                  <div className="text-xl font-bold text-purple-400">{(balance / 1e9).toFixed(4)} SOL</div>
                  <div className="text-xs text-zinc-500">{balance.toLocaleString()} lamports</div>
                </>
              ) : (
                <span className="text-sm text-zinc-500">Loading...</span>
              )}
            </div>
          </div>
          {minDeposit !== null && (
            <div className="mt-3 text-xs text-zinc-500">
              Minimum deposit: {(minDeposit / 1e9).toFixed(4)} SOL
            </div>
          )}
        </div>

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

          <Button
            onClick={loadBalance}
            disabled={loading}
            variant="ghost"
            className="w-full"
          >
            <RefreshCw className="mr-2 h-4 w-4" />
            Refresh Balance
          </Button>
        </div>

        {message && (
          <div className="rounded-lg border border-green-500/20 bg-green-500/10 p-3 text-sm text-green-400">
            {message}
          </div>
        )}

        {error && (
          <div className="rounded-lg border border-red-500/20 bg-red-500/10 p-3 text-sm text-red-400">
            {error}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
