package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/merlindrones/gomavlib/pkg/conversion"
)

var cli struct {
	Link  bool   `help:"Link included definitions instead of including them into the main definition"`
	XML   string `arg:"" help:"Path or URL pointing to a XML MAVLink dialect"`
	Chdir string `help:"Change into this directory before generating output"`
	Force bool   `help:"Remove any existing output directory before generating"`
}

func run(args []string) error {
	parser, err := kong.New(&cli,
		kong.Description("Convert MAVLink dialects from XML format to Go format."),
		kong.UsageOnError())
	if err != nil {
		return err
	}

	_, err = parser.Parse(args)
	if err != nil {
		return err
	}

	// Change working directory if requested
	if cli.Chdir != "" {
		if err := os.Chdir(cli.Chdir); err != nil {
			return fmt.Errorf("failed to chdir to %s: %w", cli.Chdir, err)
		}
	}

	// Determine output folder name from XML file name (without extension)
	baseName := filepath.Base(cli.XML)
	outDir := cli.XML
	ext := filepath.Ext(baseName)
	if ext != "" {
		outDir = baseName[:len(baseName)-len(ext)]
	} else {
		outDir = baseName
	}

	// If force is enabled, remove existing directory
	if cli.Force {
		if err := os.RemoveAll(outDir); err != nil {
			return fmt.Errorf("failed to remove existing %s: %w", outDir, err)
		}
	}

	if err := conversion.Convert(cli.XML, cli.Link); err != nil {
		return err
	}

	fmt.Printf("✅ Generated dialect %s in %s\n", outDir, filepath.Join(getCwd(), outDir))
	return nil
}

func getCwd() string {
	cwd, _ := os.Getwd()
	return cwd
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %s\n", err)
		os.Exit(1)
	}
}
