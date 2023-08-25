package internal

type Config struct {
	Address            string
	DBConnectionString string
	LogLevel           string
}

func (c *Config) Load() error {
	c.Address = Env("APP_ADDRESS", ":8080")
	c.DBConnectionString = Env("APP_DB_CONNECTION_STRING", "")
	c.LogLevel = Env("APP_LOG_LEVEL", "warn")
	return nil
}
