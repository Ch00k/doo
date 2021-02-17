package main

import "os"

func GetEnvVar(name, default_ string) string {
	value, exists := os.LookupEnv(name)
	if exists {
		return value
	} else {
		return default_
	}
}
