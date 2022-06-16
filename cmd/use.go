package cmd

import (
	"fmt"

	"github.com/albertkol/ua-qa-go/state"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Set the user",
	Long:  `Insert the user you wish to work on by email address.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Command needs authentication
		if state.Status.LoggedIn == false {
			fmt.Println(state.CGREY + "Please login" + state.CEND)

			return
		}

		// Store email from user info
		state.Status.Email = emailPromt()

		// Refresh app with the new email
		initCmd.Run(initCmd, args)

		if state.Status.Email == "" {
			// Refresh app one more time because email was invalid
			initCmd.Run(initCmd, args)
			fmt.Println(state.CRED + "Warning: Restored initial state" + state.CEND)
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func emailPromt() string {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . }} ",
		Success: "{{ . }} ",
	}

	prompt := promptui.Prompt{
		Label:     state.CBLUE + "Email:" + state.CEND,
		Templates: templates,
	}

	result, _ := prompt.Run()

	return result
}
