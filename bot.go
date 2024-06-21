package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"vs-helper/handlers"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/lmittmann/tint"
)

func main() {
	logger := tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(logger))

	slog.Info("starting the bot...", slog.String("disgo.version", disgo.Version))

	client, err := disgo.New(os.Getenv("VIRTUALSTREETS_HELPER_TOKEN"),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentsNone)),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagsNone)),
		bot.WithEventListeners(handlers.NewHandler()))
	if err != nil {
		panic(err)
	}

	defer client.Close(context.TODO())

	if err := client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}

	slog.Info("VirtualStreets Helper is now running.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}
