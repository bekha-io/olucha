package api

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type DSLValidator interface {
	Validate() error
}

type defaultDSLValidator struct {
	dsl DSL
}

func NewDSLValidator(dsl DSL) DSLValidator {
	return &defaultDSLValidator{dsl: dsl}
}

func (v *defaultDSLValidator) Validate() error {
	var err = ErrInvalidFormat

	if v.dsl.ID == "" {
		err = fmt.Errorf("%w: id is required", err)
	}

	if v.dsl.Version == "" {
		err = fmt.Errorf("%w: version is required", err)
	}

	// Each DSL must have at least one step
	if v.dsl.Steps == nil {
		err = fmt.Errorf("%w: 0 steps are defined. please define at least one step", err)
	}

	// Validate the forms
	for formID, form := range v.dsl.Forms {
		err = v.validateForm(formID, form)
		if err != nil {
			err = fmt.Errorf("%w: form %s is invalid: %w", err, formID, err)
		}
	}

	// Validate each step
	for stepID, step := range v.dsl.Steps {
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

	if _, ok := v.dsl.Forms[step.Form]; !ok {
		return fmt.Errorf("form %s is not defined in the DSL", step.Form)
	}

	if err := step.ValidateForm(v.dsl, v.dsl.Forms[step.Form]); err != nil {
		return fmt.Errorf("human task step %s has invalid form: %w", stepID, err)
	}

	return nil
}

func (v *defaultDSLValidator) validateSystemTaskStep(stepID string, step Step) error {
	return nil
}

func (v *defaultDSLValidator) validateConditionStep(stepID string, step Step) error {
	return nil
}
