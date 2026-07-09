# live-notify — usage

`live-notify` is a single bridge between the togo [`live`](https://to-go.dev/plugins/live)
system (agent conversations) and the togo [`notifications`](https://to-go.dev/plugins/notifications)
system (external delivery). Instead of writing a per-provider `live` egress
driver, you point a conversation's channel at the bridge and let an existing
`notifications-*` driver do the delivery.

## What the bridge registers

On import, `live-notify` registers three `live` egress channels via
`live.RegisterChannel`:

| live channel | Delivers via |
|---|---|
| `slack`   | the `slack` notifications channel |
| `discord` | the `discord` notifications channel |
| `notify`  | fans out to the channels named in `LIVE_NOTIFY_CHANNELS` (CSV) |

When the `live` loop calls `Deliver` for an agent reply, the bridge looks up the
`notifications` service on the kernel and calls `Send` with a push notification
carrying the conversation title (or `"New reply"` if untitled) and the reply
body. If the `notifications` plugin isn't installed, `Deliver` is a no-op.

## 1. Install the bridge and a notifications driver

The bridge only routes — the actual transport is a `notifications-*` driver, so
you must blank-import the matching driver too:

```go
import (
    _ "github.com/togo-framework/live-notify"          // slack / discord / notify live channels
    _ "github.com/togo-framework/notifications-slack"   // the Slack notifications transport
)
```

Add more drivers as needed — `notifications-discord`, `notifications-webpush`,
`notifications-fcm`, `notifications-pusher`, `notifications-mail`.

## 2. Configure the driver's env

Each `notifications-*` driver reads its own config from the environment — the
bridge passes none of this through, it just selects the channel. For example:

```bash
export SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...
export DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
```

## 3. Set the conversation channel

Set a conversation's `channel` to the name of a bridge channel:

- `slack` → every agent reply is delivered through the `slack` notifications channel.
- `discord` → delivered through the `discord` notifications channel.
- `notify` → fans out to each channel in `LIVE_NOTIFY_CHANNELS`.

For the fan-out channel, set the CSV env:

```bash
export LIVE_NOTIFY_CHANNELS=slack,discord
```

## End-to-end example

1. Install the bridge and the Slack driver, and set `SLACK_WEBHOOK_URL`.
2. Create a `live` conversation with `channel = "slack"`.
3. A user message arrives; the agent produces a reply.
4. The `live` loop calls the bridge's `Deliver`; the bridge sends a push
   notification (title = conversation title, body = reply text) to the `slack`
   notifications channel.
5. `notifications-slack` posts it to your Slack incoming webhook — the reply
   lands in Slack.

To deliver the same reply to both Slack and Discord, set the conversation
channel to `notify` and `LIVE_NOTIFY_CHANNELS=slack,discord` (with both
`notifications-slack` and `notifications-discord` imported and their webhook env
set).

## Routing note

The bridge routes recipients by the conversation's `ExternalRef`. Webhook-based
channels (Slack, Discord) read their target from env and ignore per-recipient
routing; push-token channels (FCM, web-push, Pusher) use `ExternalRef` as the
push token.
