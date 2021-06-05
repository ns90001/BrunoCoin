package test

import (
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/block/tx/txi"
	"BrunoCoin/pkg/block/tx/txo"
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"strings"
	"testing"
)

// Validation functions

func TestChkTxNilInput (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	genNd.ChkTx(nil)

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

	errorMessage := "ERROR {Node.ChkTx}: nil transaction"

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}

func TestChkTxEmptyIO (t *testing.T) {
	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	transaction := tx.Transaction{
		Version:  1,
		Inputs: make([]*txi.TransactionInput, 0),
		Outputs:  make([]*txo.TransactionOutput, 0),
		LockTime: 0,
	}

	isValid := genNd.ChkTx(&transaction)

	if isValid {
		t.Errorf("expected invalid transaction since inputs/outputs are empty")
	}

}

func TestChkTxNormalTx (t *testing.T) {
	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	inputs := []*txi.TransactionInput{
		{
			TransactionHash: "123",
			OutputIndex:     0,
			UnlockingScript: "123",
			Amount:          10,
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

	genNd.Chain.Add(&block.Block{
		Hdr:          block.Header{},
		Transactions: []*tx.Transaction{&transaction},
	})

	isValid := genNd.ChkTx(&transaction)

	if !isValid {
		t.Errorf("function returned false, should return true")
	}

}

func TestChkTxOutputOfZero (t *testing.T) {
	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	inputs := []*txi.TransactionInput{
		{
			TransactionHash: "123",
			OutputIndex:     0,
			UnlockingScript: "123",
			Amount:          10,
		}}

	outputs := []*txo.TransactionOutput{
		{
			Amount:        0,
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

	isValid := genNd.ChkTx(&transaction)

	if isValid {
		t.Errorf("function returned true, should return false: output amount is 0")
	}

}

func TestChkTxOutputsGreaterInput (t *testing.T) {
	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	inputs := []*txi.TransactionInput{
		{
			TransactionHash: "123",
			OutputIndex:     0,
			UnlockingScript: "123",
			Amount:          10,
		}}

	outputs := []*txo.TransactionOutput{
		{
			Amount:        1000,
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

	isValid := genNd.ChkTx(&transaction)

	if isValid {
		t.Errorf("function returned true, should return false: output amt greater than input amt")
	}

}

func TestChkBlk(t *testing.T) {
	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	var byte1 byte = 0

	testHdr := block.Header{
		Ver:       0,
		PrvBlkHsh: "",
		MrklRt:    "",
		Timestamp: 0,
		DiffTarg:  hex.EncodeToString([]byte{byte1}),
		Nonce:     0,
	}

	b := block.Block{
		Hdr: testHdr,
		Transactions: nil,
	}

	t1 := genNd.ChkBlk(&b)

	if t1 {
		t.Errorf("ERROR {validation.ChkBlk}: POW not satisfied")
	}

}

