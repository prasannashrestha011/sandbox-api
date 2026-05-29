package app

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"main/internal/controllers"
	"main/internal/pkg"
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
}

// Services groups all services.
type Services struct {
	SandboxService     services.SandboxService
	UserService        services.UserService
	AuthService        services.AuthService
	DockerImageService services.DockerImageService
}

// Controllers groups all controllers.
type Controllers struct {
	SandboxController     *controllers.SandboxController
	UserController        *controllers.UserController
	DockerImageController *controllers.DockerImageController
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
	}

	servicesGroup := Services{
		SandboxService:     services.NewSandboxService(repos.SandboxRepo, sandboxClient),
		UserService:        services.NewUserService(repos.UserRepo),
		AuthService:        services.NewAuthService(repos.UserRepo, repos.RefreshRepo),
		DockerImageService: services.NewDockerImageService(repos.DockerImageRepo),
	}

	controllersGroup := Controllers{
		SandboxController:     controllers.NewSandboxController(servicesGroup.SandboxService),
		UserController:        controllers.NewUserController(servicesGroup.UserService, servicesGroup.AuthService),
		DockerImageController: controllers.NewDockerImageController(servicesGroup.DockerImageService),
	}

	router := chi.NewRouter()
	routes.RegisterSandboxRoutes(router, controllersGroup.SandboxController)
	routes.RegisterUserRoutes(router, controllersGroup.UserController)
	routes.RegisterDockerImageRoutes(router, controllersGroup.DockerImageController)
	pkg.NewWebSocket(router)
	return &App{
		Repos:       repos,
		Services:    servicesGroup,
		Controllers: controllersGroup,
		Router:      router,
	}, nil
}
