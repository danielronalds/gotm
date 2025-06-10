package services

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

const CONTROLLERS_DIR string = "controllers"
const SERVICES_DIR string = "services"
const REPOSITORIES_DIR string = "repositories"
const MIGRATION_DIR string = "migrations"
const QUERY_DIR string = "queries"
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

// Utility function for dealing with the logic of generating a typical component.
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

func (s ComponentService) GenerateSqlcConfigIfNotExists() (bool, error) {
	hasFile, err := s.filesystem.HasDirectoryOrFile(s.filesystem.FromRoot("sqlc.yml"))
	if err != nil {
		return false, fmt.Errorf("unable to check if sqlc.yml already exists: %v", err.Error())
	}
	if hasFile {
		return false, nil
	}

	file, err := s.filesystem.CreateFile(s.filesystem.FromRoot("sqlc.yml"))
	if err != nil {
		return false, fmt.Errorf("unable to create sqlc.yml: %v", err.Error())
	}
	defer file.Close()

	if err := s.templates.WriteTemplate(file, "sqlc.yml.tmpl", nil); err != nil {
		return false, fmt.Errorf("unable to write template: %v", err.Error())
	}

	return true, nil
}

type columnName = string
type columnType = string
type columns = map[columnName]columnType

// type of data passed to a database specific template
type databaseTemplateData struct {
	Name                  string
	NameSentenceCase      string
	Columns               columns
	ColumnsUpperCaseTypes columns
}

func newDatabaseTemplateData(name string, cols columns) (databaseTemplateData, error) {
	colsUpperCase := make(columns, 0)

	for name, colType := range cols {
		formattedType := strings.ToUpper(colType)
		if !isValidSqliteType(formattedType) {
			return databaseTemplateData{}, fmt.Errorf("%v is not a valid sqlite type", colType)
		}
		colsUpperCase[name] = formattedType
	}

	return databaseTemplateData{
		Name:                  name,
		NameSentenceCase:      toSentenceCase(name),
		Columns:               cols,
		ColumnsUpperCaseTypes: colsUpperCase,
	}, nil
}

func (s ComponentService) GenerateTable(name string, cols columns) error {
	tableTemplateData, err := newDatabaseTemplateData(name, cols)
	if err != nil {
		return err
	}

	// Opening tempate file and writing the file
	timestamp := time.Now().UTC().Format("20060102150405")
	filename := fmt.Sprintf("%v_add_%v_table.sql", timestamp, name)

	return s.writeDatabaseTemplateFile(filename, "table.sql.tmpl", MIGRATION_DIR, tableTemplateData)
}

func (s ComponentService) GenerateQueries(name string, cols columns) error {
	tableTemplateData, err := newDatabaseTemplateData(name, cols)
	if err != nil {
		return err
	}

	// Opening tempate file and writing the file
	filename := fmt.Sprintf("%v.sql", name)

	return s.writeDatabaseTemplateFile(filename, "queries.sql.tmpl", QUERY_DIR, tableTemplateData)
}

// Utility function for handling writing a databse related template file
func (s ComponentService) writeDatabaseTemplateFile(filename, template, dir string, data databaseTemplateData) error {
	hasDir, err := s.filesystem.HasDirectoryOrFile(dir)
	if err != nil {
		return fmt.Errorf("unable to check if %v directory exists: %v", dir, err.Error())
	}

	if !hasDir {
		if err := s.filesystem.CreateDirectory(dir); err != nil {
			return fmt.Errorf("unable to create %v directory: %v", dir, err.Error())
		}
	}

	// No need to check if this is unique, as the filename contains a timestamp
	templateFilepath := fmt.Sprintf("%v/%v", dir, filename)

	file, err := s.filesystem.CreateFile(templateFilepath)
	if err != nil {
		return fmt.Errorf("unable to create template file: %v", err.Error())
	}
	defer file.Close()

	// Writing to the file
	if err := s.templates.WriteTemplate(file, template, data); err != nil {
		return fmt.Errorf("unable to write template: %v", err.Error())
	}

	return nil
}

func toSentenceCase(s string) string {
	firstLetter := strings.ToUpper(string(s[0]))
	rest := strings.ToLower(s[1:])

	return fmt.Sprintf("%v%v", firstLetter, rest)
}

// Array of all possible Sqlite Types, fetched from https://www.sqlite.org/datatype3.html
var sqliteTypes = []string{
	"INT",
	"INTEGER",
	"TINYINT",
	"SMALLINT",
	"MEDIUMINT",
	"BIGINT",
	"UNSIGNED BIG INT",
	"INT2",
	"INT8",
	"CHARACTER(20)",
	"VARCHAR(255)",
	"VARYING CHARACTER(255)",
	"NCHAR(55)",
	"NATIVE CHARACTER(70)",
	"NVARCHAR(100)",
	"TEXT",
	"CLOB",
	"BLOB",
	"REAL",
	"DOUBLE",
	"DOUBLE PRECISION",
	"FLOAT",
	"NUMERIC",
	"BOOLEAN",
	"DATE",
	"DATETIME",
}

func isValidSqliteType(t string) bool {
	return slices.Contains(sqliteTypes, t)
}
