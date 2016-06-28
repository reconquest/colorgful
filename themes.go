package colorgful

import (
	"bytes"
	"regexp"
	"text/template"

	"github.com/kovetskiy/lorg"
)

// DefaultThemeStyles represents template values for default Light and Dark
// themes.
type DefaultThemeStyles struct {
	// Trace specifies overall style for trace logs level
	Trace string

	// Debug specifies overall style for debug logs level
	Debug string

	// Info specifies overall style for info logs level
	Info string

	// Warning specifies overall style for warning log level
	Warning string

	// Error specifies overall style for error log level
	Error string

	// Fatal specifies overall style for fatal log level
	Fatal string

	// ErrorLevel specifies custom style for the ${level} placeholder for
	// error log level
	ErrorLevel string

	// FatalLevel specifies custom style for the ${level} placeholder for
	// fatal log level
	FatalLevel string
}

var (
	// Dark are the default styles, suitable for the dark shell
	// backgrounds.
	Dark = DefaultThemeStyles{
		Trace:      `{fg 243}`,
		Debug:      `{fg 250}`,
		Info:       `{fg 110}`,
		Warning:    `{fg 178}`,
		Error:      `{fg 202}`,
		Fatal:      `{bold}{fg 197}{bg 17}`,
		ErrorLevel: `{bold}{bg 52}`,
		FatalLevel: ``,
	}

	// Light are the default styles, suitable for the light shell
	// backgrounds.
	Light = DefaultThemeStyles{
		Trace:      `{fg 250}`,
		Debug:      `{fg 243}`,
		Info:       `{fg 26}`,
		Warning:    `{fg 167}{bg 230}`,
		Error:      `{bold}{fg 161}`,
		Fatal:      `{bold}{fg 231}{bg 124}`,
		ErrorLevel: `{reverse}{bold}{bg 231}`,
		FatalLevel: ``,
	}
)

var (
	levelPlaceholderRegexp = regexp.MustCompile(`\s*\${level[^}]*}\s*`)
)

// ApplyDefaultTheme applies default theme to the given lorg formatting string.
//
// `styles` can be used to specify color scheme for the default theme.
func ApplyDefaultTheme(
	formatting string,
	styles DefaultThemeStyles,
) (lorg.Formatter, error) {
	lineStyles := `` +
		`{ontrace "{{.Trace}}"}` +
		`{ondebug "{{.Debug}}"}` +
		`{oninfo "{{.Info}}"}` +
		`{onwarning "{{.Warning}}"}` +
		`{onerror "{{.Error}}"}` +
		`{onfatal "{{.Fatal}}"}` +
		`{store}`

	levelStyles := `` +
		`{onerror "{{.ErrorLevel}}"}` +
		`{onfatal "{{.FatalLevel}}"}` +
		`$0` +
		`{restore}`

	format := lineStyles + levelPlaceholderRegexp.ReplaceAllString(
		formatting,
		levelStyles,
	)

	buffer := bytes.Buffer{}

	err := template.Must(template.New(`theme`).Parse(format)).Execute(
		&buffer,
		styles,
	)
	if err != nil {
		return nil, err
	}

	theme, err := FormatWithReset(buffer.String())
	if err != nil {
		return nil, err
	}

	return theme, nil
}

// MustApplyDefaultTheme is the same, as ApplyDefaultTheme, but panics on error.
func MustApplyDefaultTheme(
	formatting string,
	styles DefaultThemeStyles,
) lorg.Formatter {
	format, err := ApplyDefaultTheme(formatting, styles)
	if err != nil {
		panic(err)
	}

	return format
}
