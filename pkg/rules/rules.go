package rules

import (
	"embed"
	"path"
)

//go:embed rego
var local embed.FS

const ruleDir = "rego"

type Rule struct {
	Name string
	Rule string
}

func FetchRegoRules() ([]Rule, error) {
	fis, _ := local.ReadDir(ruleDir)

	rules := []Rule{}
	for _, info := range fis {
		data, _ := local.ReadFile(path.Join(ruleDir, info.Name()))
		rules = append(rules, Rule{
			Name: info.Name(),
			Rule: string(data),
		})
	}

	return rules, nil
}
