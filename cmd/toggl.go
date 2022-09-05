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
	"time"

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
		// fmt.Println("toggl called")
		// getTogglTimes()
		togglUsername := os.Getenv("TOGGL_USERNAME")
		togglPassword := os.Getenv("TOGGL_PASSWORD")
		// println(togglUsername)
		// println(togglPassword)

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

		// fmt.Print(string(body))

		type TimeEntry struct {
			Start       time.Time `json:"start"`
			Stop        time.Time `json:"stop"`
			Description string    `json:"description"`
			Duration    int       `json:"duration"`
		}

		var togglApiResponse []TimeEntry

		// err := json.Unmarshal([]byte(body), &togglApiResponse)
		json.Unmarshal([]byte(body), &togglApiResponse)

		// if err != nil {
		// 	fmt.Println(err)
		// }

		fmt.Println(togglApiResponse)

	},
}

func init() {
	getCmd.AddCommand(togglCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// togglCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// togglCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTogglTimes() {
	fmt.Println("Hi mom")

}
