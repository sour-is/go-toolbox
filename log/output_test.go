package log

import (
	"bytes"
	"testing"

	"sour.is/x/toolbox/log/event"
	"sour.is/x/toolbox/log/loggers"
	"sour.is/x/toolbox/log/scheme"
	"sour.is/x/toolbox/log/tag"
)

func TestStdLog(t *testing.T) {
	type args []interface{}
	type out []string

	tests := []struct {
		name string
		args args
		out  out
	}{
		{
			name: "basic error",
			args: args{event.VerbError, "something happened"},
			out:  out{"ERR", "something happened", "output_test.go"},
		},
		{
			name: "warning with tags",
			args: args{event.VerbWarning, "something happened\nwith another line", tag.Tags{"foo": tag.Value("bar"), "bin": tag.Value("baz")}},
			out:  out{"WARN", "something happened", "with another line", "foo", "bar", "bin", "baz", "output_test.go"},
		},
	}

	for _, tt := range tests {
		Convey(tt.name, t, func() {
			var b bytes.Buffer
			logger := loggers.NewStdLogger(&b, scheme.MonoScheme, event.VerbDebug)

			Output(logger, 0, tt.args...)
			txt := b.String()
			for _, out := range tt.out {
				So(txt, ShouldContainSubstring, out)
			}
		})
	}

	Convey("setting values in tag", t, func() {
		tags := make(tag.Tags)
		var b bytes.Buffer
		logger := loggers.NewStdLogger(&b, scheme.MonoScheme, event.VerbDebug)

		tags.Set("random", "string")
		tags.Set("int", 123)
		tags.Set("event", event.Event{Message: "hello"})

		Output(logger, 2, "something happened", tags)
		txt := b.String()

		So(txt, ShouldContainSubstring, "DBUG")
		So(txt, ShouldContainSubstring, "something happened")
		So(txt, ShouldContainSubstring, "random=string")
		So(txt, ShouldContainSubstring, "int=123")
		So(txt, ShouldContainSubstring, "event={")
		So(txt, ShouldContainSubstring, "hello")
	})
}

func TestJSONLog(t *testing.T) {
	type args []interface{}
	type out []string

	tests := []struct {
		name string
		args args
		out  out
	}{
		{
			name: "basic error",
			args: args{event.VerbError, "something happened"},
			out:  out{"ERR", "something happened"},
		},
		{
			name: "warning with tags",
			args: args{event.VerbWarning, "something happened\nwith another line", tag.Tags{"foo": tag.Value("bar"), "bin": tag.Value("baz")}},
			out:  out{"WARN", "something happened", "with another line", "foo", "bar", "bin", "baz", "output_test.go"},
		},
	}

	for _, tt := range tests {
		Convey(tt.name, t, func() {
			var b bytes.Buffer
			logger := loggers.NewJSONLogger(&b, event.VerbDebug)

			Output(logger, 0, tt.args...)
			txt := b.String()
			for _, out := range tt.out {
				So(txt, ShouldContainSubstring, out)
			}
		})
	}
}

func TestFanLog(t *testing.T) {
	type args []interface{}
	type out []string

	tests := []struct {
		name string
		args args
		out  out
	}{
		{
			name: "basic error",
			args: args{event.VerbError, "something happened"},
			out:  out{"ERR", "something happened"},
		},
		{
			name: "warning with tags",
			args: args{event.VerbWarning, "something happened\nwith another line", tag.Tags{"foo": tag.Value("bar"), "bin": tag.Value("baz")}},
			out:  out{"WARN", "something happened", "with another line", "foo", "bar", "bin", "baz", "output_test.go"},
		},
	}

	for _, tt := range tests {
		Convey(tt.name, t, func() {
			var b1 bytes.Buffer
			logger1 := loggers.NewStdLogger(&b1, scheme.MonoScheme, event.VerbDebug)

			var b2 bytes.Buffer
			logger2 := loggers.NewJSONLogger(&b2, event.VerbDebug)

			logger := loggers.NewFanLogger(event.VerbDebug, logger1, logger2)

			Output(logger, 0, tt.args...)

			txt1 := b1.String()
			txt2 := b2.String()
			for _, out := range tt.out {
				So(txt1, ShouldContainSubstring, out)
				So(txt2, ShouldContainSubstring, out)
			}
		})
	}
}
