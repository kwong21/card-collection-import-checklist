package function

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseChecklist(t *testing.T) {
	actual, err := parseChecklist("test_data.xlsx")

	if assert.NoError(t, err) {
		sets := make([]string, 0, len(actual))

		for set := range actual {
			sets = append(sets, set)
		}

		assert.Equal(t, 3, len(sets))

		rookieBaseSet := actual["Base Set - Rookies"]
		assert.Equal(t, 10, len(rookieBaseSet))

		winterStormGoldAutoParallel := actual["Winter Storm Warning Golden Treasure Auto Parallel"]
		card := winterStormGoldAutoParallel[10]

		assert.Equal(t, "Bo Horvat", card.Player)
		assert.True(t, card.HasAuto)
	}
}
