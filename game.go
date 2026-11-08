package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxLives = 3
	fieldTop = 2
)

type word struct {
	text  string
	x     int
	y     float64
	speed float64
}

type popup struct {
	text  string
	color string
	x, y  int
	ttl   int
}

type drop struct {
	x     int
	y     float64
	speed float64
}

type game struct {
	w, h int
	rng  *rand.Rand

	words  []*word
	target *word
	typed  int
	popups []popup
	rain   []drop

	score, lives, combo, maxCombo int
	correct, mistakes             int
	elapsed, spawnIn              float64
	shake, inputErr               int

	over     bool
	copyReq  bool
	copied   bool
	lastCard string
	best     Best
	newBest  bool
	saveErr  error
}

func newGame(w, h int, best Best) *game {
	g := &game{rng: rand.New(rand.NewSource(time.Now().UnixNano())), best: best}
	g.resize(w, h)
	g.reset()
	return g
}

func (g *game) reset() {
	g.words, g.target, g.typed, g.popups = nil, nil, 0, nil
	g.score, g.lives, g.combo, g.maxCombo = 0, maxLives, 0, 0
	g.correct, g.mistakes = 0, 0
	g.elapsed, g.spawnIn = 0, 0.3
	g.shake, g.inputErr = 0, 0
	g.over, g.newBest, g.saveErr = false, false, nil
	g.copyReq, g.copied = false, false
}

func (g *game) resize(w, h int) {
	if w == g.w && h == g.h {
		return
	}
	g.w, g.h = w, h
	g.rain = g.rain[:0]
	for i := 0; i < w/3; i++ {
		g.rain = append(g.rain, drop{g.rng.Intn(max(w, 1)), float64(g.rng.Intn(max(h, 1))), 8 + g.rng.Float64()*10})
	}
	for _, wd := range g.words {
		wd.x = clamp(wd.x, 1, max(1, w-len(wd.text)-1))
	}
}

func (g *game) ground() int { return g.h - 3 }
func (g *game) level() int  { return 1 + int(g.elapsed/20) }

func (g *game) wpm() int {
	if g.elapsed < 1 {
		return 0
	}
	return int(float64(g.correct) / 5 / (g.elapsed / 60))
}

func (g *game) accuracy() int {
	total := g.correct + g.mistakes
	if total == 0 {
		return 100
	}
	return g.correct * 100 / total
}

func (g *game) spawn() {
	lvl := g.level()
	pool := wordsWithLen(min(2+lvl/2, 7), 3+lvl*2)
	if len(pool) == 0 {
		pool = wordList
	}
	var text string
	for try := 0; try < 12; try++ {
		text = pool[g.rng.Intn(len(pool))]
		if g.findStarting(text[0]) == nil {
			break
		}
	}
	maxX := g.w - len(text) - 1
	if maxX < 1 {
		return
	}
	speed := (0.8 + g.elapsed*0.022) * (0.75 + g.rng.Float64()*0.5)
	g.words = append(g.words, &word{text: text, x: 1 + g.rng.Intn(maxX), y: fieldTop, speed: math.Min(speed, 7)})
}

func (g *game) findStarting(c byte) *word {
	var best *word
	for _, wd := range g.words {
		if wd.text[0] == c && (best == nil || wd.y > best.y) {
			best = wd
		}
	}
	return best
}

func (g *game) remove(target *word) {
	for i, wd := range g.words {
		if wd == target {
			g.words = append(g.words[:i], g.words[i+1:]...)
			break
		}
	}
	if g.target == target {
		g.target, g.typed = nil, 0
	}
}

func (g *game) update(dt float64) {
	for i := range g.rain {
		d := &g.rain[i]
		d.y += d.speed * dt
		if int(d.y) >= g.ground() {
			d.y, d.x = fieldTop, g.rng.Intn(max(g.w, 1))
		}
	}
	for i := 0; i < len(g.popups); i++ {
		if g.popups[i].ttl--; g.popups[i].ttl <= 0 {
			g.popups = append(g.popups[:i], g.popups[i+1:]...)
			i--
		}
	}
	g.shake = max(0, g.shake-1)
	g.inputErr = max(0, g.inputErr-1)
	if g.over {
		return
	}

	g.elapsed += dt
	g.spawnIn -= dt
	if g.spawnIn <= 0 {
		g.spawn()
		g.spawnIn = math.Max(0.45, 2.1-g.elapsed*0.02)
	}

	for _, wd := range append([]*word(nil), g.words...) {
		wd.y += wd.speed * dt
		if int(wd.y) >= g.ground() {
			g.remove(wd)
			g.lives--
			g.combo = 0
			g.shake = 10
			g.popups = append(g.popups, popup{"✗ " + wd.text, red, wd.x, g.ground() - 1, 25})
			if g.lives <= 0 {
				g.gameOver()
				return
			}
		}
	}
}

func (g *game) gameOver() {
	g.over = true
	g.target, g.typed = nil, 0
	if g.score > g.best.Score {
		g.newBest = true
		g.best = Best{Score: g.score, WPM: g.wpm(), Combo: g.maxCombo, When: time.Now()}
		g.saveErr = saveBest(g.best)
	}
	g.lastCard = g.card()
}

const repoURL = "github.com/avichal-08/termrain"

func (g *game) survived() string {
	s := int(g.elapsed)
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

func (g *game) card() string {
	head := fmt.Sprintf("☂ termrain · %d pts", g.score)
	if g.newBest {
		head += " 🏆 new best"
	}
	stats := fmt.Sprintf("⌨️ %d wpm  🎯 %d%%  🔥 x%d  ⏱ %s", g.wpm(), g.accuracy(), g.maxCombo, g.survived())
	return head + "\n" + stats + "\n" + repoURL
}

func (g *game) key(c byte) bool {
	switch {
	case c == 3 || c == 27:
		return true
	case g.over:
		if c == 'q' || c == 'Q' {
			return true
		}
		if c == 'r' || c == 'R' || c == ' ' || c == '\r' {
			g.reset()
		}
		if c == 'c' || c == 'C' {
			g.copyReq = true
			g.copied = true
		}
		return false
	case c == 127 || c == 8:
		g.target, g.typed = nil, 0
		return false
	}
	if c >= 'A' && c <= 'Z' {
		c += 'a' - 'A'
	}
	if c < 'a' || c > 'z' {
		return false
	}

	if g.target == nil {
		if g.target = g.findStarting(c); g.target != nil {
			g.typed = 1
			g.correct++
		} else {
			g.miss()
		}
	} else if g.target.text[g.typed] == c {
		g.typed++
		g.correct++
	} else {
		g.miss()
	}

	if g.target != nil && g.typed == len(g.target.text) {
		g.combo++
		g.maxCombo = max(g.maxCombo, g.combo)
		mult := 1 + g.combo/5
		pts := len(g.target.text) * 10 * mult
		g.score += pts
		label := fmt.Sprintf("+%d", pts)
		if mult > 1 {
			label += fmt.Sprintf(" x%d", mult)
		}
		g.popups = append(g.popups, popup{label, green, g.target.x, int(g.target.y), 18})
		g.remove(g.target)
	}
	return false
}

func (g *game) miss() {
	g.mistakes++
	g.combo = 0
	g.inputErr = 6
}

const (
	reset   = "\x1b[0m"
	dim     = "\x1b[90m"
	white   = "\x1b[97m"
	green   = "\x1b[1;92m"
	yellow  = "\x1b[1;93m"
	red     = "\x1b[1;91m"
	cyan    = "\x1b[1;96m"
	magenta = "\x1b[1;95m"
	blue    = "\x1b[34m"
)

type cell struct {
	r     rune
	color string
}

type canvas struct {
	w, h  int
	cells []cell
}

func newCanvas(w, h int) *canvas {
	c := &canvas{w: w, h: h, cells: make([]cell, w*h)}
	for i := range c.cells {
		c.cells[i] = cell{' ', reset}
	}
	return c
}

func (c *canvas) put(x, y int, s, color string) {
	if y < 0 || y >= c.h {
		return
	}
	for _, r := range s {
		if x >= 0 && x < c.w {
			c.cells[y*c.w+x] = cell{r, color}
		}
		x++
	}
}

func (c *canvas) center(y int, s, color string) {
	c.put((c.w-utf8.RuneCountInString(s))/2, y, s, color)
}

func (c *canvas) String() string {
	var b strings.Builder
	b.Grow(c.w * c.h * 2)
	b.WriteString("\x1b[H")
	cur := ""
	for y := 0; y < c.h; y++ {
		fmt.Fprintf(&b, "\x1b[%d;1H", y+1)
		for x := 0; x < c.w; x++ {
			cl := c.cells[y*c.w+x]
			if cl.color != cur {
				b.WriteString(reset)
				b.WriteString(cl.color)
				cur = cl.color
			}
			b.WriteRune(cl.r)
		}
	}
	b.WriteString(reset)
	return b.String()
}

func (g *game) render() string {
	c := newCanvas(g.w, g.h)
	if g.w < 40 || g.h < 14 {
		c.center(g.h/2, "make the terminal bigger ☂", yellow)
		return c.String()
	}

	hearts := strings.Repeat("♥", max(g.lives, 0)) + strings.Repeat("♡", maxLives-max(g.lives, 0))
	c.put(1, 0, "☂ TERMRAIN", cyan)
	c.put(16, 0, fmt.Sprintf("score %d", g.score), white)
	c.put(30, 0, fmt.Sprintf("best %d", max(g.best.Score, g.score)), magenta)
	right := fmt.Sprintf("lvl %d  %d wpm  combo %d  ", g.level(), g.wpm(), g.combo)
	c.put(g.w-len(right)-maxLives-1, 0, right, dim)
	c.put(g.w-maxLives-1, 0, hearts, red)
	c.put(0, 1, strings.Repeat("─", g.w), dim)

	for _, d := range g.rain {
		c.put(d.x, int(d.y), "'", blue)
	}

	for _, wd := range g.words {
		y := int(wd.y)
		col := white
		if y >= g.ground()-4 {
			col = red
		}
		if wd == g.target {
			c.put(wd.x, y, wd.text[:g.typed], green)
			c.put(wd.x+g.typed, y, wd.text[g.typed:], yellow)
		} else {
			c.put(wd.x, y, wd.text, col)
		}
	}
	for _, p := range g.popups {
		c.put(p.x, p.y, p.text, p.color)
	}

	groundCol := dim
	if g.shake > 0 && g.shake%2 == 0 {
		groundCol = red
	}
	c.put(0, g.ground(), strings.Repeat("▀", g.w), groundCol)
	prompt, promptCol := "❯ ", cyan
	if g.inputErr > 0 {
		promptCol = red
	}
	c.put(1, g.h-2, prompt, promptCol)
	if g.target != nil {
		c.put(3, g.h-2, g.target.text[:g.typed], green)
		c.put(3+g.typed, g.h-2, "▏", dim)
	}
	c.put(1, g.h-1, "type the falling words · backspace = drop target · esc = quit", dim)

	if g.over {
		g.renderGameOver(c)
	}
	return c.String()
}

func (g *game) renderGameOver(c *canvas) {
	lines := []struct{ s, col string }{
		{"G A M E   O V E R", red},
		{"", ""},
		{fmt.Sprintf("score  %d", g.score), white},
		{fmt.Sprintf("%d wpm · %d%% accuracy · best combo %d · %s", g.wpm(), g.accuracy(), g.maxCombo, g.survived()), dim},
		{"", ""},
	}
	switch {
	case g.newBest && g.saveErr != nil:
		lines = append(lines, struct{ s, col string }{"★ NEW BEST! (couldn't save: " + g.saveErr.Error() + ")", yellow})
	case g.newBest:
		lines = append(lines, struct{ s, col string }{"★ NEW PERSONAL BEST · saved ★", yellow})
	default:
		lines = append(lines, struct{ s, col string }{fmt.Sprintf("personal best  %d  (%d wpm)", g.best.Score, g.best.WPM), magenta})
	}
	share := struct{ s, col string }{"[c] copy result to share", green}
	if g.copied {
		share = struct{ s, col string }{"✓ copied · paste it anywhere", green}
	}
	lines = append(lines,
		struct{ s, col string }{"", ""},
		share,
		struct{ s, col string }{"[r] play again    [q] quit", cyan})

	boxW := 0
	for _, l := range lines {
		boxW = max(boxW, utf8.RuneCountInString(l.s))
	}
	boxW += 6
	top := (g.h - len(lines)) / 2
	left := (g.w - boxW) / 2
	c.put(left, top-2, "╭"+strings.Repeat("─", boxW-2)+"╮", dim)
	for y := top - 1; y < top+len(lines)+1; y++ {
		c.put(left, y, "│"+strings.Repeat(" ", boxW-2)+"│", dim)
	}
	c.put(left, top+len(lines)+1, "╰"+strings.Repeat("─", boxW-2)+"╯", dim)
	for i, l := range lines {
		c.center(top+i, l.s, l.col)
	}
}

func clamp(v, lo, hi int) int { return max(lo, min(v, hi)) }
