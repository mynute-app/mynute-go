package handler

import (
	"agenda-kaki-go/core/config/db/model"
	"strings"
)

type PolicySubject struct {
	ID    uint
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
		subMatch := s.Attrs[rule.SubjectAttr] == rule.SubjectValue
		resMatch := r.Attrs[rule.ResourceAttr] == rule.ResourceValue
		methodMatch := strings.EqualFold(rule.Method, method)
		pathMatch := matchPath(rule.Path, path)
		attrMatch := false

		switch rule.AttrCondition {
		case "equal":
			attrMatch = subMatch && resMatch
		case "contains":
			attrMatch = strings.Contains(s.Attrs[rule.SubjectAttr], r.Attrs[rule.ResourceAttr])
		}

		if attrMatch && methodMatch && pathMatch {
			return true
		}
	}
	return false
}

func matchPath(rulePath, realPath string) bool {
	ruleSegments := strings.Split(rulePath, "/")
	realSegments := strings.Split(realPath, "/")

	if len(ruleSegments) != len(realSegments) {
		return false
	}

	for i := range ruleSegments {
		if strings.HasPrefix(ruleSegments[i], ":") {
			continue
		}
		if ruleSegments[i] != realSegments[i] {
			return false
		}
	}
	return true
}