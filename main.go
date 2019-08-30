// build-it is a tool to simplify the task of wiring up build commands (sent to a shell or something).
//
// Think of it as a Jenkinsfile that you can use from your CI/CD pipeline or from your terminal.
//
// Build scripts are written in lisp because I wanted an excuse to use it
//
// Applies POLA principle, if you have a reference to a function that means you can call it, do give functions to people that don't deserve them!
//
// shell-out is a super dangerous functions and will bit you if you give it to somebody.
package main

import (
	"flag"
	"io/ioutil"
	"os"
	"unicode/utf8"

	"github.com/pkg/errors"

	"github.com/rumlang/rum/parser"
	"github.com/rumlang/rum/runtime"
)

var (
	scriptFile = flag.String("script", "-", "Which script to run. Leave - to read from stdin (not interactive)")
	scriptCode = flag.String("code", "", "When configured, will run this code instead of script")
	allowHTTP  = flag.Bool("http", false, "Allow the library to issue simple http requests (get/post/put/delete/head)")
	runCmd     = flag.Bool("run-cmd", false, "Allow the script start other processes")
	prefix     = flag.String("prefix", "[build-it]", "Prefix used to output build-it info")
	verbose    = flag.Int("verbose", 0, "How verbose the app should be (only used for build-it related logging)")
)

func main() {
	flag.Parse()
	log.SetPrefix(*prefix)
	script := readAllScript(*scriptFile, *scriptCode)

	runScript(script)
	log.Print(4, "All good!", "(･o･)ง")
}

func runScript(script []byte) {
	if !utf8.Valid(script) {
		log.Fatal("build-it: script is not utf-8 encoded! <(｀^´)>")
	}

	root, err := parser.Parse(parser.NewSource(string(script)))
	if err != nil {
		log.Fatal(errors.Wrap(err, "build-it: parse error! (ʘᗩʘ’)"))
	}

	ctx := runtime.NewContext(nil)
	builditLib := &lib{canRunCmd: *runCmd}
	builditLib.LoadLib(ctx)
	(httpLib{}).LoadLibrary(ctx)

	if _, err = ctx.TryEval(root); err != nil {
		log.Fatal(errors.Wrap(err, "build-it: execution error! (╯°□°）╯︵ ┻━┻"))
	}
}

func readAllScript(file string, code string) []byte {
	if len(code) != 0 {
		return []byte(code)
	}
	log.Print(3, "Reading input from", file)
	fileHandle := os.Stdin
	if file != "-" {
		var err error
		fileHandle, err = os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fileHandle.Close()
	}
	data, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
