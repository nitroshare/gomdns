package vtime

import (
	"testing"
	"time"
)

func TestUnmockedTicker(t *testing.T) {
	ticker := NewTicker(1 * time.Second)
	defer ticker.Stop()
	ticker.Reset(1 * time.Second)
}

func TestTicker(t *testing.T) {
	Mock()
	defer Unmock()
	ticker := NewTicker(2 * time.Second)
	defer ticker.Stop()
	Advance(3 * time.Second)
	<-ticker.C
	Advance(2 * time.Second)
	<-ticker.C
}

func TestTickerReset(t *testing.T) {
	Mock()
	defer Unmock()
	ticker := NewTicker(4 * time.Second)
	defer ticker.Stop()
	ticker.Reset(2 * time.Second)
	Advance(3 * time.Second)
	<-ticker.C
	Advance(2 * time.Second)
	<-ticker.C
}
