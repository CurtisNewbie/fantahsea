package config

import (
	"testing"
	"fmt"
)

func TestParseProfile(t *testing.T) {

	args := make([]string, 2)
	args[0] = "profile=abc"
	args[1] = "--someflag"
	fmt.Printf("args: %v", args)
	

	profile := ParseProfile(args)
	if profile != "abc" {
		t.Errorf("Expected abc, but got: %v", profile)
	}
}



