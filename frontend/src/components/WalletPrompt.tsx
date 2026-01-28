import { ExternalLink } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';

export function WalletPrompt() {
  const wallets = [
    {
      name: 'Phantom',
      description: 'The most popular Solana wallet with millions of users',
      url: 'https://phantom.app/',
      icon: '👻',
    },
    {
      name: 'Umbra',
      description: 'Privacy-focused Solana wallet for anonymous transactions',
      url: 'https://umbra.cash/',
      icon: '🌑',
    },
  ];

  return (
    <div className="flex min-h-screen items-center justify-center p-4 bg-black">
      <Card className="w-full max-w-2xl border-zinc-800 bg-zinc-950">
        <CardHeader>
          <CardTitle className="text-2xl">Solana Wallet Required</CardTitle>
          <CardDescription>
            Install a Solana wallet to connect to ShadowPay and access privacy-focused payments
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-4">
            <p className="text-sm text-zinc-400">
              Choose a wallet that supports Solana. Both Phantom and Umbra work great with ShadowPay's
              zero-knowledge payment system.
            </p>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            {wallets.map((wallet) => (
              <div
                key={wallet.name}
                className="rounded-lg border border-zinc-800 bg-zinc-900/30 p-5 space-y-4"
              >
                <div className="flex items-center gap-3">
                  <span className="text-3xl">{wallet.icon}</span>
                  <div>
                    <h3 className="font-semibold text-white">{wallet.name}</h3>
                    <p className="text-xs text-zinc-500">Wallet</p>
                  </div>
                </div>

                <p className="text-sm text-zinc-400 leading-relaxed">
                  {wallet.description}
                </p>

                <Button
                  className="w-full bg-purple-600 hover:bg-purple-700"
                  onClick={() => window.open(wallet.url, '_blank')}
                >
                  <ExternalLink className="mr-2 h-4 w-4" />
                  Install {wallet.name}
                </Button>
              </div>
            ))}
          </div>

          <div className="space-y-2 rounded-lg border border-zinc-800 bg-zinc-900/30 p-4 text-xs text-zinc-500">
            <p className="font-medium text-zinc-400">After installing:</p>
            <ol className="list-inside list-decimal space-y-1 pl-2">
              <li>Create or import your wallet</li>
              <li>Secure your seed phrase</li>
              <li>Refresh this page</li>
              <li>Click your wallet to connect</li>
            </ol>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
