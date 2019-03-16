package models

import (
	"sort"
	"strings"

	"github.com/joakim-ribier/gttp/core"
)

const defaultValue = "default"

// Context reprensents a context structure
type Context struct {
	Env map[string][]ContextVariable
}

// ContextVariable reprensents a context variable structure
type ContextVariable struct {
	Variable string
	Value    string
}

// NewContextVariable creates new ContextVariable struct
func NewContextVariable(variable string, value string) ContextVariable {
	return ContextVariable{
		Variable: variable,
		Value:    value,
	}
}

// GetEnvsName gets all environments name
func (c Context) GetEnvsName() core.StringSlice {
	var tab core.StringSlice = []string{}
	for key := range c.Env {
		tab = append(tab, strings.ToLower(key))
	}
	sort.Strings(tab)
	if index := tab.GetIndex(defaultValue); index == -1 {
		return append([]string{defaultValue}, tab...)
	}
	return tab
}

// Add adds new variable to an environment
func (c *Context) Add(env string, variable string, value string) {
	c.add(strings.ToLower(env), strings.ToLower(variable), value)
}

func (c *Context) add(env string, variable string, value string) {
	add := func(variable string, value string, slice []ContextVariable) []ContextVariable {
		if slice == nil {
			return []ContextVariable{NewContextVariable(variable, value)}
		}
		newSlice := []ContextVariable{NewContextVariable(variable, value)}
		for _, value := range slice {
			if value.Variable != variable {
				newSlice = append(newSlice, value)
			}
		}
		return newSlice
	}
	if c.Env == nil {
		c.Env = make(map[string][]ContextVariable)
	}
	c.Env[env] = add(variable, value, c.Env[env])
}

// Remove removes a variable to an environment
func (c *Context) Remove(env string, variable string) {
	remove := func(variable string, slice []ContextVariable) []ContextVariable {
		newSlice := []ContextVariable{}
		for _, value := range slice {
			if value.Variable != variable {
				newSlice = append(newSlice, value)
			}
		}
		return newSlice
	}
	if _, is := c.Env[env]; is {
		c.Env[env] = remove(variable, c.Env[env])
		if len(c.Env[env]) == 0 && env != defaultValue {
			delete(c.Env, env)
		}
	}
}

// GetAllVariables returns all variables for an specific environment
func (c Context) GetAllVariables(env string) core.StringSlice {
	newSlice := []string{}
	for _, value := range c.Env[env] {
		newSlice = append(newSlice, value.Variable)
	}
	return newSlice
}

// GetAllKeyValue gets all ContextVariable for an specific environment
func (c Context) GetAllKeyValue(env string) map[string]string {
	newMap := make(map[string]string)
	for _, value := range c.Env[env] {
		newMap[value.Variable] = value.Value
	}
	return newMap
}

// FindVariableByEnv finds variable for an specific environment
func (c Context) FindVariableByEnv(env string, variable string) ContextVariable {
	var value ContextVariable
	for _, value := range c.Env[env] {
		if value.Variable == variable {
			return value
		}
	}
	return value
}
