package main

import (
	"log"

	"github.com/anacrolix/torrent"
)

var client *torrent.Client

// TCStart starts a torrent client and sets videos as the download dir
func TCStart() {
	var err error

	log.Print("Starting torrent client")

	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = "videos"

	client, err = torrent.NewClient(clientConfig)
	if err != nil {
		log.Fatal(err)
	}
}

// TCDownload downloads the video from a torrent file
func TCDownload(magnet string) error {
	torrent, err := client.AddMagnet(magnet)
	if err != nil {
		return err
	}
	<-torrent.GotInfo() // do this so things don't whine

	files := torrent.Files()
	if len(files) == 1 { // torrent is a single file so we download it all
		torrent.DownloadAll()
	} else { // figure out which file is the video file
		for file := range files {
			println(file)
		}
	}
	client.WaitAll()
	return nil
}
