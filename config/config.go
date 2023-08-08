package config

// Config is the database configuration
type Config struct {
	Host      string
	Name      string
	User      string
	Password  string
	Deletable bool
}

// New prepares the configuration for mxorm
func New(host string, name, user, password string, allowDelete bool) *Config {
	return &Config{
		Host:      host,
		Name:      name,
		User:      user,
		Password:  password,
		Deletable: allowDelete,
	}
}
