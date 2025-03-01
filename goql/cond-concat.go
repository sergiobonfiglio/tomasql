package goql

import (
	"fmt"
	"strings"
)

type ConditionConnector string

const (
	AndCondConnector = ConditionConnector("AND")
	OrCondConnector  = ConditionConnector("OR")
)

type ConcatCondition struct {
	conditions []Condition
	connector  ConditionConnector
}

var _ Condition = &ConcatCondition{} // Ensure ConcatCondition implements Condition

func newConcatCondition(connector ConditionConnector, conditions ...Condition) *ConcatCondition {
	return &ConcatCondition{conditions: conditions, connector: connector}
}

func (c *ConcatCondition) SQL(p ParamsMap) string {
	var innerSQLs []string
	for _, cond := range c.conditions {
		innerSQLs = append(innerSQLs, cond.SQL(p))
	}
	connector := fmt.Sprintf(" %s ", c.connector)
	return fmt.Sprintf("%s", strings.Join(innerSQLs, connector))
}

func (c *ConcatCondition) And(condition Condition) Condition {
	return newConcatCondition(AndCondConnector, c, condition)
}

func (c *ConcatCondition) Or(condition Condition) Condition {
	return newConcatCondition(OrCondConnector, c, condition)
}

type GroupedCondition struct {
	ConcatCondition
}

var _ Condition = &GroupedCondition{} // Ensure GroupedCondition implements Condition

func Grouped(conditions ...Condition) *GroupedCondition {
	return &GroupedCondition{ConcatCondition{conditions: conditions}}
}

func (g *GroupedCondition) SQL(p ParamsMap) string {
	return fmt.Sprintf("(%s)", g.ConcatCondition.SQL(p))
}

func (g *GroupedCondition) And(condition Condition) Condition {
	return newConcatCondition(AndCondConnector, g, condition)
}

func (g *GroupedCondition) Or(condition Condition) Condition {
	return newConcatCondition(OrCondConnector, g, condition)
}
