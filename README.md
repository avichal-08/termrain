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
go run .
```

- Start typing any falling word. The first letter locks onto it.
- Finish a word to clear it. Every 5 words in a row bumps your multiplier.
- A word that reaches the ground costs a life. You have 3.
- Words get longer and faster the longer you last.

Keys: `Backspace` drops the current word, `Esc` quits, `r` restarts after game over.

## Best score

Your best score, WPM and combo are saved locally:

- Windows: `%AppData%\termrain\best.json`
- macOS: `~/Library/Application Support/termrain/best.json`
- Linux: `~/.config/termrain/best.json`

Reset it with `go run . -reset`.

## Notes

- Raw keyboard input without libraries: `SetConsoleMode` from `kernel32.dll` on Windows, `stty raw` on macOS/Linux.
- Keystrokes are read on a separate goroutine and fed through a channel; the game loop selects between input and a 30 fps ticker.
- Frames are built in a cell buffer and written once per tick to the alternate screen.

Runs in Windows Terminal, PowerShell, macOS Terminal, iTerm and Linux terminals. Not supported in Git Bash (mintty).
