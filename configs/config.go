package configs

import (
	"errors"
	"os"
	"strconv"

	"go.uber.org/multierr"

	"github.com/Bazhenator/tools/src/logger"
	grpcListener "github.com/Bazhenator/tools/src/server/grpc/listener"
)

const (
	EnvBaseSpeed   = "BASE_SPEED"
	EnvTeamsAmount = "TEAMS_AMOUNT"
)

// Config is a main configuration struct for application
type Config struct {
	Environment  string
	Grpc         *grpcListener.GrpcConfig
	LoggerConfig *logger.LoggerConfig

	BaseSpeed   uint64
	TeamsAmount uint64
}

// NewConfig returns application config instance
func NewConfig() (*Config, error) {
	var errorBuilder error

	grpcConfig, err := grpcListener.NewStandardGrpcConfig()
	multierr.AppendInto(&errorBuilder, err)

	loggerConfig, err := logger.NewLoggerConfig()
	multierr.AppendInto(&errorBuilder, err)

	EnvBaseSpeedStr, ok := os.LookupEnv(EnvBaseSpeed)
	if !ok {
		multierr.AppendInto(&errorBuilder, errors.New("BASE_SPEED is not defined"))
	}

	baseSpeed, err := strconv.Atoi(EnvBaseSpeedStr)
	multierr.AppendInto(&errorBuilder, err)

	EnvTeamsAmountStr, ok := os.LookupEnv(EnvTeamsAmount)
	if !ok {
		multierr.AppendInto(&errorBuilder, errors.New("TEAMS_AMOUNT is not defined"))
	}

	teamsAmount, err := strconv.Atoi(EnvTeamsAmountStr)
	multierr.AppendInto(&errorBuilder, err)

	if errorBuilder != nil {
		return nil, errorBuilder
	}

	glCfg := &Config{
		Grpc:         grpcConfig,
		LoggerConfig: loggerConfig,

		BaseSpeed:   uint64(baseSpeed),
		TeamsAmount: uint64(teamsAmount),
	}

	return glCfg, nil
}
