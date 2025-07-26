package cmd

import "github.com/spf13/cobra"

var (
	senderusername   string
	receiverusername string
	alllinks         bool = false
)

var LinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Manage meeting links between two users",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Println("Please provide a subcommand: 'send' or 'get'")
		} else if args[0] == "send" {
			// Logic to create a link
			cmd.Println("Creating link...")
			// Createandsendlink(senderusername, receiverusername)
		} else if args[0] == "get" {
			// Logic to get a link
			cmd.Println("Retrieving link...")
			// Getlink(receiverusername, alllinks)
		}
	},
}

func init() {
	LinkCmd.Flags().StringVarP(&senderusername, "sender", "s", "", "Username of the sender")
	LinkCmd.Flags().StringVarP(&receiverusername, "receiver", "r", "", "Username of the receiver")
	LinkCmd.Flags().BoolVarP(&alllinks, "all", "a", false, "Retrieve all links for the user")
}
