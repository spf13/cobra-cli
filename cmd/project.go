package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra-cli/tpl"
)

// Project contains name, license and paths to projects.
type Project struct {
	// v2
	PkgName      string
	Copyright    string
	AbsolutePath string
	Legal        License
	Viper        bool
	AppName      string
	LocalVars    bool
}

type Command struct {
	CmdName   string
	CmdParent string
	*Project
}

func (p *Project) Create() error {
	// check if AbsolutePath exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.MkdirAll(p.AbsolutePath, 0754); err != nil {
			return err
		}
	}

	// create main.go
	mainFile, err := os.Create(fmt.Sprintf("%s/main.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.MainTemplate())))
	err = mainTemplate.Execute(mainFile, p)
	if err != nil {
		return err
	}

	// create cmd/root.go
	if _, err = os.Stat(fmt.Sprintf("%s/cmd", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.MkdirAll(fmt.Sprintf("%s/cmd", p.AbsolutePath), 0751))
	}
	rootFile, err := os.Create(fmt.Sprintf("%s/cmd/root.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer rootFile.Close()

	tmpl := tpl.RootTemplateGlobal()
	if p.LocalVars {
		tmpl = tpl.RootTemplateLocal()
	}
	rootTemplate := template.Must(template.New("root").Parse(string(tmpl)))
	err = rootTemplate.Execute(rootFile, p)
	if err != nil {
		return err
	}

	// create license
	return p.createLicenseFile()
}

func (p *Project) createLicenseFile() error {
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

func (c *Command) Create() error {
	filename := c.CmdName
	if c.CmdParent != "rootCmd" {
		filename = fmt.Sprintf("%s_%s", strings.TrimSuffix(c.CmdParent, "Cmd"), c.CmdName)
	}
	cmdFile, err := os.Create(fmt.Sprintf("%s/cmd/%s.go", c.AbsolutePath, filename))
	if err != nil {
		return err
	}
	defer cmdFile.Close()

	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"typeName": func(cmd string, parent string) string {
			name := strings.Title(cmd)
			if parent != "rootCmd" {
				name = strings.Title(strings.TrimSuffix(parent, "Cmd")) + name
			}
			name += "Cmd"
			return name
		},
		"title": strings.Title,
	}

	tmpl := tpl.AddCommandTemplateGlobal()
	if c.LocalVars {
		tmpl = tpl.AddCommandTemplateLocal()
	}
	commandTemplate := template.Must(template.New("sub").Funcs(funcMap).Parse(string(tmpl)))
	err = commandTemplate.Execute(cmdFile, c)
	if err != nil {
		return err
	}
	return nil
}
