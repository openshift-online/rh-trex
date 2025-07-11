package generator

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/openshift-online/rh-trex/pkg/core/api"
	"github.com/openshift-online/rh-trex/pkg/core/controllers"
	"github.com/openshift-online/rh-trex/pkg/core/dao"
	"github.com/openshift-online/rh-trex/pkg/core/services"
	"gorm.io/gorm"
)

// GeneratorConfig defines what to generate
type GeneratorConfig struct {
	KindName     string
	ProjectName  string
	ModulePath   string
	OutputDir    string
	Templates    []string
}

// Generator creates new resources using core patterns
type Generator struct {
	config    GeneratorConfig
	templates map[string]*template.Template
}

// NewGenerator creates a new generator
func NewGenerator(config GeneratorConfig) *Generator {
	return &Generator{
		config:    config,
		templates: make(map[string]*template.Template),
	}
}

// Generate creates a new resource with all necessary files
func (g *Generator) Generate() error {
	// 1. Generate base files from templates
	// 2. Update existing files (routes, migrations, etc.)
	// 3. Register with controller framework
	// 4. Run post-generation hooks
	
	return nil
}

// ResourceFactory creates and wires up a complete resource
type ResourceFactory struct {
	db              *gorm.DB
	controllerMgr   *controllers.ControllerManager
	eventEmitter    services.EventEmitter
}

// NewResourceFactory creates a new resource factory
func NewResourceFactory(
	db *gorm.DB,
	controllerMgr *controllers.ControllerManager,
	eventEmitter services.EventEmitter,
) *ResourceFactory {
	return &ResourceFactory{
		db:            db,
		controllerMgr: controllerMgr,
		eventEmitter:  eventEmitter,
	}
}

// CreateResource creates a complete resource with DAO, service, and controller
func CreateResource[T any](
	factory *ResourceFactory,
	kindName string,
	model T,
) services.CRUDService[T] {
	// Create DAO
	resourceDAO := dao.NewBaseDAO[T](factory.db)
	
	// Create service
	service := services.NewBaseCRUDService[T](resourceDAO, factory.eventEmitter, kindName)
	
	// Auto-register with controller manager
	if factory.controllerMgr != nil {
		controllers.AutoRegisterCRUDController(factory.controllerMgr, service, kindName)
	}
	
	return service
}

// ResourceInfo contains metadata about a resource
type ResourceInfo struct {
	Name       string
	Type       reflect.Type
	TableName  string
	Fields     []FieldInfo
}

// FieldInfo contains metadata about a field
type FieldInfo struct {
	Name     string
	Type     string
	JSONTag  string
	GORMTag  string
	Required bool
}

// AnalyzeResource analyzes a resource type and returns metadata
func AnalyzeResource[T any](model T) *ResourceInfo {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	
	info := &ResourceInfo{
		Name:      modelType.Name(),
		Type:      modelType,
		TableName: toSnakeCase(modelType.Name()),
		Fields:    make([]FieldInfo, 0),
	}
	
	// Analyze fields
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		
		fieldInfo := FieldInfo{
			Name:     field.Name,
			Type:     field.Type.String(),
			JSONTag:  field.Tag.Get("json"),
			GORMTag:  field.Tag.Get("gorm"),
			Required: !strings.Contains(field.Tag.Get("json"), "omitempty"),
		}
		
		info.Fields = append(info.Fields, fieldInfo)
	}
	
	return info
}

// ValidateResource validates that a resource implements required interfaces
func ValidateResource[T any](service services.CRUDService[T]) error {
	// Check if service implements required methods
	serviceType := reflect.TypeOf(service)
	
	requiredMethods := []string{
		"Get", "Create", "Replace", "Delete", "List", "FindByIDs",
		"OnUpsert", "OnDelete",
	}
	
	for _, methodName := range requiredMethods {
		method, exists := serviceType.MethodByName(methodName)
		if !exists {
			return fmt.Errorf("service missing required method: %s", methodName)
		}
		
		// Basic method signature validation could be added here
		_ = method
	}
	
	return nil
}

// TemplateVars contains variables for template generation
type TemplateVars struct {
	Kind                string
	KindPlural          string
	KindLowerSingular   string
	KindLowerPlural     string
	KindSnakeCasePlural string
	Project             string
	Module              string
	Fields              []FieldInfo
}

// GenerateTemplateVars creates template variables from a resource
func GenerateTemplateVars[T any](model T, project, module string) *TemplateVars {
	info := AnalyzeResource(model)
	
	return &TemplateVars{
		Kind:                info.Name,
		KindPlural:          info.Name + "s",
		KindLowerSingular:   strings.ToLower(info.Name[:1]) + info.Name[1:],
		KindLowerPlural:     strings.ToLower(info.Name[:1]) + info.Name[1:] + "s",
		KindSnakeCasePlural: info.TableName + "s",
		Project:             project,
		Module:              module,
		Fields:              info.Fields,
	}
}

// RegisterResourceType automatically registers a resource type with all systems
func RegisterResourceType[T any](
	factory *ResourceFactory,
	kindName string,
	model T,
) services.CRUDService[T] {
	// Create the complete resource
	service := CreateResource(factory, kindName, model)
	
	// Validate the resource
	if err := ValidateResource(service); err != nil {
		panic(fmt.Sprintf("Resource validation failed for %s: %v", kindName, err))
	}
	
	return service
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// EventEmitterAdapter adapts the existing event service to the EventEmitter interface
type EventEmitterAdapter struct {
	eventService interface {
		Create(ctx context.Context, event *api.Event) (*api.Event, error)
	}
}

// NewEventEmitterAdapter creates a new event emitter adapter
func NewEventEmitterAdapter(eventService interface {
	Create(ctx context.Context, event *api.Event) (*api.Event, error)
}) *EventEmitterAdapter {
	return &EventEmitterAdapter{
		eventService: eventService,
	}
}

// EmitEvent emits an event through the adapted service
func (e *EventEmitterAdapter) EmitEvent(ctx context.Context, source, sourceID string, eventType api.EventType) error {
	event := &api.Event{
		Source:    source,
		SourceID:  sourceID,
		EventType: eventType,
	}
	
	_, err := e.eventService.Create(ctx, event)
	return err
}