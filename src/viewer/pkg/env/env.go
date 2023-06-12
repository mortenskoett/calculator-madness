package env

import "os"

func GetEnvVarOrDefault(envName string, def string) string {
	envvar := os.Getenv(envName)
	if len(envvar) == 0 {
		return def
	}
	return envvar
}
