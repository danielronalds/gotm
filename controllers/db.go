package controllers

import (
	"errors"
	"fmt"
	"strings"
)

type dbInitialiser interface {
	Init() error
}

type DbController struct {
	db dbInitialiser
}

func NewDbController(db dbInitialiser) DbController {
	return DbController{db}
}

type handlerFunc = func(args []string) error

func (c DbController) Handle(args []string) error {
	if len(args) == 0 || args[0] != "db" {
		return errors.New("passed to incorrect controller! Passed to `db` controller")
	}

	slicedArgs := args[1:]

	cmd := "help"
	if len(slicedArgs) > 0 {
		cmd = strings.ToLower(slicedArgs[0])
	}

	handlerMap := map[string]handlerFunc{
		"help": c.handleHelp,
		"init": c.handleInit,
	}
	handler, ok := handlerMap[cmd]
	if !ok {
		handler = c.handleHelp
	}

	return handler(slicedArgs)
}

func (c DbController) handleHelp(args []string) error {
	fmt.Println("wip")
	return nil
}

func (c DbController) handleInit(args []string) error {
	if err := c.db.Init(); err != nil {
		return fmt.Errorf("unable to initialise database: %v", err.Error())
	}

	fmt.Println("\nInitialised database")

	return nil
}
