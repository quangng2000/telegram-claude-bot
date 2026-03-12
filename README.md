# Telegram Claude Bot

A Go-based Telegram bot that bridges messages to [Claude Code](https://claude.ai/claude-code) CLI via PTY. Acts as a personal AI assistant accessible from Telegram with session memory, auto-expiry, and chat logging.

## Features

- **Claude Code integration** — runs Claude CLI with full tool access (bash, file editing, etc.)
- **Session memory** — maintains conversation context per chat using `--continue`
- **Auto-expiry** — sessions expire after configurable inactivity (default 24h)
- **Summary carry-over** — generates a conversation summary before expiring, injected into next session
- **Chat logging** — all messages saved to `sessions/<chatId>/history.log`
- **Telegram-friendly** — ANSI stripping, message chunking, typing indicators
- **Per-chat concurrency** — mutex prevents overlapping Claude calls per user

## Prerequisites

- **Go 1.21+** — [install](https://go.dev/doc/install)
- **Claude Code CLI** — [install](https://docs.anthropic.com/en/docs/claude-code)
- **Telegram Bot Token** — create via [@BotFather](https://t.me/BotFather)

## Project Structure

```
├── cmd/bot/
│   └── main.go                 # Entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Environment-based configuration
│   ├── claude/
│   │   ├── runner.go           # Claude CLI execution
│   │   ├── pty.go              # PTY management and timeout
│   │   └── ansi.go             # ANSI escape code stripping
│   ├── session/
│   │   ├── manager.go          # Session CRUD and directory management
│   │   ├── logger.go           # Chat history logging
│   │   └── expiry.go           # Session expiry and summary generation
│   └── telegram/
│       ├── handler.go          # Message routing and chat flow
│       ├── commands.go         # Bot commands (/start, /reset, /history)
│       └── sender.go           # Message sending and chunking
├── CLAUDE.md                   # Claude system prompt (local only, gitignored)
├── go.mod
└── go.sum
```

## Setup

### 1. Clone and build

```bash
git clone https://github.com/quangng2000/telegram-claude-bot.git
cd telegram-claude-bot
go build -o telegram-claude-bot ./cmd/bot
```

### 2. Create CLAUDE.md

Create a `CLAUDE.md` file in the project root. This defines the system prompt for Claude — customize it for your use case:

```markdown
You are a personal AI assistant accessible via Telegram.

You can help with:
- General questions and conversation
- Coding and technical help
- System administration tasks

IMPORTANT FORMATTING RULES (responses go to Telegram):
- Keep responses short and concise
- DO NOT use markdown tables
- Use bullet points and numbered lists
- Use code blocks for command output
```

### 3. Configure environment variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `TELEGRAM_BOT_TOKEN` | Yes | — | Token from @BotFather |
| `BOT_DIR` | No | Current directory | Working directory for the bot |
| `SESSION_EXPIRY_HOURS` | No | `24` | Hours of inactivity before session expires |
| `CLAUDE_MODEL` | No | `sonnet` | Claude model to use |

### 4. Run

```bash
TELEGRAM_BOT_TOKEN="your-token-here" ./telegram-claude-bot
```

## Running as a systemd service

Create `/etc/systemd/system/telegram-claude-bot.service`:

```ini
[Unit]
Description=Telegram Claude Bot
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=your-username
WorkingDirectory=/path/to/telegram-claude-bot
ExecStart=/path/to/telegram-claude-bot/telegram-claude-bot
Environment=TELEGRAM_BOT_TOKEN=your-token-here
Environment=HOME=/home/your-username
Environment=PATH=/home/your-username/.local/bin:/usr/local/bin:/usr/bin:/bin
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Then enable and start:

```bash
sudo cp telegram-claude-bot.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable telegram-claude-bot
sudo systemctl start telegram-claude-bot
```

Manage the service:

```bash
sudo systemctl status telegram-claude-bot    # check status
sudo systemctl restart telegram-claude-bot   # restart
sudo journalctl -u telegram-claude-bot -f    # view logs
```

## Bot Commands

| Command | Description |
|---------|-------------|
| `/start` | Welcome message with capabilities |
| `/reset` | Clear conversation history and summary |
| `/history` | Show last 20 chat messages |

## How It Works

1. User sends a message on Telegram
2. Bot spawns `claude -p --continue` in a PTY with a per-chat working directory
3. Claude CLI runs with full tool access (bash, file editing, etc.)
4. PTY output is captured, ANSI codes are stripped, and the response is sent back
5. Session context persists via `--continue` flag and per-chat directories
6. After inactivity, sessions are summarized and the summary is injected into future conversations

## License

MIT
