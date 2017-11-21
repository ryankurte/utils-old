package gpm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {

	t.Run("Semantic version tag matching", func(t *testing.T) {
		tests := []struct {
			Name  string
			Range string
			In    []string
			Out   interface{}
		}{
			{
				"Filters non-semver tags",
				"",
				[]string{"assbd", "1.2.3"},
				[]string{"1.2.3"},
			}, {
				"Filters tags by ranges",
				"^0.1.0",
				[]string{"0.1.1", "1.2.0"},
				[]string{"0.1.1"},
			}, {
				"Re-orders tags",
				"",
				[]string{"v0.1.0", "v0.3.0", "v0.2.0"},
				[]string{"v0.1.0", "v0.2.0", "v0.3.0"},
			},
		}

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				tags := make(Tags)
				for _, v := range test.In {
					tags[v] = v
				}
				filtered, err := tags.Filter(test.Range)
				assert.Nil(t, err)
				sorted, err := filtered.Sort()
				assert.Nil(t, err)
				assert.Equal(t, test.Out, sorted)
			})
		}
	})
}
