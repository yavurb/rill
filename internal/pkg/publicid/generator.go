package publicid

import (
	"log"
	"strings"

	nanoid "github.com/matoous/go-nanoid/v2"
)

func New(prefix string, lenght int) (string, error) {
	if lenght <= 0 {
		lenght = 15
	}

	id, err := nanoid.Generate("alphabet", lenght)
	if err != nil {
		log.Printf("Error generating nanoid: %v\n", err)

		return "", err
	}

	return strings.Join([]string{prefix, id}, "_"), nil
}
