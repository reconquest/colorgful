package colorgful

import (
	"bytes"
	"testing"

	"github.com/kovetskiy/lorg"
	"github.com/reconquest/loreley"
	"github.com/stretchr/testify/assert"
)

func TestNewFormat_SetsLorgFormatIfNoStylesSpecified(t *testing.T) {
	test := assert.New(t)

	buffer, log := setupBufferedLogger(t, `${level} %s`)

	log.Error(`hello`)

	test.Equal("ERROR hello\n", buffer.String())
}

func TestNewFormat_SetsStyleForAllLevels(t *testing.T) {
	test := assert.New(t)

	buffer, log := setupBufferedLogger(t, `{bg 1}${level} %s`)

	log.Error(`hello`)

	test.Equal(
		compileExpectedStyle(t, "{bg 1}ERROR hello\n"),
		buffer.String(),
	)
}

func TestNewFormat_SetsStyleForSpecifiedLevel(t *testing.T) {
	test := assert.New(t)

	formats := []string{
		`{onlevel "error" "{bg 1}"}${level} %s`,
		`{onerror "{bg 1}"}${level} %s`,
	}

	for _, format := range formats {
		buffer, log := setupBufferedLogger(t, format)

		log.Error(`hello`)

		test.Equal(
			compileExpectedStyle(t, "{bg 1}ERROR hello\n"),
			buffer.String(),
		)

		buffer.Reset()

		log.Info(`ok`)
		test.Equal("INFO ok\n", buffer.String())
	}

}

func TestNewFormat_CanRestoreStyleAfterSuccessfullLevelMatch(t *testing.T) {
	test := assert.New(t)

	buffer, log := setupBufferedLogger(
		t,
		`{fg 1}{onlevel "error" "{bg 199}"}${level}{restore} %s`,
	)

	log.Error(`hello`)

	test.Equal(
		compileExpectedStyle(
			t,
			"{fg 1}{bg 199}ERROR{fg 1}{nobg}{nobold}{noreverse} hello\n",
		),
		buffer.String(),
	)
}

func TestNewFormat_CanRestoreStyleAfterTwoSuccessfullLevelMatch(t *testing.T) {
	test := assert.New(t)

	buffer, log := setupBufferedLogger(
		t,
		`{onerror "{fg 1}"}{store}XXX {onerror "{bg 199}"}${level}{restore} %s`,
	)

	log.Error(`hello`)

	test.Equal(
		compileExpectedStyle(
			t,
			"{fg 1}XXX {bg 199}ERROR"+
				"{nofg}{nobg}{nobold}{noreverse}{fg 1} hello\n",
		),
		buffer.String(),
	)
}

func TestNewFormat_CanRestoreStyleAfterNoLevelMatches(t *testing.T) {
	test := assert.New(t)

	buffer, log := setupBufferedLogger(
		t,
		`{fg 2}{onerror "{fg 1}"}XXX {onerror "{bg 199}"}${level}{restore} %s`,
	)

	log.Error(`change state`)
	buffer.Reset()

	log.Info(`hello`)

	test.Equal(
		compileExpectedStyle(
			t,
			"{fg 2}XXX INFO{fg 2}{nobg}{nobold}{noreverse} hello\n",
		),
		buffer.String(),
	)
}

func TestNewFormat_CanStoreConditionalStyle(t *testing.T) {
	test := assert.New(t)

	buffer, log := setupBufferedLogger(
		t,
		`{oninfo "{fg 1}"}{store}XXX {onerror "{bg 199}"}${level}{restore} %s`,
	)

	log.Info(`hello`)

	test.Equal(
		compileExpectedStyle(
			t,
			"{fg 1}XXX INFO{nofg}{nobg}{nobold}{noreverse}{fg 1} hello\n",
		),
		buffer.String(),
	)
}

func setupBufferedLogger(
	t *testing.T,
	formatting string,
) (*bytes.Buffer, *lorg.Log) {
	format, err := Format(formatting)
	assert.NoError(t, err)

	log := lorg.NewLog()
	log.SetFormat(format)

	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)

	return buffer, log
}

func compileExpectedStyle(t *testing.T, style string) string {
	expected, err := loreley.CompileAndExecuteToString(
		style, nil, nil,
	)

	assert.NoError(t, err)

	return expected
}
