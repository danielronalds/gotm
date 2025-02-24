package services

import (
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
	return s.generateComponent(name, "controller", CONTROLLER_DIR, ".go", "controller.go.tmpl")
}

func (s ComponentService) GenerateService(name string) error {
	return s.generateComponent(name, "service", SERVICES_DIR, ".go", "service.go.tmpl")
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

// general method for dealing with the logic of generating a component.
//
// `fileExtension` should include the dot, i.e. ".go"
func (s ComponentService) generateComponent(name, componentType, componentDir, fileExtension, templateName string) error {
	hasDir, err := s.filesystem.HasDirectoryOrFile(componentDir)
	if err != nil {
		return fmt.Errorf("unable to check if %v directory exists: %v", componentDir, err.Error())
	}

	if !hasDir {
		if err := s.filesystem.CreateDirectory(componentDir); err != nil {
			return fmt.Errorf("unable to create %v directory: %v", componentDir, err.Error())
		}
	}

	componentFilepath := fmt.Sprintf("%v/%v%v", componentDir, strings.ToLower(name), fileExtension)

	hasFile, err := s.filesystem.HasDirectoryOrFile(componentFilepath)
	if err != nil {
		return fmt.Errorf("unable to check if %v with that name already exists: %v", componentType, err.Error())
	}
	if hasFile {
		return fmt.Errorf("%v with that name already exists", componentType)
	}

	file, err := s.filesystem.CreateFile(componentFilepath)
	if err != nil {
		return fmt.Errorf("unable to create %v file: %v", componentType, err.Error())
	}
	defer file.Close()

	controllerName := toSentenceCase(name)

	if err := s.templates.ExecuteTemplate(file, templateName, ComponentConfig{Name: controllerName}); err != nil {
		return fmt.Errorf("unable to write template: %v", err.Error())
	}

	return nil
}

func toSentenceCase(s string) string {
	firstLetter := strings.ToUpper(string(s[0]))
	rest := strings.ToLower(s[1:])

	return fmt.Sprintf("%v%v", firstLetter, rest)
}
