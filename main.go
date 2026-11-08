package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-reset" || os.Args[1] == "--reset") {
		os.Remove(bestPath())
		fmt.Println("best score cleared:", bestPath())
		return
	}

	card, err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "termrain needs an interactive terminal (Windows Terminal, PowerShell, macOS/Linux terminal):", err)
		os.Exit(1)
	}
	if card != "" {
		fmt.Println("\n" + card + "\n")
	}
}

func run() (string, error) {
	restore, err := enableRaw()
	if err != nil {
		return "", err
	}
	out := bufio.NewWriterSize(os.Stdout, 1<<16)
	fmt.Fprint(out, "\x1b[?1049h\x1b[?25l\x1b[2J")
	out.Flush()
	defer func() {
		fmt.Fprint(out, "\x1b[0m\x1b[?25h\x1b[?1049l")
		out.Flush()
		restore()
	}()

	keys := make(chan []byte, 64)
	go func() {
		buf := make([]byte, 64)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				close(keys)
				return
			}
			k := make([]byte, n)
			copy(k, buf[:n])
			keys <- k
		}
	}()

	w, h := termSize()
	g := newGame(w, h, loadBest())
	frame := time.NewTicker(time.Second / 30)
	defer frame.Stop()
	last, lastSize := time.Now(), time.Now()

	for {
		select {
		case k, ok := <-keys:
			if !ok {
				return g.lastCard, nil
			}
			if len(k) > 1 && k[0] == 27 {
				continue
			}
			for _, c := range k {
				if g.key(c) {
					return g.lastCard, nil
				}
			}
			if g.copyReq {
				g.copyReq = false
				fmt.Fprintf(out, "\x1b]52;c;%s\a", base64.StdEncoding.EncodeToString([]byte(g.lastCard)))
			}
		case now := <-frame.C:
			if now.Sub(lastSize) > 500*time.Millisecond {
				g.resize(termSize())
				lastSize = now
			}
			g.update(now.Sub(last).Seconds())
			last = now
			out.WriteString(g.render())
			out.Flush()
		}
	}
}
