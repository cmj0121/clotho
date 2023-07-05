// The SubCommand interface is used to define the subcommand of the CLI.
package clotho

type SubCommand interface {
	// open the necessary resources for the subcommand.
	Prologue()

	// clean up the resources for the subcommand.
	Epilogue()

	// execute the subcommand and get the result.
	Execute() (resp map[string]interface{}, err error)
}
