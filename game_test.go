package main

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func isolate(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPDATA", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
}

func typeWord(g *game, s string) {
	for i := 0; i < len(s); i++ {
		g.key(s[i])
	}
}

func TestTypingAWordScores(t *testing.T) {
	isolate(t)
	g := newGame(80, 24, Best{})
	g.words = []*word{{text: "defer", x: 5, y: 5, speed: 1}}
	typeWord(g, "defer")
	if len(g.words) != 0 || g.score != 50 || g.combo != 1 {
		t.Fatalf("words=%d score=%d combo=%d", len(g.words), g.score, g.combo)
	}
}

func TestWrongKeyBreaksCombo(t *testing.T) {
	isolate(t)
	g := newGame(80, 24, Best{})
	g.combo = 4
	g.words = []*word{{text: "chan", x: 5, y: 5, speed: 1}}
	typeWord(g, "cx")
	if g.combo != 0 || g.mistakes != 1 || g.target == nil || g.typed != 1 {
		t.Fatalf("combo=%d mistakes=%d typed=%d", g.combo, g.mistakes, g.typed)
	}
}

func TestTargetsLowestMatchingWord(t *testing.T) {
	isolate(t)
	g := newGame(80, 24, Best{})
	hi := &word{text: "map", x: 5, y: 3}
	lo := &word{text: "mutex", x: 20, y: 15}
	g.words = []*word{hi, lo}
	g.key('m')
	if g.target != lo {
		t.Fatal("should lock onto the word closest to the ground")
	}
}

func TestLosingSavesBestAndReloads(t *testing.T) {
	isolate(t)
	g := newGame(80, 24, Best{})
	g.score = 1234
	for g.lives > 0 {
		g.words = []*word{{text: "bug", x: 5, y: float64(g.ground()), speed: 1}}
		g.update(0.01)
	}
	if !g.over || !g.newBest || g.saveErr != nil {
		t.Fatalf("over=%v newBest=%v err=%v", g.over, g.newBest, g.saveErr)
	}
	if b := loadBest(); b.Score != 1234 {
		t.Fatalf("reloaded best = %d, want 1234 (file %s)", b.Score, bestPath())
	}

	g2 := newGame(80, 24, loadBest())
	g2.score = 10
	g2.gameOver()
	if g2.newBest || loadBest().Score != 1234 {
		t.Fatal("lower score overwrote the best")
	}
	if _, err := os.Stat(bestPath()); err != nil {
		t.Fatal(err)
	}
}

func TestRestartAndQuitKeys(t *testing.T) {
	isolate(t)
	g := newGame(80, 24, Best{})
	g.gameOver()
	if g.key('r'); g.over || g.lives != maxLives {
		t.Fatal("r should restart")
	}
	if !g.key(27) {
		t.Fatal("esc should quit")
	}
}

var ansi = regexp.MustCompile(`\x1b\[[0-9;?]*[A-Za-z]`)

func TestRenderSnapshot(t *testing.T) {
	isolate(t)
	g := newGame(72, 20, Best{Score: 900, WPM: 61})
	for i := 0; i < 30*8; i++ {
		g.update(1.0 / 30)
	}
	if len(g.words) > 0 {
		w := g.words[len(g.words)-1].text
		g.key(w[0])
	}
	screen := ansi.ReplaceAllString(g.render(), "\n")
	screen = regexp.MustCompile(`\n+`).ReplaceAllString(screen, "\n")
	if !strings.Contains(screen, "TERMRAIN") {
		t.Fatal("HUD missing")
	}
	t.Log("\n" + screen)

	g.score = 1500
	g.gameOver()
	over := regexp.MustCompile(`\n+`).ReplaceAllString(ansi.ReplaceAllString(g.render(), "\n"), "\n")
	if !strings.Contains(over, "NEW PERSONAL BEST") {
		t.Fatal("game-over box missing")
	}
	t.Log("\n" + over)
}

func TestShareCard(t *testing.T) {
	isolate(t)
	g := newGame(80, 24, Best{})
	g.score, g.correct, g.elapsed, g.maxCombo = 1840, 552, 92, 14
	g.gameOver()
	card := g.lastCard
	for _, want := range []string{"1840 pts", "new best", "72 wpm", "x14", "1:32", repoURL} {
		if !strings.Contains(card, want) {
			t.Fatalf("card missing %q:\n%s", want, card)
		}
	}
	g.key('c')
	if !g.copyReq || !strings.Contains(g.render(), "copied") {
		t.Fatal("c should request a copy and show feedback")
	}
	t.Log("\n" + card)
}

func TestWordListIsClean(t *testing.T) {
	if len(wordList) < 500 {
		t.Fatalf("only %d words", len(wordList))
	}
	for _, w := range wordList {
		for i := 0; i < len(w); i++ {
			if w[i] < 'a' || w[i] > 'z' {
				t.Fatalf("word %q has a character players can't type", w)
			}
		}
	}
	for lvl := 1; lvl <= 10; lvl++ {
		if n := len(wordsWithLen(min(2+lvl/2, 7), 3+lvl*2)); n < 50 {
			t.Fatalf("level %d pool only %d words", lvl, n)
		}
	}
}
