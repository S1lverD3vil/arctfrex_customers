package common

import (
	"log"
	"math"
	"strings"
	"time"

	"math/rand"

	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
)

// const (
// 	lowerCharset   = "abcdefghijklmnopqrstuvwxyz"
// 	upperCharset   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	numberCharset  = "0123456789"
// 	specialCharset = "!@#$%^&*()-_+=<>?~"
// 	allCharset     = lowerCharset + upperCharset + numberCharset + specialCharset
// 	passwordLength = 12 // Minimum recommended length
// )

func UUIDNormalizer(uuid uuid.UUID) string {
	return strings.ReplaceAll(uuid.String(), DASH, STRING_EMPTY)
}
func UUIDNormalizer2(uuid uuid.UUID) string {
	return strings.ReplaceAll(uuid.String(), DASH, STRING_EMPTY)
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// Remove removes the first occurrence of a value from a slice of any type.
func Remove[T comparable](slice []T, value T) []T {
	newSlice := []T{}
	for _, v := range slice {
		if v != value {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

func GenerateRandomNumber(min, max int) int {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between min and max
	number := rand.Intn(max-min+1) + min

	// Randomly flip the sign
	if rand.Intn(2) == 0 {
		number = -number
	}

	return number
}

func RoundTo4DecimalPlaces(value float64) float64 {
	return math.Round(value*10000) / 10000
}

func GenerateShortKSUID() string {
	return ksuid.New().String()[:8] // Shorten to 8 characters
}

// GenerateSecurePassword generates a secure password matching ISO 27001 complexity requirements.
func GenerateSecurePassword() (string, error) {
	// Ensure the password meets complexity requirements
	password := make([]byte, passwordLength)

	// Add at least one character of each required type
	password[0] = lowerCharset[randomIndex(len(lowerCharset))]
	password[1] = upperCharset[randomIndex(len(upperCharset))]
	password[2] = numberCharset[randomIndex(len(numberCharset))]
	password[3] = specialCharset[randomIndex(len(specialCharset))]

	// Fill the rest with random characters from the full charset
	for i := 4; i < passwordLength; i++ {
		password[i] = allCharset[randomIndex(len(allCharset))]
	}

	// Shuffle the password to avoid predictable positions
	shuffledPassword := shuffle(password)

	return string(shuffledPassword), nil
}

// randomIndex generates a cryptographically secure random index.
func randomIndex(max int) int {
	randomByte := make([]byte, 1)
	_, err := rand.Read(randomByte)
	if err != nil {
		log.Fatal("Failed to generate random index:", err)
	}
	return int(randomByte[0]) % max
}

// shuffle shuffles a slice of bytes randomly.
func shuffle(input []byte) []byte {
	shuffled := make([]byte, len(input))
	copy(shuffled, input)
	for i := range shuffled {
		j := randomIndex(len(shuffled))
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}
