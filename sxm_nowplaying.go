package siriusxm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type messages struct {
	Code    int
	Message string
}

type artists struct {
	ID   string
	Name string
}

type album struct {
	Name string
}

type createArts struct {
	Encrypted bool
	Size      string
	Type      string
	URL       string
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
	ChannelId     interface{}
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

type channelMetadataResponseArr struct {
	Messages messages
	Status   int
	MetaData []metaData
}

type NowPlaying struct {
	ChannelId     interface{}
	ChannelNumber int
	ChannelName   string
	Song          string
	Artist        string
	Album         string
	Composer      string
	Details       currentEvent
}

// CurrentSong returns the song that is currently playing
func (playing channelMetadataResponse) CurrentSong() string {
	return playing.MetaData.CurrentEvent.Song.Name
}

// CurrentArtist returns the artist of the current song
func (playing channelMetadataResponse) CurrentArtist() string {
	return playing.MetaData.CurrentEvent.Artists.Name
}

func GetAllNowPlaying() []NowPlaying {
	var channelResponse channelMetadataResponseArr
	var nowPlaying []NowPlaying
	timeString := time.Now().Add(time.Duration(-1) * time.Minute).UTC().Format("01-02-15:04:00")
	reqString := fmt.Sprintf("https://www.siriusxm.com/metadata/pdt/en-us/json/events/timestamp/%s", timeString)
	resp, err := http.Get(reqString)

	checkErr(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	checkErr(err)

	json.Unmarshal(skipRoot(body), &channelResponse)
	for _, info := range channelResponse.MetaData {
		nowPlaying = append(nowPlaying, NowPlaying{
			ChannelId:     info.ChannelId,
			ChannelNumber: info.ChannelNumber,
			ChannelName:   info.ChannelName,
			Song:          info.CurrentEvent.Song.Name,
			Artist:        info.CurrentEvent.Artists.Name,
			Album:         info.CurrentEvent.Song.Album.Name,
			Composer:      info.CurrentEvent.Song.Composer,
			Details:       info.CurrentEvent,
		})
	}

	return nowPlaying
}

func GetNowPlaying(channel interface{}) NowPlaying {
	var channelResponse channelMetadataResponse
	timeString := time.Now().Add(time.Duration(-1) * time.Minute).UTC().Format("01-02-15:04:00")
	reqString := fmt.Sprintf("https://www.siriusxm.com/metadata/pdt/en-us/json/channels/%v/timestamp/%s", channel, timeString)
	resp, err := http.Get(reqString)

	checkErr(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	checkErr(err)

	json.Unmarshal(skipRoot(body), &channelResponse)

	nowPlaying := NowPlaying{
		Song:     channelResponse.MetaData.CurrentEvent.Song.Name,
		Artist:   channelResponse.MetaData.CurrentEvent.Artists.Name,
		Album:    channelResponse.MetaData.CurrentEvent.Song.Album.Name,
		Composer: channelResponse.MetaData.CurrentEvent.Song.Composer,
		Details:  channelResponse.MetaData.CurrentEvent,
	}

	return nowPlaying
}
