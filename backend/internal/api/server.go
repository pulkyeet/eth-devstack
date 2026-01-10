package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/pulkyeet/eth-devstack/backend/internal/api/handlers"
	"github.com/pulkyeet/eth-devstack/backend/internal/api/middleware"
	"github.com/pulkyeet/eth-devstack/backend/internal/database"
	"github.com/pulkyeet/eth-devstack/backend/internal/responses"
	"go.uber.org/zap"
)

type Server struct {
	app *fiber.App
	db *database.DB
	logger *zap.Logger
	port string
}

func NewServer(db *database.DB, logger *zap.Logger, port string) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return responses.Error(c, code, "INTERNAL_ERROR", err.Error(), nil)
		},
	})

	app.Use(middleware.Recovery(logger))
	app.Use(middleware.Logger(logger))
	app.Use(middleware.CORS())

	blockHandler := handlers.NewBlockHandler(db)
	txHandler := handlers.NewTransactionHandler(db)
	addrHandler := handlers.NewAddressHandler(db)
	chainHandler := handlers.NewChainHandler(db)
	searchHandler := handlers.NewSearchHandler(db)
	streamHandler := handlers.NewStreamHandler(db, logger)
	statsHandler := handlers.NewStatsHandler(db)

	api := app.Group("/api/v1")

	api.Get("/health", chainHandler.GetHealth)
	api.Get("/chains", chainHandler.GetChains)

	api.Get("/blocks", blockHandler.GetBlocks)
	api.Get("/blocks/:id", blockHandler.GetBlock)

	api.Get("/transactions", txHandler.GetTransactions)
	api.Get("/transactions/:hash", txHandler.GetTransaction)

	api.Get("/addresses/:address", addrHandler.GetAddress)
	api.Get("/addresses/:address/:/transactions", addrHandler.GetAddressTransactions)

	api.Get("/search", searchHandler.Search)

	api.Get("/stream/blocks", streamHandler.StreamBlocks)

	api.Get("/stats", statsHandler.GetStats)

	api.Get("/addresses/:address/tokens", addrHandler.GetAddressTokens)

	return &Server{
		app: app,
		db: db,
		logger: logger,
		port: port,
	}
}

func (s *Server) Start() error {
	s.logger.Sugar().Infow("Starting API server", "port", s.port)
	return s.app.Listen(fmt.Sprintf(":%s", s.port))
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}