package models

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/joakim-ribier/gttp/core"
)

// Test 'GetEnvsName' method
func TestGetEnvsNameIfEmpty(t *testing.T) {
	var ctx Context
	actual := ctx.GetEnvsName()

	if len(actual) != 1 || actual[0] != "default" {
		t.Error("Expected 'default', got ", actual)
	}
}

func TestGetEnvsNameSortBy(t *testing.T) {
	var ctx Context
	ctx.Add("prod", "{prod}}", "prod-value")
	ctx.Add("dev", "{dev}", "dev-value")

	actual := ctx.GetEnvsName()

	if len(actual) != 3 {
		t.Error("Expected len(3), got ", len(actual))
	}
	var expected core.StringSlice = []string{"default", "dev", "prod"}
	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Sprintf("Expected %v, got ", expected), actual)
	}
}

// Test 'Add' method
func TestAddFromEmptyContext(t *testing.T) {
	var ctx Context

	ctx.Add("dev", "{dev}", "dev-value")

	if ctx.Env["dev"][0] != NewContextVariable("{dev}", "dev-value") {
		t.Error("Expected '{dev}' env, got ", ctx.Env)
	}
}

func TestAddLowerCaseValues(t *testing.T) {
	var ctx Context

	// Transform only 'env' & 'variable' to lower case
	ctx.Add("DEV", "{dEv}", "dev-VALUE")

	if ctx.Env["dev"][0] != NewContextVariable("{dev}", "dev-VALUE") {
		t.Error("Expected '{dev}' env, got ", ctx.Env)
	}
}

func TestAdd(t *testing.T) {
	var ctx Context
	ctx.Add("dev", "{dev-1}", "dev-1-value")
	ctx.Add("default", "{default}", "default-value")
	ctx.Add("dev", "{dev-2}", "dev-2-value")
	ctx.Add("prod", "{prod}", "prod-value")

	if len(ctx.Env) != 3 {
		t.Error("Expected len(3), got ", len(ctx.Env))
	}
	t.Run("dev", testLength(ctx.Env["dev"], 2))
	t.Run("default", testLength(ctx.Env["default"], 1))
	t.Run("prod", testLength(ctx.Env["prod"], 1))
}

// Test 'Remove' method
func TestRemoveIfEnvDoesNotExist(t *testing.T) {
	var ctx Context

	// Remove nothing...
	ctx.Remove("dev", "{dev-1}")

	if len(ctx.Env) != 0 {
		t.Error("Expected len(0), got ", len(ctx.Env))
	}
}

func TestRemove(t *testing.T) {
	var ctx Context
	ctx.Add("dev", "{dev-1}", "dev-1-value")
	ctx.Add("dev", "{dev-2}", "dev-2-value")

	ctx.Remove("dev", "{dev-2}")

	t.Run("dev", testLength(ctx.Env["dev"], 1))
}

func TestRemoveEnvIfEmptyVariables(t *testing.T) {
	var ctx Context
	ctx.Add("dev", "{dev-1}", "dev-1-value")
	ctx.Add("dev", "{dev-2}", "dev-2-value")

	// Remove all data
	ctx.Remove("dev", "{dev-1}")
	ctx.Remove("dev", "{dev-2}")

	// Check if 'dev' still exists
	if _, is := ctx.Env["dev"]; is {
		t.Error("Expected false, got ", ctx.GetEnvsName())
	}
}

// Test 'GetAllVariables'
func TestGetAllVariables(t *testing.T) {
	expected := []string{"{dev-1}", "{dev-2}"}

	var ctx Context
	ctx.Add("dev", expected[0], "dev-1-value")
	ctx.Add("dev", expected[1], "dev-2-value")

	var actual []string = ctx.GetAllVariables("dev")

	// Sort actual result to match with expected
	sort.Strings(actual)

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Sprintf("Expected %v, got ", expected), actual)
	}
}

func TestGetAllVariablesEmpyIfEnvDoesNotExist(t *testing.T) {
	expected := []string{}
	var ctx Context

	var actual []string = ctx.GetAllVariables("dev")

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Sprintf("Expected %v, got ", expected), actual)
	}
}

// Test 'GetAllKeyValue'
func TestGetAllKeyValue(t *testing.T) {
	var ctx Context
	ctx.Add("dev", "{dev-1}", "dev-1-value")
	ctx.Add("dev", "{dev-2}", "dev-2-value")
	ctx.Add("prod", "{prod-1}", "prod-1-value")

	actual := ctx.GetAllKeyValue("dev")

	if len(actual) != 2 || actual["{dev-1}"] != "dev-1-value" || actual["{dev-2}"] != "dev-2-value" {
		t.Error(fmt.Sprintf("Expected %v, got ", ctx.Env["dev"]), actual)
	}
}

func TestGetAllKeyValueEmpyIfEnvDoesNotExist(t *testing.T) {
	var ctx Context

	actual := ctx.GetAllKeyValue("dev")

	if len(actual) != 0 {
		t.Error(fmt.Sprintf("Expected %v, got ", ctx.Env["dev"]), actual)
	}
}

// Test 'FindVariableByEnv'
func TestFindVariableByEnv(t *testing.T) {
	expected := NewContextVariable("{dev-1}", "dev-1-value")

	var ctx Context
	ctx.Add("dev", "{dev-1}", "dev-1-value")
	ctx.Add("dev", "{dev-2}", "dev-2-value")

	actual := ctx.FindVariableByEnv("dev", "{dev-1}")

	if actual != expected {
		t.Error(fmt.Sprintf("Expected %v, got ", expected), actual)
	}
}

func TestFindVariableByEnvEmpyIfEnvDoesNotExist(t *testing.T) {
	var expected ContextVariable
	var ctx Context

	var actual = ctx.FindVariableByEnv("dev", "{dev-1}")

	if actual != expected {
		t.Error(fmt.Sprintf("Expected %v, got ", expected), actual)
	}
}

func testLength(values []ContextVariable, expected int) func(*testing.T) {
	return func(t *testing.T) {
		actual := len(values)
		if actual != expected {
			t.Error(fmt.Sprintf("Expected length of %v to be %d but instead got %d!", values, expected, actual))
		}
	}
}
