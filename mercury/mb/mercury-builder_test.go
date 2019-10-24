package mb

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
	result := `# This is a Space Note
@is.sour.registry meta
# Key note about such and such
key-name meta-key  :green
# Key note about such and such
key-name meta-key  :blue

`

	Convey("Evaluate Mercury Builder", t, func() {
		So(tt.String(), ShouldEqual, result)
	})
}
