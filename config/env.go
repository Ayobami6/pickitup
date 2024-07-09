package config

import "syscall"

func GetEnv(envStr string, fallback string) string {
	if val, ok := syscall.Getenv(envStr); ok {
		return val
	}
	return fallback
}