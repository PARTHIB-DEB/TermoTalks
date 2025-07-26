package cmd

import "github.com/spf13/cobra"

var (
	firstname string
	lastname  string
	email     string
	password  string
	// username  string = ""
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
			cmd.Println("Viewing profile...")
		} else if args[0] == "create" {
			// Logic to create profile
			cmd.Println("Creating profile...")
			// CreateProfile(firstname, lastname, email, password)
		} else if args[0] == "update" {
			// Logic to update profile
			cmd.Println("Updating profile...")
			// UpdateProfile(firstname, lastname, email, password, username)
		} else if args[0] == "delete" {
			// Logic to delete profile
			cmd.Println("Deleting profile...")
			// DeleteProfile(username)
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
}
