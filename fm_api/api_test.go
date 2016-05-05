package fm_api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLondonTime(t *testing.T) {
	_, err := getLondonTime()
	assert.Nil(t, err)
}

func TestParseQueue(t *testing.T) {
	body := `[
	    {
	        "track": {
	            "album": {
	                "id": "01981749-3291-49fe-bdf4-7c7102650b74",
	                "images": [
	                    {
	                        "url": "https://i.scdn.co/image/545af3401e7b5bc1e81590557b7daa7c8faceb9b",
	                        "width": 640,
	                        "height": 640
	                    },
	                    {
	                        "url": "https://i.scdn.co/image/9e399897e283ae3156522fd27c99870e5bb3e82b",
	                        "width": 300,
	                        "height": 300
	                    },
	                    {
	                        "url": "https://i.scdn.co/image/5d9d645c85529fa300f1f46407a6393f8739adae",
	                        "width": 64,
	                        "height": 64
	                    }
	                ],
	                "name": "Dive",
	                "uri": "spotify:album:3I3PmRvn5iFY8i6zzvEcci"
	            },
	            "name": "Elegy",
	            "uri": "spotify:track:5MbUyUE6erY9mVgXaecQwR",
	            "play_count": 0,
	            "artists": [
	                {
	                    "id": "6f840614-5675-44aa-9a03-d74093887edc",
	                    "uri": "spotify:artist:5oOhM2DFWab8XhSdQiITry",
	                    "name": "Tycho"
	                }
	            ],
	            "duration": 263560,
	            "id": "181daf93-81ef-4538-9f72-9872b486fe4e"
	        },
	        "user": {
	            "family_name": "",
	            "display_name": "",
	            "avatar_url": "",
	            "spotify_playlists": null,
	            "given_name": "",
	            "id": ""
	        },
	        "uuid": ""
	    }
	]`

	manager := FmApiManager{}
	queue, err := manager.parseQueue([]byte(body))

	assert.Nil(t, err)
	assert.Equal(t, 1, len(queue))
}
