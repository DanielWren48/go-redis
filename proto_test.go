package main

import (
	"fmt"
	"testing"
)

func TestProtocol(t *testing.T) {
	raw := ""
	cmd, err := parseCommand(raw)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cmd)
}
