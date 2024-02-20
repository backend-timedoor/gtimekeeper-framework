package config

import "github.com/spf13/cast"

// Env Get config from env.
func (c *Config) Env(envName string, defaultValue ...any) any {
	value := c.Get(envName, defaultValue...)
	if cast.ToString(value) == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return nil
	}

	return value
}

// Add config to application.
func (c *Config) Add(name string, configuration any) {
	c.vp.Set(name, configuration)
}

// Get config from application.
func (c *Config) Get(path string, defaultValue ...any) any {
	if !c.vp.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	return c.vp.Get(path)
}

// GetString Get string type config from application.
func (c *Config) GetString(path string, defaultValue ...any) string {
	value := cast.ToString(c.Get(path, defaultValue...))
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0].(string)
		}

		return ""
	}

	return value
}

// GetInt Get int type config from application.
func (c *Config) GetInt(path string, defaultValue ...any) int {
	value := c.Get(path, defaultValue...)
	if cast.ToString(value) == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0].(int)
		}

		return 0
	}

	return cast.ToInt(value)
}

// GetBool Get bool type config from application.
func (c *Config) GetBool(path string, defaultValue ...any) bool {
	value := c.Get(path, defaultValue...)
	if cast.ToString(value) == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0].(bool)
		}

		return false
	}

	return cast.ToBool(value)
}
