package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bmt-saas/api/internal/config"
	"github.com/bmt-saas/api/internal/handler/auth"
	"github.com/bmt-saas/api/internal/handler/developer"
	"github.com/bmt-saas/api/internal/handler/ecommerce"
	"github.com/bmt-saas/api/internal/handler/finance"
	handlerform "github.com/bmt-saas/api/internal/handler/form"
	"github.com/bmt-saas/api/internal/handler/merchant"
	"github.com/bmt-saas/api/internal/handler/nasabah"
	"github.com/bmt-saas/api/internal/handler/nfc"
	"github.com/bmt-saas/api/internal/handler/platform"
	"github.com/bmt-saas/api/internal/handler/pondok"
	"github.com/bmt-saas/api/internal/handler/teller"
	"github.com/bmt-saas/api/internal/middleware"
	"github.com/bmt-saas/api/internal/repository/postgres"
	"github.com/bmt-saas/api/internal/service"
	"github.com/bmt-saas/api/internal/worker"
	"github.com/bmt-saas/api/pkg/jwt"
	"github.com/bmt-saas/api/pkg/settings"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// ── Logger ────────────────────────────────────────────────────────────────
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if os.Getenv("APP_ENV") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// ── Config ────────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("gagal load config")
	}

	log.Info().
		Str("app", cfg.App.Name).
		Str("env", cfg.App.Env).
		Str("version", cfg.App.Version).
		Msg("Starting BMT SaaS API")

	// ── PostgreSQL ────────────────────────────────────────────────────────────
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("gagal parse database URL")
	}
	poolCfg.MaxConns = int32(cfg.Database.MaxOpenConns)
	poolCfg.MinConns = int32(cfg.Database.MaxIdleConns)
	poolCfg.MaxConnLifetime = cfg.Database.ConnMaxLifetime

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("gagal connect ke PostgreSQL")
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("database ping gagal")
	}
	log.Info().Msg("PostgreSQL connected")

	// ── Redis ─────────────────────────────────────────────────────────────────
	redisOpts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("gagal parse Redis URL")
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Warn().Err(err).Msg("Redis ping gagal — lanjut tanpa cache")
	} else {
		log.Info().Msg("Redis connected")
	}

	// ── Repositories ──────────────────────────────────────────────────────────
	settingsRepo := postgres.NewSettingsRepository(dbPool)
	platformRepo := postgres.NewPlatformRepository(dbPool)
	nasabahRepo := postgres.NewNasabahRepository(dbPool)
	rekeningRepo := postgres.NewRekeningRepository(dbPool)
	autodebetRepo := postgres.NewAutodebetRepository(dbPool)
	formRepo := postgres.NewFormRepository(dbPool)
	pembiayaanRepo := postgres.NewPembiayaanRepository(dbPool)
	sesiTellerRepo := postgres.NewSesiTellerRepository(dbPool)
	akuntansiRepo := postgres.NewAkuntansiRepository(dbPool)

	// Pondok repositories
	santriRepo := postgres.NewSantriRepository(dbPool)
	kelasRepo := postgres.NewKelasRepository(dbPool)
	pengajarRepo := postgres.NewPengajarRepository(dbPool)
	jenisTagihanRepo := postgres.NewJenisTagihanRepository(dbPool)
	tagihanSPPRepo := postgres.NewTagihanSPPRepository(dbPool)

	// Ecommerce repositories
	tokoRepo := postgres.NewTokoRepository(dbPool)
	produkRepo := postgres.NewProdukRepository(dbPool)
	pesananRepo := postgres.NewPesananRepository(dbPool)

	// Suppress unused variable warnings for repos not yet used by handlers
	_ = formRepo
	_ = pembiayaanRepo
	_ = santriRepo
	_ = kelasRepo
	_ = pengajarRepo
	_ = jenisTagihanRepo
	_ = tagihanSPPRepo
	_ = tokoRepo
	_ = produkRepo
	_ = pesananRepo

	auditRepo := postgres.NewAuditRepository(dbPool)
	penggunaRepo := postgres.NewPenggunaRepository(dbPool)
	keamananRepo := postgres.NewKeamananRepository(dbPool)

	// ── Settings Resolver ─────────────────────────────────────────────────────
	settingsResolver := settings.NewResolver(settingsRepo)

	// ── JWT Manager ───────────────────────────────────────────────────────────
	jwtManager := jwt.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

	// ── Services ──────────────────────────────────────────────────────────────
	akuntansiService := service.NewAkuntansiService(akuntansiRepo)
	rekeningService := service.NewRekeningService(
		rekeningRepo,
		autodebetRepo,
		settingsResolver,
		akuntansiService,
	)
	autodebetService := service.NewAutodebetService(autodebetRepo, rekeningService)
	nasabahService := service.NewNasabahService(nasabahRepo, rekeningRepo)
	sesiTellerService := service.NewSesiTellerService(sesiTellerRepo, settingsResolver)
	featureChecker := service.NewPlatformFeatureChecker(platformRepo)
	sessionService := service.NewSessionService(keamananRepo, jwtManager, settingsResolver)
	otpService := service.NewOTPService(keamananRepo, redisClient, settingsResolver, nil)
	authService := service.NewAuthService(penggunaRepo, nasabahRepo, sessionService, otpService, settingsResolver)
	formService := service.NewFormService(formRepo, nasabahRepo, rekeningRepo, settingsResolver)
	_ = service.NewSettingsService(settingsResolver)
	_ = featureChecker

	// ── Workers ───────────────────────────────────────────────────────────────
	asynqRedisOpt := asynq.RedisClientOpt{Addr: redisOpts.Addr, Password: redisOpts.Password}

	autodebetWorker := worker.NewAutodebetWorker(autodebetService)

	asynqServer := asynq.NewServer(
		asynqRedisOpt,
		asynq.Config{
			Concurrency: cfg.Asynq.Concurrency,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	worker.RegisterWorkers(mux, autodebetWorker)

	go func() {
		log.Info().Msg("Starting asynq worker server")
		if err := asynqServer.Run(mux); err != nil {
			log.Error().Err(err).Msg("asynq server error")
		}
	}()

	// Asynq scheduler untuk periodic tasks
	asynqScheduler := asynq.NewScheduler(asynqRedisOpt, nil)
	worker.SchedulePeriodicTasks(asynqScheduler)

	go func() {
		log.Info().Msg("Starting asynq scheduler")
		if err := asynqScheduler.Run(); err != nil {
			log.Error().Err(err).Msg("asynq scheduler error")
		}
	}()

	// ── Router ────────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Idempotency-Key", "Developer-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Idempotency)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","version":"%s"}`, cfg.App.Version)
	})

	// Developer routes (no JWT, Developer-Token only)
	r.Route("/dev", func(r chi.Router) {
		r.Use(middleware.DeveloperAuth(os.Getenv("DEVELOPER_TOKEN")))
		developer.RegisterRoutes(r)
	})

	// Auth routes
	authHandler := auth.NewHandler(authService, sessionService, jwtManager)
	r.Route("/auth", func(r chi.Router) {
		authHandler.RegisterRoutes(r)
	})

	// Platform routes (BMT management)
	r.Route("/platform", func(r chi.Router) {
		r.Use(middleware.Auth(jwtManager))
		r.Use(middleware.TenantRequired)
		platform.RegisterRoutes(r)
	})

	// Teller routes
	tellerHandler := teller.NewHandler(sesiTellerService, rekeningService, nasabahService)
	r.Route("/teller", func(r chi.Router) {
		r.Use(middleware.Auth(jwtManager))
		r.Use(middleware.TenantRequired)
		r.Use(middleware.RequireRole("TELLER", "MANAJER_CABANG"))
		r.Use(middleware.AuditLog(auditRepo))
		tellerHandler.RegisterRoutes(r)
	})

	// Nasabah routes
	nasabahHandler := nasabah.NewHandler(nasabahService, rekeningService)
	r.Route("/nasabah", func(r chi.Router) {
		r.Use(middleware.Auth(jwtManager))
		nasabahHandler.RegisterRoutes(r)
	})

	// Management API routes
	formHandler := handlerform.NewHandler(formService)
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Auth(jwtManager))
		r.Use(middleware.TenantRequired)
		r.Use(middleware.AuditLog(auditRepo))
		// Form workflow routes
		r.Route("/form", func(r chi.Router) {
			formHandler.RegisterRoutes(r)
		})
		// Finance sub-routes
		r.Route("/finance", func(r chi.Router) {
			r.Use(middleware.RequireRole("FINANCE", "MANAJER_CABANG", "MANAJER_BMT"))
			finance.RegisterRoutes(r)
		})
	})

	// Finance routes (top-level for FINANCE role)
	r.Route("/finance", func(r chi.Router) {
		r.Use(middleware.Auth(jwtManager))
		r.Use(middleware.TenantRequired)
		r.Use(middleware.RequireRole("FINANCE", "MANAJER_CABANG", "MANAJER_BMT", "AUDITOR_BMT"))
		finance.RegisterRoutes(r)
	})

	// Pondok routes
	r.Route("/pondok", func(r chi.Router) {
		r.Use(middleware.Auth(jwtManager))
		r.Use(middleware.TenantRequired)
		r.Use(middleware.RequireRole("ADMIN_PONDOK", "OPERATOR_PONDOK", "BENDAHARA_PONDOK",
			"PETUGAS_UKS", "PUSTAKAWAN", "PETUGAS_PPDB", "BK"))
		pondok.RegisterRoutes(r)
	})

	// E-commerce shop routes (seller/pondok side)
	r.Route("/shop", func(r chi.Router) {
		r.Use(middleware.Auth(jwtManager))
		r.Use(middleware.TenantRequired)
		ecommerce.RegisterRoutes(r)
	})

	// OPOP marketplace (lintas pondok, publik read)
	r.Route("/opop", func(r chi.Router) {
		ecommerce.RegisterOPOPRoutes(r)
	})

	// NFC routes
	r.Route("/nfc", func(r chi.Router) {
		r.Use(middleware.Idempotency)
		nfc.RegisterRoutes(r)
	})

	// Merchant routes
	r.Route("/merchant", func(r chi.Router) {
		r.Use(middleware.Auth(jwtManager))
		r.Use(middleware.RequireRole("KASIR_MERCHANT", "OWNER_MERCHANT"))
		merchant.RegisterRoutes(r)
	})

	// Webhook routes
	r.Route("/webhook", func(r chi.Router) {
		// TODO: register Midtrans webhook
	})

	// ── HTTP Server ───────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	// ── Graceful Shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("port", cfg.HTTP.Port).Msg("Server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server error")
		}
	}()

	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	// Shutdown asynq
	asynqServer.Shutdown()
	asynqScheduler.Shutdown()

	log.Info().Msg("Server stopped gracefully")
}
