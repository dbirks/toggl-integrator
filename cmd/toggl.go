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
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
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

		togglTimeData := getTogglTimeData()

		var togglTimeEntries []TimeEntry

		err := json.Unmarshal([]byte(togglTimeData), &togglTimeEntries)

		if err != nil {
			fmt.Println(err)
		}

		printTable(togglTimeEntries)
	},
}

type TimeEntry struct {
	Start       time.Time `json:"start"`
	Stop        time.Time `json:"stop"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
}

func init() {
	getCmd.AddCommand(togglCmd)
}

func getTogglTimeData() []byte {
	togglUsername := os.Getenv("TOGGL_USERNAME")
	togglPassword := os.Getenv("TOGGL_PASSWORD")

	req, err := http.NewRequest("GET",
		"https://api.track.toggl.com/api/v9/me/time_entries", nil)
	if err != nil {
		print(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(togglUsername, togglPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	return body
}

// Print the time entries as a formatted table to stdout
func printTable(timeEntries []TimeEntry) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Start", "Stop", "Description", "Duration"})

	for _, entry := range timeEntries {
		str := []string{entry.Start.String(), entry.Stop.String(), entry.Description, strconv.Itoa(entry.Duration)}
		table.Append(str)
	}
	table.Render()
}
