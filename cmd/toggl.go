/*
Copyright © 2022 David Birks <david@birks.dev>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// togglCmd represents the toggl command
var togglCmd = &cobra.Command{
	Use:   "toggl",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {

		var togglProjects = []TogglProject{}
		var togglTimeEntries = []TogglTimeEntry{}
		var err error

		togglProjectsJson := queryTogglApi("https://api.track.toggl.com/api/v9/me/projects")
		togglTimeEntriesJson := queryTogglApi("https://api.track.toggl.com/api/v9/me/time_entries")

		err = json.Unmarshal([]byte(togglProjectsJson), &togglProjects)
		if err != nil {
			log.WithError(err).Fatal("Unable to unmarshal the projects json.")
		}

		// Load just the json we want into the togglTimeEntries object
		err = json.Unmarshal([]byte(togglTimeEntriesJson), &togglTimeEntries)
		if err != nil {
			log.WithError(err).Fatal("Unable to unmarshal the time entries json.")
		}

		// Lookup ProjectName based on the ProjectId, and add it to each object in the array
		for index, togglTimeEntry := range togglTimeEntries {

			// Find the project name for the project id in the togglTimeEntries
			for _, togglProject := range togglProjects {
				if togglProject.Id == togglTimeEntry.ProjectId {
					togglTimeEntries[index].ProjectName = togglProject.Name
				}
			}

		}

		printTable(togglTimeEntries)
	},
}

type TogglTimeEntry struct {
	Start       time.Time `json:"start"`
	Stop        time.Time `json:"stop"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	ProjectId   int32     `json:"project_id"`
	ProjectName string
}

type TogglProject struct {
	Name  string `json:"name"`
	Id    int32  `json:"id"`
	Color string `json:"color"`
}

func init() {
	log.SetLevel(log.InfoLevel)

	// Cobra
	getCmd.AddCommand(togglCmd)
}

// Print the time entries as a formatted table to stdout
func printTable(timeEntries []TogglTimeEntry) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Start", "Stop", "Description", "Duration", "Project"})

	for _, entry := range timeEntries {

		log.Trace(fmt.Sprintf("proj name::::: %s", entry.ProjectName))

		str := []string{entry.Start.String(), entry.Stop.String(), entry.Description, strconv.Itoa(entry.Duration), entry.ProjectName}
		log.Trace(str)
		table.Append(str)
	}
	table.Render()
}
