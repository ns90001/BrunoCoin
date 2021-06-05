package test

import (
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/block/tx/txi"
	"BrunoCoin/pkg/block/tx/txo"
	"BrunoCoin/pkg/wallet"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// Wallet Tests

// Add Tests
func TestAddNilInput (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	limTx := wallet.LiminalTxs{
		TxQ:          nil,
		TxRplyThresh: 0,
	}

	limTx.Add(nil)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	errorMessage := "ERROR {LiminalTxs.Add}: transaction is nil"

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}

func TestWalletAddNormal (t *testing.T) {
	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	inputs := []*txi.TransactionInput{
		{
			TransactionHash: "123",
			OutputIndex:     0,
			UnlockingScript: "123",
			Amount:          1,
		}}

	outputs := []*txo.TransactionOutput{
		{
			Amount:        1,
			LockingScript: "123",
			Liminal:       false,
		},
	}

	transaction := tx.Transaction{
		Version:  1,
		Inputs: inputs,
		Outputs:  outputs,
		LockTime: 0,
	}

	limTx := wallet.LiminalTxs{
		TxQ:          tx.NewTxHeap(),
		TxRplyThresh: 0,
	}

	limTx.Add(&transaction)

	if limTx.TxQ.Len() != 1 {
		t.Errorf("expected transaction to be added, expected length of 1, received: %v", limTx.TxQ.Len())
	}

}

// ChkTxs Tests
func TestChkTxsNilInput (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	limTx := wallet.LiminalTxs{
		TxQ:          nil,
		TxRplyThresh: 0,
	}

	limTx.Add(nil)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	errorMessage := "ERROR {LiminalTxs.Add}: transaction is nil"

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}