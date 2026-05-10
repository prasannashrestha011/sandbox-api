package app

import (
	"errors"

	"gorm.io/gorm"

	"main/internal/controllers"
	"main/internal/repository"
	"main/internal/sandbox/core"
	"main/internal/services"
)

// Repos groups all repositories.
type Repos struct {
	SandboxRepo repository.SandboxRepository
}

// Services groups all services.
type Services struct {
	SandboxService services.SandboxService
}

// Controllers groups all controllers.
type Controllers struct {
	SandboxController *controllers.SandboxController
}

// App wires repositories, services, controllers, and routes.
type App struct {
	Repos       Repos
	Services    Services
	Controllers Controllers
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
		SandboxRepo: repository.NewSandboxRepository(db),
	}

	servicesGroup := Services{
		SandboxService: services.NewSandboxService(repos.SandboxRepo, sandboxClient),
	}

	controllersGroup := Controllers{
		SandboxController: controllers.NewSandboxController(servicesGroup.SandboxService),
	}

	return &App{
		Repos:       repos,
		Services:    servicesGroup,
		Controllers: controllersGroup,
	}, nil
}
