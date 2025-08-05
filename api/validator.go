package api

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type DSLValidator interface {
	Validate(dsl DSL) error
}

type defaultDSLValidator struct {
}

func (v *defaultDSLValidator) Validate(dsl DSL) error {
	var err = ErrInvalidFormat

	if dsl.ID == "" {
		err = fmt.Errorf("%w: id is required", err)
	}

	if dsl.Version == "" {
		err = fmt.Errorf("%w: version is required", err)
	}

	// Each DSL must have at least one step
	if dsl.Steps == nil {
		err = fmt.Errorf("%w: 0 steps are defined. please define at least one step", err)
	}

	// Validate the forms
	for formID, form := range dsl.Forms {
		err = v.validateForm(formID, form)
		if err != nil {
			err = fmt.Errorf("%w: form %s is invalid: %w", err, formID, err)
		}
	}

	// Validate each step
	for stepID, step := range dsl.Steps {
		err = v.validateStep(stepID, step)
		if err != nil {
			err = fmt.Errorf("%w: step %s is invalid: %w", err, stepID, err)
		}
	}

	return err
}

func (v *defaultDSLValidator) validateForm(formID string, form map[string]any) error {
	_, err := gojsonschema.NewSchema(gojsonschema.NewGoLoader(form))
	if err != nil {
		return fmt.Errorf("failed to create schema for form %s: %w", formID, err)
	}
	return nil
}

func (v *defaultDSLValidator) validateStep(stepID string, step Step) error {

	if stepID == "" {
		return fmt.Errorf("all steps must have an ID")
	}

	if step.Type == "" {
		return fmt.Errorf("step '%s' has no type", stepID)
	}

	// Validate the step based on its type
	switch step.Type {
	case StepTypeHumanTask:
		return v.validateHumanTaskStep(stepID, step)
	case StepTypeSystemTask:
		return v.validateSystemTaskStep(stepID, step)
	case StepTypeCondition:
		return v.validateConditionStep(stepID, step)
	}

	return nil
}

func (v *defaultDSLValidator) validateHumanTaskStep(stepID string, step Step) error {
	if step.Form == "" {
		return fmt.Errorf("human task step %s has no form", stepID)
	}

	return nil
}

func (v *defaultDSLValidator) validateSystemTaskStep(stepID string, step Step) error {
	return nil
}

func (v *defaultDSLValidator) validateConditionStep(stepID string, step Step) error {
	return nil
}
