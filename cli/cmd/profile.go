package cmd

import (
	"fmt"
	"log"

	View "github.com/PARTHIB-DEB/TermoTalks/cli/view"
	"github.com/spf13/cobra"
)

var (
	firstname string = ""
	lastname  string = ""
	email     string = ""
	password  string = ""
	username  string = ""
)

var ProfCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage user profiles",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Println("Please provide a subcommand: view, create, update, or delete")
			return
		} else if args[0] == "view" {
			// Logic to view profile
			uname := cmd.Flags().Lookup("username")
			if uname != nil && uname.Value.String() != "" {
				cmd.Println("Viewing profile for user:", uname.Value.String())
			} else {
				log.Fatal("Only Username is required to view your profile")
			}
		} else if args[0] == "create" {
			var pwd2 string
			fmt.Print("Enter Password Again :")
			fmt.Scanf("%s", &pwd2)
			View.CreateProfileView(firstname, lastname, email, password, pwd2)
		} else if args[0] == "update" {
			View.UpdateProfileView(firstname, lastname, email, password, username)
		} else if args[0] == "delete" {
			uname := cmd.Flags().Lookup("username")
			if uname != nil && uname.Value.String() != "" {
				View.DeleteProfileView(uname.Value.String())
			} else {
				log.Fatal("Only Username is required to delete your profile")
			}
		} else {
			cmd.Help()
		}
	},
}

func init() {
	ProfCmd.Flags().StringVarP(&firstname, "firstname", "f", "", "First name of the user")
	ProfCmd.Flags().StringVarP(&lastname, "lastname", "l", "", "Last name of the user")
	ProfCmd.Flags().StringVarP(&email, "email", "e", "", "Email address of the user")
	ProfCmd.Flags().StringVarP(&password, "password", "p", "", "Password for the user profile")
	ProfCmd.Flags().StringVarP(&username, "username", "u", "", "Username for the user profile")
}
