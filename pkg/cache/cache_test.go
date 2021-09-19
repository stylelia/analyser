package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeyPath(t *testing.T) {

	t.Run("getKeyPath returns the correct keyPath", func(t *testing.T) {
		org := "youshy"
		name := "ENAE-The-System"
		expected := "github/youshy/ENAE-The-System"

		k := Key{
			OrgName:  org,
			RepoName: name,
		}
		path := k.getKeyPath()
		assert.Equal(t, expected, path)
	})
}
