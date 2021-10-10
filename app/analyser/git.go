package analyser

import (
	"fmt"
	"os/exec"
)

func createBranchName(cookstyleVersion string) string {
	return fmt.Sprintf("stylelia/cookstyle_%s", cookstyleVersion)
}

func buildBranchCommand(branchName string) *exec.Cmd {
	return exec.Command("git", "branch", "-b", branchName)
}

func createBranch(exec CommandRunner) error {
	// err := exec.Command("git", "branch", "-b", cmdMessage).Run()
	err := exec.Run()
	if err != nil {
		return err
	}

	return nil
}
