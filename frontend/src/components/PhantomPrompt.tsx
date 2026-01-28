import { ExternalLink } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';

export function PhantomPrompt() {
  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md border-zinc-800 bg-zinc-950">
        <CardHeader>
          <CardTitle className="text-2xl">Phantom Wallet Required</CardTitle>
          <CardDescription>
            Install Phantom to connect to ShadowPay and access privacy-focused payments
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-4">
            <p className="text-sm text-zinc-400">
              Phantom is a secure wallet for Solana that enables you to manage your assets
              and connect to dapps with privacy and security.
            </p>
          </div>

          <Button
            className="w-full bg-purple-600 hover:bg-purple-700"
            size="lg"
            onClick={() => window.open('https://phantom.app/', '_blank')}
          >
            <ExternalLink className="mr-2 h-4 w-4" />
            Install Phantom Wallet
          </Button>

          <div className="space-y-2 text-xs text-zinc-500">
            <p>After installing:</p>
            <ol className="list-inside list-decimal space-y-1 pl-2">
              <li>Create or import your wallet</li>
              <li>Refresh this page</li>
              <li>Click "Connect" to get started</li>
            </ol>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
