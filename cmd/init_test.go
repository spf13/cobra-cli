package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

type (
	testArgs struct {
		name      string
		args      []string
		force     bool
		pkgName   string
		expectErr bool
		testFunc  func(*testArgs) error
	}
)

func getProject() *Project {
	wd, _ := os.Getwd()
	return &Project{
		AbsolutePath: fmt.Sprintf("%s/testproject", wd),
		Legal:        getLicense(),
		Copyright:    copyrightLine(),
		AppName:      "cmd",
		PkgName:      "github.com/spf13/cobra-cli/cmd/cmd",
		Viper:        true,
	}
}

func TestGoldenInitCmd(t *testing.T) {

	dir, err := ioutil.TempDir("", "cobra-init")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []testArgs{
		{
			name:      "successfully creates a project based on module",
			args:      []string{"testproject"},
			force:     false,
			pkgName:   "github.com/spf13/testproject",
			expectErr: false,
			testFunc: func(tt *testArgs) error {
				projectPath, err := initializeProject(tt.force, tt.args)
				defer func() {
					if projectPath != "" {
						os.RemoveAll(projectPath)
					}
				}()
				if err != nil {
					return err
				}

				expectedFiles := []string{"LICENSE", "main.go", "cmd/root.go"}
				for _, f := range expectedFiles {
					generatedFile := fmt.Sprintf("%s/%s", projectPath, f)
					goldenFile := fmt.Sprintf("testdata/%s.golden", filepath.Base(f))
					err := compareFiles(generatedFile, goldenFile)
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			name:      "does not overwrite files without force",
			args:      []string{"testproject"},
			force:     false,
			pkgName:   "github.com/spf13/testproject",
			expectErr: true,
			testFunc: func(tt *testArgs) error {
				projectPath, err := initializeProject(tt.force, tt.args)
				defer func() {
					if projectPath != "" {
						os.RemoveAll(projectPath)
					}
				}()
				if err != nil {
					return err
				}

				_, err = initializeProject(tt.force, tt.args)
				return err
			},
		},
		{
			name:      "does overwrite files with force",
			args:      []string{"testproject"},
			force:     true,
			pkgName:   "github.com/spf13/testproject",
			expectErr: false,
			testFunc: func(tt *testArgs) error {
				projectPath, err := initializeProject(tt.force, tt.args)
				defer func() {
					if projectPath != "" {
						os.RemoveAll(projectPath)
					}
				}()
				if err != nil {
					return err
				}

				_, err = initializeProject(tt.force, tt.args)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			viper.Set("useViper", true)
			viper.Set("license", "apache")

			err := tt.testFunc(&tt)

			if !tt.expectErr && err != nil {
				t.Fatalf("did not expect an error, got %s", err)
			}
			if tt.expectErr {
				if err == nil {
					t.Fatal("expected an error but got none")
				} else {
					// got an expected error nothing more to do
					return
				}
			}
		})
	}
}
