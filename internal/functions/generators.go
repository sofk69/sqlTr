package generator

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateAge() int {
	minAge := 18
	maxAge := 90
	return minAge + rand.Intn(maxAge-minAge)
}

func GenerateID() int {
	minAge := 1
	maxAge := 29
	return minAge + rand.Intn(maxAge-minAge)
}

func GenerateNumber() int {
	minN := 1
	maxN := 999
	return minN + rand.Intn(maxN-minN)
}

func GeneratePaymentName() string {
	names := []string{
		"Membership",
		"Product 1",
		"Product 2",
		"Product 3",
		"1-week trial"}

	randomIndex := rand.Intn(len(names))
	return names[randomIndex]
}

func GenerateEmail() string {
	localPartChars := "abcdefghijklmnopqrstuvwxyz0123456789"
	domainChars := "abcdefghijklmnopqrstuvwxyz"

	localPartLength := rand.Intn(6) + 3
	localPart := make([]byte, localPartLength)
	for i := range localPart {
		localPart[i] = localPartChars[rand.Intn(len(localPartChars))]
	}

	domainLength := rand.Intn(5) + 2
	domain := make([]byte, domainLength)
	for i := range domain {
		domain[i] = domainChars[rand.Intn(len(domainChars))]
	}

	return fmt.Sprintf("%s@%s.com", string(localPart), string(domain))
}

func GenerateName() string {
	consonants := "bcdfghjklmnpqrstvwxyz"
	vowels := "aeiou"
	allLetters := "abcdefghijklmnopqrstuvwxyz"

	nameLength := rand.Intn(6) + 3
	name := make([]byte, nameLength)

	for i := 0; i < nameLength; i++ {
		if i == 0 {
			name[i] = allLetters[rand.Intn(len(allLetters))]
		} else if i%2 == 1 {
			name[i] = consonants[rand.Intn(len(consonants))]
		} else {
			name[i] = vowels[rand.Intn(len(vowels))]
		}

		if rand.Intn(100) < 30 && i > 0 {
			name[i] = allLetters[rand.Intn(len(allLetters))]
		}
	}

	if name[0] >= 'a' && name[0] <= 'z' {
		name[0] = name[0] - 32
	}

	return string(name)
}
