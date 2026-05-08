package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gettextcodegen/internal/diff"
	"github.com/gettextcodegen/internal/generator"
	"github.com/gettextcodegen/internal/parser"
)

func main() {
	dir := flag.String("dir", "", "Directory containing language subdirectories with LC_MESSAGES/default.po")
	lang := flag.String("lang", "", "Default language code for code generation (required)")
	pkg := flag.String("package", "", "Package name override (default: directory name)")

	flag.Parse()

	if *dir == "" || *lang == "" {
		flag.Usage()
		os.Exit(1)
	}

	dirAbs, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid directory: %v\n", err)
		os.Exit(1)
	}

	packageName := *pkg
	if packageName == "" {
		packageName = filepath.Base(dirAbs)
		packageName = strings.ReplaceAll(packageName, "-", "")
		packageName = strings.ReplaceAll(packageName, "_", "")
		if packageName == "" {
			packageName = "translations"
		}
	}

	defaultPOPath := parser.GetDefaultPOPath(dirAbs, *lang)
	defaultEntries, err := parser.ParsePO(defaultPOPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to parse default PO file %s: %v\n", defaultPOPath, err)
		os.Exit(1)
	}

	defaultKeys := make(map[string]bool)
	for _, entry := range defaultEntries {
		defaultKeys[entry.MsgID] = true
	}

	languages, err := parser.ListLanguages(dirAbs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to list languages: %v\n", err)
		os.Exit(1)
	}

	g := &generator.Generator{
		Package:   packageName,
		Dir:       dirAbs,
		Entries:   defaultEntries,
		Languages: languages,
	}

	code, err := g.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate code: %v\n", err)
		os.Exit(1)
	}

	for _, langName := range languages {
		if langName == *lang {
			continue
		}

		langPOPath := parser.GetLanguagePOPath(dirAbs, langName)
		langEntries, err := parser.ParsePO(langPOPath)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", langPOPath, err)
			continue
		}

		langKeys := make(map[string]bool)
		for _, entry := range langEntries {
			langKeys[entry.MsgID] = true
		}

		d := diff.Compare(defaultKeys, langKeys)
		if d.HasDiff() {
			fmt.Println(d.Format(langName, *lang))
		}
	}

	outputPath := filepath.Join(dirAbs, "gettext.go")
	if err := os.WriteFile(outputPath, []byte(code), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to write output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated: %s\n", outputPath)
}