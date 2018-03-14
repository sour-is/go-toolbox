package log // import "sour.is/x/toolbox/log"

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	Ldate         = 1 << iota            // the date in the local time zone: 2009/01/23
	Ltime                                // the time in the local time zone: 01:23:23
	Lmicroseconds                        // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                            // full file name and line number: /a/b/c/d.go:23
	Lshortfile                           // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                                 // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime | LUTC // initial values for the standard logger
)

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

const (
	Vnone     = 0
	Vcritical = 1 << iota
	Verror
	Vwarning
	Vnotice
	Vinfo
	Vdebug
)

var Freset = Creset
var Fprefix = Cprefix
var Fdebug = Cdebug
var Finfo = Cinfo
var Fnotice = Cnotice
var Fwarning = Cwarning
var Ferror = Cerror
var Fcritical = Ccritical
var Fcontinue = Ccontinue

var std = log.New(os.Stderr, Fprefix, LstdFlags)
var mu = sync.Mutex{}
var verb int = Vnone

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
func Flags() (flag int) {
	return std.Flags()
}

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
}

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
	s := strings.Split(fmt.Sprint(v...), "\n")
	std.Output(2, Ferror+s[0]+Freset)
	for _, l := range s[1:] {
		std.Output(2, Fcontinue+l+Freset)
	}
	panic(s)
}

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
