package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
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

	tests := []struct {
		name            string
		pkgName         string
		args            []string
		expectedFiles   []string
		unexpectedFiles []string
		expectErr       bool
		skipLicense     bool
	}{
		{
			name:          "successfully creates a project based on module",
			args:          []string{"testproject"},
			pkgName:       "github.com/spf13/testproject",
			expectErr:     false,
			expectedFiles: []string{"LICENSE", "main.go", "cmd/root.go"},
			skipLicense:   false,
		},
		{
			name:            "successfully creates a project based on module",
			args:            []string{"testproject"},
			pkgName:         "github.com/spf13/testproject",
			expectErr:       false,
			expectedFiles:   []string{"main.go", "cmd/root.go"},
			unexpectedFiles: []string{"LICENSE"},
			skipLicense:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			viper.Set("useViper", true)
			viper.Set("license", "apache")

			if tt.skipLicense {
				skipLicense = true
			}

			projectPath, err := initializeProject(tt.args)
			defer func() {
				if projectPath != "" {
					_ = os.RemoveAll(projectPath)
				}
			}()

			if !tt.expectErr && err != nil {
				t.Fatalf("did not expect an error, got %s", err)
			}
			if tt.expectErr && err == nil {
				t.Fatal("expected an error but got none")
			}

			if err != nil {
				t.Fatal(err)
			}

			for _, f := range tt.expectedFiles {
				generatedFile := fmt.Sprintf("%s/%s", projectPath, f)
				goldenFile := fmt.Sprintf("testdata/%s.golden", filepath.Base(f))
				err := compareFiles(generatedFile, goldenFile)
				if err != nil {
					t.Fatal(err)
				}
			}

			for _, f := range tt.unexpectedFiles {
				path := fmt.Sprintf("%s/%s", projectPath, f)
				if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
					continue
				}
				t.Fatalf("%s should not be generated", path)
			}
		})
	}
}
