package colorgful

import (
	"bytes"
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
	// DefaultThemeTemplate is the common template for the Dark and Light
	// themes.
	//
	// It should be populated using DefaultThemeStyles
	DefaultThemeTemplate = template.Must(
		template.New(`theme`).Delims(`@{`, `}`).Parse(`` +
			`{ontrace "@{.Trace}"}` +
			`{ondebug "@{.Debug}"}` +
			`{oninfo "@{.Info}"}` +
			`{onwarning "@{.Warning}"}` +
			`{onerror "@{.Error}"}` +
			`{onfatal "@{.Fatal}"}` +
			`{store}` +
			`${time} ` +
			`{onerror "@{.ErrorLevel}"}` +
			`{onfatal "@{.FatalLevel}"}` +
			`${level:[%s]:right:true}` +
			`{restore}` +
			` %s`,
		),
	)

	// DarkStyles are the default styles, suitable for the dark shell
	// backgrounds.
	DarkStyles = DefaultThemeStyles{
		Trace:      `{fg 243}`,
		Debug:      `{fg 250}`,
		Info:       `{fg 110}`,
		Warning:    `{fg 178}`,
		Error:      `{fg 202}`,
		Fatal:      `{bold}{fg 168}{bg 17}`,
		ErrorLevel: `{bold}{bg 52}`,
		FatalLevel: ``,
	}

	// LightStyles are the default styles, suitable for the dark shell
	// backgrounds.
	LightStyles = DefaultThemeStyles{
		Trace:      `{fg 250}`,
		Debug:      `{fg 243}`,
		Info:       `{fg 26}`,
		Warning:    `{fg 167}{bg 230}`,
		Error:      `{bold}{fg 161}`,
		Fatal:      `{bold}{fg 231}{bg 124}`,
		ErrorLevel: `{reverse}{bold}{bg 231}`,
		FatalLevel: ``,
	}

	// Dark is a default theme, suitable for the dark backgrounds.
	Dark lorg.Formatter

	// Light is a default theme, suitable for the light backgrounds.
	Light lorg.Formatter
)

func init() {
	var err error

	Dark, err = compileTheme(DarkStyles)
	if err != nil {
		panic(err)
	}

	Light, err = compileTheme(LightStyles)
	if err != nil {
		panic(err)
	}
}

func compileTheme(styles DefaultThemeStyles) (lorg.Formatter, error) {
	buffer := bytes.Buffer{}

	err := DefaultThemeTemplate.Execute(&buffer, styles)
	if err != nil {
		return nil, err
	}

	theme, err := FormatWithReset(buffer.String())
	if err != nil {
		return nil, err
	}

	return theme, nil
}
