/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	datastore "ToDo/data"
	debugger "ToDo/debug"
	"ToDo/structs"

	"github.com/spf13/cobra"
)

// readTasksCmd represents the readTasks command
var readTasksCmd = &cobra.Command{
	Use:   "readTasks",
	Short: "Read all stored tasks.",
	Long:  `Read all stored tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks := datastore.FetchAll()
		firstTask := &tasks[0]

		menu := structs.Menu{
			Focused_Task:  firstTask,
			Tasks:         tasks,
			OutputChannel: make(chan []structs.Task),
		}

		debugger.Trace("starting")
		go menu.Display()
		debugger.Trace("finsihed starting")

		for {
			debugger.Trace("for loop")
			tasks := <-menu.OutputChannel
			if tasks != nil {
				debugger.Trace("overwrote")
				datastore.SetCache(tasks)
			} else {
				debugger.Trace("escaped")
				break
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
