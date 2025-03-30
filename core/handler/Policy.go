package handler

import (
	"strings"
	"agenda-kaki-go/core/lib"
	"agenda-kaki-go/core/config/db/model"
)

type PolicySubject struct {
	ID    string
	Attrs map[string]string
}

type PolicyResource struct {
	Attrs map[string]string
}

type PolicyEnvironment struct {
	Attrs map[string]string
}

type PolicyEngine struct {
	Rules []model.PolicyRule
}

func Policy(rules []model.PolicyRule) *PolicyEngine {
	return &PolicyEngine{Rules: rules}
}

func (pe *PolicyEngine) CanAccess(s PolicySubject, method, path string, r PolicyResource, e PolicyEnvironment) bool {
	for _, rule := range pe.Rules {
		if !strings.EqualFold(rule.Method, method) || !lib.MatchPath(rule.Resource.Path, path) {
			continue
		}

		if s.Attrs[rule.SubjectAttr] != rule.SubjectValue {
			continue
		}

		match := true
		for _, cond := range rule.Conditions {
			val := r.Attrs[cond.Attr]
			switch cond.Op {
			case "equal":
				if val != cond.Value {
					match = false
				}
			case "contains":
				if !strings.Contains(val, cond.Value) {
					match = false
				}
			}
		}

		if match {
			return true
		}
	}
	return false
}