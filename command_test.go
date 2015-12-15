package main

import(
	"testing"
	"os/exec"
	"fmt"
)

func TestSpacedFormat(t *testing.T){
	out, err := exec.Command("./commandExecuter", "command_file_spaced.txt").Output()
	if err != nil && len(out) == 0{
		t.Error("exec.Command call failed")
	}
	fmt.Print(string(out))
}

func TestMixedFormat(t *testing.T){
	out, err := exec.Command("./commandExecuter", "command_file_mixed.txt").Output()
	if err != nil && len(out) == 0{
		t.Error("exec.Command call failed")
	}
	fmt.Print(string(out))
}

func TestInvalidFormat(t *testing.T){
	out, err := exec.Command("./commandExecuter", "command_file_invalid.txt").Output()
	if err != nil && len(out) == 0{
		t.Error("exec.Command call failed")
	}
	fmt.Print(string(out))
}