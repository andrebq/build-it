// build-it is a tool to simplify the task of wiring up build commands (sent to a shell or something).
//
// Think of it as a Jenkinsfile that you can use from your CI/CD pipeline or from your terminal.
//
// Build scripts are written in lua because I like it!
//
// Applies POLA principle, if you have a reference to a function that means you can call it, do give functions to people that don't deserve them!
//
// cmd.run is a super dangerous functions and will bit you if you give it to somebody, by default it is disabled. Run with -run-cmd true to enable it
package main

import (
	"flag"
	"io/ioutil"
	"os"
	"unicode/utf8"

	jsonmodule "github.com/layeh/gopher-json"

	lua "github.com/yuin/gopher-lua"
)

var (
	scriptFile = flag.String("script", "-", "Which script to run. Leave - to read from stdin (not interactive)")
	scriptCode = flag.String("code", "", "When configured, will run this code instead of script")
	allowHTTP  = flag.Bool("http", false, "Allow the library to issue simple http requests (get/post/put/delete/head)")
	runCmd     = flag.Bool("run-cmd", false, "Allow the script start other processes")
	prefix     = flag.String("prefix", "[build-it] ", "Prefix used to output build-it info")
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

	L := lua.NewState()
	L.PreloadModule("json", jsonmodule.Loader)
	if *allowHTTP {
		L.PreloadModule("http", (httpLib{}.load))
	}
	L.PreloadModule("cmd", lib{*runCmd}.load)

	defer L.Close()

	if err := L.DoString(string(script)); err != nil {
		panic(err)
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
