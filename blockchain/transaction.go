package blockchain

// Transaction represents a token transfer in this educational prototype.
// Nonce is the sender's next expected transaction number (starting at 1).
type Transaction struct {
    From   string `json:"from"`
    To     string `json:"to"`
    Amount uint64 `json:"amount"`
    Nonce  uint64 `json:"nonce"`
}
