package colorgful

import (
	"regexp"

	"github.com/kovetskiy/lorg"
	"github.com/reconquest/loreley"
)

var (
	placeholderRegexp = regexp.MustCompile(`\${([^}]+)}`)
)

// Format parses specified formatting (loreley based) and return
// lorg.Formatter.
//
// Following functions are available:
//
//   * {bg <color>}, {fg <color>}, {nobg}, {nofg}, {bold}, {nobold}, {reverse},
//     {noreverse}, {reset}, {from <text> <bg>}, {to <bg> <text>}.
//   * {onlevel <level> <style>} - match given level and insert given style.
//   * {ontrace <style>} - alias for {onlevel "trace" <style>}.
//   * {ondebug <style>} - alias for {onlevel "debug" <style>}.
//   * {oninfo <style>} - alias for {onlevel "info" <style>}.
//   * {onwarning <style>} - alias for {onlevel "warning" <style>}.
//   * {onerror <style>} - alias for {onlevel "error" <style>}.
//   * {onfatal <style>} - alias for {onlevel "fatal" <style>}.
//   * {store} - store all previous style
//   * {restore} - restore previous style, saved by {store}.
//
func Format(formatting string) (lorg.Formatter, error) {
	formatting = placeholderRegexp.ReplaceAllString(
		formatting,
		`{lorg "$1"}`,
	)

	var (
		inserter = &inserter{}
		restorer = &restorer{}
	)

	extensions := map[string]interface{}{
		"lorg":    inserter.insertLorg,
		"onlevel": inserter.insertOnLevel,
		"restore": inserter.insertRestore,
		"store":   inserter.insertStore,

		"ontrace": func(value string) string {
			return inserter.insertOnLevel("trace", value)
		},
		"ondebug": func(value string) string {
			return inserter.insertOnLevel("debug", value)
		},
		"oninfo": func(value string) string {
			return inserter.insertOnLevel("info", value)
		},
		"onwarning": func(value string) string {
			return inserter.insertOnLevel("warning", value)
		},
		"onerror": func(value string) string {
			return inserter.insertOnLevel("error", value)
		},
		"onfatal": func(value string) string {
			return inserter.insertOnLevel("fatal", value)
		},
	}

	style, err := loreley.Compile(
		formatting,
		extensions,
	)
	if err != nil {
		return nil, err
	}

	inserter.Style = style

	formatting, err = style.ExecuteToString(nil)
	if err != nil {
		return nil, err
	}

	format := lorg.NewFormat(formatting)
	format.SetPlaceholder("onlevel", restorer.handleOnLevel)
	format.SetPlaceholder("restore", restorer.handleRestore)
	format.SetPlaceholder("store", restorer.handleStore)

	return &formatter{
		Format: format,

		restorer: restorer,
	}, nil
}

// FormatWithReset is same as Format, but appends {reset} to the end.
func FormatWithReset(formatting string) (lorg.Formatter, error) {
	return Format(formatting + loreley.StyleReset)
}
