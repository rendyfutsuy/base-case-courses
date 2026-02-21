//	@title			Base Template API Documentation
//	@version		0.0-beta
//	@description	Welcome to the API documentation for the Base Template Web Application. This comprehensive guide is designed to help developers seamlessly integrate and interact with our platform's functionalities. Whether you're building new features, enhancing existing ones, or troubleshooting, this documentation provides all the necessary resources and information.

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description					Enter JWT token (ex: Bearer eyJhbGciOiJIU....)
package router

import (
	"net/http"
	"time"

	// "github.com/go-playground/validator/v10"

	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/redis/go-redis/v9"
	_ "github.com/rendyfutsuybase-case-courses/docs"
	"github.com/rendyfutsuybase-case-courses/utils"
	"github.com/rendyfutsuybase-case-courses/utils/services"
	"github.com/rendyfutsuybase-case-courses/worker"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"

	_homepageController "github.com/rendyfutsuybase-case-courses/modules/homepage/delivery/http"

	// middleware "github.com/rendyfutsuybase-case-courses/helpers/middleware"
	_reqContext "github.com/rendyfutsuybase-case-courses/helpers/middleware/request"
	// "github.com/rendyfutsuybase-case-courses/helpers/validations"

	_authController "github.com/rendyfutsuybase-case-courses/modules/auth/delivery/http"
	_authRepo "github.com/rendyfutsuybase-case-courses/modules/auth/repository"
	_authService "github.com/rendyfutsuybase-case-courses/modules/auth/usecase"

	authmiddleware "github.com/rendyfutsuybase-case-courses/helpers/middleware"
	roleMiddleware "github.com/rendyfutsuybase-case-courses/helpers/middleware"

	_userManagementController "github.com/rendyfutsuybase-case-courses/modules/user_management/delivery/http"
	_userManagementRepo "github.com/rendyfutsuybase-case-courses/modules/user_management/repository"
	_userManagementService "github.com/rendyfutsuybase-case-courses/modules/user_management/usecase"

	_roleManagementController "github.com/rendyfutsuybase-case-courses/modules/role_management/delivery/http"
	_roleManagementRepo "github.com/rendyfutsuybase-case-courses/modules/role_management/repository"
	_roleManagementService "github.com/rendyfutsuybase-case-courses/modules/role_management/usecase"

	_courseController "github.com/rendyfutsuybase-case-courses/modules/course/delivery/http"
	_courseRepo "github.com/rendyfutsuybase-case-courses/modules/course/repository"
	_courseService "github.com/rendyfutsuybase-case-courses/modules/course/usecase"
	_parameterController "github.com/rendyfutsuybase-case-courses/modules/parameter/delivery/http"
	_parameterRepo "github.com/rendyfutsuybase-case-courses/modules/parameter/repository"
	_parameterService "github.com/rendyfutsuybase-case-courses/modules/parameter/usecase"
)

func InitializedRouter(gormDB *gorm.DB, redisClient *redis.Client, timeoutContext time.Duration, v *validator.Validate, nrApp *newrelic.Application) *echo.Echo {
	router := echo.New()

	// queries := sqlc.New(db)

	// Config CORS
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXCSRFToken},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(nrecho.Middleware(nrApp))

	// Config Rate Limiter with configurable limit (default 1000 requests/sec)
	throttleMiddleware := authmiddleware.NewThrottleMiddleware()
	router.Use(throttleMiddleware.Throttle())

	router.GET("/", _homepageController.DefaultHomepage)
	router.GET("/health/storage", _homepageController.StorageHealth)

	// Swagger documentation - hanya tersedia di development environment
	if utils.ConfigVars.String("app_env") == "development" {
		router.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// uses to render files stored on local device.
	router.Static("/storage", "public/storage")
	// Services  ------------------------------------------------------------------------------------------------------------------------------------------------------
	emailServices, err := services.NewEmailService()
	if err != nil {
		panic(err)
	}

	// Initialize the Redis client for Asynq
	redisSetting := asynq.RedisClientOpt{
		Addr:     utils.ConfigVars.String("redis.address"),
		Password: utils.ConfigVars.String("redis.password"),
		DB:       utils.ConfigVars.Int("redis.db"),
	}

	queueClient := asynq.NewClient(redisSetting)

	// Repositories ------------------------------------------------------------------------------------------------------------------------------------------------------
	authRepo := _authRepo.NewAuthRepository(gormDB, emailServices, queueClient)   // Using GORM for auth
	roleManagementRepo := _roleManagementRepo.NewRoleManagementRepository(gormDB) // Using GORM for role_management

	userManagementRepo := _userManagementRepo.NewUserManagementRepository(gormDB) // Using GORM for user_management

	parameterRepo := _parameterRepo.NewParameterRepository(gormDB) // Using GORM for parameter
	courseRepo := _courseRepo.NewCourseRepository(gormDB)          // Using GORM for course

	// Middlewares ------------------------------------------------------------------------------------------------------------------------------------------------------
	middlewareAuth := authmiddleware.NewMiddlewareAuth()
	middlewarePermission := roleMiddleware.NewMiddlewarePermission(
		roleManagementRepo,
	)

	middlewarePageRequest := _reqContext.NewMiddlewarePageRequest()

	// Initialize race condition middleware
	raceConditionMiddleware := authmiddleware.NewRaceConditionMiddleware(redisClient)

	// Example: Add protected routes with race condition prevention
	// This demonstrates how to apply race condition middleware to specific routes
	raceProtectedGroup := router.Group("/v1/protected")
	raceProtectedGroup.Use(middlewareAuth.AuthorizationCheck)
	raceProtectedGroup.Use(raceConditionMiddleware.PreventRaceCondition("protected_operations"))
	raceProtectedGroup.GET("/safe-operation", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "This operation is protected from race conditions"})
	})

	// Auth
	authService := _authService.NewAuthUsecase(
		authRepo,
		roleManagementRepo,
		timeoutContext,
		utils.ConfigVars.String("jwt_key"),
		[]byte(utils.ConfigVars.String("jwt_key")),
		[]byte(utils.ConfigVars.String("jwt_refresh_key")),
	)
	_authController.NewAuthHandler(
		router,
		authService,
		middlewareAuth,
		middlewarePageRequest,
	)

	// role management
	roleManagementService := _roleManagementService.NewRoleManagementUsecase(
		roleManagementRepo,
		authRepo,
		timeoutContext,
	)
	_roleManagementController.NewRoleManagementHandler(
		router,
		roleManagementService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	// user management
	userManagementService := _userManagementService.NewUserManagementUsecase(
		userManagementRepo,
		roleManagementRepo,
		authRepo,
		timeoutContext,
	)
	_userManagementController.NewUserManagementHandler(
		router,
		userManagementService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	// parameter management
	parameterService := _parameterService.NewParameterUsecase(parameterRepo)
	_parameterController.NewParameterHandler(
		router,
		parameterService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	// course management (public index & detail, protected create/update/delete)
	courseService := _courseService.NewCourseUsecase(courseRepo, parameterRepo)
	_courseController.NewCourseHandler(
		router,
		courseService,
		middlewarePageRequest,
		middlewareAuth,
		middlewarePermission,
	)

	usecaseRegistry := worker.UsecaseRegistry{
		// Add any other usecases that your background jobs might need
	}

	dispatcher := worker.NewDispatcher(10, usecaseRegistry) // Using 10 workers, for example
	dispatcher.Run()

	time.Sleep(1000 * time.Millisecond)
	return router
}
