package colorgful

import (
	"fmt"
	"strings"

	"github.com/reconquest/loreley"
)

type inserter struct {
	*loreley.Style
}

func (inserter *inserter) insertLorg(
	value string,
) string {
	return fmt.Sprintf(`${%s}`, value)
}

func (inserter *inserter) insertOnLevel(
	level string,
	value string,
) string {
	style, err := loreley.Compile(value, nil)
	if err != nil {
		return fmt.Sprintf("#{COMPILE ERROR: %s}", err)
	}

	style.NoColors = NoColors

	var previous string

	if !NoColors {
		previous, _ = loreley.CompileAndExecuteToString(
			inserter.GetState().String(),
			nil,
			nil,
		)
	}

	style.SetState(inserter.GetState())

	sequence, err := style.ExecuteToString(nil)
	if err != nil {
		return fmt.Sprintf("#{EXECUTE ERROR: %s}", err)
	}

	return inserter.insertLorg(
		`onlevel` + `:` +
			strings.ToLower(level) + `:` + sequence + `:` +
			previous,
	)
}

func (inserter *inserter) insertRestore() string {
	return inserter.insertLorg(`restore`)
}

func (inserter *inserter) insertStore() string {
	return inserter.insertLorg(`store`)
}
