package utils

import (
	"log"
)

func LogInfo(msg string, args ...interface{}) {
	log.Printf("✅ INFO: "+msg, args...)
}

func LogError(msg string, args ...interface{}) {
	log.Printf("❌ ERROR: "+msg, args...)
}

func LogWarn(msg string, args ...interface{}) {
	log.Printf("⚠️ WARN: "+msg, args...)
}
