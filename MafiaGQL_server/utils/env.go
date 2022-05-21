package utils

import "os"

func SetEnvVar(name, defaultVal string) string {
	val := os.Getenv(name)

	if val == "" {
		err := os.Setenv(name, defaultVal)
		PanicOnError("Failed to set env", err)
		val = os.Getenv(name)
	}

	return val
}
