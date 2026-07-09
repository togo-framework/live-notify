// Package notify bridges togo `live` conversations to the togo `notifications`
// system, so an agent's reply can be delivered out over ANY registered
// notifications channel — slack, discord, webpush, fcm, pusher, mail, broadcast —
// without a per-provider live driver. Blank-import this alongside the
// notifications channel driver(s) you want (e.g. togo-framework/notifications-slack).
//
// A live conversation whose `channel` is "slack" or "discord" delivers through
// that notifications channel; the "notify" channel fans out to the comma list in
// LIVE_NOTIFY_CHANNELS (e.g. "slack,discord").
package notify

import (
	"context"
	"os"
	"strings"

	"github.com/togo-framework/live"
	"github.com/togo-framework/notifications"
	"github.com/togo-framework/togo"
)

func init() {
	for _, name := range []string{"slack", "discord", "notify"} {
		n := name
		live.RegisterChannel(n, func(k *togo.Kernel) live.Channel {
			return &bridge{k: k, channel: n}
		})
	}
}

type bridge struct {
	k       *togo.Kernel
	channel string
}

func (b *bridge) Name() string { return b.channel }

func (b *bridge) Deliver(ctx context.Context, conv live.Conversation, msg live.Message) error {
	svc, ok := notifications.FromKernel(b.k)
	if !ok {
		return nil // notifications plugin not installed → no-op
	}
	via := []string{b.channel}
	if b.channel == "notify" {
		via = splitCSV(os.Getenv("LIVE_NOTIFY_CHANNELS"))
	}
	if len(via) == 0 {
		return nil
	}
	title := conv.Title
	if title == "" {
		title = "New reply"
	}
	return svc.Send(ctx, recipient(conv.ExternalRef), &push{title: title, body: msg.Body, via: via})
}

// push is a notifications PushNotification carrying the agent reply.
type push struct {
	title, body string
	via         []string
}

func (p *push) Via(notifications.Notifiable) []string { return p.via }
func (p *push) ToPush(notifications.Notifiable) notifications.PushMessage {
	return notifications.PushMessage{Title: p.title, Body: p.body}
}

// recipient routes by the conversation's external ref. Slack/Discord read their
// webhook from env and ignore routing; push-token channels use RoutePushTokens.
type recipient string

func (r recipient) RouteID() string    { return string(r) }
func (r recipient) RouteEmail() string { return "" }
func (r recipient) RoutePushTokens() []string {
	if r == "" {
		return nil
	}
	return []string{string(r)}
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
