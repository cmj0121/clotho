// The Clotho library for OSINT collector.
package clotho

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cmj0121/clotho/internal/github"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// The entry point of the Clotho.
type Clotho struct {
	// the logger options
	Quiet   bool `short:"q" group:"logger" xor:"verbose,quiet" help:"Disable all logger."`
	Verbose int  `short:"v" group:"logger" xor:"verbose,quiet" type:"counter" help:"Show the verbose logger."`

	Github *github.GitHub `cmd:"" help:"The GitHub user collector."`
}

// Create the new instance of the Clotho.
func New() *Clotho {
	/// the default constructor
	return &Clotho{}
}

// run the Clotho on command-line and return the exit code.
func (c *Clotho) Run() (exitcode int) {
	/// execute the CLI parser by kong
	kong.Parse(c)

	c.prologue()
	defer c.epilogue()

	var command SubCommand

	switch {
	case c.Github != nil:
		command = c.Github
	default:
		log.Error().Msg("No command is specified.")
		exitcode = 1
		return
	}

	exitcode = c.run(command)
	return
}

// the extra options for the Kong when parsing the command-line.
func (c *Clotho) AfterApply() (err error) {
	if c.Quiet {
		c.Verbose = -1
	}

	return
}

// run the subcommand.
func (c *Clotho) run(cmd SubCommand) (exitcode int) {
	cmd.Prologue()
	defer cmd.Epilogue()

	resp, err := cmd.Execute()
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute the command.")
		exitcode = 1
		return
	}

	log.Info().Interface("data", resp).Msg("The result of the command.")

	// show the result as JSON
	payload, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(payload))

	return
}

// setup everything before the execution.
func (c *Clotho) prologue() {
	c.setupLogger()
}

// clean-up everything after the execution.
func (c *Clotho) epilogue() {
}

// setup the logger subsystem, depends on the options.
func (c *Clotho) setupLogger() {
	writter := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = zerolog.New(writter).With().Timestamp().Logger()

	switch c.Verbose {
	case -1:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case 0:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}
