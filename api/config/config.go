package config

type Config struct {
	Environment    *ENV
	OAuthProviders *Providers
}

func Load() *Config {

	env, err := LoadEnvironment()
	if err != nil {
		panic(err)
	}

	oac, err := LoadOAuthProviders()
	if err != nil {
		panic(err)
	}

	return &Config{
		Environment:    env,
		OAuthProviders: oac,
	}
}

func IsProduction(config *Config) bool {
	if config.Environment == nil {
		return false
	}
	if config.Environment.ENV == "production" {
		return true
	}
	return false
}
