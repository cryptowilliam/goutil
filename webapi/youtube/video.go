package youtube

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"google.golang.org/api/youtube/v3"
)

var (
	snippetContentDetailsStatistics = []string{"snippet", "contentDetails", "statistics"}
	snippetContentDetails           = []string{"snippet", "contentDetails"}
)

type (
	Video youtube.Video
)

func (yt Youtube) GetVideo(videoId string) (*Video, error) {
	response, err := yt.service.Videos.List(snippetContentDetailsStatistics).Id(videoId).Do()
	if err != nil {
		return nil, err
	}
	if len(response.Items) == 0 {
		return nil, gerrors.New("nil response length")
	}
	return (*Video)(response.Items[0]), nil
}

// Gets all playlists of current user - maxResult is set to 50 (default is 5)
// returns array of all playlists (id, name, count)
func (yt Youtube) GetAllPlaylists() ([]*youtube.Playlist, error) {

	// get all playlists
	call := yt.service.Playlists.List(snippetContentDetails)
	call = call.MaxResults(50).Mine(true)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}
