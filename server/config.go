package server

import (
	"log/slog"
)

// Config provides configuration for Server.
type Config struct {

	// ChanPacket sends packet data.
	ChanPacket chan<- []byte

	// Logger can be used to capture log messages.
	Logger *slog.Logger
}
