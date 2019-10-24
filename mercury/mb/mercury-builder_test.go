package mb

import (
	"fmt"
	"testing"
)

func TestBuilder(t *testing.T) {
	tt := Config(
		NewSpace("is.sour.registry",
			WithNotes("This is a Space Note"),
			WithTags("meta"),

			NewKey("key-name",
				WithNotes("Key note about such and such"),
				WithTags("meta-key"),
				WithValue("green"),
			),

			NewKey("key-name",
				WithNotes("Key note about such and such"),
				WithTags("meta-key"),
				WithValue("blue"),
			),
		),
	)

	fmt.Println(tt.String())
}
