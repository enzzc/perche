package main

import (
	"encoding/json"
	"net"
	"net/http"
	"time"
)

type RadioFranceStatus struct {
	DelayToRefresh int `json:"delayToRefresh"`
	Now            struct {
		NowTime    int     `json:"nowTime"`
		NowPercent float64 `json:"nowPercent"`
		Artist     string  `json:"secondLine"`
		Song       struct {
			Title        string   `json:"title"`
			Interpreters []string `json:"interpreters"`
			Year         int      `json:"year"`
		} `json:"song"`
		Cover struct {
			Src string `json:"src"`
		} `json:"cover"`
	} `json:"now"`
}

func minDuration(d1, d2 time.Duration) time.Duration {
	if d1 < d2 {
		return d1
	}
	return d2
}

func MonitorAPI() chan *RadioFranceStatus {
	ch := make(chan *RadioFranceStatus)
	go func() {
		defer close(ch)
		c := &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DialContext:           (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
				TLSHandshakeTimeout:   3 * time.Second,
				ResponseHeaderTimeout: 3 * time.Second,
			},
		}
		for {
			resp, err := c.Get(FipAPIURL)
			if err != nil {
				panic(err)
			}
			decoder := json.NewDecoder(resp.Body)
			var status RadioFranceStatus
			decoder.Decode(&status)
			ch <- &status
			time.Sleep(minDuration(10*time.Second, time.Duration(status.DelayToRefresh)*time.Millisecond))
		}
	}()
	return ch
}
