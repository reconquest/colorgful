package colorgful

import (
	"strings"
	"sync"

	"github.com/kovetskiy/lorg"
)

type restorer struct {
	previous string
	current  string
	stored   string

	mutex sync.Mutex
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

	restorer.mutex.Lock()
	restorer.previous = previous
	restorer.current = sequence
	restorer.mutex.Unlock()

	return sequence
}

func (restorer *restorer) handleRestore(
	_ lorg.Level,
	_ string,
) string {
	restorer.mutex.Lock()
	defer restorer.mutex.Unlock()
	if restorer.stored != "" {
		return restorer.stored
	}

	return restorer.previous
}

func (restorer *restorer) handleStore(
	_ lorg.Level,
	_ string,
) string {
	restorer.mutex.Lock()
	restorer.stored = restorer.previous + restorer.current
	restorer.mutex.Unlock()

	return ""
}

func (restorer *restorer) reset() {
	restorer.mutex.Lock()
	restorer.stored = ""
	restorer.mutex.Unlock()
}
