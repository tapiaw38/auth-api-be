package config

import "os"

const (
	prodScope = "prod"
	testScope = "test"
)

func GetScope() string {
	return os.Getenv("SCOPE")
}

func IsProdScope() bool {
	return GetScope() == prodScope
}

func IsTestScope() bool {
	return GetScope() == testScope
}

func IsLocalScope() bool {
	return GetScope() == ""
}
