package judge

import (
	"context"

	"github.com/doitintl/kube-no-trouble/pkg/rules"
	"github.com/open-policy-agent/opa/rego"
	"github.com/rs/zerolog/log"
)

type RegoJudge struct {
	preparedQuery rego.PreparedEvalQuery
}

type RegoOpts struct {
}

func NewRegoJudge(opts *RegoOpts, rules []rules.Rule) (*RegoJudge, error) {
	ctx := context.Background()

	r := rego.New(
		rego.Query("data[_].main"),
	)

	for _, info := range rules {
		rego.Module(info.Name, info.Rule)(r)
		log.Info().Str("name", info.Name).Msg("Loaded ruleset")
	}

	pq, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}

	judge := &RegoJudge{preparedQuery: pq}
	return judge, nil
}

func (j *RegoJudge) Eval(input []map[string]interface{}) ([]Result, error) {
	ctx := context.Background()

	log.Debug().Msgf("evaluating +%v", input)
	rs, err := j.preparedQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, err
	}

	results := []Result{}
	for _, r := range rs {
		for _, e := range r.Expressions {
			for _, i := range e.Value.([]interface{}) {
				m := i.(map[string]interface{})
				results = append(results, Result{
					Name:        m["Name"].(string),
					Namespace:   m["Namespace"].(string),
					Kind:        m["Kind"].(string),
					ApiVersion:  m["ApiVersion"].(string),
					ReplaceWith: m["ReplaceWith"].(string),
					RuleSet:     m["RuleSet"].(string),
					Since:       m["Since"].(string),
				})
			}
		}
	}

	return results, nil
}
