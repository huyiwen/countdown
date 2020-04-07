package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	usage = `usage
    countdown 25s
    countdown 1m50s
	countdown 2h45m50s
`
	keybinding = `Space  pause/resume
Esc    stop
b      turn on/off the bell
`
	bellon  = `The bell is on `
	belloff = `The bell is off`

	tick = time.Second

	Redtime = 5 * time.Second
)

var (
	timer          *time.Timer
	ticker         *time.Ticker
	queues         chan termbox.Event
	startDone      bool
	pause          bool
	startX, startY int
	ring           bool
)

func draw(d time.Duration) {
	w, h := termbox.Size()
	clear()

	str := format(d)
	text := toText(str)

	if !startDone {
		startDone = true
		startX, startY = w/2-text.width()/2, h/2-text.height()/2
		fmt.Print(keybinding)
		/*
		if ring {
			fmt.Print(bellon)
		} else {
			fmt.Print(belloff)
		}
		*/
	}

	if d <= Redtime {
		x, y := startX, startY
		for _, s := range text {
			echoRed(s, x, y)
			x += s.width()
		}
	} else {
		x, y := startX, startY
		for _, s := range text {
			echo(s, x, y)
			x += s.width()
		}
	}

	flush()
}

func format(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h < 1 {
		return fmt.Sprintf("%02d:%02d", m, s)
	}
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func start(d time.Duration) {
	pause = false
	timer = time.NewTimer(d)
	ticker = time.NewTicker(tick)
}

func stop() {
	pause = true
	timer.Stop()
	ticker.Stop()
}

func countdown(left time.Duration) {
	var exitCode int

	start(left)

loop:
	for {
		select {
		case ev := <-queues:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
					exitCode = 1
					break loop
				}
				if ev.Key == termbox.KeySpace {
					if pause {
						start(left)
					} else {
						stop()
					}
				}
			}
			if ev.Ch == 'b' || ev.Ch == 'B' {
				// startDone = false
				ring = !ring
			}
			if ev.Ch == 'p' || ev.Ch == 'P' {
				stop()
			}
			if ev.Ch == 'c' || ev.Ch == 'C' {
				start(left)
			}
		case <-ticker.C:
			left -= time.Duration(tick)
			draw(left)
		case <-timer.C:
			break loop
		}
	}

	termbox.Close()
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func main() {
	ring = true
	if len(os.Args) != 2 {
		stderr(usage)
		os.Exit(2)
	}

	duration, err := time.ParseDuration(os.Args[1])
	if err != nil {
		stderr("error: invalid duration: %v\n", os.Args[1])
		os.Exit(2)
	}
	left := duration

	err = termbox.Init()
	if err != nil {
		panic(err)
	}

	queues = make(chan termbox.Event)
	go func() {
		for {
			queues <- termbox.PollEvent()
		}
	}()

	draw(left)
	countdown(left)

	if ring {
		cmd := exec.Command("tput","bel")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
