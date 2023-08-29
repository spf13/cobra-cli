package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

var customFileTemplate = []byte(`/*
{{ .Project.Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// {{ .CmdName }}Cmd represents the {{ .CmdName }} command
var {{ .CmdName }}Cmd = &myStruct.Command{
	Use:   "{{ .CmdName }}",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("{{ .CmdName }} called")
	},
	CustomParam: "This is a custom parameter to extend the default cobra functionality",
}

func init() {
	{{ .CmdParent }}.AddCommand({{ .CmdName }}Cmd)

    // All commands must indepenently define this flag by default
	testCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Thank you for contributing a new command.
}
`)

func TestGoldenAddCmd(t *testing.T) {
	viper.Set("useViper", true)
	viper.Set("license", "apache")
	command := &Command{
		CmdName:   "test",
		CmdParent: parentName,
		Project:   getProject(),
	}
	defer os.RemoveAll(command.AbsolutePath)

	assertNoErr(t, command.Project.Create())
	assertNoErr(t, command.Create())

	generatedFile := fmt.Sprintf("%s/cmd/%s.go", command.AbsolutePath, command.CmdName)
	goldenFile := fmt.Sprintf("testdata/%s.go.golden", command.CmdName)
	err := compareFiles(generatedFile, goldenFile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGoldenCustomAddCmd(t *testing.T) {
	viper.Set("useViper", true)
	viper.Set("license", "apache")
	command := &Command{
		CmdName:   "test",
		CmdParent: parentName,
		Project:   getProject(),
	}
	defer os.RemoveAll(command.AbsolutePath)

	assertNoErr(t, command.Project.Create())

	templateFile := fmt.Sprintf("%s/.cobra_template.tpl", command.AbsolutePath)
	err := os.WriteFile(templateFile, customFileTemplate, 0644)
	if err != nil {
		t.Fatal(err)
	}

	assertNoErr(t, command.Create())

	generatedFile := fmt.Sprintf("%s/cmd/%s.go", command.AbsolutePath, command.CmdName)
	goldenFile := fmt.Sprintf("testdata/%s.go.golden", command.CmdName)
	err = compareFiles(generatedFile, goldenFile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidateCmdName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"cmdName", "cmdName"},
		{"cmd_name", "cmdName"},
		{"cmd-name", "cmdName"},
		{"cmd______Name", "cmdName"},
		{"cmd------Name", "cmdName"},
		{"cmd______name", "cmdName"},
		{"cmd------name", "cmdName"},
		{"cmdName-----", "cmdName"},
		{"cmdname-", "cmdname"},
	}

	for _, testCase := range testCases {
		got := validateCmdName(testCase.input)
		if testCase.expected != got {
			t.Errorf("Expected %q, got %q", testCase.expected, got)
		}
	}
}
