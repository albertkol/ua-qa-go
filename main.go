package main

import (
	"os"

	"github.com/albertkol/ua-qa-go/cmd"
)

func main() {
	os.Setenv("CONTRACTS_URL", "https://contracts.staging.canonical.com/")

	cmd.Execute()
}
