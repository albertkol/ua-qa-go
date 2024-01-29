package cmd

import (
	"fmt"
	"os/exec"

	"github.com/albertkol/ua-qa-go/state"
	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset user account",
	Long:  `Detach account from user and remove account id from state of the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Command needs authentication
		if state.Status.LoggedIn == false {
			fmt.Println(state.CGREY + "Please login" + state.CEND)

			return
		}

		// Are you sure you want to reset user
		answ := yesNoPromt("Are you sure you want to reset user? (yes/no):")

		// If no do nothing
		if answ[0:1] != "y" {
			return
		}

		// Refresh app
		initCmd.Run(initCmd, args)

		// If yes detach account
		var payload string
		payload = `{"nameHint": ` + state.Status.Email + `, "email": ` + state.Status.Email + `, "role": "none", "explanation": "how dare you demand this"}`

		var bash string
		bash = "echo '" + payload + "' | contract call set-account-user-role " + state.Status.AccountID + " - --force=true"
		fmt.Println("")
		fmt.Println(string(bash))
		var resp []byte
		resp, _ = exec.Command("bash", "-c", bash).Output()
		fmt.Println(string(resp))

		// Reset account env from the App
		state.Status.AccountID = ""
		fmt.Println(state.CGREY + "Set AccountID: ''" + state.CEND)
		state.Status.SubscriptionID = ""
		fmt.Println(state.CGREY + "Set SubscriptionID: ''" + state.CEND)
		state.Status.ContractID = ""
		fmt.Println(state.CGREY + "Set ContractID: ''" + state.CEND)
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
