package main

import (
	"errors"
	"io"
	"log"
	"time"

	"github.com/anacrolix/torrent"
)

var client *torrent.Client

// openTorrents is a map containing readers to each torrent being downloaded
var openTorrents = make(map[string]io.Reader)

// TClientStart starts a torrent client and sets videos as the download dir
func TClientStart() {
	var err error

	log.Print("Starting torrent client")

	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = "videos"

	client, err = torrent.NewClient(clientConfig)
	if err != nil {
		log.Fatal(err)
	}
}

// TClientAdd adds the magnet to the torrent client.
// returns an error if the magnet is invalid or if torrent.Info times out
// or if the magnet is already loaded into the client
func TClientAdd(magnet string, index int) (*torrent.Torrent, error) {
	const timeout = 10 * time.Second
	torrent, err := client.AddMagnet(magnet)
	if err != nil {
		return torrent, err
	}

	// if torrent info hasn't arrived by timeout return errors
	select {
	case <-time.After(timeout):
		torrent.Drop()
		return torrent, errors.New("torrent.GotInfo: timeout for torrent " + magnet)
	case <-torrent.GotInfo():
		torrentHash := torrent.InfoHash().AsString()
		torrentFile := torrent.Files()[index]
		if openTorrents[torrentHash] != nil {
			return torrent, errors.New("Torrent already loaded")
		}
		openTorrents[torrentHash] = torrentFile.NewReader()
	}

	return torrent, nil
}
