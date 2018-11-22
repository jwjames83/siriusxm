# Sirius XM Scraper

This little library can retrieve the currently playing song information with a provided station. All you need to do is
get the channel you want

``` go
package main

import (
	"fmt"
	"github.com/jwjames83/siriusxm"
)

func main() {
	channel, err := siriusxm.GetChannelLineup().FindChannelByNumber(36)

	if err != nil {
		fmt.Println(err)
	}

	nowPlaying := siriusxm.GetNowPlaying((channel.ChannelKey).(string))
	fmt.Printf("Now playing on %s: %s - %s", channel.Name, nowPlaying.Song, nowPlaying.Artist)
}
