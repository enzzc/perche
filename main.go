package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const FipStreamURL = "https://stream.radiofrance.fr/fip/fip_hifi.m3u8?id=radiofrance"
const FipAPIURL = "https://www.radiofrance.fr/api/v2.0/stations/fip/webradios/fip"
const DefaultSavedTracksFilename = ".fip-tracks.txt"

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		fmt.Println("Stop subprocess mpv")
		cancel()
		time.Sleep(250 * time.Millisecond) // grace period
	}()

	mpv, err := NewMpvPlayer(ctx)
	if err != nil {
		return
	}
	defer func() {
		fmt.Println("remove", mpv.UnixSock)
		mpv.Close()
	}()

	//mpv.Play(FipStreamURL)

	tuiChans := RunTUI(ctx)
	statusNotif := MonitorAPI()

	for {
		select {
		case <-sigs:
			fmt.Println("Bye.")
			return
		case <-tuiChans.ExitNotif:
			fmt.Println("TUI returned. Bye.")
			return
		case s := <-statusNotif:
			tuiChans.Statuses <- s
		}
	}
}
