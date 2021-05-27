package test

import (
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