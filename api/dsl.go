package api

type DSL struct {
	// ID is a unique string identifier for the DSL
	ID string
	// Version is the version of the DSL
	Version string

	// Variables is a list of variables that are used by the DSL and are shared between steps
	Variables map[string]any

	// Forms is a list of JSONSchema forms that are used by different steps
	Forms map[string]map[string]any

	// Steps is a list of steps (actions) that are executed in the mentioned order
	Steps map[string]Step
}

// StepType is the type of step (whether it's being performed by a human or a system, etc.)
type StepType string

const (
	// StepTypeHumanTask is a step that requires human input
	StepTypeHumanTask  StepType = "humanTask"
	StepTypeSystemTask StepType = "systemTask"
)

type Step struct {
	Type StepType
}
