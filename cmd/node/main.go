package main

import (
    "fmt"
    "log"

    "github.com/justin-ash/go-mini-blockchain-node/blockchain"
)

func main() {
    chain, err := blockchain.NewChain(map[string]uint64{
        "Alice": 100,
        "Bob": 0,
        "Carol": 0,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Go Mini Blockchain Node")
    fmt.Println("-----------------------")

    first := blockchain.Transaction{
        From: "Alice", To: "Bob", Amount: 80, Nonce: 1,
    }
    if err := chain.AddBlock([]blockchain.Transaction{first}); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Accepted: Alice sends Bob 80 tokens")

    second := blockchain.Transaction{
        From: "Alice", To: "Carol", Amount: 30, Nonce: 2,
    }
    if err := chain.AddBlock([]blockchain.Transaction{second}); err != nil {
        fmt.Printf("Rejected: Alice sends Carol 30 tokens: %v\n", err)
    }

    alice, _ := chain.Balance("Alice")
    bob, _ := chain.Balance("Bob")
    carol, _ := chain.Balance("Carol")
    fmt.Printf("Balances: Alice=%d Bob=%d Carol=%d\n", alice, bob, carol)

    if err := chain.Validate(); err != nil {
        log.Fatalf("chain validation failed: %v", err)
    }
    fmt.Printf("Chain valid: %t (blocks=%d)\n", true, len(chain.Blocks()))
}
