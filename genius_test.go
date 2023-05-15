package lyric_test

import (
	"testing"

	"github.com/pauloo27/lyric"
	"github.com/stretchr/testify/assert"
)

const freeSoftwareSong = "https://genius.com/Richard-stallman-free-software-song-lyrics"

func TestGenius(t *testing.T) {
	t.Run("Test fetch", func(t *testing.T) {
		lyric, err := lyric.Fetch(freeSoftwareSong)
		assert.Nil(t, err)
		assert.NotEmpty(t, lyric)
		assert.Contains(t, lyric, "You'll be free, hackers, you'll be free")
		assert.Contains(t, lyric, "We'll kick out those dirty licenses")
		assert.Contains(t, lyric, "[Chorus]")
	})
	t.Run("Test duckduckgo search", func(t *testing.T) {
		url, err := lyric.SearchDDG("free software song")
		assert.Nil(t, err)
		assert.NotNil(t, url)
		assert.Equal(t, freeSoftwareSong, url)
	})
	t.Run("Test genius search", func(t *testing.T) {
		url, err := lyric.Search("free software song")
		assert.Nil(t, err)
		assert.NotNil(t, url)
		assert.Equal(t, freeSoftwareSong, url)
	})
}
