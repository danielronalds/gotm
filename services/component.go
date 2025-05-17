package services

import (
	"errors"
	"fmt"
	"strings"
)

const CONTROLLERS_DIR string = "controllers"
const SERVICES_DIR string = "services"
const REPOSITORIES_DIR string = "repositories"
const MIDDLEWARE_DIR string = "middleware"
const MODELS_DIR string = "frontend/src/models"
const VIEWS_DIR string = "frontend/src/views"
const PAGES_DIR string = "frontend/src/views/pages"

const VIEW_COMPONENT_TYPE = "view"

type ComponentServiceFilesystem interface {
	ProjectRoot
	FileCreater
	DirCreater
	DirReader
}

type ComponentService struct {
	filesystem ComponentServiceFilesystem
	templates  TemplatesWriter
}

func NewComponentService(filesystem ComponentServiceFilesystem, templates TemplatesWriter) ComponentService {
	return ComponentService{filesystem, templates}
}

func (s ComponentService) GenerateController(name string) error {
	filename := fmt.Sprintf("%v.go", name)
	return s.generateComponent(name, "controller", s.filesystem.FromRoot(CONTROLLERS_DIR), filename, "controller.go.tmpl")
}

func (s ComponentService) GenerateService(name string) error {
	filename := fmt.Sprintf("%v.go", name)
	return s.generateComponent(name, "service", s.filesystem.FromRoot(SERVICES_DIR), filename, "service.go.tmpl")
}

func (s ComponentService) GenerateRepository(name string) error {
	filename := fmt.Sprintf("%v.go", name)
	return s.generateComponent(name, "repository", s.filesystem.FromRoot(REPOSITORIES_DIR), filename, "repository.go.tmpl")
}

func (s ComponentService) GenerateMiddleware(name string) error {
	filename := fmt.Sprintf("%v.go", name)
	return s.generateComponent(name, "middleware", s.filesystem.FromRoot(MIDDLEWARE_DIR), filename, "middleware.go.tmpl")
}

func (s ComponentService) GenerateModel(name string) error {
	filename := fmt.Sprintf("%v.ts", name)
	return s.generateComponent(name, "model", s.filesystem.FromRoot(MODELS_DIR), filename, "model.ts.tmpl")
}

func (s ComponentService) GenerateView(name string) error {
	filename := fmt.Sprintf("%v.ts", name)
	return s.generateComponent(name, VIEW_COMPONENT_TYPE, s.filesystem.FromRoot(VIEWS_DIR), filename, "view.ts.tmpl")
}

func (s ComponentService) GeneratePage(name string) error {
	filename := fmt.Sprintf("%vPage.ts", toSentenceCase(name))
	return s.generateComponent(name, "page", s.filesystem.FromRoot(PAGES_DIR), filename, "page.ts.tmpl")
}

// general method for dealing with the logic of generating a component.
//
// `fileExtension` should include the dot, i.e. ".go"
func (s ComponentService) generateComponent(name, componentType, componentDir, filename, templateName string) error {
	hasDir, err := s.filesystem.HasDirectoryOrFile(componentDir)
	if err != nil {
		return fmt.Errorf("unable to check if %v directory exists: %v", componentDir, err.Error())
	}

	if !hasDir {
		if err := s.filesystem.CreateDirectory(componentDir); err != nil {
			return fmt.Errorf("unable to create %v directory: %v", componentDir, err.Error())
		}
	}

	componentFilepath := fmt.Sprintf("%v/%v", componentDir, filename)

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

	componentName := toSentenceCase(name)
	if componentType == VIEW_COMPONENT_TYPE {
		componentName = name
	}
	if err := s.templates.WriteTemplate(file, templateName, struct {
		Name          string
		LowerCaseName string
	}{Name: componentName, LowerCaseName: strings.ToLower(componentName)}); err != nil {
		return fmt.Errorf("unable to write template: %v", err.Error())
	}

	return nil
}

func (s ComponentService) GenerateDockerfile() error {
	hasFile, err := s.filesystem.HasDirectoryOrFile(s.filesystem.FromRoot("Dockerfile"))
	if err != nil {
		return fmt.Errorf("unable to check if Dockerfile already exists: %v", err.Error())
	}
	if hasFile {
		return errors.New("Dockerfile already exists")
	}

	file, err := s.filesystem.CreateFile(s.filesystem.FromRoot("Dockerfile"))
	if err != nil {
		return fmt.Errorf("unable to create dockerfile: %v", err.Error())
	}
	defer file.Close()

	if err := s.templates.WriteTemplate(file, "Dockerfile.tmpl", nil); err != nil {
		return fmt.Errorf("unable to write template: %v", err.Error())
	}

	return nil
}

func toSentenceCase(s string) string {
	firstLetter := strings.ToUpper(string(s[0]))
	rest := strings.ToLower(s[1:])

	return fmt.Sprintf("%v%v", firstLetter, rest)
}
