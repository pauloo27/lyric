package lyric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const freeSoftwareSong = "https://genius.com/Richard-stallman-free-software-song-lyrics"

func TestGenius(t *testing.T) {
	t.Run("Test fetch", func(t *testing.T) {
		lyric, err := Fetch(freeSoftwareSong)
		assert.Nil(t, err)
		assert.NotEmpty(t, lyric)
		assert.Contains(t, lyric, "You'll be free, hackers, you'll be free")
		assert.Contains(t, lyric, "We'll kick out those dirty licenses")
		assert.Contains(t, lyric, "[Chorus]")
	})
	t.Run("Test duckduckgo search", func(t *testing.T) {
		url, err := SearchDDG("free software song")
		assert.Nil(t, err)
		assert.NotNil(t, url)
		assert.Equal(t, freeSoftwareSong, url)
	})
	// TODO: search genius
}
