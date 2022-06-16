package cmd

import (
	"fmt"
	"os/exec"

	"github.com/albertkol/ua-qa-go/state"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Echo application state",
	Long:  `Display application state information`,
	Run: func(cmd *cobra.Command, args []string) {
		// Command needs authentication
		if state.Status.LoggedIn == false {
			fmt.Println(state.CGREY + "Please login" + state.CEND)

			return
		}

		// Make user info request and echo it
		fmt.Println("user:")
		info, _ := exec.Command("contract", "call", "get-user-info", "--email="+state.Status.Email).Output()
		fmt.Println(string(info))

		// Make get accounts request and echo it
		accounts, _ := exec.Command("contract", "call", "get-accounts", "--email="+state.Status.Email).Output()
		fmt.Println(string(accounts))

		// Echo env values of the App
		fmt.Println("environment variables:")
		fmt.Println("- marketplace:", state.CGREY, state.Status.Marketplace, state.CEND)
		fmt.Println("- email:", state.CGREY, state.Status.Email, state.CEND)
		fmt.Println("- account:", state.CGREY, state.Status.AccountID, state.CEND)
		fmt.Println("- subscription:", state.CGREY, state.Status.SubscriptionID, state.CEND)
		fmt.Println("- contract:", state.CGREY, state.Status.ContractID, state.CEND)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
