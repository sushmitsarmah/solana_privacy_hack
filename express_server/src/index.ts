import express, { Request, Response } from 'express';
import cors from 'cors';
import dotenv from 'dotenv';
import umbraRoutes from './routes/umbra';

// Load environment variables
dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;

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
      'POST /api/umbra/balance',
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

// Start server
app.listen(PORT, () => {
  console.log('========================================');
  console.log('🔒 Umbra Privacy API Server');
  console.log('========================================');
  console.log(`Server running on http://localhost:${PORT}`);
  console.log(`\nAvailable endpoints:`);
  console.log(`  GET  http://localhost:${PORT}/`);
  console.log(`  GET  http://localhost:${PORT}/api/umbra/info`);
  console.log(`\nWallet Management:`);
  console.log(`  POST http://localhost:${PORT}/api/umbra/wallet/create`);
  console.log(`\nAccount Registration:`);
  console.log(`  POST http://localhost:${PORT}/api/umbra/register`);
  console.log(`  POST http://localhost:${PORT}/api/umbra/register-confidentiality`);
  console.log(`  POST http://localhost:${PORT}/api/umbra/register-full`);
  console.log(`\nPrivacy Pool Operations:`);
  console.log(`  POST http://localhost:${PORT}/api/umbra/deposit`);
  console.log(`  POST http://localhost:${PORT}/api/umbra/balance`);
  console.log('========================================');
  console.log('Powered by Umbra SDK - Zero-Knowledge Privacy on Solana');
  console.log('========================================');
});

export default app;
