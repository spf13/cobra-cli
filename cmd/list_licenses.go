/*
Copyright © 2026 Matthew Stobbs <matthew@stobbs.ca>

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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// listLicensesCmd represents the listLicenses command
var listLicensesCmd = &cobra.Command{
	Use:     "list_licenses",
	Aliases: []string{"list", "licenses"},
	Short:   "List available licenses",
	Long:    `List licenses that are registered`,
	Run: func(_ *cobra.Command, args []string) {
		var listSB strings.Builder
		listSB.WriteString("List of defined licenses:\n")
		for k, v := range Licenses {
			listSB.WriteString("\nName: ")
			listSB.WriteString(k)
			listSB.WriteString("\n  Proper name: ")
			listSB.WriteString(v.Name)
			listSB.WriteString("\n  Matches: ")
			listSB.WriteString(strings.Join(v.PossibleMatches, " "))
			listSB.WriteRune('\n')
		}
		fmt.Println(listSB.String())
	},
}

func init() {
	rootCmd.AddCommand(listLicensesCmd)
}
