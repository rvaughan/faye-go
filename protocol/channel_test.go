package protocol_test

import (
	. "git.xuvasi.com/gocode/faye-go/protocol"
	"git.xuvasi.com/gocode/faye-go/utils"
	"testing"
)

func TestExpandSimpleChannel(t *testing.T) {
	chan1 := NewChannel("/foo/bar")

	expected := []string{
		"/**",
		"/foo/**",
		"/foo/*",
		"/foo/bar",
	}

	patterns := chan1.Expand()

	if !utils.CompareStringSlices(expected, patterns) {
		t.Fatal("Expected ", expected, " got ", patterns)
	}
}
