/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	datastore "ToDo/data"
	"ToDo/structs"

	"github.com/spf13/cobra"
)

// readTasksCmd represents the readTasks command
var readTasksCmd = &cobra.Command{
	Use:   "readTasks",
	Short: "Read all stored tasks.",
	Long:  `Read all stored tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		menu := structs.Menu{
			Focused_Task_Id: 0,
			Tasks:           datastore.FetchAll(),
		}

		menu.Display()
	},
}

func init() {
	rootCmd.AddCommand(readTasksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readTasksCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readTasksCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
