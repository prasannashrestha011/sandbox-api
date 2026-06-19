package lang

import "main/internal/domain"

var lang = map[string][]string{
	"python":     {"python3", "-c"},
	"javascript": {"node", "-e"},
	"typescript": {"ts-node", "-e"},
	"deno":       {"deno", "eval"},
	"ruby":       {"ruby", "-e"},
	"php":        {"php", "-r"},
	"perl":       {"perl", "-e"},
	"lua":        {"lua", "-e"},
	"groovy":     {"groovy", "-e"},
	"r":          {"Rscript", "-e"},
	"julia":      {"julia", "-e"},
	"bash":       {"bash", "-c"},
	"sh":         {"sh", "-c"},
	"zsh":        {"zsh", "-c"},
	"fish":       {"fish", "-c"},
	"powershell": {"pwsh", "-c"},
	"racket":     {"racket", "-e"},
	"guile":      {"guile", "-c"},
	"awk":        {"awk"},
	"octave":     {"octave", "--eval"},
	"prolog":     {"swipl", "-g"},
	"lisp":       {"sbcl", "--non-interactive", "--eval"},
	"crystal":    {"crystal", "eval"},
	"scala":      {"scala", "-e"},
}

func MapLangToRuntime(language string) ([]string, error) {
	if _, exists := lang[language]; !exists {
		return []string{}, domain.InvalidRequestError("unsupported language", nil)
	}
	return lang[language], nil
}
func BuildCommand(language string, code string) ([]string, error) {
	execCmd, exists := lang[language]
	if !exists {
		return []string{}, domain.InvalidRequestError("unsupported language", nil)
	}
	return append(execCmd, code), nil
}
