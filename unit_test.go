package main

import (
	"strings"
	"testing"
)

func TestToUpper(t *testing.T) {
	word := "a3bcD"
	ans := strings.ToUpper(word)
	t.Logf("%v\n", ans)
}
