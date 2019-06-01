package bitrixentitygenerator

import (
	"bufio"
	"fmt"
)

type EntityGenerator interface {
	CollectEntityParameters(scanner *bufio.Scanner)
	CreateEntity()
}

type EntityPathData struct {
	Path string
}

const (
	defaultModuleName      = "my.module.name"
	defaultNamespace       = "myns"
	defaultComponentName   = "my.component"
	defaultLanguage        = "ru"
	defaultTemplateName    = ".default"
	defaultDirPermissions  = 0755
	defaultFilePermissions = 0644

	entityTypeComponent         = "c"
	entityTypeLangFile          = "l"
	entityTypeModule            = "m"
	entityTypeComponentTemplate = "t"
	// TODO - should be implemented later
	//entityTypeSiteTemplate = "s"
)

type generatorFactory func(path string) (EntityGenerator, error)

var registeredFactories = make(map[string]generatorFactory)

// RegisterFactories registers all available factories
func RegisterFactories() {
	registeredFactories[entityTypeComponent] = NewComponentGenerator
	registeredFactories[entityTypeLangFile] = NewLangFileGenerator
	registeredFactories[entityTypeModule] = NewModuleGenerator
	registeredFactories[entityTypeComponentTemplate] = NewTemplateGenerator
}

// ScanAndGetEntityType scans input and returns a code for new entity
func ScanAndGetEntityType(scanner *bufio.Scanner) string {
	fmt.Print("Enter entity type to create, ('c' => 'component'): ")
	scanner.Scan()
	entityType := clearTextValue(scanner.Text())
	if entityType == "" {
		entityType = entityTypeComponent
	}

	return entityType
}

// GetGenerator creates and returns new Generator entity or error
func GetGenerator(entityType, createPath string) (EntityGenerator, error) {
	eg, ok := registeredFactories[entityType]
	if !ok {
		return nil, fmt.Errorf("Invalid generator type '%s'", entityType)
	}

	return eg(createPath)
}
