package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zeebo/bencode"
)

var (
	path = flag.String("path", "", "path to directory or file. if empty - current directory will be used.")
)

// Torrent represents a .torrent file structure
type Torrent struct {
	Announce string                 `bencode:"announce"` // Tracker URL
	Info     map[string]interface{} `bencode:"info"`     // Torrent file info section
	Other    map[string]interface{} `bencode:",inline"`  // Other fields (e.g., announce-list, etc.)
}

func main() {
	flag.Parse()
	targetPath := *path
	if targetPath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("can't get current directory: %s", err)
		}
		targetPath = pwd
	}

	items, err := os.ReadDir(targetPath)
	if err != nil {
		log.Fatal("Error reading directory:", err.Error())
	}
	for _, item := range items {
		if item.IsDir() || !strings.HasSuffix(item.Name(), ".torrent") {
			continue
		}
		cleanFile(targetPath + "/" + item.Name())
	}
}

func cleanFile(torrentFile string) {
	fileData, err := os.ReadFile(torrentFile)
	if err != nil {
		log.Fatal("Error reading torrent file:", err.Error())
	}

	var torrent Torrent
	if err = bencode.DecodeBytes(fileData, &torrent); err != nil {
		log.Fatal("Error decoding torrent file:", err.Error())
	}

	torrent.Announce = ""
	delete(torrent.Other, "announce-list")

	modifiedTorrentData, err := bencode.EncodeBytes(torrent)
	if err != nil {
		log.Fatal("Error encoding modified torrent:", err.Error())
	}

	if err = os.WriteFile(torrentFile, modifiedTorrentData, 0644); err != nil {
		log.Fatal("Error writing modified torrent file:", err.Error())
	}
	fmt.Println("Modified torrent file replaced: " + torrentFile)
}
