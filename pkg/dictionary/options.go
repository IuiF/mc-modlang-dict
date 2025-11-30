package dictionary

// Config holds the client configuration.
type Config struct {
	TargetLang         string
	IncludeGlobalTerms bool
	CacheEnabled       bool
}

// Option is a function that configures the client.
type Option func(*Config)

// defaultConfig returns the default configuration.
func defaultConfig() *Config {
	return &Config{
		TargetLang:         "ja_jp",
		IncludeGlobalTerms: true,
		CacheEnabled:       true,
	}
}

// WithTargetLang sets the default target language.
func WithTargetLang(lang string) Option {
	return func(c *Config) {
		c.TargetLang = lang
	}
}

// WithGlobalTerms enables or disables global terms inclusion.
func WithGlobalTerms(include bool) Option {
	return func(c *Config) {
		c.IncludeGlobalTerms = include
	}
}

// WithCache enables or disables caching.
func WithCache(enabled bool) Option {
	return func(c *Config) {
		c.CacheEnabled = enabled
	}
}
