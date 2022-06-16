package cmd

import (
	"fmt"
	"os/exec"

	"github.com/albertkol/ua-qa-go/state"
	"github.com/spf13/cobra"
)

// attachOfferCmd represents the attachOffer command
var attachOfferCmd = &cobra.Command{
	Use:   "attach-offer",
	Short: "Attach offer",
	Long:  `Command to attach offer to the user in 'use'.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Command needs authentication
		if state.Status.LoggedIn == false {
			fmt.Println(state.CGREY + "Please login" + state.CEND)

			return
		}

		// Offers require purchase account
		if state.Status.AccountID == "" {
			fmt.Println(state.CRED + "Warning: Offer attach requires purchase account" + state.CEND)
			fmt.Println(state.CRED + "Warning: Starting init process" + state.CEND)

			// Start purchase account init
			initCmd.Run(initCmd, args)

			// Purchase account init failed
			if state.Status.AccountID == "" {
				fmt.Println(state.CRED + "Error: User has no purchase account" + state.CEND)

				return
			}
		}

		// Choose offer type
		option := selectPromt([]string{"single", "multiple"})

		file := "offer"
		if option == "multiple" {
			file = "offer-multi"
		}

		// Make create offer request
		payload := getPayload(file, map[string]string{"account_id": state.Status.AccountID})
		bash := "echo '" + payload + "' | contract call EnsureOffer " + state.Status.Marketplace + " -"
		resp, err := exec.Command("bash", "-c", bash).Output()

		if err != nil {
			fmt.Println(err)

			return
		}

		fmt.Println(string(resp))
		fmt.Println(state.CGREEN + "Offer was attached" + state.CEND)
	},
}

func init() {
	rootCmd.AddCommand(attachOfferCmd)
}
