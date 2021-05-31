package test

import (
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/block/tx/txi"
	"BrunoCoin/pkg/block/tx/txo"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestHndlTxNilInput (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	genNd.Mnr.HndlTx(nil)

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

	errorMessage := "ERROR {Miner.HndlTx}: transaction is nil"

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}

func TestHndlTxActiveMiner (t *testing.T) {
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

	genNd.Mnr.Pause()
	genNd.Mnr.HndlTx(&transaction)
	genNd.Mnr.StartMiner()

	isUpdated := <- genNd.Mnr.PoolUpdated
	
	if isUpdated{
		t.Errorf("mining pool should not have updated")
	}

}