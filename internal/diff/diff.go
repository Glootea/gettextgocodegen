package diff

import (
	"fmt"
	"sort"
	"strings"
)

type Diff struct {
	Missing []string
	Extra   []string
}

func Compare(defaultKeys, langKeys map[string]bool) *Diff {
	var missing, extra []string

	for key := range defaultKeys {
		if !langKeys[key] {
			missing = append(missing, key)
		}
	}

	for key := range langKeys {
		if !defaultKeys[key] {
			extra = append(extra, key)
		}
	}

	sort.Strings(missing)
	sort.Strings(extra)

	return &Diff{
		Missing: missing,
		Extra:   extra,
	}
}

func (d *Diff) HasDiff() bool {
	return len(d.Missing) > 0 || len(d.Extra) > 0
}

func (d *Diff) Format(langName, defaultName string) string {
	if !d.HasDiff() {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("%s diff with default %s:", langName, defaultName))

	for _, key := range d.Missing {
		lines = append(lines, fmt.Sprintf("  - %s", key))
	}

	for _, key := range d.Extra {
		lines = append(lines, fmt.Sprintf("  + %s", key))
	}

	return strings.Join(lines, "\n")
}