package main

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

type MpvPlayer struct {
	UnixSock string
}

func (p *MpvPlayer) Play(uri string) error {
	conn, err := net.Dial("unix", p.UnixSock)
	if err != nil {
		return err
	}
	defer conn.Close()
	encoder := json.NewEncoder(conn)
	encoder.Encode(struct {
		Command []string `json:"command"`
	}{
		[]string{"loadfile", uri},
	})
	buf := make([]byte, 1024)
	_, err = conn.Read(buf[:])
	if err != nil {
		return err
	}
	return err
}

func (p *MpvPlayer) Close() {
	os.Remove(p.UnixSock)
}

func NewMpvPlayer(ctx context.Context) (*MpvPlayer, error) {
	socket := "/tmp/mpv_" + uuid.New().String() + ".sock"
	cmd := exec.CommandContext(ctx,
		"mpv",
		"--idle",
		"--input-ipc-server="+socket,
	)
	go cmd.Run()
	time.Sleep(500 * time.Millisecond) // FIXME: smell
	return &MpvPlayer{UnixSock: socket}, nil
}
