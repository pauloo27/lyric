package lyric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenius(t *testing.T) {
	t.Run("Test fetch", func(t *testing.T) {
		lyric, err := Fetch("https://genius.com/Richard-stallman-free-software-song-lyrics")
		assert.Nil(t, err)
		assert.NotEmpty(t, lyric)
		assert.Contains(t, lyric, "You'll be free, hackers, you'll be free")
		assert.Contains(t, lyric, "We'll kick out those dirty licenses")
		assert.Contains(t, lyric, "[Chorus]")
	})
	// TODO: search duckduckgo
	// TODO: search genius
}
