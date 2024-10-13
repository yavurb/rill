package testhelpers

import (
	"encoding/json"
	"log"

	"github.com/google/go-cmp/cmp"
)

func CompareByteMaps(a, b []byte) bool {
	aT := make(map[string]any)
	bT := make(map[string]any)

	err := json.Unmarshal(a, &aT)
	if err != nil {
		log.Fatal("Error unmarshalling left side bytes")
	}

	err = json.Unmarshal(b, &bT)
	if err != nil {
		log.Fatal("Error unmarshalling right side bytes")
	}

	return cmp.Equal(aT, bT)
}

func CompareMaps(a, b any) bool {
	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)
	aT := make(map[string]any)
	bT := make(map[string]any)

	err := json.Unmarshal(aBytes, &aT)
	if err != nil {
		log.Fatal("Error unmarshalling aBytes")
	}

	err = json.Unmarshal(bBytes, &bT)
	if err != nil {
		log.Fatal("Error unmarshalling bBytes")
	}

	return cmp.Equal(aT, bT)
}
