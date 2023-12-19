package middleware

import (
	"github.com/Yangiboev/todo/config"
	"github.com/Yangiboev/todo/pkg/logger"
)

// Middleware manager
type MiddlewareManager struct {
	cfg     *config.Config
	origins []string
	logger  logger.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, origins: origins, logger: logger}
}
