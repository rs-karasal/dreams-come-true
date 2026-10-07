package config

func EnvOrDefault(key, def string) string {
	return envOrDefault(key, def)
}
