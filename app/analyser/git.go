package analyser

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
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

func createClientWithAuth(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
