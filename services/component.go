package services

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const CONTROLLERS_DIR string = "controllers"
const SERVICES_DIR string = "services"
const REPOSITORIES_DIR string = "repositories"
const MIGRATION_DIR string = "migrations"
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

func (s ComponentService) GenerateMigration(name string) error {
	// Generating timestamp that matches how goose generates timestamps
	timestamp := time.Now().UTC().Format("20060102150405")

	filename := fmt.Sprintf("%v_%v.sql", timestamp, name)
	return s.generateComponent(name, "migration", s.filesystem.FromRoot(MIGRATION_DIR), filename, "migration.sql.tmpl")
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

// Generates a migraton for creating a table with the given columns
//
// Columns should be a map with column names as the key and their type being the value
func (s ComponentService) GenerateTable(name string, columns map[string]string) error {
	tableTemplateData := struct {
		Name    string
		Columns map[string]string
	}{Name: name, Columns: columns}

	// Opening tempate file and writing the file
	timestamp := time.Now().UTC().Format("20060102150405")
	filename := fmt.Sprintf("%v_add_%v_table.sql", timestamp, name)

	hasDir, err := s.filesystem.HasDirectoryOrFile(MIGRATION_DIR)
	if err != nil {
		return fmt.Errorf("unable to check if %v directory exists: %v", MIGRATION_DIR, err.Error())
	}

	if !hasDir {
		if err := s.filesystem.CreateDirectory(MIGRATION_DIR); err != nil {
			return fmt.Errorf("unable to create %v directory: %v", MIGRATION_DIR, err.Error())
		}
	}

	// No need to check if this is unique, as the filename contains a timestamp
	migrationFilepath := fmt.Sprintf("%v/%v", MIGRATION_DIR, filename)

	file, err := s.filesystem.CreateFile(migrationFilepath)
	if err != nil {
		return fmt.Errorf("unable to create migraiton file: %v", err.Error())
	}
	defer file.Close()

	// Writing to the file
	if err := s.templates.WriteTemplate(file, "table.sql.tmpl", tableTemplateData); err != nil {
		return fmt.Errorf("unable to write template: %v", err.Error())
	}

	return nil
}

func toSentenceCase(s string) string {
	firstLetter := strings.ToUpper(string(s[0]))
	rest := strings.ToLower(s[1:])

	return fmt.Sprintf("%v%v", firstLetter, rest)
}
