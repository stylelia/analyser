package analyser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockCreateBranchCommand struct{}

func (m *MockCreateBranchCommand) Run() error {
	return nil
}

func (m *MockCreateBranchCommand) Output() ([]byte, error) {
	return nil, nil
}

type MockCreateBranchCommand_Error struct{}

func (m *MockCreateBranchCommand_Error) Run() error {
	return errors.New("test error")
}

func (m MockCreateBranchCommand_Error) Output() ([]byte, error) {
	return nil, nil
}

func TestCreateBranchName(t *testing.T) {
	version := "v10.10.10"
	expected := "stylelia/cookstyle_v10.10.10"

	actual := createBranchName(version)
	assert.Equal(t, expected, actual)
}

func TestBuildBranchCommand(t *testing.T) {
	message := "stylelia/cookstyle_v10.10.10"

	cmd := buildBranchCommand(message)

	expectedPath := "/usr/bin/git"
	assert.Equal(t, expectedPath, cmd.Path)

	expectedArgs := []string{"git", "branch", "-b", message}
	assert.Equal(t, expectedArgs, cmd.Args)
}

func TestCreateBranch(t *testing.T) {
	t.Run("createBranch throws an error on a faulty command", func(t *testing.T) {
		faulty := &MockCreateBranchCommand_Error{}

		err := createBranch(faulty)
		assert.Error(t, err)
	})

	t.Run("createBranch doesn't return any error on a valid command", func(t *testing.T) {
		runner := &MockCreateBranchCommand{}

		err := createBranch(runner)
		assert.NoError(t, err)
	})
}
