<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/live-notify</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/live-notify"><img src="https://pkg.go.dev/badge/github.com/togo-framework/live-notify.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Bridge <a href="https://to-go.dev">togo</a> live agent replies out through the notifications system — one bridge, zero per-provider live drivers.</strong></p>
</div>

## Install

```bash
togo install togo-framework/live-notify
```

## Usage

`live-notify` registers three [`live`](https://to-go.dev/plugins/live) egress channels — **`slack`**, **`discord`**, and **`notify`** — that route an agent's reply through the togo [`notifications`](https://to-go.dev/plugins/notifications) system. A conversation on one of these channels is delivered by the existing `notifications-*` driver, so there's no per-provider live code to write.

Blank-import this bridge alongside the notifications channel driver(s) you want:

```go
import (
    _ "github.com/togo-framework/live-notify"          // registers slack / discord / notify live channels
    _ "github.com/togo-framework/notifications-slack"   // the matching notifications driver
)
```

Then set a conversation's `channel` and provide the driver's webhook env:

- Set a conversation's `channel` to **`"slack"`** or **`"discord"`** to deliver each agent reply through that notifications channel.
- Set it to **`"notify"`** to fan out to every channel listed in `LIVE_NOTIFY_CHANNELS` (CSV), e.g. `LIVE_NOTIFY_CHANNELS=slack,discord`.
- Configure the underlying driver as usual — `SLACK_WEBHOOK_URL` for `notifications-slack`, `DISCORD_WEBHOOK_URL` for `notifications-discord`, etc.

Any registered notifications channel works the same way — `notifications-webpush`, `-fcm`, `-pusher`, `-mail` — with no extra live code. If the `notifications` plugin isn't installed, delivery is a no-op.

See [`docs/usage.md`](docs/usage.md) for an end-to-end example.

---

<div align="center"><sub>💎 Premium sponsors — <b>ID8 Media</b> · <b>One Studio</b></sub></div>
