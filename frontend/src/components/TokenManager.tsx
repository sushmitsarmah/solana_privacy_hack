import { useState, useEffect } from 'react';
import api from '../services/api';

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
    <div className="space-y-4 rounded-2xl border border-border-low bg-card p-6 shadow-[0_20px_80px_-50px_rgba(0,0,0,0.35)]">
      <div className="flex items-start justify-between">
        <div className="space-y-1">
          <h2 className="text-xl font-semibold">Supported Tokens</h2>
          <p className="text-sm text-muted">
            SPL tokens supported by ShadowPay
          </p>
        </div>
        <button
          onClick={loadTokens}
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

      <div className="space-y-2">
        {tokens.length === 0 && !loading && (
          <p className="py-8 text-center text-sm text-muted">
            No tokens configured
          </p>
        )}

        {tokens.map((token) => (
          <div
            key={token.mint}
            className="rounded-lg border border-border-low bg-cream p-4"
          >
            <div className="flex items-start justify-between">
              <div className="space-y-1">
                <div className="flex items-center gap-2">
                  <span className="font-semibold">{token.symbol}</span>
                  <span
                    className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                      token.enabled
                        ? 'bg-green-500/10 text-green-600'
                        : 'bg-red-500/10 text-red-600'
                    }`}
                  >
                    {token.enabled ? 'Enabled' : 'Disabled'}
                  </span>
                </div>
                <p className="font-mono text-xs text-muted">{token.mint}</p>
                <p className="text-xs text-muted">Decimals: {token.decimals}</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
