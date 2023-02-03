package cmd

import (
	"fmt"
	"testing"
)

func assertNoErr(t *testing.T, e error) {
	if e != nil {
		t.Error(e)
	}
}

func assertErr(t *testing.T, e error) {
	if e == nil {
		t.Error(fmt.Errorf("expected error but not error occurred"))
	}
}
