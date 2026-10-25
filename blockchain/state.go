package blockchain

import (
    "errors"
    "math"
)

var (
    ErrUnknownSender      = errors.New("unknown sender account")
    ErrUnknownRecipient   = errors.New("unknown recipient account")
    ErrZeroAmount         = errors.New("amount must be greater than zero")
    ErrSameAccount        = errors.New("sender and recipient must differ")
    ErrInsufficientFunds  = errors.New("insufficient funds")
    ErrInvalidNonce       = errors.New("invalid transaction nonce")
    ErrBalanceOverflow    = errors.New("recipient balance overflow")
    ErrNonceOverflow      = errors.New("sender nonce overflow")
)

// State is the current account state. Keep its maps private so callers cannot
// accidentally mutate balances without going through ApplyTransaction.
type State struct {
    balances map[string]uint64
    nonces   map[string]uint64
}

func NewState(initialBalances map[string]uint64) *State {
    balances := make(map[string]uint64, len(initialBalances))
    for account, balance := range initialBalances {
        balances[account] = balance
    }
    return &State{
        balances: balances,
        nonces:   make(map[string]uint64),
    }
}

func (s *State) Balance(account string) (uint64, bool) {
    balance, exists := s.balances[account]
    return balance, exists
}

func (s *State) Nonce(account string) uint64 {
    return s.nonces[account]
}

func (s *State) Clone() *State {
    cloned := NewState(s.balances)
    for account, nonce := range s.nonces {
        cloned.nonces[account] = nonce
    }
    return cloned
}

// ApplyTransaction validates first, then mutates state only after all checks pass.
func (s *State) ApplyTransaction(tx Transaction) error {
    senderBalance, senderExists := s.balances[tx.From]
    if !senderExists {
        return ErrUnknownSender
    }
    recipientBalance, recipientExists := s.balances[tx.To]
    if !recipientExists {
        return ErrUnknownRecipient
    }
    if tx.Amount == 0 {
        return ErrZeroAmount
    }
    if tx.From == tx.To {
        return ErrSameAccount
    }
    if senderBalance < tx.Amount {
        return ErrInsufficientFunds
    }
    if s.nonces[tx.From] == math.MaxUint64 {
        return ErrNonceOverflow
    }
    if tx.Nonce != s.nonces[tx.From]+1 {
        return ErrInvalidNonce
    }
    if recipientBalance > math.MaxUint64-tx.Amount {
        return ErrBalanceOverflow
    }

    // All checks passed. State mutation begins here.
    s.balances[tx.From] = senderBalance - tx.Amount
    s.balances[tx.To] = recipientBalance + tx.Amount
    s.nonces[tx.From] = tx.Nonce
    return nil
}
