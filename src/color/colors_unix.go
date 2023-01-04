//go:build !windows

package color

import "github.com/jandedobbeleer/oh-my-posh/src/platform"

func GetAccentColor(env platform.Environment) (*RGB, error) {
	return nil, &platform.NotImplemented{}
}

func (d *DefaultColors) SetAccentColor(env platform.Environment, defaultColor string) {
	if len(defaultColor) == 0 {
		return
	}
	d.accent = &Color{
		Foreground: string(d.AnsiColorFromString(defaultColor, false)),
		Background: string(d.AnsiColorFromString(defaultColor, true)),
	}
}
