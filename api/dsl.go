package api

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

type HumanTaskStep struct {
	// ID of the form to be used for this step. Forms are used only for those steps that are of type StepTypeHumanTask (humanTask)
	Form string `json:"form"`
}

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

type Hook struct {
	// Script is a script to be executed as a hook
	// Access variables using the syntax: {{variableName}}
	// Or the filled forms using the syntax: {{formName.propertyName}}
	Script string `json:"script"`
}
