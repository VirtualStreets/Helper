package handlers

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/lmittmann/tint"
)

func NewHandler() *Handler {
	mux := handler.New()
	mux.Error(func(e *handler.InteractionEvent, err error) {
		i := e.Interaction.(discord.ApplicationCommandInteraction)
		slog.Error("error while handling a command", slog.String("command.name", i.Data.CommandName()), tint.Err(err))
		_ = e.Respond(discord.InteractionResponseTypeCreateMessage, discord.NewMessageCreateBuilder().
			SetContentf("There was an error while handling the command: %v", err).
			SetEphemeral(true).
			Build())
	})
	h := &Handler{
		Router: handler.New(),
	}
	h.SlashCommand("/timeout", h.HandleTimeout)
	return h
}

type Handler struct {
	handler.Router
}
