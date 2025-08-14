package passwordmanagerv2_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func checkAttributeMinLength(resourceName, attribute string, minLength int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		attrKey := fmt.Sprintf("%s.#", attribute)
		attr, ok := rs.Primary.Attributes[attrKey]
		if !ok {
			return fmt.Errorf("attribute %s not found", attrKey)
		}

		count, err := strconv.Atoi(attr)
		if err != nil {
			return fmt.Errorf("error converting %s to int: %s", attrKey, err)
		}

		if count < minLength {
			return fmt.Errorf("expected %s to be >= %d, got %d", attrKey, minLength, count)
		}

		return nil
	}
}

func checkAttributeLengthEqual(resourceName, attribute string, expectedLength int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		attrKey := fmt.Sprintf("%s.#", attribute)
		attr, ok := rs.Primary.Attributes[attrKey]
		if !ok {
			return fmt.Errorf("attribute %s not found", attrKey)
		}

		count, err := strconv.Atoi(attr)
		if err != nil {
			return fmt.Errorf("error converting %s to int: %s", attrKey, err)
		}

		if count != expectedLength {
			return fmt.Errorf("expected %s to be %d, got %d", attrKey, expectedLength, count)
		}

		return nil
	}
}

// Character sets
var (
	lowerLetters = []rune("abcdefghijklmnopqrstuvwxyz")
	upperLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	digits       = []rune("0123456789")
	specials     = []rune("@#$%^&*()-_+=!~")
)

// allChars is the union of all allowed characters.
var allChars = append(append(append(lowerLetters, upperLetters...), digits...), specials...)

// secureRandomRune returns a cryptographically secure random rune from a given set.
func secureRandomRune(set []rune) (rune, error) {
	idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
	if err != nil {
		return 0, err
	}
	return set[idx.Int64()], nil
}

// shuffleRunes shuffles a slice of runes in place.
func shuffleRunes(runes []rune) error {
	for i := len(runes) - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return err
		}
		j := int(jBig.Int64())
		runes[i], runes[j] = runes[j], runes[i]
	}
	return nil
}

// hasConsecutiveDuplicates checks for more than two identical characters in a row.
func hasConsecutiveDuplicates(p []rune) bool {
	count := 1
	for i := 1; i < len(p); i++ {
		if p[i] == p[i-1] {
			count++
			if count > 2 {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}

// meetsCharClassRequirements checks that p contains at least one rune from each class.
func meetsCharClassRequirements(p []rune) bool {
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, c := range p {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsDigit(c):
			hasDigit = true
		case strings.ContainsRune(string(specials), c):
			hasSpecial = true
		}
	}
	return hasLower && hasUpper && hasDigit && hasSpecial
}

// hammingDistance computes the number of differing positions between two strings.
func hammingDistance(a, b string) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	diff := abs(len(a) - len(b))
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			diff++
		}
	}
	return diff
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// inLastN checks if pwd is exactly equal to any of the last N passwords.
func inLastN(pwd string, history []string, n int) bool {
	limit := n
	if len(history) < n {
		limit = len(history)
	}
	for i := 0; i < limit; i++ {
		if history[i] == pwd {
			return true
		}
	}
	return false
}

// isDictionaryWord rejects passwords that exactly match a dictionary word.
func isDictionaryWord(pwd string, dict map[string]struct{}) bool {
	low := strings.ToLower(pwd)
	if _, found := dict[low]; found {
		return true
	}
	return false
}

// isSequential checks for runs of 4 increasing or decreasing characters.
func isSequential(p []rune) bool {
	if len(p) < 4 {
		return false
	}
	for i := 0; i <= len(p)-4; i++ {
		inc, dec := true, true
		for j := 1; j < 4; j++ {
			if p[i+j] != p[i+j-1]+1 {
				inc = false
			}
			if p[i+j] != p[i+j-1]-1 {
				dec = false
			}
		}
		if inc || dec {
			return true
		}
	}
	return false
}

// GeneratePassword creates a password meeting all requirements.
// previous list holds most recent passwords (index 0 = most recent).
func GeneratePassword(prev []string) (string, error) {
	const (
		minLen       = 8
		maxLen       = 16
		diffRequired = 4
		historyCount = 5
	)
	// Attempt generation up to 200 times
	for attempt := 0; attempt < 200; attempt++ {
		lengthBig, err := rand.Int(rand.Reader, big.NewInt(maxLen-minLen+1))
		if err != nil {
			return "", err
		}
		length := int(lengthBig.Int64()) + minLen

		// Build initial slice ensuring one of each class
		password := make([]rune, 0, length)
		sets := [][]rune{lowerLetters, upperLetters, digits, specials}
		for _, set := range sets {
			r, err := secureRandomRune(set)
			if err != nil {
				return "", err
			}
			password = append(password, r)
		}

		// Fill remaining runes
		for len(password) < length {
			r, err := secureRandomRune(allChars)
			if err != nil {
				return "", err
			}
			password = append(password, r)
		}

		// Shuffle to remove predictability
		if err := shuffleRunes(password); err != nil {
			return "", err
		}

		pwdStr := string(password)

		// Validate constraints
		if hasConsecutiveDuplicates(password) {
			continue
		}
		if !meetsCharClassRequirements(password) {
			continue
		}
		if hammingDistance(pwdStr, prev[0]) < diffRequired {
			continue
		}
		if inLastN(pwdStr, prev, historyCount) {
			continue
		}
		if isSequential(password) {
			continue
		}

		return pwdStr, nil
	}

	return "", fmt.Errorf("failed to generate valid password after multiple attempts")
}
