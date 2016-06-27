package colorgful

import (
	"strings"

	"github.com/kovetskiy/lorg"
)

type restorer struct {
	previous string
	current  string
	stored   string
}

func (restorer *restorer) handleOnLevel(
	level lorg.Level,
	value string,
) string {
	var (
		parts = strings.SplitN(value, `:`, 3)

		targetLevel = parts[0]

		sequence = parts[1]
		previous = parts[2]
	)

	if targetLevel != strings.ToLower(level.String()) {
		return ""
	}

	restorer.previous = previous
	restorer.current = sequence

	return sequence
}

func (restorer *restorer) handleRestore(
	_ lorg.Level,
	_ string,
) string {
	if restorer.stored != "" {
		return restorer.stored
	}

	return restorer.previous
}

func (restorer *restorer) handleStore(
	_ lorg.Level,
	_ string,
) string {
	restorer.stored = restorer.previous + restorer.current

	return ""
}

func (restorer *restorer) reset() {
	restorer.stored = ""
}
