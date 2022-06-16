package cmd

import (
	"fmt"
	"os/exec"

	"github.com/albertkol/ua-qa-go/state"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to ua-contracts",
	Long: `This command promts the user to the SSO link to login for the 
ua-contracts API.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Make login request
		login, err := exec.Command("contract", "login").Output()

		if err != nil {
			fmt.Println(err)

			return
		}

		fmt.Println(string(login))

		// Store LoggedIn status
		state.Status.LoggedIn = true
		fmt.Println(state.CGREEN + "Set LoggedIn: true" + state.CEND)

		// Init data for logged in user and store it in the App env
		initCmd.Run(initCmd, args)

		// Echo user info and App env
		statusCmd.Run(statusCmd, args)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
