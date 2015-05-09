package sde

import (
	"time"
)

func ApplyTypeToType(tone, ttwo SDEType) (SDEType, error) {
	defer Debug(time.Now())
	out := tone

	return out, nil

}
