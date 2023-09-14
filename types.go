package apis

import (
	"github.com/google/uuid"
	"strings"
)

// StringUUID is used to cast query values to uuid.UUID. Because gin BindQuery does not recognize uuids, we use a
// string instead, with a method to dynamically convert it to uuid.UUID. If the content of the string is not
// a valid uuid, it will return uuid.Nil.
// https://github.com/gin-gonic/gin/issues/1516#issuecomment-1269846541
type StringUUID string

func (s StringUUID) Value() uuid.UUID {
	parsed, err := uuid.Parse(string(s))
	if err != nil {
		return uuid.Nil
	}

	return parsed
}

// StringUUIDs represents an array of StringUUID.
type StringUUIDs string

func (s StringUUIDs) Value() []uuid.UUID {
	var uuids []uuid.UUID

	for _, id := range strings.Split(string(s), ",") {
		parsed, err := uuid.Parse(id)
		if err != nil {
			continue
		}

		uuids = append(uuids, parsed)
	}

	return uuids
}
