package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/merlindrones/gomavlib/pkg/conversion"
)

var cli struct {
	Link  bool   `help:"Link included definitions instead of including them into the main definition"`
	XML   string `arg:"" help:"Path or URL pointing to a XML MAVLink dialect"`
	Out   string `help:"Output base directory for generated files (default: pkg/dialects).The dialect will be placed in a subfolder named after the XML (e.g.,swarmos.xml → pkg/dialects/swarmos)" default:"pkg/dialects"`
	Force bool   `help:"Remove any existing output directory before generating"`
}

func isURL(s string) bool {
	return strings.Contains(s, "://")
}

func run(args []string) error {
	parser, err := kong.New(&cli,
		kong.Description("Convert MAVLink dialects from XML format to Go format."),
		kong.UsageOnError())
	if err != nil {
		return err
	}
	if _, err = parser.Parse(args); err != nil {
		return err
	}

	// Capture current working dir (where the user ran the command).
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working dir: %w", err)
	}

	// Resolve XML input:
	// - If it's a URL, leave it as-is.
	// - If it's a relative path, make it absolute against the *original* cwd.
	xmlInput := cli.XML
	if !isURL(xmlInput) && !filepath.IsAbs(xmlInput) {
		xmlInput = filepath.Join(cwd, xmlInput)
	}

	// Dialect directory name = xml filename without extension.
	baseName := filepath.Base(cli.XML)
	dialectDir := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	// Resolve output base relative to cwd.
	outBase := cli.Out
	if !filepath.IsAbs(outBase) {
		outBase = filepath.Join(cwd, outBase)
	}
	targetDir := filepath.Join(outBase, dialectDir)

	// Prepare output dir.
	if cli.Force {
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("failed to remove existing %s: %w", targetDir, err)
		}
	}
	if err := os.MkdirAll(outBase, 0o755); err != nil {
		return fmt.Errorf("failed to create output base %s: %w", outBase, err)
	}

	// Temporarily chdir to outBase so the converter writes into ./<dialectDir>.
	prevCwd, _ := os.Getwd()
	if err := os.Chdir(outBase); err != nil {
		return fmt.Errorf("failed to chdir to %s: %w", outBase, err)
	}
	defer os.Chdir(prevCwd)

	// Use the absolute XML input so chdir doesn't break relative resolution.
	if err := conversion.Convert(xmlInput, cli.Link); err != nil {
		return err
	}

	fmt.Printf("✅ Generated dialect %s in %s\n", dialectDir, targetDir)
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %s\n", err)
		os.Exit(1)
	}
}
