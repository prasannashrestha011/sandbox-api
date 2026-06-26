package app

import (
	"context"
	"errors"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gorm.io/gorm"

	"main/internal/controllers"
	lab_handler "main/internal/controllers/lab"
	"main/internal/jobs"
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
	SandboxRepo         repository.SandboxTemplateRepository
	UserRepo            repository.UserRepository
	RefreshRepo         repository.RefreshTokenRepository
	DockerImageRepo     repository.DockerImageRepository
	LabRepo             lab_repo.LabRepository
	ChapterRepo         lab_repo.ChapterRepository
	ExerciseRepo        lab_repo.ExerciseRepository
	EnrollmentRepo      lab_repo.EnrollmentRepository
	SubmissionRepo      lab_repo.SubmissionRepository
	SandboxTemplateRepo repository.SandboxTemplateRepository
	SandboxSessionRepo  repository.SandboxInstanceRepository
	WarmPoolRepo        repository.WarmPoolRepository
}

// Services groups all services.
type Services struct {
	SandboxService     services.SandboxTemplateService
	UserService        services.UserService
	AuthService        services.AuthService
	DockerImageService services.DockerImageService
	LabService         lab_services.LabService
	ChapterService     lab_services.ChapterService
	ExerciseService    lab_services.ExerciseService
	EnrollmentService  lab_services.EnrollmentService
	SubmissionService  lab_services.SubmissionService
	WarmPoolService    services.WarmpoolService

	SandboxTemplateService services.SandboxTemplateService
	SandboxInstanceService services.SandboxInstanceService
}

// Controllers groups all controllers.

type Controllers struct {
	SandboxController     *controllers.SandboxController
	UserController        *controllers.UserController
	DockerImageController *controllers.DockerImageController
	PingerController      *controllers.PingerController
	// WebSocketController   *websocket.WebSocketController
	LabController          *lab_handler.LabController
	ChapterController      *lab_handler.ChapterController
	ExerciseController     *lab_handler.ExerciseHandler
	EnrollmentController   *lab_handler.EnrollmentController
	SubmissionController   *lab_handler.SubmissionHandler
	SandboxTemplateHandler *controllers.SandboxController
	SandboxInstanceHandler *controllers.SandboxInstanceHandler
	WarmPoolHandler        *controllers.WarmPoolHandler
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
	asynqClient := jobs.InitAsynq()

	repos := Repos{
		SandboxRepo:         repository.NewSandboxTemplateRepository(db),
		UserRepo:            repository.NewUserRepository(db),
		RefreshRepo:         repository.NewRefreshTokenRepository(db),
		DockerImageRepo:     repository.NewDockerImageRepository(db),
		LabRepo:             lab_repo.NewLabRepository(db),
		ChapterRepo:         lab_repo.NewChapterRepository(db),
		ExerciseRepo:        lab_repo.NewExerciseRepository(db),
		EnrollmentRepo:      lab_repo.NewEnrollmentRepository(db),
		SubmissionRepo:      lab_repo.NewSubmissionRepository(db),
		SandboxTemplateRepo: repository.NewSandboxTemplateRepository(db),
		SandboxSessionRepo:  repository.NewSandboxInstanceRepository(db),
		WarmPoolRepo:        repository.NewWarmPoolRepository(db),
	}

	servicesGroup := Services{
		SandboxService:         services.NewSandboxTemplateService(repos.SandboxRepo, repos.DockerImageRepo, sandboxClient),
		UserService:            services.NewUserService(repos.UserRepo),
		AuthService:            services.NewAuthService(repos.UserRepo, repos.RefreshRepo),
		DockerImageService:     services.NewDockerImageService(repos.DockerImageRepo, sandboxClient),
		LabService:             lab_services.NewLabService(repos.LabRepo),
		ChapterService:         lab_services.NewChapterService(repos.ChapterRepo),
		ExerciseService:        lab_services.NewExerciseService(repos.ExerciseRepo),
		EnrollmentService:      lab_services.NewEnrollmentService(repos.EnrollmentRepo),
		SubmissionService:      lab_services.NewSubmissionService(repos.SubmissionRepo, repos.ExerciseRepo),
		SandboxTemplateService: services.NewSandboxTemplateService(repos.SandboxTemplateRepo, repos.DockerImageRepo, sandboxClient),
		SandboxInstanceService: services.NewSandboxInstanceService(repos.SandboxSessionRepo, repos.SandboxTemplateRepo, sandboxClient, asynqClient),
		WarmPoolService:        services.NewWarmpoolService(repos.WarmPoolRepo, asynqClient),
	}

	controllersGroup := Controllers{
		SandboxController:     controllers.NewSandboxController(servicesGroup.SandboxService),
		UserController:        controllers.NewUserController(servicesGroup.UserService, servicesGroup.AuthService),
		DockerImageController: controllers.NewDockerImageController(servicesGroup.DockerImageService),
		PingerController:      controllers.NewPingerController(),
		// WebSocketController:   websocket.NewWebSocketController(servicesGroup.SandboxService),
		LabController:          lab_handler.NewLabController(servicesGroup.LabService),
		ChapterController:      lab_handler.NewChapterController(servicesGroup.ChapterService),
		ExerciseController:     lab_handler.NewExerciseHandler(servicesGroup.ExerciseService),
		EnrollmentController:   lab_handler.NewEnrollmentController(servicesGroup.EnrollmentService),
		SubmissionController:   lab_handler.NewSubmissionHandler(servicesGroup.SubmissionService),
		SandboxTemplateHandler: controllers.NewSandboxController(servicesGroup.SandboxTemplateService),
		SandboxInstanceHandler: controllers.NewSandboxInstanceHandler(servicesGroup.SandboxInstanceService),
		WarmPoolHandler:        controllers.NewWarmPoolHandler(servicesGroup.WarmPoolService),
	}

	//listeners for sandbox events (e.g., cleanup after timeout)
	go sandboxClient.ListenContainerEvents(context.Background(), func(containerID string) {
		log.Println("Received container event for container ID: ", containerID)
		// Find the sandbox associated with the container ID
		// err := repos.SandboxRepo.UpdateStatus(context.Background(), containerID, "inactive")
		// if err != nil {
		// 	log.Println("Error updating sandbox status: ", err)
		// }
	})
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://yourdomain.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Use(proxy.ResponseWriterMiddleware)
	router.Use(proxy.ErrorMiddleware)
	router.Use(proxy.RateLimiterMiddleware)
	router.Get("/", controllersGroup.PingerController.Ping)
	routes.RegisterSandboxRoutes(router, controllersGroup.SandboxController, controllersGroup.SandboxInstanceHandler)
	routes.RegisterUserRoutes(router, controllersGroup.UserController)
	routes.RegisterDockerImageRoutes(router, controllersGroup.DockerImageController)
	routes.RegisterLabRoutes(router,
		controllersGroup.LabController,
		controllersGroup.ChapterController,
		controllersGroup.ExerciseController,
		controllersGroup.EnrollmentController,
		controllersGroup.SubmissionController)
	routes.RegisterWarmPoolRoutes(router, controllersGroup.WarmPoolHandler)
	// websocket.RegisterWebSocketRoutes(router, controllersGroup.WebSocketController)
	return &App{
		Repos:       repos,
		Services:    servicesGroup,
		Controllers: controllersGroup,
		Router:      router,
	}, nil
}
