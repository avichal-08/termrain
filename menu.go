package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

const version = "v0.3.0"

type line struct{ s, col string }

var commands = []struct{ name, desc string }{
	{"/play", "start a round"},
	{"/stats", "show your personal best"},
	{"/reset", "clear your saved best score"},
	{"/help", "how to play"},
	{"/quit", "exit termrain"},
}

func (g *game) toMenu() {
	g.reset()
	g.inMenu = true
	g.cmd = ""
	g.menuOut = nil
}

func (g *game) matches() []int { return matchesFor(g.cmd) }

func matchesFor(prefix string) []int {
	var out []int
	for i, c := range commands {
		if strings.HasPrefix(c.name, strings.ToLower(prefix)) {
			out = append(out, i)
		}
	}
	return out
}

func (g *game) menuKey(c byte) bool {
	switch {
	case c == 27:
		return true
	case c == '\r' || c == '\n':
		return g.runCommand()
	case c == 127 || c == 8:
		if len(g.cmd) > 0 {
			g.cmd = g.cmd[:len(g.cmd)-1]
		}
	case c == '\t':
		if strings.HasPrefix(g.cmd, "/") {
			if m := g.matches(); len(m) > 0 {
				g.cmd = commands[m[0]].name
			}
		}
	case c >= 32 && c < 127 && len(g.cmd) < 40:
		g.cmd += string(c)
	}
	return false
}

func (g *game) runCommand() bool {
	input := strings.ToLower(strings.TrimSpace(g.cmd))
	g.cmd = ""
	if input == "" {
		input = "/play"
	}
	if !strings.HasPrefix(input, "/") {
		input = "/" + input
	}
	name := ""
	for _, c := range commands {
		if c.name == input {
			name = c.name
		}
	}
	if m := matchesFor(input); name == "" && len(m) == 1 {
		name = commands[m[0]].name
	}

	echo := line{"❯ " + input, dim}
	armed := g.resetArmed
	g.resetArmed = false

	switch name {
	case "/play":
		g.reset()
		g.inMenu = false
	case "/quit":
		return true
	case "/stats":
		if g.best.Score == 0 {
			g.menuOut = []line{echo, {"⎿  no runs saved yet · /play to set one", dim}}
			break
		}
		g.menuOut = []line{echo,
			{fmt.Sprintf("⎿  best score   %d", g.best.Score), magenta},
			{fmt.Sprintf("   speed        %d wpm", g.best.WPM), white},
			{fmt.Sprintf("   best combo   x%d", g.best.Combo), white},
			{"   set on       " + g.best.When.Local().Format("2 Jan 2006, 15:04"), dim},
		}
	case "/reset":
		if !armed {
			g.resetArmed = true
			g.menuOut = []line{echo, {"⎿  this deletes your best score · run /reset again to confirm", yellow}}
			break
		}
		if err := os.Remove(bestPath()); err != nil && !os.IsNotExist(err) {
			g.menuOut = []line{echo, {"⎿  couldn't clear: " + err.Error(), red}}
			break
		}
		g.best = Best{}
		g.menuOut = []line{echo, {"⎿  best score cleared", green}}
	case "/help":
		g.menuOut = []line{echo,
			{"⎿  words fall · type one to lock on, finish it to clear it", white},
			{"   5 clears in a row raises your multiplier · 3 misses and you're out", white},
			{"   backspace drops a word · esc goes back to this prompt", dim},
			{"   on game over: [c] copies a result card to share", dim},
		}
	default:
		g.menuOut = []line{echo, {"⎿  unknown command · try /help", red}}
	}
	return false
}

func prettyPath(p string) string {
	if ad := os.Getenv("APPDATA"); ad != "" && strings.HasPrefix(p, ad) {
		return "%AppData%" + p[len(ad):]
	}
	if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(p, home) {
		return "~" + p[len(home):]
	}
	return p
}

func fit(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	r := []rune(s)
	return string(r[:max(n-1, 0)]) + "…"
}

func (c *canvas) box(x, y, w, h int, col string) {
	c.put(x, y, "╭"+strings.Repeat("─", w-2)+"╮", col)
	for i := 1; i < h-1; i++ {
		c.put(x, y+i, "│"+strings.Repeat(" ", w-2)+"│", col)
	}
	c.put(x, y+h-1, "╰"+strings.Repeat("─", w-2)+"╯", col)
}

func (g *game) renderMenu(c *canvas) {
	for _, d := range g.rain {
		c.put(d.x, int(d.y), "'", blue)
	}

	bw := min(g.w-2, 66)
	inner := bw - 6
	c.box(1, 1, bw, 8, dim)
	c.put(3, 2, "☂ termrain", cyan)
	c.put(15, 2, version, dim)
	c.put(5, 4, fit("type the rain before it lands.", inner), white)
	if g.best.Score > 0 {
		c.put(5, 5, "best", dim)
		c.put(12, 5, fit(fmt.Sprintf("%d pts · %d wpm · x%d", g.best.Score, g.best.WPM, g.best.Combo), inner-7), magenta)
	} else {
		c.put(5, 5, "best", dim)
		c.put(12, 5, fit("no runs yet · go set one", inner-7), dim)
	}
	c.put(5, 6, "words", dim)
	c.put(12, 6, fit(fmt.Sprintf("%d · saved to %s", len(wordList), prettyPath(bestPath())), inner-7), dim)

	row := 10
	tips := []string{
		"start typing any falling word to lock on",
		"5 in a row = score multiplier",
		"press enter to play, or / for commands",
	}
	if len(g.menuOut) == 0 {
		c.put(2, row, "Tips:", dim)
		for i, t := range tips {
			c.put(3, row+1+i, fit("· "+t, g.w-5), dim)
		}
	}

	promptTop := g.h - 5
	for _, l := range g.menuOut {
		if row >= promptTop-1 {
			break
		}
		c.put(2, row, fit(l.s, g.w-4), l.col)
		row++
	}

	c.box(1, promptTop, bw, 3, dim)
	c.put(3, promptTop+1, "❯ ", cyan)
	c.put(5, promptTop+1, fit(g.cmd, bw-8), white)
	c.put(5+utf8.RuneCountInString(fit(g.cmd, bw-8)), promptTop+1, "▏", dim)
	if g.cmd == "" {
		c.put(6, promptTop+1, "press enter to play", dim)
	}

	x := 3
	m := g.matches()
	if !strings.HasPrefix(g.cmd, "/") {
		m = []int{0, 1, 2, 3, 4}
	}
	for i, idx := range m {
		col := dim
		if i == 0 && strings.HasPrefix(g.cmd, "/") {
			col = cyan
		}
		c.put(x, g.h-2, commands[idx].name, col)
		x += len(commands[idx].name) + 2
	}
	if strings.HasPrefix(g.cmd, "/") && len(m) > 0 {
		c.put(x, g.h-2, fit(commands[m[0]].desc, max(g.w-x-1, 0)), dim)
	} else if strings.HasPrefix(g.cmd, "/") {
		c.put(3, g.h-2, "no matching command", red)
	}
	c.put(3, g.h-1, "tab to complete · esc to quit", dim)
}
