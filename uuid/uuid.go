package uuid

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GenPrefixed uuid
func GenPrefixed(length int, prefix ...string) string {
	if len(prefix) > 0 && prefix[0] != "" {
		return fmt.Sprintf("%s-%s", prefix[0], genV4UUID()[:length])
	}
	return genV4UUID()[:length]
}

// Generate uuid
func Generate() string {
	return uuid.New().String()
}

func genV4UUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
