package logtag

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
			args: args{Verror, "something happened"},
			out:  out{"ERR", "something happened", "log_test.go", "logtag.TestStdLog.func1"},
		},
		{
			name: "warning with tags",
			args: args{Vwarning, "something happened\nwith another line", Tags{"foo": value("bar"), "bin": value("baz")}},
			out:  out{"WARN", "something happened", "with another line", "foo", "bar", "bin", "baz", "log_test.go", "logtag.TestStdLog.func1"},
		},
	}

	for _, tt := range tests {
		Convey(tt.name, t, func() {
			var b bytes.Buffer
			logger := &StdLogger{out: &b, scheme: MonoScheme}

			Output(logger, 2, tt.args...)
			txt := b.String()
			for _, out := range tt.out {
				So(txt, ShouldContainSubstring, out)
			}
		})
	}

	Convey("setting values in tag", t, func() {
		tags := make(Tags)
		var b bytes.Buffer
		logger := &StdLogger{out: &b, scheme: MonoScheme}

		tags.Set("random", "string")
		tags.Set("int", 123)
		tags.Set("event", Event{Message: "hello"})

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
			args: args{Verror, "something happened"},
			out:  out{"ERR", "something happened"},
		},
		{
			name: "warning with tags",
			args: args{Vwarning, "something happened\nwith another line", Tags{"foo": value("bar"), "bin": value("baz")}},
			out:  out{"WARN", "something happened", "with another line", "foo", "bar", "bin", "baz", "log_test.go", "logtag.TestJSONLog.func1"},
		},
	}

	for _, tt := range tests {
		Convey(tt.name, t, func() {
			var b bytes.Buffer
			logger := &JSONLogger{out: &b}

			Output(logger, 2, tt.args...)
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
			args: args{Verror, "something happened"},
			out:  out{"ERR", "something happened"},
		},
		{
			name: "warning with tags",
			args: args{Vwarning, "something happened\nwith another line", Tags{"foo": value("bar"), "bin": value("baz")}},
			out:  out{"WARN", "something happened", "with another line", "foo", "bar", "bin", "baz", "log_test.go", "logtag.TestFanLog.func1"},
		},
	}

	for _, tt := range tests {
		Convey(tt.name, t, func() {
			var b1 bytes.Buffer
			logger1 := &StdLogger{out: &b1, scheme: MonoScheme}

			var b2 bytes.Buffer
			logger2 := &JSONLogger{out: &b2}

			logger := &FanLogger{outs: []Logger{logger1, logger2}}

			Output(logger, 2, tt.args...)

			txt1 := b1.String()
			txt2 := b2.String()
			for _, out := range tt.out {
				So(txt1, ShouldContainSubstring, out)
				So(txt2, ShouldContainSubstring, out)
			}
		})
	}
}
