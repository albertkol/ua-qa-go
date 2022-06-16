package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/albertkol/ua-qa-go/state"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ua-qa-go",
	Short: "UA QA Tool - Making UA QA easier since last Friday",
	Long: `
                 ,---,                     ,----..        ,---, 
         ,--,   '  .' \                   /   /   \      '  .' \ 
       ,'_ /|  /  ;    '.        ,---,.  /   .     :    /  ;    '. 
  .--. |  | : :  :       \     ,'  .' | .   /   ;.  \  :  :       \ 
,'_ /| :  . | :  |   /\   \  ,---.'   ,.   ;   /  ' ;  :  |   /\   \ 
|  ' | |  . . |  :  ' ;.   : |   |    |;   |  ; \ ; |  |  :  ' ;.   : 
|  | ' |  | | |  |  ;/  \   \:   :  .' |   :  | ; | '  |  |  ;/  \   \ 
:  | | :  ' ; '  :  | \  \ ,':   |.'   .   |  ' ' ' :  '  :  | \  \ ,' 
|  ; ' |  | ' |  |  '  '--'  '---'     '   ;  \; /  |  |  |  '  '--' 
:  | : ;  ; | |  :  :                   \   \  ',  . \ |  :  : 
'  :  '--'   \|  | ,'                    ;   :      ; ||  | ,' 
:  ,      .-./'--''                       \   \ .'---" '--'' 
 '--'----'                                 '--- 


----------------------------------------------------

Welcome!

----------------------------------------------------

UA QA tool was designed with the goal to make 
difficult UA QA scenarios a child's play.

This is the Go Lang conversion of the 1st version of the UA QA Tool
written in Python.

Tool options:
    - login                 SSO login
    - logout                Logout
    - status                Echos user in use info and app envs
    - refresh               Re-initialize app with current user data
    - use                   Set user you want to work on
    - switch                Alias for 'use'
    - attach                Attach offer or renewal to user in use
    - clear                 Clear user in use accounts
    - exit                  Close ua-qa tool

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(cmd.Long)
		initCmd.Run(initCmd, args)

		for {
			command := uaCammandPromt()

			if command == "" {
				continue
			}

			if command == "login" {
				loginCmd.Run(loginCmd, args)

				continue
			}

			if command == "logout" {
				logoutCmd.Run(logoutCmd, args)

				continue
			}

			if command == "use" || command == "switch" {
				useCmd.Run(useCmd, args)

				continue
			}

			if command == "refresh" {
				initCmd.Run(initCmd, args)

				continue
			}

			if command == "status" {
				statusCmd.Run(statusCmd, args)

				continue
			}

			if command == "attach" {
				option := selectPromt([]string{"offer", "renewal"})

				if option == "offer" {
					attachOfferCmd.Run(attachOfferCmd, args)
				}

				if option == "renewal" {
					attachRenewalCmd.Run(attachRenewalCmd, args)
				}

				continue
			}

			if command == "reset" {
				resetCmd.Run(resetCmd, args)

				continue
			}

			if command == "exit" {
				fmt.Println(state.CGREEN + "Goodbye!")

				break
			}

			fmt.Println(state.CRED + "Error: Invalid command")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func uaCammandPromt() string {
	templates := &promptui.PromptTemplates{
		Prompt:  "🦗 {{ . }} ",
		Valid:   "🦗 {{ . }} ",
		Invalid: "🦗 {{ . }} ",
		Success: "🦗 {{ . }} ",
	}

	prompt := promptui.Prompt{
		Label:     state.CYELLOW + "ua-qa [" + state.Status.Email + "] >" + state.CEND,
		Templates: templates,
	}

	result, _ := prompt.Run()

	return result
}

func yesNoPromt(question string) string {
	prompt := promptui.Prompt{
		Label: state.CGREY + question + state.CEND,
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . }} ",
			Success: "{{ . }} ",
		},
	}

	result, _ := prompt.Run()

	return result
}

func selectPromt(items []string) string {
	prompt := promptui.Select{
		Label: "Type",
		Items: items,
	}

	_, result, _ := prompt.Run()

	return result
}

func getPayload(payload string, replaces map[string]string) string {
	input, err := ioutil.ReadFile("payloads/" + payload + ".json")
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		for from, with := range replaces {
			if strings.Contains(line, "${"+from+"}") {
				lines[i] = strings.Replace(line, "${"+from+"}", with, -1)
			}
		}
	}

	output := strings.Join(lines, "\n")

	return output
}
