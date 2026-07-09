---
name: live-notify
description: Bridge togo live agent replies out through the notifications system (Slack, Discord, web-push, FCM, Pusher, mail) — one bridge, no per-provider live driver.
---

# togo live-notify

Use this skill to deliver a togo `live` agent's replies through the togo
`notifications` system. The bridge registers three `live` egress channels:
`slack`, `discord`, and `notify` (fan-out).

## Wire it up
```go
import (
    _ "github.com/togo-framework/live-notify"         // slack / discord / notify live channels
    _ "github.com/togo-framework/notifications-slack"  // the matching notifications transport
)
```
Set the driver's env, e.g. `SLACK_WEBHOOK_URL` / `DISCORD_WEBHOOK_URL`.

## Route a conversation
- Conversation `channel = "slack"` or `"discord"` → deliver through that notifications channel.
- Conversation `channel = "notify"` → fan out to `LIVE_NOTIFY_CHANNELS` (CSV), e.g. `LIVE_NOTIFY_CHANNELS=slack,discord`.

## Notes
- Egress only: the bridge pushes the agent reply (title = conversation title or "New reply", body = reply text) to the notifications channel; it does not ingest inbound messages or store in-app.
- Needs the `notifications` plugin + the matching `notifications-*` driver installed; without them delivery is a no-op.
- Works with any notifications channel (`-slack`, `-discord`, `-webpush`, `-fcm`, `-pusher`, `-mail`) with no extra live code.
- Prefer a native `live` driver (e.g. `live-whatsapp`) when you need two-way/threaded delivery.
