package siriusxm

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type logos struct {
	Height       int
	ResourceType string
	URL          string
	Width        int
}

type busCodes struct {
	Mobile    string
	SiriusIP  string
	SiriusSat string
	XMIP      string
	XMSat     string
}

type channels struct {
	BusCodes          busCodes
	ChannelKey        interface{}
	ContentID         interface{}
	Description       string
	DisplayName       string
	FullDescription   string
	GeoRestriction    int
	IsAvailable       bool
	IsBizMature       bool
	IsMature          bool
	IsMyXM            bool
	Logos             [6]logos
	Name              string
	Order             int
	RelatedContentIDs string
	Replay            int
	Select            int
	ServiceID         int
	ServiceTypes      []string
	Shadow            int
	SiriusChannelNo   int
	SpanishContent    bool
	Type              int
	URL               string
	XMChannelNo       int
	XMServiceID       int
}

type genres struct {
	Channels    []channels
	Description string
	Key         string
	Name        string
	Order       int
	URL         string
	XMSatID     int
}

type categories struct {
	Description string
	Genres      []genres
	Key         string
	Name        string
	Order       int
	URL         string
}

type lineup struct {
	Categories      [7]categories
	ConsumerID      string
	LineupID        int
	PackageID       string
	UpsellPackageID string
}

type ChannelResponse struct {
	LastModified string
	Messages     messages
	Status       int
	Lineup       lineup
}

type ChannelDetails struct {
	Name            string
	XMChannelNo     int
	SiriusChannelNo int
	ChannelKey      interface{}
	ContentID       interface{}
	DisplayName     string
	Genre           string
	Category        string
}

type ChannelLineup struct {
	Channels     []ChannelDetails
	LastModified string
}

func (lineup ChannelLineup) FindChannel(channel interface{}) (ChannelDetails, error) {
	var c ChannelDetails
	var err error = nil
	for _, v := range lineup.Channels {
		if reflect.TypeOf(v.ChannelKey).String() != "string" {
			v.ChannelKey = fmt.Sprintf("%d", v.ChannelKey)
		}
		if strings.ToLower(v.Name) == strings.ToLower(channel.(string)) ||
			strings.ToLower(v.ChannelKey.(string)) == strings.ToLower(channel.(string)) ||
			v.SiriusChannelNo == channel ||
			v.XMChannelNo == channel {
			c = v
			break
		}
	}

	if c.ChannelKey == nil {
		err = errors.New("Could not find key by name")
	}

	return c, err
}

func (lineup ChannelLineup) FindChannelByNumber(channel int) (ChannelDetails, error) {
	var c ChannelDetails
	var err error = nil
	for _, v := range lineup.Channels {
		if v.SiriusChannelNo == channel {
			c = v
			break
		}
	}

	if c.ChannelKey == nil {
		err = errors.New("Could not find key by name")
	}

	return c, err
}

func GetChannelLineup() (ChannelLineup, ChannelResponse, error) {
	var response ChannelResponse
	var lineup ChannelLineup
	resp, err := http.Get("https://www.siriusxm.com/userservices/cl/en-us/json/lineup/350/client/ump")

	if err != nil {
		fmt.Println(err)
		if resp == nil {
			return lineup, response, err
		}

	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {

	}

	json.Unmarshal(skipRoot(body), &response)
	lineup = response.createLineup()
	return lineup, response, err
}

func (response ChannelResponse) createLineup() ChannelLineup {
	var lineup ChannelLineup
	lineup.LastModified = response.LastModified

	// Iterate through categories, genres and channels
	for _, cat := range response.Lineup.Categories {
		for _, genre := range cat.Genres {
			for _, channel := range genre.Channels {
				lineup.Channels = append(lineup.Channels, ChannelDetails{
					Name:            channel.Name,
					XMChannelNo:     channel.XMChannelNo,
					SiriusChannelNo: channel.SiriusChannelNo,
					ChannelKey:      channel.ChannelKey,
					ContentID:       channel.ContentID,
					DisplayName:     channel.DisplayName,
					Genre:           genre.Name,
					Category:        cat.Name})
			}
		}
	}

	return lineup
}
