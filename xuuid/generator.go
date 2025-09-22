package xuuid

import (
	"encoding/base32"
	"encoding/base64"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/google/uuid"
	"strings"
)

func Base32UUID() string {
	uid, _ := uuid.NewV7()

	encoder := base32.StdEncoding.WithPadding(base64.NoPadding)
	return strings.ToLower(encoder.EncodeToString(uid[:]))
}

func HumanReadableID() string {
	return petname.Generate(2, "-")
}
