package colorgful

import "github.com/kovetskiy/lorg"

type formatter struct {
	*lorg.Format

	restorer *restorer
}

func (format *formatter) Render(level lorg.Level, prefix string) string {
	format.restorer.reset()

	return format.Format.Render(level, prefix)
}
