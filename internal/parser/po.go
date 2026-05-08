package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/leonelquinteros/gotext"
)

type POEntry struct {
	MsgID     string
	MsgStr    string
	MsgCtx    string
	IsPlural  bool
	PluralID  string
	Variables []Variable
}

type Variable struct {
	Name string
	Type string
}

func (p *POEntry) HasVariables() bool {
	return len(p.Variables) > 0
}

func (p *POEntry) HasPluralVariables() bool {
	return p.IsPlural && extractVariables(p.PluralID) != nil && len(extractVariables(p.PluralID)) > 0
}

var printfSpecifiers = map[string]string{
	"%s":  "string",
	"%d":  "int",
	"%v":  "int",
	"%f":  "float64",
	"%x":  "int",
	"%b":  "int",
	"%c":  "rune",
	"%e":  "float64",
	"%g":  "float64",
	"%o":  "int",
	"%q":  "string",
	"%t":  "bool",
	"%U":  "int",
}

func ListLanguages(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var languages []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		languages = append(languages, entry.Name())
	}
	return languages, nil
}

func GetDefaultPOPath(dir, lang string) string {
	return filepath.Join(dir, lang, "LC_MESSAGES", "default.po")
}

func GetLanguagePOPath(dir, lang string) string {
	return filepath.Join(dir, lang, "LC_MESSAGES", "default.po")
}

func ParsePO(path string) ([]*POEntry, error) {
	po := gotext.NewPo()
	po.ParseFile(path)

	domain := po.GetDomain()
	translations := domain.GetTranslations()

	var entries []*POEntry
	for _, trans := range translations {
		entry := &POEntry{
			MsgID:     trans.ID,
			MsgStr:    trans.Get(),
			IsPlural:  trans.PluralID != "",
			PluralID:  trans.PluralID,
			Variables: extractVariables(trans.Get()),
		}
		entries = append(entries, entry)
	}

	ctxTranslations := domain.GetCtxTranslations()
	for ctx, ctxMap := range ctxTranslations {
		for _, trans := range ctxMap {
			entry := &POEntry{
				MsgID:     trans.ID,
				MsgCtx:    ctx,
				MsgStr:    trans.Get(),
				IsPlural:  trans.PluralID != "",
				PluralID:  trans.PluralID,
				Variables: extractVariables(trans.Get()),
			}
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func extractVariables(str string) []Variable {
	var vars []Variable
	re := regexp.MustCompile(`%[sdvfxbcogtU]`)
	matches := re.FindAllStringIndex(str, -1)

	paramCounter := 1
	for _, match := range matches {
		spec := str[match[0]:match[1]]
		if spec == "%" {
			continue
		}
		goType, ok := printfSpecifiers[spec]
		if !ok {
			goType = "interface{}"
		}
		vars = append(vars, Variable{
			Name: fmt.Sprintf("param%d", paramCounter),
			Type: goType,
		})
		paramCounter++
	}
	return vars
}