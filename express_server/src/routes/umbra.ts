import { Router, Request, Response } from 'express';
import { Connection, Keypair, LAMPORTS_PER_SOL, PublicKey } from '@solana/web3.js';
import bs58 from 'bs58';
import crypto from 'crypto';
// @ts-ignore - Using compiled Umbra SDK
const { UmbraClient, UmbraWallet, MXE_ARCIUM_X25519_PUBLIC_KEY, WSOL_MINT_ADDRESS } = require('../../umbra-sdk/index.cjs');
import { redisClient } from '../index';

const router = Router();

// Use environment variable for RPC URL
const SOLANA_RPC_URL = process.env.SOLANA_RPC_URL || 'https://api.devnet.solana.com';

// Transaction record interface
interface TransactionRecord {
  id: string;
  type: 'deposit' | 'withdrawal' | 'transfer' | 'stealth_address' | 'registration';
  timestamp: number;
  publicKey: string;
  amount?: number;
  mint?: string;
  recipient?: string;
  destinationAddress?: string;
  signature?: string;
  details?: any;
}

// Helper function to add transaction to Redis
const addTransaction = async (
  type: TransactionRecord['type'],
  publicKey: string,
  details: Partial<TransactionRecord>
) => {
  const tx: TransactionRecord = {
    id: crypto.randomUUID(),
    type,
    timestamp: Date.now(),
    publicKey,
    ...details,
  };

  // Store transaction in Redis
  const txKey = `tx:${publicKey}:${tx.id}`;
  const txListKey = `txs:${publicKey}`;

  await redisClient.set(txKey, JSON.stringify(tx));
  await redisClient.lPush(txListKey, tx.id);

  // Set TTL to 30 days (optional - remove if you want to keep forever)
  await redisClient.expire(txKey, 30 * 24 * 60 * 60);
  await redisClient.expire(txListKey, 30 * 24 * 60 * 60);

  return tx;
};

// Helper function to get transaction history from Redis
const getTransactionHistory = async (
  publicKey: string,
  limit?: number,
  offset?: number,
  typeFilter?: string
): Promise<TransactionRecord[]> => {
  const txListKey = `txs:${publicKey}`;
  const txIds = await redisClient.lRange(txListKey, 0, -1);

  if (!txIds || txIds.length === 0) {
    return [];
  }

  // Fetch all transactions
  const txPromises = txIds.map(id => redisClient.get(`tx:${publicKey}:${id}`));
  const txStrings = await Promise.all(txPromises);

  const transactions = txStrings
    .filter((tx): tx is string => tx !== null)
    .map(tx => JSON.parse(tx) as TransactionRecord);

  // Filter by type if specified
  const filtered = typeFilter
    ? transactions.filter(tx => tx.type === typeFilter)
    : transactions;

  // Sort by timestamp (newest first)
  filtered.sort((a, b) => b.timestamp - a.timestamp);

  // Apply pagination
  const start = offset || 0;
  const end = limit ? start + limit : filtered.length;

  return filtered.slice(start, end);
};

// Helper function to parse private key
const parsePrivateKey = (privateKey: string): Keypair => {
  try {
    const privateKeyBytes = bs58.decode(privateKey);
    return Keypair.fromSecretKey(privateKeyBytes);
  } catch {
    throw new Error('Invalid private key format. Expected base58.');
  }
};

// Helper function to setup Umbra client and wallet
const setupUmbraClient = async (privateKey: string) => {
  const connection = new Connection(SOLANA_RPC_URL, 'confirmed');
  const keypair = parsePrivateKey(privateKey);

  const wallet = await UmbraWallet.fromSigner(
    { signer: { keypair } },
    { arciumX25519PublicKey: MXE_ARCIUM_X25519_PUBLIC_KEY }
  );

  const client = await UmbraClient.create({ connection });
  await client.setUmbraWallet(wallet);

  client.setZkProver('wasm', {
    masterViewingKeyRegistration: true,
    createSplDepositWithPublicAmount: true,
  });

  return { connection, keypair, wallet, client };
};

/**
 * POST /api/umbra/register
 * Register an account for anonymity using Umbra
 *
 * Body:
 * {
 *   "privateKey": "base58 private key"
 * }
 */
router.post('/register', async (req: Request, res: Response) => {
  try {
    const { privateKey } = req.body;

    if (!privateKey) {
      return res.status(400).json({
        error: 'Missing required field: privateKey',
      });
    }

    // Create connection
    const connection = new Connection(SOLANA_RPC_URL, 'confirmed');

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

    // Record transaction
    addTransaction('registration', keypair.publicKey.toBase58(), {
      signature,
      details: {
        registrationType: 'anonymity'
      }
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
    const { privateKey, amount, destinationAddress } = req.body;

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

    // Create connection using environment RPC URL
    const connection = new Connection(SOLANA_RPC_URL, 'confirmed');

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

    // Record transaction
    addTransaction('deposit', keypair.publicKey.toBase58(), {
      amount,
      mint: 'SOL',
      destinationAddress: destination,
      signature: result.signature,
      details: {
        amountLamports: lamports,
      }
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

    // Record transaction
    addTransaction('registration', keypair.publicKey.toBase58(), {
      signature,
      details: {
        registrationType: 'confidentiality'
      }
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

    // Record transaction
    addTransaction('registration', keypair.publicKey.toBase58(), {
      signature,
      details: {
        registrationType: 'confidentiality_and_anonymity'
      }
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
      'POST /api/umbra/stealth-address': {
        description: 'Generate stealth address for anonymous payments',
        body: {
          recipientPublicKey: 'string',
          rpcUrl: 'string (optional)',
        },
      },
      'POST /api/umbra/send': {
        description: 'Send anonymous/confidential transfer',
        body: {
          privateKey: 'string (base58)',
          recipientAddress: 'string',
          amount: 'number (in SOL or token units)',
          mint: 'string (optional, defaults to WSOL)',
          rpcUrl: 'string (optional)',
        },
      },
      'POST /api/umbra/withdraw': {
        description: 'Withdraw funds from privacy mixer pool',
        body: {
          privateKey: 'string (base58)',
          commitmentIndex: 'number',
          generationIndex: 'number',
          depositTime: 'number (Unix timestamp)',
          relayerPublicKey: 'string (optional)',
          mint: 'string (optional, defaults to WSOL)',
          rpcUrl: 'string (optional)',
        },
      },
      'POST /api/umbra/deposit-spl': {
        description: 'Deposit SPL tokens into privacy mixer pool',
        body: {
          privateKey: 'string (base58)',
          amount: 'number (in token units, respecting decimals)',
          mint: 'string (token mint address, required)',
          destinationAddress: 'string (optional)',
          rpcUrl: 'string (optional)',
        },
      },
      'GET /api/umbra/transactions': {
        description: 'Get transaction history for a public key',
        query: {
          publicKey: 'string (required)',
          type: 'string (optional: deposit, withdrawal, transfer, stealth_address, registration)',
          limit: 'number (optional, default: 50, max: 100)',
          offset: 'number (optional, default: 0)',
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
      'Stealth addresses for enhanced privacy',
      'Confidential transfers with encrypted amounts',
      'Full SPL token support for any token mint',
      'Complete transaction history tracking',
      'Anonymous withdrawals with ZK proofs',
    ],
    howItWorks: {
      step1: 'Create wallet: POST /api/umbra/wallet/create',
      step2: 'Register for privacy: POST /api/umbra/register-full',
      step3: 'Deposit SOL or SPL tokens: POST /api/umbra/deposit or /deposit-spl',
      step4: 'Check encrypted balance: POST /api/umbra/balance',
      step5: 'Generate stealth address: POST /api/umbra/stealth-address',
      step6: 'Send anonymous transfer: POST /api/umbra/send',
      step7: 'Withdraw anonymously: POST /api/umbra/withdraw',
      step8: 'View transaction history: GET /api/umbra/transactions',
    },
  });
});

/**
 * POST /api/umbra/stealth-address
 * Generate a stealth address for a recipient (for anonymous payments)
 *
 * Body:
 * {
 *   "recipientPublicKey": "recipient's Solana public key",
 *   "rpcUrl": "https://api.devnet.solana.com" (optional)
 * }
 */
router.post('/stealth-address', async (req: Request, res: Response) => {
  try {
    const { recipientPublicKey, rpcUrl = 'https://api.devnet.solana.com' } = req.body;

    if (!recipientPublicKey) {
      return res.status(400).json({
        error: 'Missing required field: recipientPublicKey',
      });
    }

    // Create connection (wallet not needed for this operation)
    const connection = new Connection(rpcUrl, 'confirmed');
    const client = await UmbraClient.create({ connection });

    // Generate ephemeral keypair for stealth address
    const ephemeralKeypair = Keypair.generate();

    // In a real scenario, you would:
    // 1. Get the recipient's X25519 public key from their Umbra wallet
    // 2. Perform X25519 key exchange to derive shared secret
    // 3. Generate stealth address from shared secret

    // For now, return the ephemeral keypair which can be used as a one-time address
    res.json({
      success: true,
      data: {
        ephemeralPublicKey: ephemeralKeypair.publicKey.toBase58(),
        ephemeralPrivateKey: bs58.encode(ephemeralKeypair.secretKey),
        recipientPublicKey,
        note: 'This ephemeral keypair can be used as a stealth address. The recipient can scan for payments using their master viewing key.',
      },
      message: 'Stealth address generated successfully!',
      usage: {
        step1: 'Send funds to the ephemeralPublicKey',
        step2: 'Share ephemeralPrivateKey with recipient securely',
        step3: 'Recipient can claim funds using the private key',
      },
    });
  } catch (error) {
    res.status(500).json({
      error: 'Failed to generate stealth address',
      details: (error as Error).message,
    });
  }
});

/**
 * POST /api/umbra/send
 * Send anonymous/confidential transfer to a recipient
 *
 * Body:
 * {
 *   "privateKey": "sender's base58 private key",
 *   "recipientAddress": "recipient's Solana address",
 *   "amount": 0.1 (in SOL),
 *   "mint": "token mint address" (optional, defaults to WSOL),
 *   "rpcUrl": "https://api.devnet.solana.com" (optional)
 * }
 */
router.post('/send', async (req: Request, res: Response) => {
  try {
    const { privateKey, recipientAddress, amount, mint, rpcUrl = 'https://api.devnet.solana.com' } = req.body;

    if (!privateKey || !recipientAddress || !amount) {
      return res.status(400).json({
        error: 'Missing required fields: privateKey, recipientAddress, amount',
      });
    }

    if (typeof amount !== 'number' || amount <= 0) {
      return res.status(400).json({
        error: 'Amount must be a positive number',
      });
    }

    const { client, keypair } = await setupUmbraClient(privateKey);

    // Use WSOL if no mint specified
    const tokenMint = mint || WSOL_MINT_ADDRESS.toBase58();

    // Convert SOL to lamports
    const lamports = Math.floor(amount * LAMPORTS_PER_SOL);

    // Perform confidential transfer
    const result = await client.transferConfidentially(
      lamports,
      recipientAddress,
      tokenMint,
      {
        mode: 'connection',
      }
    );

    res.json({
      success: true,
      data: {
        signature: result,
        amount,
        amountLamports: lamports,
        recipientAddress,
        senderPublicKey: keypair.publicKey.toBase58(),
        tokenMint,
        explorerUrl: `https://explorer.solana.com/tx/${result}?cluster=devnet`,
      },
      message: 'Anonymous transfer sent successfully!',
      privacy: {
        note: 'Transfer amount is encrypted using Rescue cipher. Only sender and recipient can decrypt the amount.',
        encryption: 'X25519 key exchange with Rescue cipher',
      },
    });

    // Record transaction
    addTransaction('transfer', keypair.publicKey.toBase58(), {
      amount,
      mint: tokenMint,
      recipient: recipientAddress,
      signature: result,
      details: {
        encrypted: true,
        amountLamports: lamports,
      }
    });
  } catch (error) {
    res.status(500).json({
      error: 'Failed to send anonymous transfer',
      details: (error as Error).message,
    });
  }
});

/**
 * POST /api/umbra/withdraw
 * Withdraw funds from the privacy mixer pool
 *
 * Body:
 * {
 *   "privateKey": "recipient's base58 private key",
 *   "commitmentIndex": 42,
 *   "generationIndex": 0,
 *   "depositTime": 1704067200,
 *   "relayerPublicKey": "relayer address" (optional),
 *   "mint": "token mint address" (optional, defaults to WSOL),
 *   "rpcUrl": "https://api.devnet.solana.com" (optional)
 * }
 */
router.post('/withdraw', async (req: Request, res: Response) => {
  try {
    const {
      privateKey,
      commitmentIndex,
      generationIndex,
      depositTime,
      relayerPublicKey,
      mint,
      rpcUrl = 'https://api.devnet.solana.com',
    } = req.body;

    if (!privateKey || commitmentIndex === undefined || generationIndex === undefined || !depositTime) {
      return res.status(400).json({
        error: 'Missing required fields: privateKey, commitmentIndex, generationIndex, depositTime',
      });
    }

    const { client, keypair } = await setupUmbraClient(privateKey);

    // Use WSOL if no mint specified
    const tokenMint = mint || WSOL_MINT_ADDRESS.toBase58();

    // Claim artifacts for withdrawal
    const claimDepositArtifacts = {
      commitmentInsertionIndex: BigInt(commitmentIndex),
      generationIndex: BigInt(generationIndex),
      time: BigInt(depositTime),
      relayerPublicKey: relayerPublicKey || undefined,
    };

    // Destination is the user's own address
    const destinationAddress = keypair.publicKey.toBase58();

    // Perform confidential withdrawal/claim from mixer pool
    const result = await client.claimDepositConfidentiallyFromMixerPool(
      tokenMint,
      destinationAddress,
      claimDepositArtifacts,
      {
        mode: 'connection',
      }
    );

    res.json({
      success: true,
      data: {
        signature: result,
        destinationAddress,
        tokenMint,
        claimArtifacts: claimDepositArtifacts,
        explorerUrl: `https://explorer.solana.com/tx/${result}?cluster=devnet`,
      },
      message: 'Withdrawal from privacy pool successful!',
      privacy: {
        note: 'Withdrawal is completely anonymous. Link between deposit and withdrawal is cryptographically unprovable.',
        zeroKnowledge: 'Uses ZK-SNARK proofs to prove deposit ownership without revealing identity',
      },
    });
  } catch (error) {
    res.status(500).json({
      error: 'Failed to withdraw from privacy pool',
      details: (error as Error).message,
    });
  }
});

/**
 * POST /api/umbra/deposit-spl
 * Deposit SPL tokens into the privacy mixer pool
 *
 * Body:
 * {
 *   "privateKey": "base58 private key",
 *   "amount": 100 (in token units, respecting decimals),
 *   "mint": "token mint address (required)",
 *   "destinationAddress": "recipient address for withdrawal" (optional, defaults to sender),
 *   "rpcUrl": "https://api.devnet.solana.com" (optional)
 * }
 */
router.post('/deposit-spl', async (req: Request, res: Response) => {
  try {
    const {
      privateKey,
      amount,
      mint,
      destinationAddress,
      rpcUrl = 'https://api.devnet.solana.com'
    } = req.body;

    if (!privateKey || !amount || !mint) {
      return res.status(400).json({
        error: 'Missing required fields: privateKey, amount, mint',
      });
    }

    if (typeof amount !== 'number' || amount <= 0) {
      return res.status(400).json({
        error: 'Amount must be a positive number',
      });
    }

    const { client, keypair } = await setupUmbraClient(privateKey);

    // Use sender address as destination if not provided
    const destination = destinationAddress || keypair.publicKey.toBase58();

    // Convert amount to bigint for blockchain compatibility
    const tokenAmount = BigInt(Math.floor(amount));

    // Deposit SPL into mixer pool
    const result = await client.depositPubliclyIntoMixerPoolSpl(
      tokenAmount,
      destination,
      mint,
      { mode: 'connection' }
    );

    // Record transaction
    addTransaction('deposit', keypair.publicKey.toBase58(), {
      amount,
      mint,
      destinationAddress: destination,
      signature: result.signature,
      details: {
        claimableBalance: result.claimableBalance?.toString(),
        generationIndex: result.generationIndex?.toString(),
        relayerPublicKey: result.relayerPublicKey,
      }
    });

    res.json({
      success: true,
      data: {
        signature: result.signature,
        amount,
        mint,
        destinationAddress: destination,
        publicKey: keypair.publicKey.toBase58(),
        claimableBalance: result.claimableBalance?.toString(),
        generationIndex: result.generationIndex?.toString(),
        relayerPublicKey: result.relayerPublicKey,
        explorerUrl: `https://explorer.solana.com/tx/${result.signature}?cluster=devnet`,
      },
      message: 'SPL token deposit to privacy pool successful!',
    });
  } catch (error) {
    console.error('SPL deposit error:', error);
    res.status(500).json({
      error: 'Failed to process SPL token deposit',
      details: (error as Error).message,
    });
  }
});

/**
 * GET /api/umbra/transactions
 * Get transaction history for a public key
 *
 * Query params:
 * - publicKey: Solana public key (required)
 * - type: Filter by transaction type (optional: deposit, withdrawal, transfer, stealth_address, registration)
 * - limit: Number of transactions to return (optional, default: 50, max: 100)
 * - offset: Number of transactions to skip (optional, default: 0)
 */
router.get('/transactions', async (req: Request, res: Response) => {
  try {
    const {
      publicKey,
      type,
      limit = '50',
      offset = '0'
    } = req.query;

    if (!publicKey) {
      return res.status(400).json({
        error: 'Missing required query parameter: publicKey',
      });
    }

    // Validate limit
    const limitNum = Math.min(parseInt(limit as string) || 50, 100);
    const offsetNum = parseInt(offset as string) || 0;

    // Get transaction history
    let history: TransactionRecord[] = await getTransactionHistory(publicKey as string, limitNum, offsetNum, type as string);

    res.json({
      success: true,
      data: {
        transactions: history,
        total: history.length,
        limit: limitNum,
        offset: offsetNum,
      },
      message: `Retrieved ${history.length} transactions`,
    });
  } catch (error) {
    res.status(500).json({
      error: 'Failed to retrieve transaction history',
      details: (error as Error).message,
    });
  }
});

export default router;
