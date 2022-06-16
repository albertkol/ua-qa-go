package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/albertkol/ua-qa-go/state"
	"github.com/spf13/cobra"
)

// attachRenewalCmd represents the attachRenewal command
var attachRenewalCmd = &cobra.Command{
	Use:   "attach-renewal",
	Short: "Attach renewal",
	Long:  `Command to attach renewal to the user in 'use'.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Command needs authentication
		if state.Status.LoggedIn == false {
			fmt.Println(state.CGREY + "Please login" + state.CEND)

			return
		}

		var bash string
		var resp []byte

		// To attach renewals user mustn't have purchase account
		if state.Status.AccountID != "" {
			fmt.Println(state.CRED + "Warning: To attach renewals user mustn't have purchase account" + state.CEND)
			fmt.Println(state.CRED + "Warning: Starting reset prcess" + state.CEND)

			// Are you sure you want to reset user
			answ := yesNoPromt("Are you sure you want to reset user? (yes/no)")
			if answ[0:1] != "y" {
				fmt.Println(state.CRED + "Error: User was not reset" + state.CEND)

				return
			}

			// Reset user
			resetCmd.Run(resetCmd, args)
		}

		// Choose offer type
		option := selectPromt([]string{"single", "multiple", "not actionable", "expired"})

		pdata := map[string]string{
			"random_account_id": "",
			"random_asset_id":   "",
			"random_asset_id_2": "",
			"contract_start":    "",
			"contract_start_2":  "",
			"renewal_start":     "",
			"renewal_start_2":   "",
			"contract_end":      "",
			"contract_end_2":    "",
			"renewal_end":       "",
			"renewal_end_2":     "",
		}

		now := time.Now()
		pdata["random_account_id"] = "ua_qa_account_" + now.Format("2006-01-02T15:04:05Z")
		pdata["random_asset_id"] = "ua_qa_asset_" + now.Format("2006-01-02T15:04:05Z")
		pdata["random_asset_id_2"] = ""

		pdata["contract_start"] = now.AddDate(0, -6, 0).Format("2006-01-02T15:04:05Z")
		pdata["contract_end"] = now.AddDate(0, 6, 0).Format("2006-01-02T15:04:05Z")
		pdata["renewal_start"] = now.AddDate(0, -3, 0).Format("2006-01-02T15:04:05Z")
		pdata["renewal_end"] = now.AddDate(0, 3, 0).Format("2006-01-02T15:04:05Z")

		file := "renewal"
		if option == "not actionable" {
			file = "renewal-no-action"
		} else if option == "multiple" {
			file = "renewal-multi"
			pdata["contract_start_2"] = now.AddDate(0, -4, 0).Format("2006-01-02T15:04:05Z")
			pdata["contract_end_2"] = now.AddDate(0, 8, 0).Format("2006-01-02T15:04:05Z")
			pdata["renewal_start_2"] = now.AddDate(0, -1, 0).Format("2006-01-02T15:04:05Z")
			pdata["renewal_end_2"] = now.AddDate(0, 5, 0).Format("2006-01-02T15:04:05Z")
		} else if option == "expired" {
			file = "renewal-expired"
			pdata["contract_start"] = now.AddDate(-3, 0, 0).Format("2006-01-02T15:04:05Z")
			pdata["contract_end"] = now.AddDate(-2, 0, 0).Format("2006-01-02T15:04:05Z")
			pdata["renewal_start"] = now.AddDate(-2, -10, 0).Format("2006-01-02T15:04:05Z")
			pdata["renewal_end"] = now.AddDate(2, 10, 0).Format("2006-01-02T15:04:05Z")
		}

		// Make create renewal request
		var payload string
		payload = getPayload(file, pdata)

		bash = "echo '" + payload + "' | contract call EnsureContractForExternalAccount " + pdata["random_account_id"] + " -"
		resp, _ = exec.Command("bash", "-c", bash).Output()
		fmt.Println(string(resp))

		// Unmarshal is ideal but it won't work :(
		// var acdata = new(params.AccountContractInfo)
		// err0 := yaml.Unmarshal(resp, &acdata)

		r, _ := regexp.Compile("id: a([a-zA-Z0-9-_]*)")
		matched_strings := r.FindAllString(string(resp), -1)

		// Because mashalling response was erroring out we regex for the accountID
		state.Status.AccountID = strings.Replace(matched_strings[0], "id: ", "", 1)
		fmt.Println(state.CGREEN + "Set AccountID: " + state.Status.AccountID + state.CEND)

		payload = `{"email": ` + state.Status.Email + `, "role": "admin"}`
		bash = "echo '" + payload + "' | contract call SetAccountUserRole " + state.Status.AccountID + " -"
		resp, _ = exec.Command("bash", "-c", bash).Output()

		fmt.Println(string(resp))
		fmt.Println(state.CGREEN + "Rewnal was attached" + state.CEND)
	},
}

func init() {
	rootCmd.AddCommand(attachOfferCmd)
}
