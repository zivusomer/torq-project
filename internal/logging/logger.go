package logging

import (
	"log/slog"
	"os"
)

var (
	Logger = &Service{
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	}
)

type Service struct {
	logger *slog.Logger
}

func (s *Service) Info(message string) {
	s.logger.Info(message)
}

func (s *Service) Warn(message string) {
	s.logger.Warn(message)
}

func (s *Service) Error(message string) {
	s.logger.Error(message)
}
