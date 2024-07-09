package application

type Config struct {
	Host     string
	Supabase SupabaseConfig
	Database DatabaseConfig
}

type GetEnvFunc func(string) string

func NewConfig(getenv GetEnvFunc) *Config {
	return &Config{
		Host: getenv("HOST"),
		Supabase: SupabaseConfig{
			URL:    getenv("SUPABASE_URL"),
			Secret: getenv("SUPABASE_SECRET"),
		},
		Database: DatabaseConfig{
			Name:     getenv("DB_NAME"),
			Host:     getenv("DB_HOST"),
			Port:     getenv("DB_PORT"),
			Username: getenv("DB_USER"),
			Password: getenv("DB_PASSWORD"),
		},
	}
}

type DatabaseConfig struct {
	Name     string
	Host     string
	Port     string
	Username string
	Password string
}

type SupabaseConfig struct {
	URL    string
	Secret string
}
