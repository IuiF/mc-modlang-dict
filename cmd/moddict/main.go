// Package main provides the CLI for minecraft-mod-dictionary.
package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error
	switch command {
	case "import":
		err = runImport(args)
	case "import-dir":
		err = runImportDir(args)
	case "translate":
		err = runTranslate(args)
	case "export":
		err = runExport(args)
	case "view":
		err = runView(args)
	case "build":
		err = runBuild(args)
	case "migrate":
		err = runMigrate(args)
	case "repair":
		err = runRepair(args)
	case "analyze":
		err = runAnalyze(args)
	case "fix-schema":
		err = runFixSchema(args)
	case "version", "-v", "--version":
		fmt.Printf("moddict version %s\n", version)
		return
	case "help", "-h", "--help":
		printUsage()
		return
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`moddict - Minecraft Mod Dictionary CLI

Usage:
  moddict <command> [options]

Commands:
  import      Import translations from a mod JAR file
  import-dir  Import translations from a directory (cloned repo)
  translate   Add/update translations, show status
  export      Export translations to various formats
  view        View translations in the database
  build       Build translation database from YAML files
  migrate     Migrate existing data to new source-based schema
  repair      Repair database inconsistencies
  analyze     Analyze translation consistency and discover patterns
  version     Show version information
  help        Show this help message

Use "moddict <command> --help" for more information about a command.
`)
}
