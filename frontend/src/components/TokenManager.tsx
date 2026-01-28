import { useState, useEffect } from 'react';
import { Coins, RefreshCw, CheckCircle, XCircle } from 'lucide-react';
import api from '../services/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';

interface Token {
  mint: string;
  symbol: string;
  decimals: number;
  enabled: boolean;
}

export function TokenManager() {
  const [tokens, setTokens] = useState<Token[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadTokens();
  }, []);

  const loadTokens = async () => {
    try {
      setLoading(true);
      const data = await api.listTokens();
      setTokens(data.tokens);
      setError('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load tokens');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="border-zinc-800 bg-zinc-950">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Coins className="h-5 w-5 text-purple-500" />
            <CardTitle>Supported Tokens</CardTitle>
          </div>
          <Button
            onClick={loadTokens}
            disabled={loading}
            variant="ghost"
            size="sm"
          >
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          </Button>
        </div>
        <CardDescription>
          SPL tokens supported by ShadowPay for private transactions
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {error && (
          <div className="rounded-lg border border-red-500/20 bg-red-500/10 p-3 text-sm text-red-400">
            {error}
          </div>
        )}

        {tokens.length === 0 && !loading && (
          <div className="flex flex-col items-center justify-center py-10 text-center">
            <Coins className="mb-3 h-12 w-12 text-zinc-600" />
            <p className="text-sm text-zinc-400">No tokens configured</p>
          </div>
        )}

        <div className="space-y-3">
          {tokens.map((token) => (
            <div
              key={token.mint}
              className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-4 transition-colors hover:bg-zinc-900"
            >
              <div className="flex items-start justify-between">
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <span className="text-lg font-semibold text-purple-400">{token.symbol}</span>
                    {token.enabled ? (
                      <CheckCircle className="h-4 w-4 text-green-500" />
                    ) : (
                      <XCircle className="h-4 w-4 text-red-500" />
                    )}
                  </div>
                  <p className="font-mono text-xs text-zinc-500 break-all">{token.mint}</p>
                  <p className="text-xs text-zinc-400">Decimals: {token.decimals}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
