package tomasql

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

func (c *ConcatCondition) Columns() []Column {
	var cols []Column
	for _, cond := range c.conditions {
		cols = append(cols, cond.Columns()...)
	}
	return cols
}

var _ Condition = &ConcatCondition{} // Ensure ConcatCondition implements Condition

func NewConcatCondition(connector ConditionConnector, conditions ...Condition) *ConcatCondition {
	return &ConcatCondition{conditions: conditions, connector: connector}
}

func (c *ConcatCondition) SQL(p ParamsMap) string {
	var innerSQLs []string
	for _, cond := range c.conditions {
		innerSQLs = append(innerSQLs, cond.SQL(p))
	}
	connector := fmt.Sprintf(" %s ", c.connector)
	return strings.Join(innerSQLs, connector)
}

func (c *ConcatCondition) And(condition Condition) Condition {
	return NewConcatCondition(AndCondConnector, c, condition)
}

func (c *ConcatCondition) Or(condition Condition) Condition {
	return NewConcatCondition(OrCondConnector, c, condition)
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
	return NewConcatCondition(AndCondConnector, g, condition)
}

func (g *GroupedCondition) Or(condition Condition) Condition {
	return NewConcatCondition(OrCondConnector, g, condition)
}
