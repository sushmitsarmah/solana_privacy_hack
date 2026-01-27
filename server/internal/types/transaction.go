package types

// UnsignedTxResponse contains the base64 encoded transaction to be signed by the user.
type UnsignedTxResponse struct {
	UnsignedTxBase64     string `json:"unsigned_tx_base64"`
	RecentBlockhash      string `json:"recent_blockhash"`
	LastValidBlockHeight int64  `json:"last_valid_block_height"`
}
