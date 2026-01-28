import { useWalletConnection } from "@solana/react-hooks";
import { Shield, Wallet as WalletIcon, LogOut } from 'lucide-react';
import { PhantomPrompt } from "./components/PhantomPrompt";
import { PoolManager } from "./components/PoolManager";
import { TokenManager } from "./components/TokenManager";
import { MerchantDashboard } from "./components/MerchantDashboard";
import { Card, CardContent } from "./components/ui/card";
import { Button } from "./components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "./components/ui/tabs";

export default function App() {
  const { connectors, connect, disconnect, wallet, status } =
    useWalletConnection();

  const address = wallet?.account.address.toString();

  // Check if Phantom is available
  const hasPhantom = connectors.some(c => c.name.toLowerCase().includes('phantom'));

  // Show install prompt if no Phantom
  if (!hasPhantom) {
    return <PhantomPrompt />;
  }

  return (
    <div className="min-h-screen bg-black">
      {/* Header */}
      <header className="border-b border-zinc-800 bg-zinc-950/50 backdrop-blur-xl">
        <div className="mx-auto max-w-6xl px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-purple-600">
                <Shield className="h-6 w-6 text-white" />
              </div>
              <div>
                <h1 className="text-xl font-bold text-white">ShadowPay</h1>
                <p className="text-xs text-zinc-400">Privacy-focused payments</p>
              </div>
            </div>

            {status === "connected" && address ? (
              <div className="flex items-center gap-3">
                <div className="hidden rounded-lg border border-zinc-800 bg-zinc-900/50 px-4 py-2 sm:block">
                  <p className="font-mono text-xs text-zinc-400">
                    {address.slice(0, 4)}...{address.slice(-4)}
                  </p>
                </div>
                <Button
                  onClick={() => disconnect()}
                  variant="outline"
                  size="sm"
                  className="border-zinc-800"
                >
                  <LogOut className="mr-2 h-4 w-4" />
                  Disconnect
                </Button>
              </div>
            ) : (
              <div className="flex items-center gap-2">
                {connectors.filter(c => c.name.toLowerCase().includes('phantom')).map((connector) => (
                  <Button
                    key={connector.id}
                    onClick={() => connect(connector.id)}
                    disabled={status === "connecting"}
                    className="bg-purple-600 hover:bg-purple-700"
                  >
                    <WalletIcon className="mr-2 h-4 w-4" />
                    {status === "connecting" ? "Connecting..." : "Connect Phantom"}
                  </Button>
                ))}
              </div>
            )}
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="mx-auto max-w-6xl px-4 py-8">
        {status !== "connected" ? (
          <Card className="border-zinc-800 bg-zinc-950">
            <CardContent className="flex flex-col items-center justify-center py-16 text-center">
              <div className="mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-purple-600/10">
                <WalletIcon className="h-8 w-8 text-purple-500" />
              </div>
              <h2 className="mb-2 text-2xl font-bold text-white">Connect Your Wallet</h2>
              <p className="mb-6 max-w-md text-zinc-400">
                Connect your Phantom wallet to access privacy pools, manage tokens,
                and track merchant earnings with zero-knowledge proofs.
              </p>
              <div className="grid gap-3 sm:grid-cols-2">
                <div className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-4 text-left">
                  <Shield className="mb-2 h-5 w-5 text-purple-500" />
                  <p className="text-sm font-medium text-white">Private Transactions</p>
                  <p className="text-xs text-zinc-500">ZK proofs for anonymity</p>
                </div>
                <div className="rounded-lg border border-zinc-800 bg-zinc-900/50 p-4 text-left">
                  <WalletIcon className="mb-2 h-5 w-5 text-purple-500" />
                  <p className="text-sm font-medium text-white">Secure Storage</p>
                  <p className="text-xs text-zinc-500">Non-custodial privacy pool</p>
                </div>
              </div>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-6">
            {/* Welcome Section */}
            <div className="rounded-xl border border-zinc-800 bg-gradient-to-br from-purple-950/50 to-zinc-950 p-6">
              <h2 className="mb-2 text-2xl font-bold text-white">
                Welcome back! 👋
              </h2>
              <p className="text-zinc-400">
                Manage your privacy-focused payments on Solana
              </p>
            </div>

            {/* Tabs for organized content */}
            <Tabs defaultValue="pool" className="space-y-4">
              <TabsList className="bg-zinc-900">
                <TabsTrigger value="pool">Privacy Pool</TabsTrigger>
                <TabsTrigger value="merchant">Merchant</TabsTrigger>
                <TabsTrigger value="tokens">Tokens</TabsTrigger>
              </TabsList>

              <TabsContent value="pool" className="space-y-4">
                <PoolManager />
              </TabsContent>

              <TabsContent value="merchant" className="space-y-4">
                <MerchantDashboard />
              </TabsContent>

              <TabsContent value="tokens" className="space-y-4">
                <TokenManager />
              </TabsContent>
            </Tabs>
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="mt-16 border-t border-zinc-800 bg-zinc-950/50 py-8">
        <div className="mx-auto max-w-6xl px-4 text-center text-xs text-zinc-500">
          <p>Powered by ShadowPay API • Zero-Knowledge Payments on Solana</p>
        </div>
      </footer>
    </div>
  );
}
