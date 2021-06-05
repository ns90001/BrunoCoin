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

	limTx.ChkTxs(nil)

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

	errorMessage := "ERROR {LiminalTxs.ChkTxs}: nil input list"

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}

func TestWalletChkTxsNormal (t *testing.T) {
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

	inputs2 := []*txi.TransactionInput{
		{
			TransactionHash: "123",
			OutputIndex:     0,
			UnlockingScript: "123",
			Amount:          1,
		}}

	outputs2 := []*txo.TransactionOutput{
		{
			Amount:        1,
			LockingScript: "123",
			Liminal:       false,
		},
	}

	transaction2 := tx.Transaction{
		Version:  1,
		Inputs: inputs2,
		Outputs:  outputs2,
		LockTime: 0,
	}

	limTx := wallet.LiminalTxs{
		TxQ:          tx.NewTxHeap(),
		TxRplyThresh: 0,
	}

	limTx.Add(&transaction)
	limTx.Add(&transaction2)
	limTx.ChkTxs([]*tx.Transaction{&transaction, &transaction2})

	if limTx.TxQ.Len() != 0 {
		t.Errorf("heap size should have changed. Expected length of 0, received: %v", limTx.TxQ.Len())
	}

}

// HndlTxReq Tests
func TestHndlTxReqZeroAmt (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	genNd.Wallet.HndlTxReq(&wallet.TxReq{
		PubK: nil,
		Amt:  0,
		Fee:  0,
	})

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

	errorMessage := "ERROR {Wallet.HndlTxReq}: no amount requested"

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}

func TestHndlTxReqLargeAmt (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	genNd.Wallet.HndlTxReq(&wallet.TxReq{
		PubK: []byte("123"),
		Amt:  100000,
		Fee:  1000,
	})

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

	errorMessage := "ERROR {Wallet.HndlTxReq}: could not get utxo for this request"

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}

func TestHndlTxReqNilInput (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	genNd.Wallet.HndlTxReq(nil)

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

	errorMessage := "ERROR {Wallet.HndlTxReq}: nil input"

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}

// HndlBlk Tests
func TestWalletHndlBlkNilInput (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	genNd.Wallet.HndlBlk(nil)

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

	errorMessage := "ERROR {Wallet.HndlBlk}: Inputted block is nil"

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage but received: %v", out)
	}
}