package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Spinner struct {
	frame  int
	mspf   int
	frames []rune

	done chan struct{}
	wg   *sync.WaitGroup
}

func NewSpinner() *Spinner {
	return &Spinner{
		frame:  0,
		mspf:   200,
		frames: []rune{'.', 'o', 'O', '@', '*'},
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
	fmt.Print(HideCursor)
	defer fmt.Print(EraseEOL + ShowCursor + SGR(Reset))

	for {
		select {
		case <-s.done:
			return
		default:
			fmt.Printf("%s%c%s", SGR(FgGreen), s.frames[s.frame], Left(1)+SGR(FgReset))
			s.frame = (s.frame + 1) % len(s.frames)
			<-time.After(time.Duration(s.mspf) * time.Millisecond)
		}
	}
}

const (
	EscRune = '\u001b'
	Esc     = string(EscRune)
	CSI     = Esc + "[" // Control Sequence Introducer

	EraseEOL = CSI + "K" // erase to end of current line

	// DECTCEM commands

	ShowCursor = CSI + "?25h"
	HideCursor = CSI + "?25l"
)

func Right(n int) string {
	return fmt.Sprintf(CSI+"%dC", n)
}

func Left(n int) string {
	return fmt.Sprintf(CSI+"%dD", n)
}

const (
	Reset         = "0"
	FgBlack       = "30"
	FgDarkRed     = "31"
	FgDarkGreen   = "32"
	FgDarkYellow  = "33"
	FgDarkBlue    = "34"
	FgDarkMagenta = "35"
	FgDarkCyan    = "36"
	FgLightGray   = "37"
	FgReset       = "39" // sets foreground color to default
	BgBlack       = "40"
	BgDarkRed     = "41"
	BgDarkGreen   = "42"
	BgDarkYellow  = "43"
	BgDarkBlue    = "44"
	BgDarkMagenta = "45"
	BgDarkCyan    = "46"
	BgLightGray   = "47"
	BgReset       = "49" // sets background color to default
	FgDarkGray    = "90"
	FgRed         = "91"
	FgGreen       = "92"
	FgYellow      = "93"
	FgBlue        = "94"
	FgMagenta     = "95"
	FgCyan        = "96"
	FgWhite       = "97"
	BgDarkGray    = "100"
	BgRed         = "101"
	BgGreen       = "102"
	BgYellow      = "103"
	BgBlue        = "104"
	BgMagenta     = "105"
	BgCyan        = "106"
	BgWhite       = "107"
)

// SGR applies the above sgr params in the order specified (later commands may override earlier commands)
func SGR(params ...string) string {
	var sb strings.Builder
	sb.Grow(len(params)*4 + 2) // worst case: n 3-char params, n-1 semicolons, <esc>, '[', and 'm'
	sb.WriteString(CSI)
	for i, v := range params {
		if i != 0 {
			sb.WriteRune(';')
		}
		sb.WriteString(string(v))
	}
	sb.WriteRune('m')
	return sb.String()
}
