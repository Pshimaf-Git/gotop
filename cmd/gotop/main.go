package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Pshimaf-Git/gotop/internal/config"
	"github.com/Pshimaf-Git/gotop/internal/process"
	"github.com/Pshimaf-Git/gotop/internal/ui"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	cfg, err := config.Load(config.FetchConfigPath())
	if err != nil {
		log.Fatalf("config parse failed: %s", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Fprintf(os.Stderr, "Stoped by signal %s\n", sig.String())
		cancel()
	}()

	app := ui.NewApp()

	processInfoChan := process.GetProcessInfoChan(ctx, cfg.RefreshInterval)

	err = app.Run(ctx, processInfoChan, cfg)
	if err != nil {
		log.Fatalf("UI appication failed: %s", err.Error())
	}

	log.Println("gotop finished")
}
