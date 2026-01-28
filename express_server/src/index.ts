import express, { Request, Response } from 'express';
import cors from 'cors';
import dotenv from 'dotenv';
import { createClient, RedisClientType } from 'redis';
import umbraRoutes from './routes/umbra';

// Load environment variables
dotenv.config();

// Validate required environment variables
const requiredEnvVars = ['SOLANA_RPC_URL', 'REDIS_URL'];
const missingEnvVars = requiredEnvVars.filter(envVar => !process.env[envVar]);

if (missingEnvVars.length > 0) {
  console.error('❌ Missing required environment variables:', missingEnvVars.join(', '));
  console.error('Please check your .env file and ensure all required variables are set.');
  process.exit(1);
}

// Initialize Redis client
export let redisClient: RedisClientType;

const initializeRedis = async () => {
  try {
    redisClient = createClient({
      url: process.env.REDIS_URL,
    });

    redisClient.on('error', (err) => {
      console.error('Redis Client Error:', err);
    });

    redisClient.on('connect', () => {
      console.log('✅ Connected to Redis successfully');
    });

    await redisClient.connect();
    console.log('🔌 Redis client initialized');
  } catch (error) {
    console.error('❌ Failed to connect to Redis:', error);
    console.error('Please ensure Redis is running and REDIS_URL is configured correctly.');
    process.exit(1);
  }
};

const app = express();
const PORT = process.env.PORT || 3000;
const SOLANA_RPC_URL = process.env.SOLANA_RPC_URL;

// Middleware
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Request logging middleware
app.use((req: Request, _res: Response, next) => {
  console.log(`[${new Date().toISOString()}] ${req.method} ${req.path}`);
  next();
});

// Health check endpoint
app.get('/', (_req: Request, res: Response) => {
  res.json({
    name: 'Umbra Privacy API',
    version: '1.0.0',
    status: 'running',
    description: 'API for privacy-focused payments on Solana using Umbra SDK',
    endpoints: {
      umbra: '/api/umbra',
    },
    documentation: {
      umbra: 'GET /api/umbra/info',
    },
  });
});

// API Routes
app.use('/api/umbra', umbraRoutes);

// 404 handler
app.use((_req: Request, res: Response) => {
  res.status(404).json({
    error: 'Not Found',
    message: 'The requested endpoint does not exist',
    availableEndpoints: [
      'GET /',
      'GET /api/umbra/info',
      'POST /api/umbra/wallet/create',
      'POST /api/umbra/register',
      'POST /api/umbra/register-confidentiality',
      'POST /api/umbra/register-full',
      'POST /api/umbra/deposit',
      'POST /api/umbra/deposit-spl',
      'POST /api/umbra/balance',
      'GET /api/umbra/transactions',
      'POST /api/umbra/stealth-address',
      'POST /api/umbra/send',
      'POST /api/umbra/withdraw',
    ],
  });
});

// Error handler
app.use((err: Error, _req: Request, res: Response, _next: any) => {
  console.error('Error:', err);
  res.status(500).json({
    error: 'Internal Server Error',
    message: err.message,
  });
});

// Start server with Redis initialization
const startServer = async () => {
  try {
    console.log('========================================');
    console.log('🔒 Umbra Privacy API Server - Starting...');
    console.log('========================================');

    // Initialize Redis
    await initializeRedis();

    // Start Express server
    app.listen(PORT, () => {
      console.log('✅ Server started successfully!');
      console.log(`📡 Server running on http://localhost:${PORT}`);
      console.log(`\n🔐 Environment: ${process.env.NODE_ENV || 'development'}`);
      console.log(`⛓️  Solana RPC: ${SOLANA_RPC_URL}`);
      console.log(`🗄️  Redis: ${process.env.REDIS_URL}`);
      console.log('========================================');
      console.log('\n📚 Available endpoints:');
      console.log(`  GET  http://localhost:${PORT}/`);
      console.log(`  GET  http://localhost:${PORT}/api/umbra/info`);
      console.log(`\n👛 Wallet Management:`);
      console.log(`  POST http://localhost:${PORT}/api/umbra/wallet/create`);
      console.log(`\n📝 Account Registration:`);
      console.log(`  POST http://localhost:${PORT}/api/umbra/register`);
      console.log(`  POST http://localhost:${PORT}/api/umbra/register-confidentiality`);
      console.log(`  POST http://localhost:${PORT}/api/umbra/register-full`);
      console.log(`\n🏊 Privacy Pool Operations:`);
      console.log(`  POST http://localhost:${PORT}/api/umbra/deposit (SOL)`);
      console.log(`  POST http://localhost:${PORT}/api/umbra/deposit-spl (Tokens)`);
      console.log(`  POST http://localhost:${PORT}/api/umbra/balance`);
      console.log(`  GET  http://localhost:${PORT}/api/umbra/transactions`);
      console.log(`\n🕵️ Stealth Addresses & Anonymous Transfers:`);
      console.log(`  POST http://localhost:${PORT}/api/umbra/stealth-address`);
      console.log(`  POST http://localhost:${PORT}/api/umbra/send`);
      console.log(`  POST http://localhost:${PORT}/api/umbra/withdraw`);
      console.log('========================================');
      console.log('⚡ Powered by Umbra SDK - Zero-Knowledge Privacy on Solana');
      console.log('========================================');
    });
  } catch (error) {
    console.error('❌ Failed to start server:', error);
    process.exit(1);
  }
};

startServer();
