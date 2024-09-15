package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gelleson/codemerge/codemerge/pkg/walker"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

var diffCmd = &cli.Command{
	Name:    "diff",
	Aliases: []string{"d"},
	Usage:   "compare two commits and merge changed files",
	Action:  diff,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "commit1",
			Aliases:  []string{"c1"},
			Usage:    "first commit hash",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "commit2",
			Aliases:  []string{"c2"},
			Usage:    "second commit hash",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "output",
			Aliases:  []string{"o"},
			Usage:    "output file",
			Required: true,
			EnvVars:  []string{"OUTPUT_FILE"},
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "verbose output",
			Value:   false,
		},
	},
}

func diff(c *cli.Context) error {
	commit1 := c.String("commit1")
	commit2 := c.String("commit2")
	output := c.String("output")

	// Step 1: Get the list of changed files using git diff
	cmd := exec.Command("git", "diff", "--name-only", commit1, commit2)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get changed files: %v", err)
	}

	changedFiles := strings.Split(string(out), "\n")

	// Step 2: Call CodeMerge to merge the changed files
	return mergeFiles(changedFiles, output)
}

func mergeFiles(files []string, outputFile string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	currentDir = currentDir + "/"

	fs := afero.NewOsFs()
	writer, err := fs.Create(outputFile)
	if err != nil {
		return err
	}

	wk := walker.New(afero.NewBasePathFs(afero.NewOsFs(), currentDir), ".", writer, true)
	for _, file := range files {
		if file == "" {
			continue
		}
		err = wk.ProcessFile(file) // Assuming this processes the file and writes it to output
		if err != nil {
			return err
		}
	}

	fmt.Printf("Merged %d files into %s\n", len(files), outputFile)
	return nil
}
