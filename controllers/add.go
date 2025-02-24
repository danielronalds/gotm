package controllers

import (
	"errors"
	"fmt"
	"strings"
)

type ComponentGenerator interface {
	GenerateController(name string) error
	GenerateService(name string) error
	GenerateRepository(name string) error
	GenerateModel(name string) error
	GenerateView(name string) error
}

type Generator = func(name string) error

type AddController struct {
	generatorMap map[string]Generator
}

func NewAddController(generator ComponentGenerator) AddController {
	generatorMap := map[string]Generator{
		"controllor": generator.GenerateController,
		"service":    generator.GenerateService,
		"repository": generator.GenerateRepository,
		"model":      generator.GenerateModel,
		"view":       generator.GenerateView,
	}
	return AddController{generatorMap}
}

func (c AddController) HandleCmd(args []string) error {
	if len(args) == 0 || args[0] != "add" {
		return errors.New("passed to incorrect controller! Passed to `add` controller")
	}

	if len(args) < 3 || args[1] == "" || args[2] == "" {
		return errors.New("expected argument [component-type] [component-name]")
	}

	componentType := strings.ToLower(args[1])
	componentName := args[2]

	generator, ok := c.generatorMap[componentType]
	if !ok {
		return fmt.Errorf("\"%v\" is not a valid component", componentType)
	}

	if err := generator(componentName); err != nil {
		return fmt.Errorf("failed to generate %v component: %v", componentType, err.Error())
	}

	fmt.Printf("Added \"%v\" %v\n", componentName, componentType)

	return nil
}
