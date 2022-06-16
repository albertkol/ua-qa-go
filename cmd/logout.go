package cmd

import (
	"fmt"
	"os/exec"

	"github.com/albertkol/ua-qa-go/state"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout from ua-contracts",
	Long:  `This command logs the user out from the ua-contracts API.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Command needs authentication
		if state.Status.LoggedIn == false {
			fmt.Println(state.CGREY + "Please login" + state.CEND)

			return
		}

		// Make logout request
		logout := exec.Command("contract", "logout")
		logout.Run()

		// Reset App env
		state.Status.LoggedIn = false
		fmt.Println(state.CGREY + "Set LoggedIn: false" + state.CEND)
		state.Status.Email = ""
		fmt.Println(state.CGREY + "Set Email: ''" + state.CEND)
		state.Status.AccountID = ""
		fmt.Println(state.CGREY + "Set AccountID: ''" + state.CEND)
		state.Status.SubscriptionID = ""
		fmt.Println(state.CGREY + "Set SubscriptionID: ''" + state.CEND)
		state.Status.ContractID = ""
		fmt.Println(state.CGREY + "Set ContractID: ''" + state.CEND)
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
