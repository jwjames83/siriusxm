package siriusxm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type messages struct {
	Code int
	Message string
}

type artists struct {
	ID string
	Name string
}

type album struct {
	Name string
}

type createArts struct {
	Encrypted bool
	Size string
	Type string
	URL string
}

type song struct {
	Album      album
	Composer   string
	CreateArts [12]createArts
	ID         string
	Name       string

}

type currentEvent struct {
	Artists    artists
	BaseURL    string
	KeyIndex   string
	SiriusXMID int
	Song       song
	StartTime  string
}

type metaData struct {
	ChannelId     string
	ChannelName   string
	ChannelNumber int
	CurrentEvent  currentEvent
	DateTime      string
	Version       float32
}

type channelMetadataResponse struct {
	Messages messages
	Status   int
	MetaData metaData
}

type nowPlaying struct {
	ChannelMetadataResponse channelMetadataResponse
}

type NowPlaying struct {
	Song string
	Artist string
	Album string
	Composer string
}

// CurrentSong returns the song that is currently playing
func (playing channelMetadataResponse) CurrentSong() string {
	return playing.MetaData.CurrentEvent.Song.Name
}

// CurrentArtist returns the artist of the current song
func (playing channelMetadataResponse) CurrentArtist() string {
	return playing.MetaData.CurrentEvent.Artists.Name
}

func GetNowPlaying (channelId string) NowPlaying {
	var channelResponse channelMetadataResponse
	var nowPlaying NowPlaying
	timeString := time.Now().UTC().Format("01-02-15:04:00")
	resp, err := http.Get(
		fmt.Sprintf("https://www.siriusxm.com/metadata/pdt/en-us/json/channels/%s/timestamp/%s", channelId, timeString))

	if err != nil {

	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {

	}

	json.Unmarshal(skipRoot(body), &channelResponse)
	nowPlaying.Song = channelResponse.MetaData.CurrentEvent.Song.Name
	nowPlaying.Artist = channelResponse.MetaData.CurrentEvent.Artists.Name
	return nowPlaying
}