package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/japannext/cosi-powerscale/pkg/config"
	"github.com/japannext/cosi-powerscale/pkg/driver"
	log "k8s.io/klog/v2"
)

func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.New()
	d, err := driver.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.InfoS("Signal received", "type", sig)
		cancel()
		os.Exit(1)
	}()

	d.Run(ctx)
}
