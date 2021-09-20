package analyser

import (
	"fmt"
	"os/exec"
)

func createBranch(cookbookVersion string) error {
	cmdMessage := fmt.Sprintf("stylelia/cookstyle_%s", cookbookVersion)

	err := exec.Command("git", "branch", "-b", cmdMessage).Run()
	if err != nil {
		return err
	}

	return nil
}
