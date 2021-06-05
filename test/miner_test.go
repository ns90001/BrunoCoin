package test

import (
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/block/tx/txi"
	"BrunoCoin/pkg/block/tx/txo"
	"BrunoCoin/pkg/miner"
	"bytes"
	"go.uber.org/atomic"
	"io"
	"os"
	"strings"
	"testing"
)

// Mining Tests

//HndlTx Tests
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

//GenCBTx Tests
func TestGenCBTxEmptyList(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	// test with empty tx list input
	genNd.Mnr.GenCBTx(make([]*tx.Transaction, 0))

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

	errorMessage1 := "ERROR {Miner.GenCBTx}: no transactions given"

	if !strings.Contains(out, errorMessage1) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}

func TestGenCBTxNilTx(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	// test with nil transaction in input list
	genNd.Mnr.GenCBTx([]*tx.Transaction{nil})

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

	errorMessage1 := "ERROR {Miner.GenCBTx}: nil transaction"

	if !strings.Contains(out, errorMessage1) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}

func TestGenCBTxNormalTxs(t *testing.T) {

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

	transaction1 := tx.Transaction{
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
			Amount:          100,
		}}

	outputs2 := []*txo.TransactionOutput{
		{
			Amount:        99,
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

	CBTx := genNd.Mnr.GenCBTx([]*tx.Transaction{&transaction1, &transaction2})

	if !(CBTx.IsCoinbase() && CBTx.Outputs[0].Amount == 11) {
		t.Errorf("expected coinbase tx to be generated with output amount 10")
	}

}

// HndlBlk Tests
func TestHndlBlkNilInput (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	genNd.Mnr.HndlBlk(nil)

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

	errorMessage := "ERROR {Miner.HndlBlk}: nil block"

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}

func TestHndlBlkNormalInput (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	genNd.Mnr.HndlBlk(&block.Block{
		Hdr:          block.Header{},
		Transactions: nil,
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

	errorMessage := "ERROR {Miner.HndlBlk}: nil block"

	if strings.Contains(out, errorMessage) {
		t.Errorf("unexpected error: %v", out)
	}
	if genNd.Mnr.ChnLen.Load() != 2 {
		t.Errorf("chain length expected to be 2, received: %v", genNd.Mnr.ChnLen.Load())
	}

}

// CalcPri Tests
func TestCalcPriZeroPriority (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

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

	// fees are zero in this case, so priority will be calculated as 0, so 1 should return
	priority := miner.CalcPri(&transaction)

	errorMessage := "transaction is nil"

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

	if strings.Contains(out, errorMessage) {
		t.Errorf("unexpected error: %v", out)
	}
	if priority != 1 {
		t.Errorf("priority expected to be 1, received: %v", priority)
	}

}

func TestCalcPriNilInput (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	miner.CalcPri(nil)

	errorMessage := "transaction is nil"

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

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage, received: %v", out)
	}

}

// Add Tests
func TestAddNilTx (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	txpool := miner.TxPool{
		CurPri: nil,
		PriLim: 10,
		TxQ:    nil,
		Ct:     nil,
		Cap:    2,
	}

	txpool.Add(nil)

	errorMessage := "ERROR {TxPool.Add}: transaction is nil"

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

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage, received: %v", out)
	}

}

func TestAddPoolFull (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

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

	txpool := miner.TxPool{
		CurPri: nil,
		PriLim: 10,
		TxQ:    nil,
		Ct:     atomic.NewUint32(10),
		Cap:    10,
	}

	// should not add transaction since cap has already been reached
	txpool.Add(&transaction)

	errorMessage := "ERROR {TxPool.Add}: transaction pool is full"

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

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage, received: %v", out)
	}

}

func TestAddNormal (t *testing.T) {
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

	txpool := miner.TxPool{
		CurPri: atomic.NewUint32(1),
		PriLim: 10,
		TxQ:    tx.NewTxHeap(),
		Ct:     atomic.NewUint32(1),
		Cap:    10,
	}

	txpool.Add(&transaction)

	if txpool.TxQ.Len() != 1 {
		t.Errorf("expected transaction to be added, expected length of 1, received: %v", txpool.TxQ.Len())
	}

}

// ChkTxs Tests
func TestChkTxsEmptyRemover (t *testing.T) {
	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	txpool := miner.TxPool{
		CurPri: nil,
		PriLim: 10,
		TxQ:    tx.NewTxHeap(),
		Ct:     atomic.NewUint32(1),
		Cap:    10,
	}

	txpool.ChkTxs(make([]*tx.Transaction, 0))

	if txpool.TxQ.Len() != 0 {
		t.Errorf("heap length should not have changed. Expected length 0, received: %v", txpool.TxQ.Len())
	}

}

func TestChkTxsNormal (t *testing.T) {
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

	txpool := miner.TxPool{
		CurPri: atomic.NewUint32(1),
		PriLim: 10,
		TxQ:    tx.NewTxHeap(),
		Ct:     atomic.NewUint32(1),
		Cap:    10,
	}

	txpool.Add(&transaction)
	txpool.ChkTxs([]*tx.Transaction{&transaction})

	if txpool.TxQ.Len() != 0 {
		t.Errorf("expected duplicates to be removed, expected length of 0, received: %v", txpool.TxQ.Len())
	}

}

func TestChkTxsNilInput (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	txpool := miner.TxPool{
		CurPri: nil,
		PriLim: 10,
		TxQ:    nil,
		Ct:     atomic.NewUint32(10),
		Cap:    10,
	}

	// should not add since cap has been reached already
	txpool.ChkTxs([]*tx.Transaction{nil})

	errorMessage := "ERROR {TxPool.ChkTxs}: transaction is nil"

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

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage, received: %v", out)
	}

}

// HndlChkBlk Tests
func TestHndlChkBlkNilInput (t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	genNd := NewGenNd()
	genNd.Start()
	genNd.StartMiner()

	genNd.Mnr.HndlChkBlk(nil)

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

	errorMessage := "ERROR {Miner.HndlChkBlk}: nil block"

	if !strings.Contains(out, errorMessage) {
		t.Errorf("expected errorMessage but received: %v", out)
	}

}
