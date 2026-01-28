declare module '../../umbra-sdk/index.cjs' {
  export * from '../../umbra-sdk/index';
}

declare module '../../umbra-sdk/index' {
  import { Connection, Keypair } from '@solana/web3.js';

  export class UmbraClient {
    static create(config: { connection: Connection }): Promise<UmbraClient>;
    setUmbraWallet(wallet: UmbraWallet): Promise<void>;
    setZkProver(type: 'wasm', config: any): void;
    registerAccountForAnonymity(data: any, options: { mode: string }): Promise<string>;
  }

  export class UmbraWallet {
    static fromSigner(
      signer: { signer: { keypair: Keypair } },
      config: { arciumX25519PublicKey: string }
    ): Promise<UmbraWallet>;
  }

  export const MXE_ARCIUM_X25519_PUBLIC_KEY: string;
}
