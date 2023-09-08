/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	datastore "ToDo/data"
	style "ToDo/defs"
	"ToDo/structs"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var scanner = bufio.NewScanner(os.Stdin)

// addTasksCmd represents the addTask command
var addTasksCmd = &cobra.Command{
	Use:   "addTasks",
	Short: "Add to the list of tasks.",
	Long:  `Add a task to the list of tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		inputs := []string{}

		fmt.Println(style.Default_Style.Render("enter tasks"))
		fmt.Println(style.Remark_Style.Render("enter !done to finish"))
		var input string
		for scanner.Scan() {
			input = scanner.Text()
			if input == "!done" {
				break
			}
			inputs = append(inputs, input)
		}

		if len(inputs) == 0 {
			return
		}

		fmt.Println(style.Default_Style.Render("you entered:"))
		for _, element := range inputs {
			fmt.Println(style.Remark_Style.Render("• ", element))
		}

		fmt.Println(style.Default_Style.Render("save these tasks? (y/n)"))
		var confirmed bool
		for {
			fmt.Printf("save? : ")
			scanner.Scan()
			input = scanner.Text()
			if strings.ToLower(input) == "y" {
				confirmed = true
				break
			}
			if strings.ToLower(input) == "n" {
				confirmed = false
				break
			}
		}

		if confirmed {
			for _, element := range inputs {
				datastore.AddToCache(structs.Task{
					Name:        element,
					Description: "",
				})
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addTasksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addTaskCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addTaskCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
