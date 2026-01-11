package cmd

import (
	"bufio"
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
}

type Command struct {
	CmdName   string
	CmdParent string
	*Project
}

const (
	mainFilePathTemplate    = "%s/main.go"
	rootFilePathTemplate    = "%s/cmd/root.go"
	licenseFilePathTemplate = "%s/LICENSE"
	cmdDirPathTemplate      = "%s/cmd"
)

func (p *Project) Create() error {
	existingFiles := p.checkExistingFiles()
	if len(existingFiles) > 0 {
		fmt.Println("The following files already exist and will be overwritten:")
		for _, file := range existingFiles {
			fmt.Printf("  - %s\n", file)
		}
		fmt.Println()
		if !readConfirmation() {
			return fmt.Errorf("operation cancelled")
		}
	}

	if err := p.createDirectory(); err != nil {
		return err
	}

	if err := p.createMainFile(); err != nil {
		return err
	}

	if err := p.createRootFile(); err != nil {
		return err
	}

	if err := p.createLicenseFile(); err != nil {
		return err
	}

	return nil
}

func (p *Project) checkExistingFiles() []string {
	var files []string

	mainPath := fmt.Sprintf(mainFilePathTemplate, p.AbsolutePath)
	rootPath := fmt.Sprintf(rootFilePathTemplate, p.AbsolutePath)
	licensePath := fmt.Sprintf(licenseFilePathTemplate, p.AbsolutePath)

	if _, err := os.Stat(mainPath); err == nil {
		files = append(files, mainPath)
	}
	if _, err := os.Stat(rootPath); err == nil {
		files = append(files, rootPath)
	}
	if _, err := os.Stat(licensePath); err == nil {
		files = append(files, licensePath)
	}

	return files
}


func readConfirmation() bool {
	fmt.Println("Do you want to continue? This will overwrite the existing files.")
	fmt.Print("Continue? (Y/N): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "yes" {
		fmt.Println("Operation cancelled.")
		return false
	}

	return true
}

func (p *Project) createDirectory() error {
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		return os.Mkdir(p.AbsolutePath, 0754)
	}
	return nil
}

func (p *Project) createMainFile() error {
	mainPath := fmt.Sprintf(mainFilePathTemplate, p.AbsolutePath)
	mainFile, err := os.Create(mainPath)
	if err != nil {
		return err
	}
	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.MainTemplate())))
	return mainTemplate.Execute(mainFile, p)
}

func (p *Project) createRootFile() error {
	cmdPath := fmt.Sprintf(cmdDirPathTemplate, p.AbsolutePath)
	rootPath := fmt.Sprintf(rootFilePathTemplate, p.AbsolutePath)

	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(cmdPath, 0751))
	}

	rootFile, err := os.Create(rootPath)
	if err != nil {
		return err
	}
	defer rootFile.Close()

	rootTemplate := template.Must(template.New("root").Parse(string(tpl.RootTemplate())))
	return rootTemplate.Execute(rootFile, p)
}

func (p *Project) createLicenseFile() error {
	licensePath := fmt.Sprintf(licenseFilePathTemplate, p.AbsolutePath)
	licenseFile, err := os.Create(licensePath)
	if err != nil {
		return err
	}
	defer licenseFile.Close()

	data := map[string]interface{}{
		"copyright": copyrightLine(),
	}
	licenseTemplate := template.Must(template.New("license").Parse(p.Legal.Text))
	return licenseTemplate.Execute(licenseFile, data)
}

func (c *Command) Create() error {
	cmdDirPath := fmt.Sprintf(cmdDirPathTemplate, c.AbsolutePath)
	cmdFilePath := fmt.Sprintf("%s/%s.go", cmdDirPath, c.CmdName)

	if _, err := os.Stat(cmdFilePath); err == nil {
		fmt.Printf("File %s already exists and will be overwritten.\n", cmdFilePath)
		if !readConfirmation() {
			return  fmt.Errorf("operation cancelled")
		}
	}

	cmdFile, err := os.Create(cmdFilePath)
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
