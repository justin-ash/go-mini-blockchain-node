package blockchain

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
)

type Block struct {
    Index        uint64        `json:"index"`
    Timestamp    int64         `json:"timestamp"`
    Transactions []Transaction `json:"transactions"`
    PreviousHash string        `json:"previous_hash"`
    Hash         string        `json:"hash"`
}

func (b Block) CalculateHash() (string, error) {
    // Hash only the block payload, not the Hash field itself.
    payload := struct {
        Index        uint64        `json:"index"`
        Timestamp    int64         `json:"timestamp"`
        Transactions []Transaction `json:"transactions"`
        PreviousHash string        `json:"previous_hash"`
    }{
        Index: b.Index,
        Timestamp: b.Timestamp,
        Transactions: b.Transactions,
        PreviousHash: b.PreviousHash,
    }
    encoded, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }
    sum := sha256.Sum256(encoded)
    return hex.EncodeToString(sum[:]), nil
}
