package handlers

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/json"
)

var (
	durationMap = map[string]time.Duration{
		"minute": time.Minute,
		"hour":   time.Hour,
		"day":    Day,
		"week":   Week,
	}
)

const (
	Day  = 24 * time.Hour
	Week = 7 * Day
)

func (h *Handler) HandleTimeout(data discord.SlashCommandInteractionData, event *handler.CommandEvent) error {
	messageBuilder := discord.NewMessageCreateBuilder().SetEphemeral(true)
	if _, ok := data.OptMember("member"); !ok {
		return event.CreateMessage(messageBuilder.
			SetContent("User is not a member of this server.").
			Build())
	}
	now := time.Now()
	length := data.Int("length")
	unit := durationMap[data.String("unit")]
	expiry := now.Add(time.Duration(length) * unit).Add(time.Duration(-3) * time.Second) // wtf even
	if expiry.After(now.Add(28 * Day)) {
		return event.CreateMessage(messageBuilder.
			SetContent("Length of the timeout exceeds 28 days.").
			Build())
	}
	reason := "Moderator: " + event.User().Tag()
	if optReason, ok := data.OptString("reason"); ok {
		reason = optReason + " | " + reason
	}
	userID := data.User("member").ID
	client := event.Client().Rest()
	_, err := client.UpdateMember(*event.GuildID(), userID, discord.MemberUpdate{
		CommunicationDisabledUntil: json.NewNullablePtr(expiry),
	}, rest.WithReason(reason))
	if err != nil {
		return err
	}
	return event.CreateMessage(messageBuilder.
		SetContentf("Member <@%d> has been timed out until %s.", userID, discord.TimestampStyleLongDateTime.FormatTime(expiry)).
		Build())
}
