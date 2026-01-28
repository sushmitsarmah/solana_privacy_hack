declare module '../../umbra-sdk/index.cjs' {
  export * from '../../umbra-sdk/index';
}

declare module '../../umbra-sdk/index' {
  import { Connection, Keypair, PublicKey, TransactionSignature } from '@solana/web3.js';

  export class UmbraClient {
    static create(config: { connection: Connection }): Promise<UmbraClient>;
    static create(connection: Connection): Promise<UmbraClient>;
    static create(rpcUrl: string): Promise<UmbraClient>;

    setUmbraWallet(wallet: UmbraWallet): Promise<void>;
    setZkProver(type: 'wasm', config: any): void;

    // Account registration
    registerAccountForAnonymity(data: any, options: { mode: string }): Promise<string>;
    registerAccountForConfidentiality(data: any, options: { mode: string }): Promise<string>;
    registerAccountForConfidentialityAndAnonymity(
      data1: any,
      data2: any,
      options: { mode: string }
    ): Promise<string>;

    // Deposits
    depositPubliclyIntoMixerPoolSol(
      amount: bigint,
      destinationAddress: string,
      options: { mode: string }
    ): Promise<{ signature: string }>;

    // SPL token deposits
    depositPubliclyIntoMixerPoolSpl(
      amount: bigint,
      destinationAddress: string,
      mintAddress: string,
      options: { mode: string }
    ): Promise<{
      signature: string;
      claimableBalance?: bigint;
      generationIndex?: bigint;
      relayerPublicKey?: string;
    }>;

    // Confidential transfers (anonymous payments)
    transferConfidentially(
      amount: number,
      destinationAddress: string,
      mintAddress: string,
      opts?: { mode: string }
    ): Promise<string>;

    // Withdrawals from mixer pool
    claimDepositConfidentiallyFromMixerPool(
      mintAddress: string,
      destinationAddress: string,
      claimArtifacts: any,
      opts?: { mode: string }
    ): Promise<string>;

    // Encrypted balance queries
    getEncryptedTokenBalance(mintAddress: string): Promise<bigint>;
  }

  export class UmbraWallet {
    static fromSigner(
      signer: { signer: { keypair: Keypair } },
      config: { arciumX25519PublicKey: string }
    ): Promise<UmbraWallet>;

    getRescueCipherForPublicKey(publicKey: string): any;
    generateRandomSecret(index: bigint): bigint;
  }

  export const MXE_ARCIUM_X25519_PUBLIC_KEY: string;
  export const WSOL_MINT_ADDRESS: PublicKey;
}
