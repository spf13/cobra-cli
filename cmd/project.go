package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/tpl"
)

type (
	// Project contains name, license and paths to projects.
	Project struct {
		// v2
		PkgName      string
		Copyright    string
		AbsolutePath string
		Legal        License
		Viper        bool
		AppName      string
	}

	Command struct {
		CmdName   string
		CmdParent string
		*Project
	}

	createFileFunc func(p *Project) error
)

var (
	projectFiles = map[string]createFileFunc{
		"%s/main.go":     createMain,
		"%s/cmd/root.go": createRootCmd,
		"%s/LICENSE":     createLicenseFile,
	}
)

func createMain(p *Project) error {
	// create main.go
	mainFile, err := os.Create(fmt.Sprintf("%s/main.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.MainTemplate())))
	return mainTemplate.Execute(mainFile, p)
}

func createRootCmd(p *Project) error {
	// create cmd/root.go
	if _, err := os.Stat(fmt.Sprintf("%s/cmd", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/cmd", p.AbsolutePath), 0751))
	}
	rootFile, err := os.Create(fmt.Sprintf("%s/cmd/root.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer rootFile.Close()

	rootTemplate := template.Must(template.New("root").Parse(string(tpl.RootTemplate())))
	return rootTemplate.Execute(rootFile, p)
}

func createLicenseFile(p *Project) error {
	data := map[string]interface{}{
		"copyright": copyrightLine(),
	}
	licenseFile, err := os.Create(fmt.Sprintf("%s/LICENSE", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer licenseFile.Close()

	licenseTemplate := template.Must(template.New("license").Parse(p.Legal.Text))
	return licenseTemplate.Execute(licenseFile, data)
}

func (p *Project) Create(force bool) error {
	// check if AbsolutePath exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, 0754); err != nil {
			return err
		}
	}

	// Check to make sure we don't overwrite things unless we have --force
	if !force {
		for path, _ := range projectFiles {
			abspath := fmt.Sprintf(path, p.AbsolutePath)
			if _, err := os.Stat(abspath); err == nil {
				return fmt.Errorf("%s already exists; use --force to overwrite", abspath)
			}
		}
	}

	for _, createFunc := range projectFiles {
		if err := createFunc(p); err != nil {
			return err
		}
	}

	return nil
}

func (c *Command) Create(force bool) error {
	abspath := fmt.Sprintf("%s/cmd/%s.go", c.AbsolutePath, c.CmdName)
	if _, err := os.Stat(abspath); err == nil && !force {
		return fmt.Errorf("%s already exists; use --force to overwrite", abspath)
	}
	cmdFile, err := os.Create(fmt.Sprintf("%s/cmd/%s.go", c.AbsolutePath, c.CmdName))
	if err != nil {
		return err
	}
	defer cmdFile.Close()

	commandTemplate := template.Must(template.New("sub").Parse(string(tpl.AddCommandTemplate())))
	err = commandTemplate.Execute(cmdFile, c)
	if err != nil {
		return err
	}
	return nil
}
