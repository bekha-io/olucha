package api

import (
	"fmt"
	"slices"

	"github.com/xeipuuv/gojsonschema"
)

type DSL struct {
	// ID is a unique string identifier for the DSL
	ID string `json:"id"`

	// Secrets is a list of secrets that are used by the workflow
	// Defined by the system's users and uploaded separately to each workflow instance.
	// If not explicitly provided, will try to be fetched from the environment variables.
	Secrets map[string]string `json:"secrets"`

	// Version is the version of the DSL
	Version string `json:"version"`

	// Variables is a list of variables that are used by the DSL and are shared between steps.
	// Variables are used to store data that is shared between steps.
	Variables map[string]any `json:"variables"`

	// Forms is a list of JSONSchema forms that are used by different steps.
	// Forms are used to validate the data that is submitted by the user.
	Forms map[string]map[string]any `json:"forms"`

	// Steps is a list of steps (actions) that are executed in the mentioned order
	Steps map[string]Step `json:"steps"`
}

// StepType is the type of step (whether it's being performed by a human or a system, etc.)
type StepType string

const (
	// StepTypeHumanTask is a step that requires human input
	StepTypeHumanTask  StepType = "humanTask"
	StepTypeSystemTask StepType = "systemTask"
	StepTypeCondition  StepType = "condition"
)

// Step is a single step in the workflow
type Step struct {
	// Label is a human-readable name for the step. If empty, the step will be named after its type
	Label string `json:"label"`

	// Description is a human-readable description for the step. If empty, the step will not have a description
	Description string `json:"description"`

	Type StepType `json:"type"`

	// Embedding different step types
	ConditionStep
	HumanTaskStep

	// ID of the next step to be executed
	Next string `json:"next"`

	// Before is a hook to be executed before the step is executed
	Before Hook `json:"before"`

	// After is a hook to be executed after the step is executed
	After Hook `json:"after"`
}

func (s Step) ValidateForm(dsl DSL, form map[string]any) error {
	// If no form is provided, the step is allowed to be executed by anyone
	if s.Form == "" {
		return nil
	}

	// If the form is not defined in the DSL, the step is not allowed to be executed
	if _, ok := dsl.Forms[s.Form]; !ok {
		return fmt.Errorf("form %s is not defined in the DSL", s.Form)
	}

	schema, err := gojsonschema.NewSchema(gojsonschema.NewGoLoader(form))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	res, err := schema.Validate(gojsonschema.NewGoLoader(form))
	if err != nil {
		return fmt.Errorf("validation attempt failed: %w", err)
	}

	if !res.Valid() {
		return fmt.Errorf("validation error: %s", res.Errors())
	}

	return nil
}

func (s Step) IsAllowedFor(roles []string) bool {
	// If any of the roles are in the Any list, the step is allowed
	if len(s.RBAC.Any) > 0 {
		for _, role := range roles {
			if slices.Contains(s.RBAC.Any, role) {
				return true
			}
		}
	}

	// If all of the roles are in the All list, the step is allowed
	if len(s.RBAC.All) > 0 {
		for _, role := range roles {
			if !slices.Contains(s.RBAC.All, role) {
				return false
			}
		}
	}

	// If no roles are provided, the step is allowed to be executed by anyone
	return true
}

// HumanTaskStep is a step that requires human input or human interaction
type HumanTaskStep struct {
	// ID of the form to be used for this step. Forms are used only for those steps that are of type StepTypeHumanTask (humanTask)
	Form string `json:"form"`

	// RBAC is a list of roles that are allowed to perform the step
	// If not provided, the step will be available to anyone
	RBAC RBAC `json:"rbac"`
}

// ConditionStep is a step that evaluates an expression and executes a different step based on the result
type ConditionStep struct {
	// Expr is an expression to be evaluated to determine the next step to be executed
	// Access variables using the syntax: {{variableName}}
	// Or the filled forms using the syntax: {{formName.propertyName}}
	Expr string `json:"expr"`

	// ID of the step to be executed if the condition is true
	If string `json:"if"`

	// ID of the step to be executed if the condition is false
	Else string `json:"else"`
}

// Hook is a script to be executed as a hook
type Hook struct {
	// Script is a script to be executed as a hook
	// Access variables using the syntax: {{variableName}}
	// Or the filled forms using the syntax: {{formName.propertyName}}
	Script string `json:"script"`
}

// RBAC is a list of roles (permissions) that are allowed to perform the step
// Passed to the step while executing
type RBAC struct {
	// Any is a list of roles that are allowed to perform the step
	Any []string `json:"any"`

	// All is a list of roles that are required to perform the step
	All []string `json:"all"`
}
