package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"text/template"
	"time"
)

func init() {
	// Mute commands.
	addCmd.SetOut(new(bytes.Buffer))
	addCmd.SetErr(new(bytes.Buffer))
	initCmd.SetOut(new(bytes.Buffer))
	initCmd.SetErr(new(bytes.Buffer))
}

// ensureLF converts any \r\n to \n
func ensureLF(content []byte) []byte {
	return bytes.Replace(content, []byte("\r\n"), []byte("\n"), -1)
}

// compareFiles compares the content of files with pathA and pathB.
// If contents are equal, it returns nil.
// If not, it returns which files are not equal
// and diff (if system has diff command) between these files.
func compareFiles(generatedFile, goldenFile string) error {
	contentA, err := ioutil.ReadFile(generatedFile)
	if err != nil {
		return err
	}
	tmpl, err := template.ParseFiles(goldenFile)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	templateData := map[string]string{
		"Year": fmt.Sprintf("%d", time.Now().Year()),
	}

	if err := tmpl.Execute(&buf, templateData); err != nil {
		return err
	}
	contentB := buf.Bytes()

	if !bytes.Equal(ensureLF(contentA), ensureLF(contentB)) {
		output := new(bytes.Buffer)
		output.WriteString(fmt.Sprintf("%q and %q are not equal!\n\n", generatedFile, goldenFile))

		diffPath, err := exec.LookPath("diff")
		if err != nil {
			// Don't execute diff if it can't be found.
			return nil
		}
		diffCmd := exec.Command(diffPath, "-u", "--strip-trailing-cr", generatedFile, goldenFile)
		diffCmd.Stdout = output
		diffCmd.Stderr = output

		output.WriteString("$ diff -u " + generatedFile + " " + goldenFile + "\n")
		if err := diffCmd.Run(); err != nil {
			output.WriteString("\n" + err.Error())
		}
		return errors.New(output.String())
	}
	return nil
}
