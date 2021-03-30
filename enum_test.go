package main

import (
	"fmt"
	"testing"
)

func TestGeneration(t *testing.T) {
	output := generateSetLength(2, []string{"a", "b"})
	fmt.Println(output)
	if output != nil {
		t.Fail()
	}
}
