package main

import "time"

type GateSpectator struct {
	ticker *time.Ticker
	gates []Gate
}
func NewGateSpectator() *GateSpectator {
	g := GateSpectator{
		ticker: time.NewTicker(1 * time.Second), // Ticks every 1 second
	}
	go func () {
	}
	return &g
}

func (g *GateSpectator) Shutdown() {
	g.ticker.Stop()
}
