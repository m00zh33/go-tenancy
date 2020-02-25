package application

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"github.com/qor/admin"
	"github.com/qor/assetfs"
	"github.com/qor/middlewares"
	"github.com/qor/wildcard_router"
)

// MicroAppInterface micro app interface
type MicroAppInterface interface {
	ConfigureApplication(*Application)
}

// Application main application
type Application struct {
	*Config
}

// Config application config
type Config struct {
	IrisApp  *iris.Application
	Handlers []http.Handler
	AssetFS  assetfs.Interface
	Admin    *admin.Admin
	DB       *gorm.DB
}

// New new application
func New(cfg *Config) *Application {
	if cfg == nil {
		cfg = &Config{}
	}

	if cfg.IrisApp == nil {
		cfg.IrisApp = iris.New()
	}

	if cfg.AssetFS == nil {
		cfg.AssetFS = assetfs.AssetFS()
	}

	return &Application{
		Config: cfg,
	}
}

// Use mount router into micro app
func (application *Application) Use(app MicroAppInterface) {
	app.ConfigureApplication(application)
}

// NewServeMux allocates and returns a new ServeMux.
func (application *Application) NewServeMux() http.Handler {
	if len(application.Config.Handlers) == 0 {
		return middlewares.Apply(application.Config.IrisApp)
	}

	wildcardRouter := wildcard_router.New()
	for _, handler := range application.Config.Handlers {
		wildcardRouter.AddHandler(handler)
	}
	wildcardRouter.AddHandler(application.Config.IrisApp)

	return middlewares.Apply(wildcardRouter)
}
