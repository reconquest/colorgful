package colorgful

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"text/template"

	"github.com/kovetskiy/lorg"
)

type Theme struct {
	lorg.Formatter
	lorg.SmartOutput
}

// DefaultThemeLevel describes how to highlight given level.
type DefaultThemeLevel struct {
	// First describes style for first line.
	First string

	// Trail describes style for all other lines.
	Trail string

	// Level describes style for level substring.
	Level string
}

// DefaultThemeStyles represents template values for default Light and Dark
// themes.
type DefaultThemeStyles struct {
	// Trace specifies overall style for trace logs level
	Trace DefaultThemeLevel

	// Debug specifies overall style for debug logs level
	Debug DefaultThemeLevel

	// Info specifies overall style for info logs level
	Info DefaultThemeLevel

	// Warning specifies overall style for warning log level
	Warning DefaultThemeLevel

	// Error specifies overall style for error log level
	Error DefaultThemeLevel

	// Fatal specifies overall style for fatal log level
	Fatal DefaultThemeLevel
}

type DefaultOutput struct {
	io.Writer

	Trailer lorg.Formatter
}

func (output *DefaultOutput) WriteWithLevel(
	data []byte,
	level lorg.Level,
) (int, error) {
	return output.Write(
		bytes.Replace(
			data,
			[]byte("\n"),
			[]byte(output.Trailer.Render(level, ``)+"\n"),
			1,
		),
	)
}

var (
	// Dark are the default styles, suitable for the dark shell
	// backgrounds.
	Dark = DefaultThemeStyles{
		Trace:   DefaultThemeLevel{First: `{fg 243}`},
		Debug:   DefaultThemeLevel{First: `{fg 250}`},
		Info:    DefaultThemeLevel{First: `{fg 110}`},
		Warning: DefaultThemeLevel{First: `{fg 178}`},
		Error: DefaultThemeLevel{
			First: `{fg 202}`,
			Level: `{bold}{bg 52}`,
		},
		Fatal: DefaultThemeLevel{First: `{bold}{fg 197}{bg 17}`},
	}

	// Light are the default styles, suitable for the light shell
	// backgrounds.
	Light = DefaultThemeStyles{
		Trace:   DefaultThemeLevel{First: `{fg 250}`},
		Debug:   DefaultThemeLevel{First: `{fg 243}`},
		Info:    DefaultThemeLevel{First: `{fg 26}`},
		Warning: DefaultThemeLevel{First: `{fg 167}{bg 230}`},
		Error: DefaultThemeLevel{
			First: `{bold}{fg 161}`,
			Level: `{reverse}{bold}{bg 231}`,
		},
		Fatal: DefaultThemeLevel{First: `{bold}{fg 231}{bg 124}`},
	}

	// Default are the default styles, suitable for the both
	// light and dark shell backgrounds.
	Default = DefaultThemeStyles{
		Trace: DefaultThemeLevel{First: `{nofg}`},
		Debug: DefaultThemeLevel{First: `{fg 31}`},
		Info:  DefaultThemeLevel{First: `{fg 33}`},
		Warning: DefaultThemeLevel{
			First: `{bold}{fg 172}`,
			Trail: `{nobold}`,
		},
		Error: DefaultThemeLevel{
			First: `{bold}{fg 9}`,
			Trail: `{reset}{fg 9}`,
			Level: `{bold}{fg 231}{bg 196}`,
		},
		Fatal: DefaultThemeLevel{
			First: `{bold}{fg 231}{bg 124}`,
			Trail: `{reset}{bold}{fg 231}`,
			Level: `{bold}{fg 231}{bg 196}`,
		},
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
) (*Theme, error) {
	lineStyles := `` +
		`{ontrace "{{.Trace.First}}"}` +
		`{ondebug "{{.Debug.First}}"}` +
		`{oninfo "{{.Info.First}}"}` +
		`{onwarning "{{.Warning.First}}"}` +
		`{onerror "{{.Error.First}}"}` +
		`{onfatal "{{.Fatal.First}}"}` +
		`{store}`

	levelStyles := `` +
		`{onerror "{{.Error.Level}}"}` +
		`{onfatal "{{.Fatal.Level}}"}` +
		`$0` +
		`{reset}` +
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

	formatter, err := FormatWithReset(buffer.String())
	if err != nil {
		return nil, err
	}

	trailStyles := `` +
		`{ontrace "{{.Trace.Trail}}"}` +
		`{ondebug "{{.Debug.Trail}}"}` +
		`{oninfo "{{.Info.Trail}}"}` +
		`{onwarning "{{.Warning.Trail}}"}` +
		`{onerror "{{.Error.Trail}}"}` +
		`{onfatal "{{.Fatal.Trail}}"}`

	buffer.Reset()
	err = template.Must(template.New(`theme`).Parse(trailStyles)).Execute(
		&buffer,
		styles,
	)
	if err != nil {
		return nil, err
	}

	trailer, err := Format(buffer.String())
	if err != nil {
		return nil, err
	}

	return &Theme{
		Formatter: formatter,
		SmartOutput: &DefaultOutput{
			Writer:  os.Stderr,
			Trailer: trailer,
		},
	}, nil
}

// MustApplyDefaultTheme is the same, as ApplyDefaultTheme, but panics on error.
func MustApplyDefaultTheme(
	formatting string,
	styles DefaultThemeStyles,
) *Theme {
	theme, err := ApplyDefaultTheme(formatting, styles)
	if err != nil {
		panic(err)
	}

	return theme
}
