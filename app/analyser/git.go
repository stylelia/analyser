package analyser

import (
	"fmt"
	"os/exec"
)

func createBranchName(cookbookVersion string) string {
	return fmt.Sprintf("stylelia/cookstyle_%s", cookbookVersion)
}

func buildBranchCommand(message string) *exec.Cmd {
	return exec.Command("git", "branch", "-b", message)
}

func createBranch(exec CommandRunner) error {
	// err := exec.Command("git", "branch", "-b", cmdMessage).Run()
	err := exec.Run()
	if err != nil {
		return err
	}

	return nil
}