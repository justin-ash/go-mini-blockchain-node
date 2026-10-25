package blockchain

import (
    "errors"
    "testing"
)

func newTestChain(t *testing.T) *Chain {
    t.Helper()
    chain, err := NewChain(map[string]uint64{
        "Alice": 100,
        "Bob": 0,
        "Carol": 0,
    })
    if err != nil {
        t.Fatal(err)
    }
    return chain
}

func TestValidTransfer(t *testing.T) {
    chain := newTestChain(t)
    if err := chain.AddBlock([]Transaction{{
        From: "Alice", To: "Bob", Amount: 80, Nonce: 1,
    }}); err != nil {
        t.Fatalf("AddBlock returned error: %v", err)
    }

    alice, _ := chain.Balance("Alice")
    bob, _ := chain.Balance("Bob")
    if alice != 20 || bob != 80 {
        t.Fatalf("unexpected balances: Alice=%d Bob=%d", alice, bob)
    }
    if err := chain.Validate(); err != nil {
        t.Fatalf("Validate returned error: %v", err)
    }
}

func TestInsufficientFundsDoesNotChangeState(t *testing.T) {
    chain := newTestChain(t)
    if err := chain.AddBlock([]Transaction{{
        From: "Alice", To: "Bob", Amount: 80, Nonce: 1,
    }}); err != nil {
        t.Fatal(err)
    }

    err := chain.AddBlock([]Transaction{{
        From: "Alice", To: "Carol", Amount: 30, Nonce: 2,
    }})
    if !errors.Is(err, ErrInsufficientFunds) {
        t.Fatalf("expected ErrInsufficientFunds, got %v", err)
    }

    alice, _ := chain.Balance("Alice")
    carol, _ := chain.Balance("Carol")
    if alice != 20 || carol != 0 {
        t.Fatalf("state changed after rejected tx: Alice=%d Carol=%d", alice, carol)
    }
    if len(chain.Blocks()) != 2 {
        t.Fatalf("rejected transaction should not create a block")
    }
}

func TestWrongNonceRejected(t *testing.T) {
    chain := newTestChain(t)
    err := chain.AddBlock([]Transaction{{
        From: "Alice", To: "Bob", Amount: 10, Nonce: 2,
    }})
    if !errors.Is(err, ErrInvalidNonce) {
        t.Fatalf("expected ErrInvalidNonce, got %v", err)
    }
}

func TestTamperingDetected(t *testing.T) {
    chain := newTestChain(t)
    if err := chain.AddBlock([]Transaction{{
        From: "Alice", To: "Bob", Amount: 10, Nonce: 1,
    }}); err != nil {
        t.Fatal(err)
    }

    chain.blocks[1].Transactions[0].Amount = 99
    if err := chain.Validate(); err == nil {
        t.Fatal("expected tampering to be detected")
    }
}
