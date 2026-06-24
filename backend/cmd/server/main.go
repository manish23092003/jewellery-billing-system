package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"

	"jewellery-billing/database"
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

	// ── 5. Dependency Injection ────────────────────────────────────────
	// Repositories
	orgRepo := repository.NewOrganizationRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool)
	metalRateRepo := repository.NewMetalRateRepository(pool)
	billRepo := repository.NewBillRepository(pool)
	settingRepo := repository.NewSettingRepository(pool)
	expenseRepo := repository.NewExpenseRepository(pool)
	analyticsRepo := repository.NewAnalyticsRepository(pool)
	customerRepo := repository.NewPostgresCustomerRepository(pool)
	verificationLogRepo := repository.NewVerificationLogRepository(pool)
	_ = repository.NewAuditRepository(pool) // Available for future use

	// Email Sender — use SMTP if configured, otherwise console logger
	var emailSender service.EmailSender
	if cfg.IsSMTPConfigured() {
		emailSender = service.NewSMTPEmailSender(cfg)
		log.Info().Msg("✓ Using SMTP email sender")
	} else {
		emailSender = service.NewConsoleEmailSender(cfg.AppURL)
		log.Info().Msg("✓ Using console email sender (development mode)")
	}

	// Services
	authService := service.NewAuthService(userRepo, orgRepo, cfg)
	registrationService := service.NewRegistrationService(orgRepo, userRepo, tokenRepo, settingRepo, authService, emailSender)
	passwordResetService := service.NewPasswordResetService(userRepo, tokenRepo, emailSender)
	userService := service.NewUserService(userRepo)
	metalRateService := service.NewMetalRateService(metalRateRepo)
	billService := service.NewBillService(billRepo, customerRepo)
	settingService := service.NewSettingService(settingRepo)
	expenseService := service.NewExpenseService(expenseRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	customerService := service.NewCustomerService(customerRepo)
	pdfService := service.NewPDFService()

	// Handlers
	authHandler := handler.NewAuthHandler(authService, orgRepo)
	registrationHandler := handler.NewRegistrationHandler(registrationService)
	passwordResetHandler := handler.NewPasswordResetHandler(passwordResetService)
	emailVerificationHandler := handler.NewEmailVerificationHandler(authService, tokenRepo, emailSender, userRepo)
	userHandler := handler.NewUserHandler(userService)
	metalRateHandler := handler.NewMetalRateHandler(metalRateService)
	billHandler := handler.NewBillHandler(billService, settingService, pdfService, verificationLogRepo)
	settingHandler := handler.NewSettingHandler(settingService)
	expenseHandler := handler.NewExpenseHandler(expenseService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	customerHandler := handler.NewCustomerHandler(customerService)

	// ── 6. Fiber App ───────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName:      "Jewellery Billing API v2.0",
		ErrorHandler: globalErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "${time} │ ${status} │ ${latency} │ ${method} ${path}\n",
		TimeFormat: "15:04:05",
	}))
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	// ── 7. Routes ──────────────────────────────────────────────────────

	// Serve uploaded files
	app.Static("/uploads", "./uploads")

	// API Routes
	api := app.Group("/api")

	// Health check (public)
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "jewellery-billing-api",
			"version": "2.0.0",
		})
	})

	// Public routes (rate-limited)
	public := api.Group("/public")
	public.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
	}))
	public.Get("/verify/:token", billHandler.VerifyInvoice)

	// Auth routes (public, rate-limited)
	auth := api.Group("/auth")
	auth.Post("/login", middleware.LoginRateLimiter(), authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/register", middleware.RegistrationRateLimiter(), registrationHandler.Register)
	auth.Post("/forgot-password", middleware.PasswordResetRateLimiter(), passwordResetHandler.ForgotPassword)
	auth.Post("/reset-password", middleware.PasswordResetRateLimiter(), passwordResetHandler.ResetPassword)
	auth.Post("/verify-email", emailVerificationHandler.VerifyEmail)

	// Protected routes — require valid JWT
	protected := api.Group("/", middleware.AuthRequired(authService))

	// Auth (protected)
	protected.Post("auth/logout", authHandler.Logout)
	protected.Get("auth/me", authHandler.Me)
	protected.Post("auth/resend-verification", emailVerificationHandler.ResendVerification)

	// User management (admin only)
	adminMid := middleware.AdminOnly()
	protected.Get("users", adminMid, userHandler.GetAll)
	protected.Post("users", adminMid, userHandler.Create)
	protected.Get("users/:id", adminMid, userHandler.GetByID)
	protected.Put("users/:id", adminMid, userHandler.Update)
	protected.Delete("users/:id", adminMid, userHandler.Delete)

	// Metal Rates
	protected.Get("metal-rates/latest", metalRateHandler.GetCurrentRates)
	protected.Post("metal-rates", adminMid, metalRateHandler.Create)
	protected.Get("metal-rates", metalRateHandler.GetHistory)
	protected.Delete("metal-rates/:id", adminMid, metalRateHandler.Delete)

	// Bills
	protected.Post("bills", billHandler.Create)
	protected.Get("bills", billHandler.GetAll)
	protected.Get("bills/:id", billHandler.GetByID)
	protected.Get("bills/:id/pdf", billHandler.DownloadPDF)
	protected.Post("bills/:id/payment", billHandler.AddPayment)
	protected.Delete("bills/:id", adminMid, billHandler.Delete)

	// Settings (admin only for mutations)
	protected.Get("settings", settingHandler.Get)
	protected.Put("settings", adminMid, settingHandler.Update)
	protected.Post("settings/logo", adminMid, settingHandler.UploadLogo)

	// Expenses
	protected.Post("expenses", expenseHandler.Create)
	protected.Get("expenses", expenseHandler.GetAll)
	protected.Delete("expenses/:id", adminMid, expenseHandler.Delete)

	// Analytics
	protected.Get("analytics/dashboard", analyticsHandler.GetDashboard)

	// Customers
	protected.Post("customers", customerHandler.Create)
	protected.Get("customers", customerHandler.GetAll)
	protected.Get("customers/:id", customerHandler.GetByID)
	protected.Put("customers/:id", customerHandler.Update)
	protected.Delete("customers/:id", adminMid, customerHandler.Delete)

	// ── 8. Graceful Shutdown ───────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info().Msg("Received shutdown signal — draining connections…")
		if err := app.Shutdown(); err != nil {
			log.Error().Err(err).Msg("Server forced to shutdown")
		}
	}()

	// ── 9. Start Server ────────────────────────────────────────────────
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
