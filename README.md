# ☂ termrain

Words fall from the sky. Type them before they hit the ground.

A small terminal typing game in pure Go. Standard library only, no dependencies.

```
☂ TERMRAIN    score 340    best 900          lvl 2  48 wpm  combo 4  ♥♥♡
────────────────────────────────────────────────────────────────────────
        '          goroutine              '
   defer      '                 mutex
                    '      channel           '
                                     '          select
  +80 x2        context
▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
 ❯ chan▏
```

## Play

```bash
go install github.com/avichal-08/termrain@latest
termrain
```

Or from a clone: `go run .`

`go install` puts the binary in `$(go env GOPATH)/bin` (usually `~/go/bin`). If `termrain` isn't found, add that folder to your `PATH`.

termrain opens on a prompt. Press Enter to play, or type a command:

```
╭────────────────────────────────────────────────────────────────╮
│ ☂ termrain  v0.3.0                                             │
│                                                                │
│   type the rain before it lands.                               │
│   best   1840 pts · 72 wpm · x14                               │
│   words  615 · saved to %AppData%\termrain\best.json           │
╰────────────────────────────────────────────────────────────────╯

╭────────────────────────────────────────────────────────────────╮
│ ❯ /s▏                                                          │
╰────────────────────────────────────────────────────────────────╯
  /stats  show your personal best
```

| Command  | What it does |
|----------|--------------|
| `/play`  | start a round (same as pressing Enter) |
| `/stats` | show your personal best |
| `/reset` | clear your saved best score (asks twice) |
| `/help`  | how to play |
| `/quit`  | exit |

`Tab` completes a command. `termrain --help` and `termrain --version` work from the shell too.

- Start typing any falling word. The first letter locks onto it.
- Finish a word to clear it. Every 5 words in a row bumps your multiplier.
- A word that reaches the ground costs a life. You have 3.
- Words get longer and faster the longer you last.

Keys: `Backspace` drops the current word, `Esc` goes back to the prompt, `Ctrl+C` quits. After a game over: `r` plays again, `m` opens the prompt, `q` quits.

Around 600 words: everyday words like animals, food, weather and places, mixed with programming terms.

## Share your run

Press `c` on the game over screen to copy a result card:

```
☂ termrain · 1840 pts 🏆 new best
⌨️ 72 wpm  🎯 98%  🔥 x14  ⏱ 1:32
github.com/avichal-08/termrain
```

Copying uses OSC 52, which works in Windows Terminal, iTerm2, kitty, WezTerm and most modern terminals. The card is also printed when you quit, so you can always copy it from there.

## Best score

Your best score, WPM and combo are saved locally:

- Windows: `%AppData%\termrain\best.json`
- macOS: `~/Library/Application Support/termrain/best.json`
- Linux: `~/.config/termrain/best.json`

Reset it with `/reset` at the prompt, or `termrain --reset` from the shell.

## Notes

- Raw keyboard input without libraries: `SetConsoleMode` from `kernel32.dll` on Windows, `stty raw` on macOS/Linux.
- Keystrokes are read on a separate goroutine and fed through a channel; the game loop selects between input and a 30 fps ticker.
- Frames are built in a cell buffer and written once per tick to the alternate screen.

Runs in Windows Terminal, PowerShell, macOS Terminal, iTerm and Linux terminals. Not supported in Git Bash (mintty).
