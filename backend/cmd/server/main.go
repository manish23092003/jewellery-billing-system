package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"

	"jewellery-billing/database"
	"jewellery-billing/database/seed"
	"jewellery-billing/internal/config"
	"jewellery-billing/internal/handler"
	"jewellery-billing/internal/middleware"
	"jewellery-billing/internal/repository"
	"jewellery-billing/internal/service"
	"jewellery-billing/pkg/logger"
)

func main() {
	// ── 1. Logger ──────────────────────────────────────────────────────
	logger.Init()
	log.Info().Msg("Starting Jewellery Billing API…")

	// ── 2. Configuration ───────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// ── 3. Database ────────────────────────────────────────────────────
	pool, err := database.NewPostgresPool(cfg.DatabaseURL())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer pool.Close()

	// ── 4. Migrations ──────────────────────────────────────────────────
	if err := database.RunMigrations(cfg.DatabaseURL()); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// ── 5. Seed Data ───────────────────────────────────────────────────
	if err := seed.SeedAdminUser(pool); err != nil {
		log.Warn().Err(err).Msg("Failed to seed admin user")
	}

	// ── 6. Dependency Injection ────────────────────────────────────────
	// Repositories
	userRepo := repository.NewUserRepository(pool)
	metalRateRepo := repository.NewMetalRateRepository(pool)
	billRepo := repository.NewBillRepository(pool)
	settingRepo := repository.NewSettingRepository(pool)
	expenseRepo := repository.NewExpenseRepository(pool)
	analyticsRepo := repository.NewAnalyticsRepository(pool)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	metalRateService := service.NewMetalRateService(metalRateRepo)
	billService := service.NewBillService(billRepo)
	settingService := service.NewSettingService(settingRepo)
	expenseService := service.NewExpenseService(expenseRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	pdfService := service.NewPDFService()

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	metalRateHandler := handler.NewMetalRateHandler(metalRateService)
	billHandler := handler.NewBillHandler(billService, settingService, pdfService)
	settingHandler := handler.NewSettingHandler(settingService)
	expenseHandler := handler.NewExpenseHandler(expenseService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// Middlewares
	adminMid := middleware.AdminOnly()

	// ── 7. Fiber App ───────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName:      "Jewellery Billing API v1.0",
		ErrorHandler: globalErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "${time} │ ${status} │ ${latency} │ ${method} ${path}\n",
		TimeFormat: "15:04:05",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// ── 8. Routes ──────────────────────────────────────────────────────
	
	// Serve uploaded files
	app.Static("/uploads", "./uploads")

	// API Routes
	api := app.Group("/api")

	// Health check (public)
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "jewellery-billing-api",
		})
	})

	// Auth routes (public)
	auth := app.Group("/api/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes — require valid JWT
	protected := app.Group("/api", middleware.AuthRequired(authService))

	// Auth (protected)
	protected.Post("/auth/logout", authHandler.Logout)
	protected.Get("/auth/me", authHandler.Me)

	// User management (admin only)
	users := protected.Group("/users", middleware.AdminOnly())
	users.Get("/", userHandler.GetAll)


	protected.Get("/users", adminMid, userHandler.GetAll)
	protected.Post("/users", adminMid, userHandler.Create)
	protected.Get("/users/:id", adminMid, userHandler.GetByID)
	protected.Put("/users/:id", adminMid, userHandler.Update)
	protected.Delete("/users/:id", adminMid, userHandler.Delete)

	// Metal Rates
	protected.Get("/metal-rates/latest", metalRateHandler.GetCurrentRates)
	protected.Post("/metal-rates", adminMid, metalRateHandler.Create)
	protected.Get("/metal-rates", metalRateHandler.GetHistory)
	protected.Delete("/metal-rates/:id", adminMid, metalRateHandler.Delete)

	// Bills
	protected.Post("/bills", billHandler.Create)
	protected.Get("/bills", billHandler.GetAll)
	protected.Get("/bills/:id", billHandler.GetByID)
	protected.Get("/bills/:id/pdf", billHandler.DownloadPDF)
	protected.Delete("/bills/:id", adminMid, billHandler.Delete)

	// Settings
	protected.Get("/settings", settingHandler.Get)
	protected.Put("/settings", adminMid, settingHandler.Update)
	protected.Post("/settings/logo", adminMid, settingHandler.UploadLogo)

	// Expenses
	protected.Post("/expenses", expenseHandler.Create)
	protected.Get("/expenses", expenseHandler.GetAll)
	protected.Delete("/expenses/:id", adminMid, expenseHandler.Delete)

	// Analytics
	protected.Get("/analytics/dashboard", analyticsHandler.GetDashboard)

	// ── 9. Graceful Shutdown ───────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info().Msg("Received shutdown signal — draining connections…")
		if err := app.Shutdown(); err != nil {
			log.Error().Err(err).Msg("Server forced to shutdown")
		}
	}()

	// ── 10. Start Server ───────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Info().Str("address", addr).Msg("🚀 Jewellery Billing API is live")
	if err := app.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("Server crashed")
	}
}

// globalErrorHandler catches unhandled panics and Fiber errors.
func globalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	log.Error().Err(err).Int("status", code).Str("path", c.Path()).Msg("Unhandled error")

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}
