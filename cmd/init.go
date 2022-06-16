package cmd

import (
	"fmt"
	"os/exec"

	"github.com/CanonicalLtd/ua-contracts/params"
	"github.com/albertkol/ua-qa-go/state"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init user in the app state",
	Long: `If you are logged in already set the user in the state of the state. 
If you are logged in already get the user and store it in the state of the state.`,
	Run: func(cmd *cobra.Command, args []string) {
		var bash string
		var resp []byte

		// Make user info request
		bash = "contract call get-user-info --email=" + state.Status.Email
		resp, _ = exec.Command("bash", "-c", bash).Output()
		udata := make(map[string]string)
		yaml.Unmarshal(resp, &udata)

		// If not logged in or user has no email
		if udata["code"] == "unauthorized" || udata["email"] == "" {
			if state.Status.LoggedIn == true {
				fmt.Println(state.CRED + "Error: Unauthroized or invalid email" + state.CEND)
				state.Status.Email = ""
				fmt.Println(state.CGREY + "Set Email: ''" + state.CEND)
			} else {
				fmt.Println(state.CGREY + "Please login" + state.CEND)
			}

			return
		}

		// Store LoggedIn status
		if state.Status.LoggedIn == false {
			state.Status.LoggedIn = true
			fmt.Println(state.CGREEN + "Set LoggedIn: true" + state.CEND)
		}

		// Store email from user info
		if state.Status.Email == "" {
			state.Status.Email = udata["email"]
			fmt.Println(state.CGREEN + "Set Email: " + state.Status.Email + state.CEND)
		}

		// Clear state first in case init is called later as refresh
		state.Status.AccountID = ""
		fmt.Println(state.CGREY + "Set AccountID: ''" + state.CEND)
		state.Status.SubscriptionID = ""
		fmt.Println(state.CGREY + "Set SubscriptionID: ''" + state.CEND)
		state.Status.ContractID = ""
		fmt.Println(state.CGREY + "Set ContractID: ''" + state.CEND)

		// Make get accounts request
		bash = "contract call get-accounts --email=" + state.Status.Email
		resp, _ = exec.Command("bash", "-c", bash).Output()
		var accdata = new(params.GetAccountsResponse)
		yaml.Unmarshal(resp, &accdata)

		// Store paid account in app env
		for _, v := range accdata.Accounts {
			if v.Type == "paid" {
				state.Status.AccountID = v.Id
				fmt.Println(state.CGREEN + "Set AccountID: " + state.Status.AccountID + state.CEND)

				return
			}
		}

		fmt.Println(state.CRED + "Error: No purchase account found" + state.CEND)

		// If use has no purchase account
		var answ string

		// Initialize purchase account?
		answ = yesNoPromt("Initialize purchase account? (yes/no):")
		if answ[0:1] != "y" {
			fmt.Println(state.CGREY + "Purchase account was not initialised" + state.CEND)

			return
		}

		var payload string
		payload = getPayload("ensure_account", map[string]string{"email": state.Status.Email})
		bash = "echo '" + payload + "' | contract call EnsureAccountForMarketplace " + state.Status.Marketplace + " -"
		resp, _ = exec.Command("bash", "-c", bash).Output()
		ensdata := make(map[string]string)
		yaml.Unmarshal(resp, &ensdata)

		state.Status.AccountID = ensdata["accountID"]
		fmt.Println(state.CGREEN + "Purchase account was initialised" + state.CEND)
		fmt.Println(state.CGREEN + "Set AccountID: " + state.Status.AccountID + state.CEND)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
