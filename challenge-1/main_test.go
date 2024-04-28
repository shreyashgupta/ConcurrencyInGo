package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_Main(t *testing.T) {

	stdOut := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := io.ReadAll(r)

	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "universe") {
		t.Errorf("universe statement is missing")
	}
	if !strings.Contains(output, "cosmos") {
		t.Errorf("cosmos statement is missing")
	}
	if !strings.Contains(output, "world") {
		t.Errorf("world statement is missing")
	}

	id_universe := strings.Index(output, "universe")
	id_cosmos := strings.Index(output, "cosmos")
	id_world := strings.Index(output, "world")

	if id_universe > id_cosmos || id_universe > id_world || id_cosmos > id_world {
		t.Errorf("order is incorrect")
	}
}

func Test_UpdateMessage(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)

	updateMessage("shreyash", &wg)
	wg.Wait()

	if msg != "shreyash" {
		t.Errorf("Failed to update message, got: " + msg)
	}
}

func Test_PrintMessage(t *testing.T) {

	stdOut := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w
	msg = "shreyash"

	printMessage()

	_ = w.Close()

	result, _ := io.ReadAll(r)

	output := string(result)
	os.Stdout = stdOut

	if !strings.Contains(output, msg) {
		t.Errorf("output of print message is incorrect")
	}
}
