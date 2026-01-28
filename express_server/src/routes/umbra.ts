import { Router, Request, Response } from 'express';
import { Connection, Keypair, LAMPORTS_PER_SOL } from '@solana/web3.js';
// @ts-ignore - Using compiled Umbra SDK
const { UmbraClient, UmbraWallet, MXE_ARCIUM_X25519_PUBLIC_KEY, WSOL_MINT_ADDRESS } = require('../../umbra-sdk/index.cjs');

const router = Router();

/**
 * POST /api/umbra/register
 * Register an account for anonymity using Umbra
 *
 * Body:
 * {
 *   "privateKey": "base58 private key",
 *   "rpcUrl": "https://api.devnet.solana.com" (optional)
 * }
 */
router.post('/register', async (req: Request, res: Response) => {
  try {
    const { privateKey, rpcUrl = 'https://api.devnet.solana.com' } = req.body;

    if (!privateKey) {
      return res.status(400).json({
        error: 'Missing required field: privateKey',
      });
    }

    // Create connection
    const connection = new Connection(rpcUrl, 'confirmed');

    // Parse keypair
    let keypair: Keypair;
    try {
      const bs58 = require('bs58');
      const privateKeyBytes = bs58.decode(privateKey);
      keypair = Keypair.fromSecretKey(privateKeyBytes);
    } catch {
      return res.status(400).json({
        error: 'Invalid private key format. Expected base58.',
      });
    }

    // Create Umbra wallet
    const wallet = await UmbraWallet.fromSigner(
      { signer: { keypair } },
      { arciumX25519PublicKey: MXE_ARCIUM_X25519_PUBLIC_KEY }
    );

    // Create Umbra client
    const client = await UmbraClient.create({ connection });
    await client.setUmbraWallet(wallet);

    // Set ZK prover
    client.setZkProver('wasm', {
      masterViewingKeyRegistration: true,
      createSplDepositWithPublicAmount: true,
    });

    // Register account for anonymity
    const signature = await client.registerAccountForAnonymity(undefined, {
      mode: 'connection',
    });

    res.json({
      success: true,
      data: {
        signature,
        publicKey: keypair.publicKey.toBase58(),
        explorerUrl: `https://explorer.solana.com/tx/${signature}?cluster=devnet`,
      },
      message: 'Account registered for anonymity successfully!',
    });
  } catch (error) {
    res.status(500).json({
      error: 'Failed to register account',
      details: (error as Error).message,
    });
  }
});

/**
 * POST /api/umbra/wallet/create
 * Create a new Umbra wallet
 *
 * Body:
 * {
 *   "privateKey": "base58 private key" (optional, generates new if not provided)
 * }
 */
router.post('/wallet/create', async (req: Request, res: Response) => {
  try {
    const { privateKey } = req.body;

    let keypair: Keypair;
    if (privateKey) {
      const bs58 = require('bs58');
      const privateKeyBytes = bs58.decode(privateKey);
      keypair = Keypair.fromSecretKey(privateKeyBytes);
    } else {
      keypair = Keypair.generate();
    }

    // Create Umbra wallet
    const wallet = await UmbraWallet.fromSigner(
      { signer: { keypair } },
      { arciumX25519PublicKey: MXE_ARCIUM_X25519_PUBLIC_KEY }
    );

    const bs58 = require('bs58');
    res.json({
      success: true,
      data: {
        publicKey: keypair.publicKey.toBase58(),
        privateKey: bs58.encode(keypair.secretKey),
        walletCreated: true,
      },
      message: 'Umbra wallet created successfully!',
      security: {
        warning: 'Keep your private key secure and never share it!',
      },
    });
  } catch (error) {
    res.status(500).json({
      error: 'Failed to create wallet',
      details: (error as Error).message,
    });
  }
});

/**
 * POST /api/umbra/deposit
 * Deposit SOL into the privacy mixer pool
 *
 * Body:
 * {
 *   "privateKey": "base58 private key",
 *   "amount": 0.1 (in SOL),
 *   "destinationAddress": "recipient address for withdrawal" (optional, defaults to sender),
 *   "rpcUrl": "https://api.devnet.solana.com" (optional)
 * }
 */
router.post('/deposit', async (req: Request, res: Response) => {
  try {
    const { privateKey, amount, destinationAddress, rpcUrl = 'https://api.devnet.solana.com' } = req.body;

    if (!privateKey || !amount) {
      return res.status(400).json({
        error: 'Missing required fields: privateKey, amount',
      });
    }

    if (typeof amount !== 'number' || amount <= 0) {
      return res.status(400).json({
        error: 'Amount must be a positive number',
      });
    }

    // Create connection
    const connection = new Connection(rpcUrl, 'confirmed');

    // Parse keypair
    const bs58 = require('bs58');
    const privateKeyBytes = bs58.decode(privateKey);
    const keypair = Keypair.fromSecretKey(privateKeyBytes);

    // Use sender address as destination if not provided
    const destination = destinationAddress || keypair.publicKey.toBase58();

    // Create Umbra wallet
    const wallet = await UmbraWallet.fromSigner(
      { signer: { keypair } },
      { arciumX25519PublicKey: MXE_ARCIUM_X25519_PUBLIC_KEY }
    );

    // Create Umbra client
    const client = await UmbraClient.create({ connection });
    await client.setUmbraWallet(wallet);

    // Set ZK prover
    client.setZkProver('wasm', {
      masterViewingKeyRegistration: true,
      createSplDepositWithPublicAmount: true,
    });

    // Convert SOL to lamports
    const lamports = Math.floor(amount * LAMPORTS_PER_SOL);

    // Deposit into mixer pool
    const result = await client.depositPubliclyIntoMixerPoolSol(
      lamports,
      destination,
      { mode: 'connection' }
    );

    res.json({
      success: true,
      data: {
        signature: result.signature,
        amount,
        amountLamports: lamports,
        destinationAddress: destination,
        publicKey: keypair.publicKey.toBase58(),
        explorerUrl: `https://explorer.solana.com/tx/${result.signature}?cluster=devnet`,
      },
      message: 'Deposit to privacy pool successful!',
    });
  } catch (error) {
    res.status(500).json({
      error: 'Failed to process deposit',
      details: (error as Error).message,
    });
  }
});

/**
 * POST /api/umbra/balance
 * Get encrypted balance for a token
 *
 * Body:
 * {
 *   "privateKey": "base58 private key",
 *   "mint": "token mint address" (optional, defaults to WSOL),
 *   "rpcUrl": "https://api.devnet.solana.com" (optional)
 * }
 */
router.post('/balance', async (req: Request, res: Response) => {
  try {
    const { privateKey, mint, rpcUrl = 'https://api.devnet.solana.com' } = req.body;

    if (!privateKey) {
      return res.status(400).json({
        error: 'Missing required field: privateKey',
      });
    }

    // Create connection
    const connection = new Connection(rpcUrl, 'confirmed');

    // Parse keypair
    const bs58 = require('bs58');
    const privateKeyBytes = bs58.decode(privateKey);
    const keypair = Keypair.fromSecretKey(privateKeyBytes);

    // Use WSOL if no mint specified
    const tokenMint = mint || WSOL_MINT_ADDRESS.toBase58();

    // Create Umbra wallet
    const wallet = await UmbraWallet.fromSigner(
      { signer: { keypair } },
      { arciumX25519PublicKey: MXE_ARCIUM_X25519_PUBLIC_KEY }
    );

    // Create Umbra client
    const client = await UmbraClient.create({ connection });
    await client.setUmbraWallet(wallet);

    // Get encrypted balance
    const balance = await client.getEncryptedTokenBalance(tokenMint);

    res.json({
      success: true,
      data: {
        balance: balance.toString(),
        balanceSOL: Number(balance) / LAMPORTS_PER_SOL,
        mint: tokenMint,
        publicKey: keypair.publicKey.toBase58(),
      },
      message: 'Balance retrieved successfully',
    });
  } catch (error) {
    res.status(500).json({
      error: 'Failed to get balance',
      details: (error as Error).message,
    });
  }
});

/**
 * POST /api/umbra/register-confidentiality
 * Register account for confidentiality (encrypted balances)
 *
 * Body:
 * {
 *   "privateKey": "base58 private key",
 *   "rpcUrl": "https://api.devnet.solana.com" (optional)
 * }
 */
router.post('/register-confidentiality', async (req: Request, res: Response) => {
  try {
    const { privateKey, rpcUrl = 'https://api.devnet.solana.com' } = req.body;

    if (!privateKey) {
      return res.status(400).json({
        error: 'Missing required field: privateKey',
      });
    }

    // Create connection
    const connection = new Connection(rpcUrl, 'confirmed');

    // Parse keypair
    const bs58 = require('bs58');
    const privateKeyBytes = bs58.decode(privateKey);
    const keypair = Keypair.fromSecretKey(privateKeyBytes);

    // Create Umbra wallet
    const wallet = await UmbraWallet.fromSigner(
      { signer: { keypair } },
      { arciumX25519PublicKey: MXE_ARCIUM_X25519_PUBLIC_KEY }
    );

    // Create Umbra client
    const client = await UmbraClient.create({ connection });
    await client.setUmbraWallet(wallet);

    // Set ZK prover
    client.setZkProver('wasm', {
      masterViewingKeyRegistration: true,
      createSplDepositWithPublicAmount: true,
    });

    // Register account for confidentiality
    const signature = await client.registerAccountForConfidentiality(undefined, {
      mode: 'connection',
    });

    res.json({
      success: true,
      data: {
        signature,
        publicKey: keypair.publicKey.toBase58(),
        explorerUrl: `https://explorer.solana.com/tx/${signature}?cluster=devnet`,
      },
      message: 'Account registered for confidentiality successfully!',
    });
  } catch (error) {
    res.status(500).json({
      error: 'Failed to register for confidentiality',
      details: (error as Error).message,
    });
  }
});

/**
 * POST /api/umbra/register-full
 * Register account for both confidentiality and anonymity
 *
 * Body:
 * {
 *   "privateKey": "base58 private key",
 *   "rpcUrl": "https://api.devnet.solana.com" (optional)
 * }
 */
router.post('/register-full', async (req: Request, res: Response) => {
  try {
    const { privateKey, rpcUrl = 'https://api.devnet.solana.com' } = req.body;

    if (!privateKey) {
      return res.status(400).json({
        error: 'Missing required field: privateKey',
      });
    }

    // Create connection
    const connection = new Connection(rpcUrl, 'confirmed');

    // Parse keypair
    const bs58 = require('bs58');
    const privateKeyBytes = bs58.decode(privateKey);
    const keypair = Keypair.fromSecretKey(privateKeyBytes);

    // Create Umbra wallet
    const wallet = await UmbraWallet.fromSigner(
      { signer: { keypair } },
      { arciumX25519PublicKey: MXE_ARCIUM_X25519_PUBLIC_KEY }
    );

    // Create Umbra client
    const client = await UmbraClient.create({ connection });
    await client.setUmbraWallet(wallet);

    // Set ZK prover
    client.setZkProver('wasm', {
      masterViewingKeyRegistration: true,
      createSplDepositWithPublicAmount: true,
    });

    // Register account for both confidentiality and anonymity
    const signature = await client.registerAccountForConfidentialityAndAnonymity(
      undefined,
      undefined,
      { mode: 'connection' }
    );

    res.json({
      success: true,
      data: {
        signature,
        publicKey: keypair.publicKey.toBase58(),
        explorerUrl: `https://explorer.solana.com/tx/${signature}?cluster=devnet`,
      },
      message: 'Account registered for full privacy (confidentiality + anonymity) successfully!',
    });
  } catch (error) {
    res.status(500).json({
      error: 'Failed to register account',
      details: (error as Error).message,
    });
  }
});

/**
 * GET /api/umbra/info
 * Get information about Umbra SDK integration
 */
router.get('/info', (_req: Request, res: Response) => {
  res.json({
    name: 'Umbra SDK Integration',
    description: 'Privacy-focused payments on Solana using zero-knowledge proofs',
    version: '1.0.0',
    sdk: {
      name: '@umbra-defi/sdk',
      github: 'https://github.com/umbra-defi/sdk',
    },
    endpoints: {
      'POST /api/umbra/register': {
        description: 'Register account for anonymity (privacy mixer)',
        body: {
          privateKey: 'string (base58)',
          rpcUrl: 'string (optional)',
        },
      },
      'POST /api/umbra/register-confidentiality': {
        description: 'Register account for confidentiality (encrypted balances)',
        body: {
          privateKey: 'string (base58)',
          rpcUrl: 'string (optional)',
        },
      },
      'POST /api/umbra/register-full': {
        description: 'Register account for full privacy (confidentiality + anonymity)',
        body: {
          privateKey: 'string (base58)',
          rpcUrl: 'string (optional)',
        },
      },
      'POST /api/umbra/wallet/create': {
        description: 'Create new Umbra wallet',
        body: {
          privateKey: 'string (base58, optional)',
        },
      },
      'POST /api/umbra/deposit': {
        description: 'Deposit SOL into privacy mixer pool',
        body: {
          privateKey: 'string (base58)',
          amount: 'number (SOL)',
          destinationAddress: 'string (optional)',
          rpcUrl: 'string (optional)',
        },
      },
      'POST /api/umbra/balance': {
        description: 'Get encrypted token balance',
        body: {
          privateKey: 'string (base58)',
          mint: 'string (optional, defaults to WSOL)',
          rpcUrl: 'string (optional)',
        },
      },
    },
    features: [
      'Zero-knowledge proofs for privacy (Groth16 ZK-SNARKs)',
      'Encrypted balances using Arcium MXE',
      'Anonymous transactions via mixer pool',
      'Privacy pool mixing for unlinkability',
      'Rescue/Poseidon/SHA-3 based commitments',
      'Master viewing keys for balance scanning',
    ],
    howItWorks: {
      step1: 'Create wallet: POST /api/umbra/wallet/create',
      step2: 'Register for privacy: POST /api/umbra/register-full',
      step3: 'Deposit funds: POST /api/umbra/deposit',
      step4: 'Check balance: POST /api/umbra/balance',
      step5: 'Make anonymous transfers (coming soon)',
    },
  });
});

export default router;
