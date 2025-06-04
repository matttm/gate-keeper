package main

import (
	"fmt"
	"time"
)

type GateSpectator struct {
	_ticker     *time.Ticker
	gatesUpdate chan []*Gate
	gates       []*Gate
}

func NewGateSpectator(c *GateConfig, year int) *GateSpectator {
	g := GateSpectator{
		_ticker:     time.NewTicker(1 * time.Second), // Ticks every 1 second
		gatesUpdate: make(chan []*Gate),
	}
	go func() {
		fmt.Println("GateSpectator began spectating")
		for range g._ticker.C {
			g.gatesUpdate <- selectAllGates(c, year)
		}
	}()
	return &g
}

func (g *GateSpectator) Shutdown() {
	g._ticker.Stop()
	close(g.gatesUpdate)
}
