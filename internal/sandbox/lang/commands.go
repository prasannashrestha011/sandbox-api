package lang

import "fmt"

var LangCommands = map[string][]string{
	"python":     {"python", "-c"},
	"javascript": {"node", "-e"},
}

func BuildCommand(lang string, code string) ([]string, error) {
	cmdTemplate, ok := LangCommands[lang]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
	return append(cmdTemplate, code), nil
}
