package controllers

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type columnName = string
type columnType = string
type columns = map[columnName]columnType

type componentGenerator interface {
	GenerateController(name string) error
	GenerateService(name string) error
	GenerateRepository(name string) error
	GenerateMigration(name string) error
	GenerateModel(name string) error
	GenerateView(name string) error
	GeneratePage(name string) error
	GenerateDockerfile() error
	GenerateSqlcConfigIfNotExists() (bool, error)
	GenerateTable(name string, cols columns) error
	GenerateQueries(name string, cols columns) error
}

type generatorFunc = func(name string) error
type dockerGeneratorFunc = func() error
type tableGeneratorFunc = func(name string, cols columns) error
type queryGeneratorFunc = func(name string, cols columns) error
type sqlcGeneratorFunc = func() (bool, error)

type sqlcRunner interface {
	RunSqlc() error
}

type AddController struct {
	generatorMap    map[string]generatorFunc
	dockerGenerator dockerGeneratorFunc
	tableGenerator  tableGeneratorFunc
	queryGenerator  queryGeneratorFunc
	sqlcGenerator   sqlcGeneratorFunc
	sqlc            sqlcRunner
}

func NewAddController(gen componentGenerator, sqlc sqlcRunner) AddController {
	generatorMap := map[string]generatorFunc{
		"controller": gen.GenerateController,
		"service":    gen.GenerateService,
		"repository": gen.GenerateRepository,
		"migration":  gen.GenerateMigration,
		"model":      gen.GenerateModel,
		"view":       gen.GenerateView,
		"page":       gen.GeneratePage,
	}

	dockerGenerator := gen.GenerateDockerfile
	tableGeneratorFunc := gen.GenerateTable
	queryGeneraterFunc := gen.GenerateQueries
	sqlcGeneratorFunc := gen.GenerateSqlcConfigIfNotExists

	return AddController{generatorMap, dockerGenerator, tableGeneratorFunc, queryGeneraterFunc, sqlcGeneratorFunc, sqlc}
}

func (c AddController) Handle(args []string) error {
	if len(args) == 0 || args[0] != "add" {
		return errors.New("passed to incorrect controller! Passed to `add` controller")
	}

	// Handling special case of the dockerfile, no named part of the component so only 2 args
	if len(args) == 2 && args[1] == "dockerfile" {
		return c.handleDockerFile()
	}

	// Handling special case of the table, as there are more than 3 args
	if len(args) < 3 || args[1] == "" || args[2] == "" {
		return errors.New("expected argument [component-type] [component-name]")
	}

	componentType := strings.ToLower(args[1])
	componentName := args[2]

	if componentType == "table" {
		return c.handleTable(args)
	}

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

func (c AddController) handleDockerFile() error {
	if err := c.dockerGenerator(); err != nil {
		return fmt.Errorf("failed to generate dockerfile: %v", err.Error())
	}

	fmt.Println("Added Dockerfile")
	return nil
}

func (c AddController) handleTable(args []string) error {
	if len(args) < 4 {
		return errors.New("no columns supplied")
	}

	tableName := args[2]
	keyValuePairs := args[3:]
	columns, err := generateKeyValuePairs(keyValuePairs, "=")
	if err != nil {
		return fmt.Errorf("failed to parse columns: %v", err.Error())
	}

	configCreated, err := c.sqlcGenerator()
	if err != nil {
		return fmt.Errorf("failed to generate sqlc config file: %v", err.Error())
	}
	if configCreated {
		log.Println("created sqlc.yml")
	}

	if err := c.queryGenerator(tableName, columns); err != nil {
		return fmt.Errorf("failed to generate queries: %v", err.Error())
	}

	if err := c.tableGenerator(tableName, columns); err != nil {
		return fmt.Errorf("failed to generate table: %v", err.Error())
	}

	if err := c.sqlc.RunSqlc(); err != nil {
		return fmt.Errorf("failed to run sqlc: %v", err.Error())
	}

	fmt.Printf("Added \"%v\" table\n", tableName)

	return nil
}

func generateKeyValuePairs(pairs []string, sep string) (map[string]string, error) {
	keyValueMap := make(map[string]string, 0)

	for _, pair := range pairs {
		splitPair := strings.Split(pair, sep)
		if len(splitPair) != 2 {
			return keyValueMap, fmt.Errorf("incorrectly formated pair: %v", splitPair)
		}

		key := splitPair[0]
		Value := splitPair[1]

		_, ok := keyValueMap[key]
		if ok {
			return keyValueMap, fmt.Errorf("repeated pair decleration: %v", key)
		}
		keyValueMap[key] = Value
	}

	return keyValueMap, nil
}
