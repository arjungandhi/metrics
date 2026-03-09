package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user profiles",
}

var userAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new user profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := s.AddUser(args[0]); err != nil {
			return err
		}

		defaultUser, _ := s.DefaultUser()
		if defaultUser == args[0] {
			fmt.Printf("Added user %q (set as default)\n", args[0])
		} else {
			fmt.Printf("Added user %q\n", args[0])
		}
		return nil
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all user profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := s.ListUsers()
		if err != nil {
			return err
		}
		if len(users) == 0 {
			fmt.Println("No users configured. Run 'metrics user add <name>' to create one.")
			return nil
		}

		defaultUser, _ := s.DefaultUser()
		for _, u := range users {
			marker := "  "
			if u.Name == defaultUser {
				marker = "* "
			}
			fmt.Printf("%s%s\n", marker, u.Name)
		}
		return nil
	},
}

var userSetDefaultCmd = &cobra.Command{
	Use:   "set-default <name>",
	Short: "Set the default user profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := s.SetDefaultUser(args[0]); err != nil {
			return err
		}
		fmt.Printf("Default user set to %q\n", args[0])
		return nil
	},
}

func init() {
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userSetDefaultCmd)
	rootCmd.AddCommand(userCmd)
}
