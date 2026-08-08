package tflog

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type tflogWriter struct {
	ctx context.Context
}

func (w *tflogWriter) Write(p []byte) (n int, err error) {
	// Strip trailing newline added by log.Logger
	msg := string(p)
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}

	// Route to tflog
	tflog.Debug(w.ctx, msg)
	return len(p), nil
}

func StandardLogger(ctx context.Context) *log.Logger {
	// Pass 0 for flags if you don't want stdlib log to append duplicate timestamps/date prefixes
	return log.New(&tflogWriter{ctx: ctx}, "", 0)
}
