package log // import "sour.is/x/toolbox/log"

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/bouk/monkey"
	"sour.is/x/toolbox/log/event"
	"sour.is/x/toolbox/log/loggers"
	"sour.is/x/toolbox/log/scheme"
)

func TestDefaultLog(t *testing.T) {
	var b bytes.Buffer
	var w = bufio.NewWriter(&b)
	var debug = loggers.NewStdLogger(w, scheme.MonoScheme, event.VerbDebug)
	var color = loggers.NewStdLogger(w, scheme.ColorScheme, event.VerbDebug)
	var none = loggers.NewStdLogger(w, scheme.MonoScheme, event.VerbNone)

	Convey("Given the standard logger", t, func() {
		Convey("Startup Banner", func() {
			SetOutput(debug)

			StartupBanner()
			w.Flush()
			So(b.String(), ShouldNotBeBlank)
			b.Reset()

		})
		Convey("With Verbose set to Debug", func() {
			SetOutput(debug)

			Debug("Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "DBUG")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Info("Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "INFO")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Warning("Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "WARN")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Notice("Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "NOTE")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Error("Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "ERR ")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Critical("Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "CRIT")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Debugf("%s", "Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "DBUG")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Infof("%s", "Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "INFO")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Warningf("%s", "Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "WARN")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Noticef("%s", "Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "NOTE")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Errorf("%s", "Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "ERR ")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Criticalf("%s", "Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "CRIT")
			So(b.String(), ShouldContainSubstring, "Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()
		})

		Convey("With Verbose set to None", func() {
			SetOutput(none)

			Debug("Test")
			w.Flush()
			So(b.String(), ShouldEqual, "")
			b.Reset()

			Info("Test")
			w.Flush()
			So(b.String(), ShouldEqual, "")
			b.Reset()

			Warning("Test")
			w.Flush()
			So(b.String(), ShouldEqual, "")
			b.Reset()

			Notice("Test")
			w.Flush()
			So(b.String(), ShouldEqual, "")
			b.Reset()

			Error("Test")
			w.Flush()
			So(b.String(), ShouldEqual, "")
			b.Reset()

			Critical("Test")
			w.Flush()
			So(b.String(), ShouldEqual, "")
			b.Reset()

			Debugf("%s", "Test")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "")
			b.Reset()

			Infof("%s", "Test")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "")
			b.Reset()

			Warningf("%s", "Test")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "")
			b.Reset()

			Noticef("%s", "Test")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "")
			b.Reset()

			Errorf("%s", "Test")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "")
			b.Reset()

			Criticalf("%s", "Test")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "")
			b.Reset()
		})

		Convey("With Color and testing Print", func() {
			SetOutput(color)
			b.Reset()

			Print("Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "\x1B[34mINFO\x1B[90m ")
			So(b.String(), ShouldContainSubstring, "\x1B[0mTest\x1B[0m")
			So(b.String(), ShouldContainSubstring, "\x1B[90m....\x1B[90m ")
			So(b.String(), ShouldContainSubstring, "\x1B[0mMultiline\x1B[0m")
			b.Reset()

			Printf("%s", "Test\nMultiline")
			w.Flush()
			So(b.String(), ShouldContainSubstring, "\x1B[34mINFO\x1B[90m ")
			So(b.String(), ShouldContainSubstring, "\x1B[0mTest\x1B[0m")
			So(b.String(), ShouldContainSubstring, "\x1B[90m....\x1B[90m ")
			So(b.String(), ShouldContainSubstring, "\x1B[0mMultiline\x1B[0m")
			b.Reset()

		})

		Convey("Testing fatal and panic outputs", func() {
			SetOutput(debug)

			So(func() { Panic("Test\nMultiline") }, ShouldPanic)
			w.Flush()
			So(b.String(), ShouldContainSubstring, "ERR ")
			So(b.String(), ShouldContainSubstring, " Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			So(func() { Panicf("%s", "Test\nMultiline") }, ShouldPanic)
			w.Flush()
			So(b.String(), ShouldContainSubstring, "ERR ")
			So(b.String(), ShouldContainSubstring, " Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			exitCode := 0
			monkey.Patch(os.Exit, func(code int) { exitCode = code })
			defer monkey.UnpatchAll()

			Fatal("Test\nMultiline")
			So(exitCode, ShouldEqual, 1)
			w.Flush()
			So(b.String(), ShouldContainSubstring, "CRIT")
			So(b.String(), ShouldContainSubstring, " Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

			Fatalf("%s", "Test\nMultiline")
			So(exitCode, ShouldEqual, 1)
			w.Flush()
			So(b.String(), ShouldContainSubstring, "CRIT")
			So(b.String(), ShouldContainSubstring, " Test")
			So(b.String(), ShouldContainSubstring, ".... Multiline")
			b.Reset()

		})

		Convey("Nil Functions should do nothing.", func() {
			SetOutput(none)

			NilPrint("PRINT")
			NilPrintf("PRINTF")
			NilPrintln("PRINTLN")

			NilDebug("DBUG")
			NilInfo("INFO")
			NilWarning("WARN")
			NilNotice("NOTE")
			NilError("ERRO")
			NilCritical("CRIT")

			NilDebugf("%s", "DBUGF")
			NilInfof("%s", "INFOF")
			NilWarningf("%s", "WARNF")
			NilNoticef("%s", "NOTEF")
			NilErrorf("%s", "ERROF")
			NilCriticalf("%s", "CRITF")

			w.Flush()
			So(b.String(), ShouldContainSubstring, "")
			b.Reset()

		})

	})
}
