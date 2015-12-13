package main

import(
	"testing"
	"os/exec"
	"fmt"
)

func TestFormat(t *testing.T){
	out, err := exec.Command("./commandExecuter", "command_file.txt").Output()
	if err == nil && len(out) != 0{
		t.Error("exec.Command call failed")
	}
	fmt.Print(string(out[:]))

	out, err = exec.Command("./commandExecuter", "command_file_spaced.txt").Output()
	if err == nil && len(out) != 0{
		t.Error("exec.Command call failed")
	}
	fmt.Print(string(out[:]))
}