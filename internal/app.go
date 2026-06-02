package app

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"main/internal/controllers"
	"main/internal/proxy"
	"main/internal/repository"
	"main/internal/routes"
	"main/internal/sandbox/core"
	"main/internal/services"
)

// Repos groups all repositories.
type Repos struct {
	SandboxRepo     repository.SandboxRepository
	UserRepo        repository.UserRepository
	RefreshRepo     repository.RefreshTokenRepository
	DockerImageRepo repository.DockerImageRepository
	LabRepo         repository.LabRepository
}

// Services groups all services.
type Services struct {
	SandboxService     services.SandboxService
	UserService        services.UserService
	AuthService        services.AuthService
	DockerImageService services.DockerImageService
	LabService         services.LabService
}

// Controllers groups all controllers.

type Controllers struct {
	SandboxController     *controllers.SandboxController
	UserController        *controllers.UserController
	DockerImageController *controllers.DockerImageController
	PingerController      *controllers.PingerController
	// WebSocketController   *websocket.WebSocketController
	LabController *controllers.LabController
}

// App wires repositories, services, controllers, and routes.
type App struct {
	Repos       Repos
	Services    Services
	Controllers Controllers
	Router      *chi.Mux
}

// New constructs the application wiring and router.
func New(db *gorm.DB, sandboxClient core.SandboxClient) (*App, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	if sandboxClient == nil {
		return nil, errors.New("sandbox client is nil")
	}

	repos := Repos{
		SandboxRepo:     repository.NewSandboxRepository(db),
		UserRepo:        repository.NewUserRepository(db),
		RefreshRepo:     repository.NewRefreshTokenRepository(db),
		DockerImageRepo: repository.NewDockerImageRepository(db),
		LabRepo:         repository.NewLabRepository(db),
	}

	servicesGroup := Services{
		SandboxService:     services.NewSandboxService(repos.SandboxRepo, repos.DockerImageRepo, sandboxClient),
		UserService:        services.NewUserService(repos.UserRepo),
		AuthService:        services.NewAuthService(repos.UserRepo, repos.RefreshRepo),
		DockerImageService: services.NewDockerImageService(repos.DockerImageRepo),
		LabService:         services.NewLabService(repos.LabRepo),
	}

	controllersGroup := Controllers{
		SandboxController:     controllers.NewSandboxController(servicesGroup.SandboxService),
		UserController:        controllers.NewUserController(servicesGroup.UserService, servicesGroup.AuthService),
		DockerImageController: controllers.NewDockerImageController(servicesGroup.DockerImageService),
		PingerController:      controllers.NewPingerController(),
		// WebSocketController:   websocket.NewWebSocketController(servicesGroup.SandboxService),
		LabController: controllers.NewLabController(servicesGroup.LabService),
	}

	router := chi.NewRouter()
	router.Use(proxy.ResponseWriterMiddleware)
	router.Use(proxy.ErrorMiddleware)
	router.Use(proxy.RateLimiterMiddleware)
	router.Get("/", controllersGroup.PingerController.Ping)
	routes.RegisterSandboxRoutes(router, controllersGroup.SandboxController)
	routes.RegisterUserRoutes(router, controllersGroup.UserController)
	routes.RegisterDockerImageRoutes(router, controllersGroup.DockerImageController)
	routes.RegisterLabRoutes(router, controllersGroup.LabController)
	// websocket.RegisterWebSocketRoutes(router, controllersGroup.WebSocketController)
	return &App{
		Repos:       repos,
		Services:    servicesGroup,
		Controllers: controllersGroup,
		Router:      router,
	}, nil
}
