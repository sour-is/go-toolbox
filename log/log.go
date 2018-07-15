package log // import "sour.is/x/toolbox/log"

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

// Debug flags
const (
	Ldate         = 1 << iota            // the date in the local time zone: 2009/01/23
	Ltime                                // the time in the local time zone: 01:23:23
	Lmicroseconds                        // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                            // full file name and line number: /a/b/c/d.go:23
	Lshortfile                           // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                                 // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime | LUTC // initial values for the standard logger
)

// Debug error levels and color coding
const (
	Tdebug    = "DBUG"
	Tinfo     = "INFO"
	Tnotice   = "NOTE"
	Twarning  = "WARN"
	Terror    = "ERR "
	Tcritical = "CRIT"
	Tcontinue = "...."

	Creset    = "\x1B[0m"
	Cprefix   = "\x1B[90m"
	Cdebug    = "\x1B[90m" + Tdebug + " " + Cprefix + "] "
	Cinfo     = "\x1B[34m" + Tinfo + " " + Cprefix + "] "
	Cnotice   = "\x1B[32m" + Tnotice + " " + Creset + "] "
	Cwarning  = "\x1B[93m" + Twarning + " " + Cprefix + "] "
	Cerror    = "\x1B[91m" + Terror + " " + Creset + "] "
	Ccritical = "\x1B[7;91;49m" + Tcritical + " " + Creset + "] "
	Ccontinue = "\x1B[90m" + Tcontinue + " " + Cprefix + "] "

	Mreset    = ""
	Mprefix   = ""
	Mdebug    = Tdebug + " ] "
	Minfo     = Tinfo + " ] "
	Mnotice   = Tnotice + " ] "
	Mwarning  = Twarning + " ] "
	Merror    = Terror + " ] "
	Mcritical = Tcritical + " ] "
	Mcontinue = Tcontinue + " ] "
)

// Debug message levels
const (
	Vnone     = 0
	Vcritical = 1 << iota
	Verror
	Vwarning
	Vnotice
	Vinfo
	Vdebug
)

// Set default formatting
var (
   Freset = Creset
   Fprefix = Cprefix
   Fdebug = Cdebug
   Finfo = Cinfo
   Fnotice = Cnotice
   Fwarning = Cwarning
   Ferror = Cerror
   Fcritical = Ccritical
   Fcontinue = Ccontinue
)

var std = log.New(os.Stderr, Fprefix, LstdFlags)
var mu = sync.Mutex{}
var verb = Vnone

// StartupBanner displays a random banner.
func StartupBanner() {
	rand.Seed(time.Now().UnixNano())
	i := rand.Int()
	Print(strings.Join(souris[i%len(souris)], "\n"))
}
// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	std.SetOutput(w)
}
// SetFlags sets the output flags for the standard logger.
func SetFlags(flag int) {
	std.SetFlags(flag)
}
// Flags returns the current flag values
func Flags() (flag int) {
	return std.Flags()
}
// SetColor enables or disables the display of color
func SetColor(on bool) {
	mu.Lock()
	defer mu.Unlock()
	if on {
		Freset = Creset
		Fprefix = Cprefix
		Fdebug = Cdebug
		Finfo = Cinfo
		Fnotice = Cnotice
		Fwarning = Cwarning
		Ferror = Cerror
		Fcritical = Ccritical
		Fcontinue = Ccontinue
	} else {
		Freset = Mreset
		Fprefix = Mprefix
		Fdebug = Mdebug
		Finfo = Minfo
		Fnotice = Mnotice
		Fwarning = Mwarning
		Ferror = Merror
		Fcritical = Mcritical
		Fcontinue = Mcontinue
	}
	std.SetPrefix(Fprefix)
}
// SetVerbose level to output
func SetVerbose(v int) {
	mu.Lock()
	defer mu.Unlock()
	verb = v
}

// These functions write to the standard logger.

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	s := strings.Split(fmt.Sprint(v...), "\n")
	std.Output(2, Finfo+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	s := strings.Split(fmt.Sprintf(format, v...), "\n")
	std.Output(2, Finfo+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	s := strings.Split(fmt.Sprint(v...), "\n")
	std.Output(2, Finfo+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Write outputs contents of io.Reader to standard logger.
func Write(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	s := strings.Split(buf.String(), "\n")

	std.Output(2, Finfo+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Tee outputs contents of io.Reader to standard logger. and returns a new io.Reader
func Tee(r io.ReadCloser) (w io.ReadCloser) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	r.Close()

	str := buf.String()
	w = ioutil.NopCloser(strings.NewReader(str))

	s := strings.Split(str, "\n")
	std.Output(2, Finfo+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}

	return
}
// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	s := strings.Split(fmt.Sprint(v...), "\n")
	std.Output(2, Fcritical+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}

	os.Exit(1)
}
// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	s := strings.Split(fmt.Sprintf(format, v...), "\n")
	std.Output(2, Fcritical+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
	os.Exit(1)
}
// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := strings.Split(fmt.Sprint(v...), "\n")
	std.Output(2, Ferror+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
	panic(s)
}
// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := strings.Split(fmt.Sprintf(format, v...), "\n")
	std.Output(2, Ferror+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
	panic(s)
}
// Debug outputs to logger with DEBUG level.
func Debug(v ...interface{}) {
	if verb < Vdebug {
		return
	}
	s := strings.Split(fmt.Sprint(v...), "\n")
	std.Output(2, Fdebug+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Debugf formats output to logger with DEBUG level.
func Debugf(format string, v ...interface{}) {
	if verb < Vdebug {
		return
	}
	s := strings.Split(fmt.Sprintf(format, v...), "\n")
	std.Output(2, Fdebug+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Debugw outputs io.Reader to logger with DEBUG level.
func Debugw(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	s := strings.Split(buf.String(), "\n")

	std.Output(2, Fdebug+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Info outputs to logger with INFO level.
func Info(v ...interface{}) {
	if verb < Vinfo {
		return
	}
	s := strings.Split(fmt.Sprint(v...), "\n")
	std.Output(2, Finfo+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Infof formatted outputs to logger with INFO level.
func Infof(format string, v ...interface{}) {
	if verb < Vinfo {
		return
	}
	s := strings.Split(fmt.Sprintf(format, v...), "\n")
	std.Output(2, Finfo+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Infow outputs io.Reader to logger with INFO level.
func Infow(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	s := strings.Split(buf.String(), "\n")

	std.Output(2, Finfo+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Notice outputs to logger with NOTICE level.
func Notice(v ...interface{}) {
	if verb < Vnotice {
		return
	}
	s := strings.Split(fmt.Sprint(v...), "\n")
	std.Output(2, Fnotice+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Noticef formatted outputs to logger with NOTICE level.
func Noticef(format string, v ...interface{}) {
	if verb < Vnotice {
		return
	}
	s := strings.Split(fmt.Sprintf(format, v...), "\n")
	std.Output(2, Fnotice+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Noticew outputs io.Reader to logger with NOTICE level.
func Noticew(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	s := strings.Split(buf.String(), "\n")

	std.Output(2, Fnotice+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Warning outputs to logger with WARNING level.
func Warning(v ...interface{}) {
	if verb < Vwarning {
		return
	}
	s := strings.Split(fmt.Sprint(v...), "\n")
	std.Output(2, Fwarning+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Warningf formatted outputs to logger with WARNING level.
func Warningf(format string, v ...interface{}) {
	if verb < Vwarning {
		return
	}
	s := strings.Split(fmt.Sprintf(format, v...), "\n")
	std.Output(2, Fwarning+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Warningw outputs io.Reader to logger with WARNING level.
func Warningw(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	s := strings.Split(buf.String(), "\n")

	std.Output(2, Fwarning+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Error outputs to logger with ERROR level.
func Error(v ...interface{}) {
	if verb < Verror {
		return
	}
	s := strings.Split(fmt.Sprint(v...), "\n")
	std.Output(2, Ferror+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Errorf formatted outputs to logger with ERROR level.
func Errorf(format string, v ...interface{}) {
	if verb < Verror {
		return
	}
	s := strings.Split(fmt.Sprintf(format, v...), "\n")
	std.Output(2, Ferror+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Errorw outputs io.Reader to logger with ERROR level.
func Errorw(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	s := strings.Split(buf.String(), "\n")

	std.Output(2, Ferror+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Critical outputs to logger with CRITICAL level.
func Critical(v ...interface{}) {
	if verb < Vcritical {
		return
	}
	s := strings.Split(fmt.Sprint(v...), "\n")
	std.Output(2, Fcritical+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Criticalf formatted outputs to logger with CRITICAL level.
func Criticalf(format string, v ...interface{}) {
	if verb < Vcritical {
		return
	}
	s := strings.Split(fmt.Sprintf(format, v...), "\n")
	std.Output(2, Fcritical+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}
// Criticalw outputs io.Reader to logger with CRITICAL level.
func Criticalw(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	s := strings.Split(buf.String(), "\n")

	std.Output(2, Fcritical+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
}

// These functions do nothing. It makes it easy to comment out
// debug lines without having to remove the import.

// NilPrint does nothing.
func NilPrint(_ ...interface{})               {}
// NilPrintf does nothing.
func NilPrintf(_ string, _ ...interface{})    {}
// NilPrintln does nothing.
func NilPrintln(_ ...interface{})             {}
// NilDebug does nothing.
func NilDebug(_ ...interface{})               {}
// NilDebugf does nothing.
func NilDebugf(_ string, _ ...interface{})    {}
// NilDebugw does nothing.
func NilDebugw(_ io.Reader)                   {}
// NilInfo does nothing.
func NilInfo(_ ...interface{})                {}
// NilInfof does nothing.
func NilInfof(_ string, _ ...interface{})     {}
// NilInfow does nothing.
func NilInfow(_ io.Reader)                    {}
// NilNotice does nothing.
func NilNotice(_ ...interface{})              {}
// NilNoticef does nothing.
func NilNoticef(_ string, _ ...interface{})   {}
// NilNoticew does nothing.
func NilNoticew(_ io.Reader)                  {}
// NilWarning does nothing.
func NilWarning(_ ...interface{})             {}
// NilWarningf does nothing.
func NilWarningf(_ string, _ ...interface{})  {}
// NilWarningw does nothing.
func NilWarningw(_ io.Reader)                 {}
// NilError does nothing.
func NilError(_ ...interface{})               {}
// NilErrorf does nothing.
func NilErrorf(_ string, _ ...interface{})    {}
// NilErrorw does nothing.
func NilErrorw(_ io.Reader)                   {}
// NilCritical does nothing.
func NilCritical(_ ...interface{})            {}
// NilCriticalf does nothing.
func NilCriticalf(_ string, _ ...interface{}) {}
// NilCriticalw does nothing.
func NilCriticalw(_ io.Reader)                {}

var souris = [][]string{
	{
		`  _________                     .__`,
		` /   _____/ ____  __ _________  |__| ______`,
		` \_____  \ /  _ \|  |  \_  __ \ |  |/  ___/`,
		` /        (  <_> |  |  /|  | \/ |  |\___ \`,
		`/_______  /\____/|____/ |__| /\ |__/____  >`,
		`\/                           \/         \/`},

	{
		`  ________  ______   ____  ____  _______         __     ________`,
		` /"       )/    " \ ("  _||_ " |/"      \       |" \   /"       )`,
		`(:   \___/// ____  \|   (  ) : |:        |      ||  | (:   \___/ `,
		` \___  \ /  /    ) :(:  |  | . |_____/   )      |:  |  \___  \`,
		`  __/  \(: (____/ // \\ \__/ // //      /  _____|.  |   __/  \\  `,
		` /" \   :\        /  /\\ __ //\|:  __   \ ))_  "/\  |\ /" \   :)  `,
		`(_______/ \"_____/  (__________|__|  \___(_____(__\_|_(_______/`},

	{
		` _____                  _`,
		`/  ___|                (_)`,
		"\\ `--.  ___  _   _ _ __ _ ___ ",
		" `--. \\/ _ \\| | | | '__| / __|",
		`/\__/ | (_) | |_| | |_ | \__ \`,
		`\____/ \___/ \__,_|_(_)|_|___/`},

	{
		`  ██████ ▒█████  █    ██ ██▀███        ██▓ ██████`,
		`▒██    ▒▒██▒  ██▒██  ▓██▓██ ▒ ██▒     ▓██▒██    ▒`,
		`░ ▓██▄  ▒██░  ██▓██  ▒██▓██ ░▄█ ▒     ▒██░ ▓██▄`,
		`  ▒   ██▒██   ██▓▓█  ░██▒██▀▀█▄       ░██░ ▒   ██▒`,
		`▒██████▒░ ████▓▒▒▒█████▓░██▓ ▒██▒ ██▓ ░██▒██████▒▒`,
		`▒ ▒▓▒ ▒ ░ ▒░▒░▒░░▒▓▒ ▒ ▒░ ▒▓ ░▒▓░ ▒▓▒ ░▓ ▒ ▒▓▒ ▒ ░`,
		`░ ░▒  ░ ░ ░ ▒ ▒░░░▒░ ░ ░  ░▒ ░ ▒░ ░▒   ▒ ░ ░▒  ░ ░`,
		`░  ░  ░ ░ ░ ░ ▒  ░░░ ░ ░  ░░   ░  ░    ▒ ░  ░  ░`,
		`     ░     ░ ░    ░       ░       ░   ░       ░`,
		`                                   ░`},

	{
		` .▄▄ ·     ▄• ▄▄▄▄ ▪ .▄▄ ·`,
		`▐█ ▀.▪    █▪██▀▄ ███▐█ ▀.`,
		`▄▀▀▀█▄▄█▀▄█▌▐█▐▀▀▄▐█▄▀▀▀█▄`,
		`▐█▄▪▐▐█▌.▐▐█▄█▐█•█▐█▐█▄▪▐█`,
		` ▀▀▀▀ ▀█▄▀▪▀▀▀.▀  ▀▀▀▀▀▀▀`},

	{
		"  .--.--.",
		" /  /    '.                                 ,--,",
		"|  :  /`. /   ,---.          ,--,  __  ,-.,--.'|",
		";  |  |--`   '   ,'\\       ,'_ /|,' ,'/ /||  |,     .--.--.",
		"|  :  ;_    /   /   | .--. |  | :'  | |' |`--'_    /  /    '",
		" \\  \\    `..   ; ,. ,'_ /| :  . ||  |   ,',' ,'|  |  :  /`./",
		"  `----.   '   | |: |  ' | |  . .'  :  /  '  | |  |  :  ;_",
		"  __ \\  \\  '   | .; |  | ' |  | ||  | '   |  | :   \\  \\    `.",
		" /  /`--'  |   :    :  | : ;  ; |;  : |   '  : |__  `----.   \\",
		"'--'.     / \\   \\  /'  :  `--'   |  , ___ |  | '.'|/  /`--'  /",
		"  `--'---'   `----' :  ,      .-./---/  .\\;  :    '--'.     /",
		"                     `--`----'       \\  ; |  ,   /  `--'---'",
		"                                      `--\" ---`-'"},
}
