package services

import (
	"errors"
	"fmt"
	"strings"
)

const CONTROLLER_DIR string = "controllers"
const SERVICES_DIR string = "services"

type ComponentConfig struct {
	Name string
}

type ComponentService struct {
	filesystem FilesystemReaderWriter
	templates  TemplatesWriter
}

func NewComponentService(filesystem FilesystemReaderWriter, templates TemplatesWriter) ComponentService {
	return ComponentService{filesystem, templates}
}

func (s ComponentService) GenerateController(name string) error {
	hasDir, err := s.filesystem.HasDirectoryOrFile(CONTROLLER_DIR)
	if err != nil {
		return fmt.Errorf("unable to check if controller directory exists: %v", err.Error())
	}

	if !hasDir {
		if err := s.filesystem.CreateDirectory(CONTROLLER_DIR); err != nil {
			return fmt.Errorf("unable to create controller directory: %v", err.Error())
		}
	}

	controllerFilepath := fmt.Sprintf("%v/%v.go", CONTROLLER_DIR, strings.ToLower(name))

	hasFile, err := s.filesystem.HasDirectoryOrFile(controllerFilepath)
	if err != nil {
		return fmt.Errorf("unable to check if controller with that name already exists: %v", err.Error())
	}
	if hasFile {
		return errors.New("controller with that name already exists")
	}

	file, err := s.filesystem.CreateFile(controllerFilepath)
	if err != nil {
		return fmt.Errorf("unable to create controller file: %v", err.Error())
	}
	defer file.Close()

	controllerName := toSentenceCase(name)

	if err := s.templates.ExecuteTemplate(file, "controller.go.tmpl", ComponentConfig{Name: controllerName}); err != nil {
		return fmt.Errorf("unable to write template: %v", err.Error())
	}

	return nil
}

func (s ComponentService) GenerateService(name string) error {
	hasDir, err := s.filesystem.HasDirectoryOrFile(SERVICES_DIR)
	if err != nil {
		return fmt.Errorf("unable to check if services directory exists: %v", err.Error())
	}

	if !hasDir {
		if err := s.filesystem.CreateDirectory(SERVICES_DIR); err != nil {
			return fmt.Errorf("unable to create services directory: %v", err.Error())
		}
	}

	servicesFilepath := fmt.Sprintf("%v/%v.go", SERVICES_DIR, strings.ToLower(name))

	hasFile, err := s.filesystem.HasDirectoryOrFile(servicesFilepath)
	if err != nil {
		return fmt.Errorf("unable to check if service with that name already exists: %v", err.Error())
	}
	if hasFile {
		return errors.New("service with that name already exists")
	}

	file, err := s.filesystem.CreateFile(servicesFilepath)
	if err != nil {
		return fmt.Errorf("unable to create service file: %v", err.Error())
	}
	defer file.Close()

	serviceName := toSentenceCase(name)

	if err := s.templates.ExecuteTemplate(file, "service.go.tmpl", ComponentConfig{Name: serviceName}); err != nil {
		return fmt.Errorf("unable to write template: %v", err.Error())
	}

	return nil
}

func (s ComponentService) GenerateRepository(name string) error {
	return nil
}

func (s ComponentService) GenerateModel(name string) error {
	return nil
}

func (s ComponentService) GenerateView(name string) error {
	return nil
}

func toSentenceCase(s string) string {
	firstLetter := strings.ToUpper(string(s[0]))
	rest := strings.ToLower(s[1:])

	return fmt.Sprintf("%v%v", firstLetter, rest)
}
