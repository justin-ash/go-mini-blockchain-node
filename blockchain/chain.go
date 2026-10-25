package blockchain

import (
    "errors"
    "time"
)

var (
    ErrEmptyChain       = errors.New("chain must not be empty")
    ErrInvalidGenesis   = errors.New("invalid genesis block")
    ErrInvalidBlockHash = errors.New("invalid block hash")
    ErrInvalidLink      = errors.New("invalid previous-hash link")
    ErrInvalidIndex     = errors.New("invalid block index")
)

type Chain struct {
    blocks          []Block
    initialBalances map[string]uint64
    state           *State
}

func NewChain(initialBalances map[string]uint64) (*Chain, error) {
    copied := make(map[string]uint64, len(initialBalances))
    for account, balance := range initialBalances {
        copied[account] = balance
    }

    genesis := Block{
        Index: 0,
        Timestamp: time.Now().Unix(),
        Transactions: []Transaction{},
        PreviousHash: "",
    }
    hash, err := genesis.CalculateHash()
    if err != nil {
        return nil, err
    }
    genesis.Hash = hash

    return &Chain{
        blocks: []Block{genesis},
        initialBalances: copied,
        state: NewState(copied),
    }, nil
}

func (c *Chain) Blocks() []Block {
    result := make([]Block, len(c.blocks))
    copy(result, c.blocks)
    for i := range result {
        result[i].Transactions = append([]Transaction(nil), c.blocks[i].Transactions...)
    }
    return result
}

func (c *Chain) Balance(account string) (uint64, bool) {
    return c.state.Balance(account)
}

func (c *Chain) Nonce(account string) uint64 {
    return c.state.Nonce(account)
}

// AddBlock applies all transactions to a cloned state first. If any transaction
// fails, neither the state nor the chain is changed.
func (c *Chain) AddBlock(txs []Transaction) error {
    nextState := c.state.Clone()
    for _, tx := range txs {
        if err := nextState.ApplyTransaction(tx); err != nil {
            return err
        }
    }

    previous := c.blocks[len(c.blocks)-1]
    block := Block{
        Index: previous.Index + 1,
        Timestamp: time.Now().Unix(),
        Transactions: append([]Transaction(nil), txs...),
        PreviousHash: previous.Hash,
    }
    hash, err := block.CalculateHash()
    if err != nil {
        return err
    }
    block.Hash = hash

    c.blocks = append(c.blocks, block)
    c.state = nextState
    return nil
}

// Validate replays every transaction from the original balances and verifies
// the hashes and links. It does not provide consensus or protect against an
// attacker rewriting the entire chain.
func (c *Chain) Validate() error {
    if len(c.blocks) == 0 {
        return ErrEmptyChain
    }

    genesis := c.blocks[0]
    genesisHash, err := genesis.CalculateHash()
    if err != nil {
        return err
    }
    if genesis.Index != 0 || genesis.PreviousHash != "" ||
        len(genesis.Transactions) != 0 || genesis.Hash != genesisHash {
        return ErrInvalidGenesis
    }

    replayState := NewState(c.initialBalances)
    for i := 1; i < len(c.blocks); i++ {
        current := c.blocks[i]
        previous := c.blocks[i-1]

        if current.Index != uint64(i) {
            return ErrInvalidIndex
        }
        if current.PreviousHash != previous.Hash {
            return ErrInvalidLink
        }
        currentHash, hashErr := current.CalculateHash()
        if hashErr != nil {
            return hashErr
        }
        if current.Hash != currentHash {
            return ErrInvalidBlockHash
        }
        for _, tx := range current.Transactions {
            if applyErr := replayState.ApplyTransaction(tx); applyErr != nil {
                return applyErr
            }
        }
    }
    return nil
}
