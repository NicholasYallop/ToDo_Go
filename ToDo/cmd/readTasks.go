/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	datastore "ToDo/data"
	"ToDo/structs"
	"fmt"

	"github.com/eiannone/keyboard"
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

		keys, err := keyboard.GetKeys(1)
		if err != nil {
			fmt.Println("error setting up keyboard buffer")
			panic(err)
		}
		defer keyboard.Close()

		for {
			menu.Print()
			fmt.Printf("\033[%dA", 2*len(menu.Tasks)-2*menu.Focused_Task_Id)

			keyEvent := <-keys
			if keyEvent.Err != nil {
				fmt.Println("error getting key")
				panic(err)
			}

			if menu.Focused_Task_Id != 0 {
				fmt.Printf("\033[%dA", 2*menu.Focused_Task_Id)
			}

			switch keyEvent.Key {
			case keyboard.KeyArrowUp:
				menu.MoveCursorUp()
			case keyboard.KeyArrowDown:
				menu.MoveCursorDown()
			case keyboard.KeyEsc:
				fmt.Printf("\033[%dB", 2*len(menu.Tasks))
				return
			}
		}
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
