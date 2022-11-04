/*
Copyright © 2022 SplitJoin <support@splitjoin.com>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/janeczku/go-spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "This command will write you a commit message",
	Long: `This commit command will call the SplitJoin API and
write you a commit message based on the contents of your
changes staged for commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		endpoint, _ := cmd.Flags().GetString("endpoint")

		if endpoint == "" {
			endpoint = "https://www.splitjoin.com/api/v0"
		}

		accessToken := viper.GetString("SJ_ACCESS_TOKEN")

		if accessToken == "" {
			log.Fatal("No access token is configured.\n\n🛠  To fix: Create an access token at https://www.splitjoin.com, then set the following environment variables:\n\n  SJ_ACCESS_TOKEN=(your access token)\n\n   To verify that it's working, run this command again.")
		}

		// Make sure we're in a git repo
		isRepo, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()

		if strings.TrimSpace(string(isRepo)) != "true" {
			log.Fatal("This is not a git repository. Please run this command from the root of a git repository.")
		}

		// Check if there are any staged changes
		staged, err := exec.Command("git", "diff", "--staged").Output()

		if len(string(staged)) == 0 {
			log.Fatal("There are no staged changes. Please stage your changes before running this command.")
		}

		diff, err := exec.Command("git", "diff", "--staged", "--no-color").Output()

		// Now we need to send the staged changes to the SplitJoin API
		// and get back a commit message
		s := spinner.NewSpinner("Thinking a good commit message...")
		s.SetCharset([]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"})
		s.SetSpeed(100 * time.Millisecond)
		s.Start()

		values := map[string]string{"diff": string(diff)}
		json_data, err := json.Marshal(values)

		if err != nil {
			log.Fatal(err)
		}

		// Create a Bearer string by appending string access token
		var bearer = "Bearer " + accessToken

		// Create a new request using http
		req, err := http.NewRequest("POST", endpoint+"/prompt", bytes.NewBuffer(json_data))

		// add authorization header to the req
		req.Header.Add("Authorization", bearer)
		req.Header.Add("Content-Type", "application/json")

		// Send req using http Client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var res map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&res)

		// Stop the spinner
		s.Stop()

		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			fmt.Println("Unauthorized, please check your personal access token.")
		} else {
			fmt.Println(res["message"])
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
