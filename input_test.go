package survey

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestInputRender(t *testing.T) {

	tests := []struct {
		title    string
		prompt   Input
		data     InputTemplateData
		expected string
	}{
		{
			"Test Input question output without default",
			Input{Message: "What is your favorite month:"},
			InputTemplateData{Icons: &defaultIconSet},
			fmt.Sprintf("%s What is your favorite month: ", defaultIconSet.Question),
		},
		{
			"Test Input question output with default",
			Input{Message: "What is your favorite month:", Default: "April"},
			InputTemplateData{Icons: &defaultIconSet},
			fmt.Sprintf("%s What is your favorite month: (April) ", defaultIconSet.Question),
		},
		{
			"Test Input answer output",
			Input{Message: "What is your favorite month:"},
			InputTemplateData{Answer: "October", ShowAnswer: true, Icons: &defaultIconSet},
			fmt.Sprintf("%s What is your favorite month: October\n", defaultIconSet.Question),
		},
		{
			"Test Input question output without default but with help hidden",
			Input{Message: "What is your favorite month:", Help: "This is helpful"},
			InputTemplateData{Icons: &defaultIconSet},
			fmt.Sprintf("%s What is your favorite month: [%s for help] ", defaultIconSet.Question, string(defaultIconSet.HelpInput)),
		},
		{
			"Test Input question output with default and with help hidden",
			Input{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			InputTemplateData{Icons: &defaultIconSet},
			fmt.Sprintf("%s What is your favorite month: [%s for help] (April) ", defaultIconSet.Question, string(defaultIconSet.HelpInput)),
		},
		{
			"Test Input question output without default but with help shown",
			Input{Message: "What is your favorite month:", Help: "This is helpful"},
			InputTemplateData{ShowHelp: true, Icons: &defaultIconSet},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: ", defaultIconSet.Help, defaultIconSet.Question),
		},
		{
			"Test Input question output with default and with help shown",
			Input{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			InputTemplateData{ShowHelp: true, Icons: &defaultIconSet},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: (April) ", defaultIconSet.Help, defaultIconSet.Question),
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.Input = test.prompt
		err = test.prompt.Render(
			InputQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}

func TestInputPrompt(t *testing.T) {
	tests := []PromptTest{
		{
			"Test Input prompt interaction",
			&Input{
				Message: "What is your name?",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird")
				c.ExpectEOF()
			},
			"Larry Bird",
		},
		{
			"Test Input prompt interaction with default",
			&Input{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("")
				c.ExpectEOF()
			},
			"Johnny Appleseed",
		},
		{
			"Test Input prompt interaction overriding default",
			&Input{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird")
				c.ExpectEOF()
			},
			"Larry Bird",
		},
		{
			"Test Input prompt interaction and prompt for help",
			&Input{
				Message: "What is your name?",
				Help:    "It might be Satoshi Nakamoto",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("?")
				c.ExpectString("It might be Satoshi Nakamoto")
				c.SendLine("Satoshi Nakamoto")
				c.ExpectEOF()
			},
			"Satoshi Nakamoto",
		},
		{
			// https://en.wikipedia.org/wiki/ANSI_escape_code
			// Device Status Report - Reports the cursor position (CPR) to the
			// application as (as though typed at the keyboard) ESC[n;mR, where n is the
			// row and m is the column.
			"Test Input prompt with R matching DSR",
			&Input{
				Message: "What is your name?",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("R")
				c.ExpectEOF()
			},
			"R",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}
