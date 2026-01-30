package config

type Config struct {
	AppEnv string       `mapstructure:"APP_ENV"`
	Grpc   GrpcConfig   `mapstructure:",squash"`
	Logger LoggerConfig `mapstructure:",squash"`
}

type LoggerConfig struct {
	LogLevel     string `mapstructure:"LOG_LEVEL" validate:"required,oneof=debug info warn error"`
	LogFormatter string `mapstructure:"LOG_FORMATTER" validate:"required,oneof=json console"`
}

type GrpcConfig struct {
	Port int `mapstructure:"GRPC_PORT" validate:"required"`
}
