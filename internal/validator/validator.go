package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var validIdentifier = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

var printfPattern = regexp.MustCompile(`%[sdvfxbcogtU]|%[0-9]*\.?[0-9]*[a-zA-Z%]`)

func ToValidIdentifier(msgID, msgCtx string) (string, error) {
	if msgID == "" {
		return "", fmt.Errorf("empty msgid")
	}

	var name string
	if msgCtx != "" {
		name = msgCtx + "_" + msgID
	} else {
		name = msgID
	}

	name = stripPrintfSpecifiers(name)
	name = sanitize(name)
	name = toCamelCase(name)

	if !validIdentifier.MatchString(name) {
		return "", fmt.Errorf("cannot convert %q to valid identifier", msgID)
	}

	return name, nil
}

func stripPrintfSpecifiers(s string) string {
	return printfPattern.ReplaceAllString(s, "")
}

func sanitize(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			result.WriteRune(r)
		} else if unicode.IsSpace(r) || r == ',' || r == '.' || r == '-' || r == '/' || r == ':' {
			if result.Len() > 0 && result.String()[result.Len()-1] != '_' {
				result.WriteRune('_')
			}
		} else if i > 0 && result.Len() > 0 && result.String()[result.Len()-1] != '_' {
			result.WriteRune('_')
		}
	}

	str := result.String()
	str = strings.Trim(str, "_")
	str = strings.ReplaceAll(str, "__", "_")

	return str
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
	}
	return strings.Join(parts, "")
}