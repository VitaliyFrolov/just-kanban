package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"just-kanban/internal/config"
	"just-kanban/internal/handlers"
	"just-kanban/internal/middlewares"
	repositorysql "just-kanban/internal/repositories/sql"
	"just-kanban/internal/services"
	"just-kanban/pkg/database"
	"just-kanban/pkg/router"
	"just-kanban/pkg/validation"
)

type App struct {
	*http.ServeMux
	*sql.DB
	*config.AllowedHTTPMethods
	*config.URLPaths
	*config.Env
	*validation.Validate
	services.UserService
	*services.TaskService
	*services.AuthService
	*services.TokenService
	*services.BoardService
	*services.BoardMemberService
}

func NewApp() *App {
	env := config.NewEnv()
	app := App{Env: env}
	app.initDatabase()
	app.initValidator()
	app.initServices()
	app.initRouter()
	app.runListen()
	return &app
}

func (app *App) initDatabase() {
	sqlDB := database.NewPostgresConnection(
		app.Env.DBUser,
		app.Env.DBPassword,
		app.Env.DBHost,
		app.Env.DBPort,
		app.Env.DBName,
	)
	app.DB = sqlDB
}

func (app *App) initValidator() {
	validator := validation.NewValidator()
	registerErr := validation.RegisterValidationTagTrimmed(validator)
	if registerErr != nil {
		panic("Can't register validator" + registerErr.Error())
	}
	app.Validate = validator
}

func (app *App) initServices() {
	// WARNING! Right services init order is required
	app.UserService = services.NewUserService(repositorysql.NewUserRepository(app.DB))
	app.TaskService = services.NewTaskService(repositorysql.NewTaskRepository(app.DB))
	app.TokenService = services.NewTokenService(
		repositorysql.NewRefreshTokenRepository(app.DB),
		app.Env.JWTSecret,
	)
	app.AuthService = services.NewAuthService(app.TokenService, app.UserService)
	app.BoardService = services.NewBoardService(
		repositorysql.NewBoardRepository(app.DB),
		app.TaskService,
	)
	app.BoardMemberService = services.NewBoardMemberService(
		repositorysql.NewBoardMemberRepository(app.DB),
		app.BoardService,
		app.UserService,
	)
}

func (app *App) initPaths() {
	paths, allowedMethods := config.NewHTTPPaths()
	app.URLPaths = paths
	app.AllowedHTTPMethods = allowedMethods
}

func (app *App) initSecureHandlers() {
	secureRoutes := router.NewGroup(app.ServeMux, "")
	secureRoutes.Use(
		func(handler http.Handler) http.Handler {
			return middlewares.Auth(handler, app.Env)
		},
	)
	secureRoutes.Handle(app.URLPaths.LogoutHandler, handlers.NewLogoutHandler(app.AuthService))
	secureRoutes.Handle(app.URLPaths.UsersHandler, handlers.NewUserHandler(app.UserService, app.Validate))
	secureRoutes.Handle(
		app.URLPaths.BoardMembersHandler,
		handlers.NewBoardMemberHandler(app.BoardMemberService, app.Validate),
	)
	secureRoutes.Handle(
		app.URLPaths.BoardMemberHandler,
		handlers.NewBoardMemberHandler(app.BoardMemberService, app.Validate),
	)
	secureRoutes.Handle(
		app.URLPaths.BoardsHandler,
		handlers.NewBoardHandler(
			app.TaskService,
			app.BoardService,
			app.BoardMemberService,
			app.Validate,
		),
	)
	secureRoutes.Handle(
		app.URLPaths.BoardHandler,
		handlers.NewBoardHandler(
			app.TaskService,
			app.BoardService,
			app.BoardMemberService,
			app.Validate,
		),
	)
	secureRoutes.Handle(
		app.URLPaths.TasksHandler,
		handlers.NewTaskHandler(
			app.TaskService,
			app.Validate,
		),
	)
	secureRoutes.Handle(
		app.URLPaths.TaskHandler,
		handlers.NewTaskHandler(
			app.TaskService,
			app.Validate,
		),
	)
}

func (app *App) initPublicHandlers() {
	publicRoutes := router.NewGroup(app.ServeMux, "")
	publicRoutes.Handle(app.URLPaths.LoginHandler, handlers.NewLoginHandler(app.AuthService, app.Validate))
	publicRoutes.Handle(
		app.URLPaths.RegistrationHandler,
		handlers.NewRegistrationHandler(app.AuthService, app.Validate),
	)
}

func (app *App) initRouter() {
	app.ServeMux = http.NewServeMux()
	app.initPaths()
	app.initPublicHandlers()
	app.initSecureHandlers()
}

func (app *App) runListen() {
	jsonHandler := middlewares.JSONResponse(app.ServeMux)
	logHandler := middlewares.Log(jsonHandler)
	corsHandler := middlewares.CORS(logHandler, map[string][]string{
		app.URLPaths.RegistrationHandler:  app.AllowedHTTPMethods.RegistrationHandler,
		app.URLPaths.LoginHandler:         app.AllowedHTTPMethods.LoginHandler,
		app.URLPaths.LogoutHandler:        app.AllowedHTTPMethods.LogoutHandler,
		app.URLPaths.RefreshAccessHandler: app.AllowedHTTPMethods.RefreshAccessHandler,
		app.URLPaths.UsersHandler:         app.AllowedHTTPMethods.UsersHandler,
		app.URLPaths.UserHandler:          app.AllowedHTTPMethods.UserHandler,
		app.URLPaths.BoardsHandler:        app.AllowedHTTPMethods.BoardsHandler,
		app.URLPaths.BoardHandler:         app.AllowedHTTPMethods.BoardHandler,
		app.URLPaths.BoardMembersHandler:  app.AllowedHTTPMethods.BoardMembersHandler,
		app.URLPaths.BoardMemberHandler:   app.AllowedHTTPMethods.BoardMemberHandler,
	})
	log.Println("Start listening on " + "0.0.0.0:" + app.ServerPort)
	runErr := http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%s", app.ServerPort),
		corsHandler,
	)
	if runErr != nil {
		panic(runErr)
	}
}
