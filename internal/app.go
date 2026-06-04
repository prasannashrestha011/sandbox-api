package app

import (
	"context"
	"errors"
	"log"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"main/internal/controllers"
	lab_handler "main/internal/controllers/lab"
	"main/internal/proxy"
	"main/internal/repository"
	lab_repo "main/internal/repository/lab"
	"main/internal/routes"
	"main/internal/sandbox/core"
	"main/internal/services"
	lab_services "main/internal/services/lab"
)

// Repos groups all repositories.
type Repos struct {
	SandboxRepo     repository.SandboxRepository
	UserRepo        repository.UserRepository
	RefreshRepo     repository.RefreshTokenRepository
	DockerImageRepo repository.DockerImageRepository
	LabRepo         lab_repo.LabRepository
	ChapterRepo     lab_repo.ChapterRepository
}

// Services groups all services.
type Services struct {
	SandboxService     services.SandboxService
	UserService        services.UserService
	AuthService        services.AuthService
	DockerImageService services.DockerImageService
	LabService         lab_services.LabService
	ChapterService     lab_services.ChapterService
}

// Controllers groups all controllers.

type Controllers struct {
	SandboxController     *controllers.SandboxController
	UserController        *controllers.UserController
	DockerImageController *controllers.DockerImageController
	PingerController      *controllers.PingerController
	// WebSocketController   *websocket.WebSocketController
	LabController     *lab_handler.LabController
	ChapterController *lab_handler.ChapterController
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
		LabRepo:         lab_repo.NewLabRepository(db),
		ChapterRepo:     lab_repo.NewChapterRepository(db),
	}

	servicesGroup := Services{
		SandboxService:     services.NewSandboxService(repos.SandboxRepo, repos.DockerImageRepo, sandboxClient),
		UserService:        services.NewUserService(repos.UserRepo),
		AuthService:        services.NewAuthService(repos.UserRepo, repos.RefreshRepo),
		DockerImageService: services.NewDockerImageService(repos.DockerImageRepo),
		LabService:         lab_services.NewLabService(repos.LabRepo),
		ChapterService:     lab_services.NewChapterService(repos.ChapterRepo),
	}

	controllersGroup := Controllers{
		SandboxController:     controllers.NewSandboxController(servicesGroup.SandboxService),
		UserController:        controllers.NewUserController(servicesGroup.UserService, servicesGroup.AuthService),
		DockerImageController: controllers.NewDockerImageController(servicesGroup.DockerImageService),
		PingerController:      controllers.NewPingerController(),
		// WebSocketController:   websocket.NewWebSocketController(servicesGroup.SandboxService),
		LabController:     lab_handler.NewLabController(servicesGroup.LabService),
		ChapterController: lab_handler.NewChapterController(servicesGroup.ChapterService),
	}

	//listeners for sandbox events (e.g., cleanup after timeout)
	go sandboxClient.ListenContainerEvents(context.Background(), func(containerID string) {
		log.Println("Received container event for container ID: ", containerID)
		// Find the sandbox associated with the container ID
		err := repos.SandboxRepo.UpdateStatus(context.Background(), containerID, "inactive")
		if err != nil {
			log.Println("Error updating sandbox status: ", err)
		}
	})
	router := chi.NewRouter()
	router.Use(proxy.ResponseWriterMiddleware)
	router.Use(proxy.ErrorMiddleware)
	router.Use(proxy.RateLimiterMiddleware)
	router.Get("/", controllersGroup.PingerController.Ping)
	routes.RegisterSandboxRoutes(router, controllersGroup.SandboxController)
	routes.RegisterUserRoutes(router, controllersGroup.UserController)
	routes.RegisterDockerImageRoutes(router, controllersGroup.DockerImageController)
	routes.RegisterLabRoutes(router, controllersGroup.LabController, controllersGroup.ChapterController)
	// websocket.RegisterWebSocketRoutes(router, controllersGroup.WebSocketController)
	return &App{
		Repos:       repos,
		Services:    servicesGroup,
		Controllers: controllersGroup,
		Router:      router,
	}, nil
}
