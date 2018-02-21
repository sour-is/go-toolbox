package log // sour.is/go/log

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestDefaultLog(t *testing.T) {
	var b bytes.Buffer
	var w = bufio.NewWriter(&b)

	convey.Convey("Given the standard logger", t, func() {
		SetOutput(w)

		convey.Convey("Setting Flags", func() {
			SetFlags(Ldate | Ltime | Lmicroseconds)
			convey.So(Flags(), convey.ShouldEqual, Ldate|Ltime|Lmicroseconds)
		})

		convey.Convey("With Verbose set to Debug", func() {
			SetVerbose(Vdebug)
			SetColor(false)

			Debug("Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "DBUG ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Info("Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "INFO ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Warning("Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "WARN ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Notice("Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "NOTE ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Error("Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "ERR  ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Critical("Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "CRIT ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Debugf("%s", "Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "DBUG ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Infof("%s", "Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "INFO ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Warningf("%s", "Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "WARN ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Noticef("%s", "Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "NOTE ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Errorf("%s", "Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "ERR  ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Criticalf("%s", "Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "CRIT ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()
		})

		convey.Convey("With Verbose set to None", func() {
			SetVerbose(Vnone)
			SetColor(false)
			Debug("Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldEqual, "")
			b.Reset()

			Info("Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldEqual, "")
			b.Reset()

			Warning("Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldEqual, "")
			b.Reset()

			Notice("Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldEqual, "")
			b.Reset()

			Error("Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldEqual, "")
			b.Reset()

			Critical("Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldEqual, "")
			b.Reset()

			Debugf("%s", "Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "")
			b.Reset()

			Infof("%s", "Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "")
			b.Reset()

			Warningf("%s", "Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "")
			b.Reset()

			Noticef("%s", "Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "")
			b.Reset()

			Errorf("%s", "Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "")
			b.Reset()

			Criticalf("%s", "Test")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "")
			b.Reset()
		})

		convey.Convey("With Color and testing Print", func() {
			SetVerbose(Vdebug)
			SetColor(true)

			Print("Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "\x1B[34mINFO \x1B[90m] Test\x1B[0m")
			convey.So(b.String(), convey.ShouldContainSubstring, "\x1B[90m.... \x1B[90m] Multiline\x1B[0m")
			b.Reset()

			Println("Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "\x1B[34mINFO \x1B[90m] Test\x1B[0m")
			convey.So(b.String(), convey.ShouldContainSubstring, "\x1B[90m.... \x1B[90m] Multiline\x1B[0m")
			b.Reset()

			Printf("%s", "Test\nMultiline")
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "\x1B[34mINFO \x1B[90m] Test\x1B[0m")
			convey.So(b.String(), convey.ShouldContainSubstring, "\x1B[90m.... \x1B[90m] Multiline\x1B[0m")
			b.Reset()

		})

		convey.Convey("Testing fatal and panic outputs", func() {
			SetVerbose(Vnone)
			SetColor(false)

			convey.So(func() { Panic("Test\nMultiline") }, convey.ShouldPanic)
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "ERR  ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			convey.So(func() { Panicf("%s", "Test\nMultiline") }, convey.ShouldPanic)
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "ERR  ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			exitCode := 0
			monkey.Patch(os.Exit, func(code int) { exitCode = code })
			defer monkey.UnpatchAll()

			Fatal("Test\nMultiline")
			convey.So(exitCode, convey.ShouldEqual, 1)
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "CRIT ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

			Fatalf("%s", "Test\nMultiline")
			convey.So(exitCode, convey.ShouldEqual, 1)
			w.Flush()
			convey.So(b.String(), convey.ShouldContainSubstring, "CRIT ] Test")
			convey.So(b.String(), convey.ShouldContainSubstring, "... ] Multiline")
			b.Reset()

		})

	})
}
