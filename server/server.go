package server

import (
	"context"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/thoas/picfit"
	"github.com/thoas/picfit/config"
)

func New(ctx context.Context, cfg *config.Config) (*HTTPServer, error) {
	processor, err := picfit.NewProcessor(ctx, cfg)
	if err != nil {
		return nil, err
	}

	httpServer, err := NewHTTPServer(cfg, processor)
	if err != nil {
		return nil, err
	}

	return httpServer, nil
}

// Run runs the application and launch servers
func Run(ctx context.Context, path string) error {
	cfg, err := config.Load(path)
	if err != nil {
		return err
	}

	server, err := New(ctx, cfg)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := server.Run(ctx); err != nil {
		return err
	}

	go func() {
		for range time.Tick(time.Duration(cfg.Options.FreeMemoryInterval) * time.Second) {
			debug.FreeOSMemory()
		}
	}()

	select { // nolint:gosimple
	case <-ctx.Done():
		stop()
	}

	return nil
}
