package logging

import (
	"log/slog"
	"os"
	"sync"
)

var (
	Logger = &Service{
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	}
)

type Service struct {
	mu     sync.RWMutex
	logger *slog.Logger
}

func (s *Service) Info(message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.logger.Info(message)
}

func (s *Service) Warn(message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.logger.Warn(message)
}

func (s *Service) Error(message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.logger.Error(message)
}
