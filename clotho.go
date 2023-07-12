// The Clotho library for OSINT collector.
package clotho

import (
	"encoding/json"
	"os"

	"github.com/cmj0121/clotho/internal/demo"
	"github.com/cmj0121/clotho/internal/github"
	"github.com/cmj0121/clotho/internal/linkedin"

	"github.com/alecthomas/kong"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// The entry point of the Clotho.
type Clotho struct {
	// show version and exit
	Version VersionFlag `short:"V" name:"version" help:"Print version info and quit"`

	Table bool `negatable:"" short:"t" help:"Print the result as table." default:"true"`

	// the logger options
	Quiet   bool `short:"q" group:"logger" xor:"verbose,quiet" help:"Disable all logger."`
	Verbose int  `short:"v" group:"logger" xor:"verbose,quiet" type:"counter" help:"Show the verbose logger."`

	Github   *github.GitHub     `cmd:"" help:"The GitHub user collector."`
	LinkedIn *linkedin.LinkedIn `cmd:"" name:"linkedin" help:"The LinkedIn user collector."`
	Demo     *demo.Demo         `cmd:"" name:"demo" help:"The Demo collector."`
}

// Create the new instance of the Clotho.
func New() *Clotho {
	/// the default constructor
	return &Clotho{}
}

// run the Clotho on command-line and return the exit code.
func (c *Clotho) Run() (exitcode int) {
	/// execute the CLI parser by kong
	ctx := kong.Parse(c)

	c.prologue()
	defer c.epilogue()

	var command SubCommand

	switch sub := ctx.Command(); sub {
	case "github <username>":
		command = c.Github
	case "linkedin <username>":
		command = c.LinkedIn
	case "demo <link>":
		command = c.Demo
	default:
		log.Error().Str("subcmd", sub).Msg("Sub-command not implemented.")
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

	switch c.Table {
	case true:
		// show the result as Table
		table := tablewriter.NewWriter(os.Stdout)

		// set table style
		table.SetHeader([]string{"Key", "Value"})
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)
		// table.SetAutoMergeCells(true)

		switch resp.(type) {
		case nil:
			// empty result and nothing to do
			return
		case [][]string:
			table.AppendBulk(resp.([][]string))
		case map[string]interface{}:
			for key, value := range resp.(map[string]interface{}) {
				switch value.(type) {
				case string:
					table.Append([]string{key, value.(string)})
				default:
					data, _ := json.Marshal(value)
					table.Append([]string{key, string(data)})
				}
			}
		default:
			log.Error().Interface("data", resp).Msg("The result is not a map.")
			exitcode = 1
			return
		}

		table.Render()
	case false:
		// show the result as JSON
		data, _ := json.Marshal(resp)
		os.Stdout.Write(data)
	}

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
