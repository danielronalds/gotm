package controllers

import (
	"errors"
	"fmt"
	"strings"
)

type componentGenerator interface {
	GenerateController(name string) error
	GenerateService(name string) error
	GenerateRepository(name string) error
	GenerateMigration(name string) error
	GenerateModel(name string) error
	GenerateView(name string) error
	GeneratePage(name string) error
	GenerateDockerfile() error
}

type generator = func(name string) error
type dockerfileGenerator = func() error

type AddController struct {
	generatorMap    map[string]generator
	dockerGenerator dockerfileGenerator
}

func NewAddController(gen componentGenerator) AddController {
	generatorMap := map[string]generator{
		"controller": gen.GenerateController,
		"service":    gen.GenerateService,
		"repository": gen.GenerateRepository,
		"migration":  gen.GenerateMigration,
		"model":      gen.GenerateModel,
		"view":       gen.GenerateView,
		"page":       gen.GeneratePage,
	}

	dockerGenerator := gen.GenerateDockerfile

	return AddController{generatorMap, dockerGenerator}
}

func (c AddController) Handle(args []string) error {
	if len(args) == 0 || args[0] != "add" {
		return errors.New("passed to incorrect controller! Passed to `add` controller")
	}

	// Handling special case of the dockerfile, no named part of the component so only 2 args
	if len(args) == 2 && args[1] == "dockerfile" {
		if err := c.dockerGenerator(); err != nil {
			return fmt.Errorf("failed to generate dockerfile: %v", err.Error())
		}

		fmt.Println("Added Dockerfile")
		return nil
	}

	if len(args) < 3 || args[1] == "" || args[2] == "" {
		return errors.New("expected argument [component-type] [component-name]")
	}

	componentType := strings.ToLower(args[1])
	componentName := args[2]

	gen, ok := c.generatorMap[componentType]
	if !ok {
		return fmt.Errorf("\"%v\" is not a valid component", componentType)
	}

	if err := gen(componentName); err != nil {
		return fmt.Errorf("failed to generate %v component: %v", componentType, err.Error())
	}

	fmt.Printf("Added \"%v\" %v\n", componentName, componentType)

	return nil
}
