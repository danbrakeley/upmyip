package main

import (
	"fmt"
	"strings"

	"github.com/danbrakeley/ansi"
)

type Printer struct {
	NoColor bool
}

func (p Printer) Header(a ...string) {
	var s strings.Builder
	s.Grow(256)
	if !p.NoColor {
		s.WriteString(ansi.SGR(ansi.FgBlue))
	}
	for _, v := range a {
		s.WriteString(v)
	}
	if !p.NoColor {
		s.WriteString(ansi.SGR(ansi.FgReset))
	}
	fmt.Println(s.String())
}

func (p Printer) Print(a ...any) {
	fmt.Print(a...)
}

func (p Printer) BrightPrintln(msg string) {
	var s strings.Builder
	s.Grow(256)
	if !p.NoColor {
		s.WriteString(ansi.SGR(ansi.FgWhite))
	}
	s.WriteString(msg)
	if !p.NoColor {
		s.WriteString(ansi.SGR(ansi.FgReset))
	}
	fmt.Println(s.String())
}

func (p Printer) Error(prefix string, err error) {
	var s strings.Builder
	s.Grow(256)
	if !p.NoColor {
		s.WriteString(ansi.SGR(ansi.FgRed))
	}
	s.WriteString(prefix)
	s.WriteString(": ")
	if !p.NoColor {
		s.WriteString(ansi.SGR(ansi.FgYellow))
	}
	s.WriteString(err.Error())
	if !p.NoColor {
		s.WriteString(ansi.SGR(ansi.FgReset))
	}
	fmt.Println(s.String())
}

func (p Printer) BrightIPln(info *PublicInfo) {
	var s strings.Builder
	s.Grow(256)
	if !p.NoColor {
		s.WriteString(ansi.SGR(ansi.FgWhite))
	}
	s.WriteString(info.IP)
	s.WriteRune(' ')
	if !p.NoColor {
		s.WriteString(ansi.SGR(ansi.FgCyan))
	}
	s.WriteRune('(')
	s.WriteString(info.ISP)
	s.WriteRune(')')
	if !p.NoColor {
		s.WriteString(ansi.SGR(ansi.FgReset))
	}
	fmt.Println(s.String())
}
