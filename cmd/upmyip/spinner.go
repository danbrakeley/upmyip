package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/danbrakeley/ansi"
)

type Spinner struct {
	frame   int
	mspf    int
	frames  []rune
	noColor bool
	done    chan struct{}
	wg      *sync.WaitGroup
}

func NewSpinner(noColor bool) *Spinner {
	return &Spinner{
		frame:   0,
		mspf:    20,
		frames:  []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'},
		noColor: noColor,
	}
}

func (s *Spinner) Start() {
	if s.done != nil {
		return
	}
	s.done = make(chan struct{})
	s.wg = &sync.WaitGroup{}
	s.wg.Add(1)
	go s.Run()
}

func (s *Spinner) Stop() {
	if s.done == nil {
		return
	}
	close(s.done)
	s.wg.Wait()
	s.done = nil
	s.wg = nil
}

func (s *Spinner) Run() {
	defer s.wg.Done()
	fmt.Print(ansi.HideCursor)
	defer fmt.Print(ansi.EraseEOL + ansi.ShowCursor + ansi.SGR(ansi.Reset))

	green := ""
	if !s.noColor {
		green = ansi.SGR(ansi.FgYellow)
	}
	reset := ""
	if !s.noColor {
		reset = ansi.SGR(ansi.FgReset)
	}

	for {
		select {
		case <-s.done:
			return
		default:
			fmt.Printf("%s%c%s", green, s.frames[s.frame], ansi.Left(1)+reset)
			s.frame = (s.frame + 1) % len(s.frames)
			<-time.After(time.Duration(s.mspf) * time.Millisecond)
		}
	}
}
