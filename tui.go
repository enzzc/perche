package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func drawTimer(s tcell.Screen, status *RadioFranceStatus) {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)
	offsetX := 0
	offsetY := 0
	percent := status.Now.NowPercent

	for {
		var bar strings.Builder
		for i := 0; i < int(percent); i++ {
			bar.WriteString("-")
		}
		var empty strings.Builder
		for i := 0; i < 100-int(percent); i++ {
			empty.WriteString(" ")
		}
		str := "[" + bar.String() + empty.String() + "]"

		writeLine(s, offsetY+13, offsetX+6, str, defStyle)
		s.Show()
		time.Sleep(1 * time.Second)
		percent++
	}
}

func writeLine(s tcell.Screen, row, col int, text string, style tcell.Style) {
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
	}
}

func draw(s tcell.Screen, status *RadioFranceStatus) {
	go drawTimer(s, status)
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	offsetX := 0
	offsetY := 0

	writeLine(s, offsetY+0, offsetX+0, "Prev.", defStyle)
	writeLine(s, offsetY+0, offsetX+80, "save previous song [p]", defStyle)
	writeLine(s, offsetY+2, offsetX+0, "TITLE_PREV", defStyle)
	writeLine(s, offsetY+3, offsetX+0, "ARTIST_PREV", defStyle)

	writeLine(s, offsetY+8, offsetX+0, "Now Playing", defStyle)
	writeLine(s, offsetY+10, offsetX+0, status.Now.Song.Title, defStyle)

	writeLine(s, offsetY+10, offsetX+80, "save current song [s]", defStyle)
	writeLine(s, offsetY+11, offsetX+0, status.Now.Artist, defStyle)
	writeLine(s, offsetY+13, offsetX+0, "__:__", defStyle)
	writeLine(s, offsetY+13, offsetX+6, "[--------------------------]", defStyle)
	writeLine(s, offsetY+18, offsetX+0, "Up Next", defStyle)
	writeLine(s, offsetY+20, offsetX+0, "Artist", defStyle)
	writeLine(s, offsetY+21, offsetX+0, "Song", defStyle)
	s.Show()
}

type TUICommunication struct {
	Statuses  chan *RadioFranceStatus
	ExitNotif chan struct{}
}

func RunTUI(ctx context.Context) TUICommunication {
	statuses := make(chan *RadioFranceStatus)
	exitNotif := make(chan struct{})

	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}

	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventResize:
				w, h := ev.Size()
				log.Println(w, h)
				s.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					s.Fini()
					exitNotif <- struct{}{}
				}
			}
		}
	}()

	go func() {
		defer close(statuses)
		for {
			select {
			case status := <-statuses:
				draw(s, status)
			case <-ctx.Done():
				return
			}
		}
	}()

	return TUICommunication{statuses, exitNotif}
}
