// The main function of the Clotho.
package main

import (
	"os"

	"github.com/cmj0121/clotho"
)

func main() {
	clotho := clotho.New()
	os.Exit(clotho.Run())
}
