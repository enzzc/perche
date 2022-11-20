package main

import (
	"fmt"
	"os"
	"path"
	"time"
)

func SaveTrack(status *RadioFranceStatus) error {
	line := fmt.Sprintf("%s; %s; %s\n", time.Now(), status.Now.Artist, status.Now.Song.Title)
	home, _ := os.UserHomeDir()
	filename := path.Join(home, DefaultSavedTracksFilename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(line))
	return err
}
