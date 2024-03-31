package clients

import (
	"os/exec"
)

var (
	copyCmdArgs = "pbcopy"
)

func getCopyCommand() *exec.Cmd {
	return exec.Command(copyCmdArgs)
}

func WriteAll(text string) error {
	copyCmd := getCopyCommand()
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return copyCmd.Wait()
}
