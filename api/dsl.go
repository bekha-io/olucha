package api

type DSL struct {
	// ID is a unique string identifier for the DSL
	ID string `json:"id"`

	// Version is the version of the DSL
	Version string `json:"version"`

	// Variables is a list of variables that are used by the DSL and are shared between steps
	Variables map[string]any `json:"variables"`

	// Forms is a list of JSONSchema forms that are used by different steps
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
	Type StepType `json:"type"`

	// Embedding different step types
	ConditionStep
	HumanTaskStep

	// ID of the next step to be executed
	Next string `json:"next"`
}

type HumanTaskStep struct {
	// ID of the form to be used for this step. Forms are used only for those steps that are of type StepTypeHumanTask (humanTask)
	Form string `json:"form"`
}

type ConditionStep struct {
	// Expression to be evaluated to determine the next step to be executed
	// Access variables using the syntax: {{variableName}}
	// Or the filled forms using the syntax: {{formName.propertyName}}
	Condition string `json:"condition"`

	// ID of the step to be executed if the condition is true
	If string `json:"if"`

	// ID of the step to be executed if the condition is false
	Else string `json:"else"`
}
