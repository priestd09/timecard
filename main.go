package main

////////////////////////////////////////////////////////////////////////////////

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"

	"github.com/sabhiram/timecard/git"
	"github.com/sabhiram/timecard/timecard"
)

////////////////////////////////////////////////////////////////////////////////

const (
	timecardFile = ".timecard"
	version      = "0.0.1"
	usage        = `usage: timecard [--version] [--help] <command> [<args>]

Valid Timecard commands include:
    init        Create an empty timecard or re-initialize an existing one
    start       Start or re-start the timecard for the current commit
    checkpoint  Create a checkpoint within a given interval
    end         End a timestamp with a given tag (usually a commit hash)
`
)

////////////////////////////////////////////////////////////////////////////////

var (
	CLI = struct {
		version bool   // print application version
		help    bool   // print application usages
		cwd     string // application's working directory
	}{}
)

////////////////////////////////////////////////////////////////////////////////

type cmdFn func(args []string) error

func initFunc(args []string) error {
	g, err := git.New(CLI.cwd)
	if err != nil {
		log.Fatalf("Error: Could not find a valid git repository at %s. Did you \"git init\"?\n", CLI.cwd)
	}

	tcfp := path.Join(CLI.cwd, timecardFile)
	if _, err := os.Stat(tcfp); os.IsNotExist(err) {
		// Create a default timecard for this project
		if _, err := timecard.Init(g, tcfp); err != nil {
			return err
		}
		log.Printf("Initialized new timecard for %s in %s.", "user", tcfp)
		return nil
	}

	// .timecard file already exists, do nothing.
	log.Printf("Timecard already setup for %s, use --force to re-initialize.\n", tcfp)
	return nil
}

func startFunc(args []string) error {
	g, err := git.New(CLI.cwd)
	if err != nil {
		log.Fatalf("Error: Could not find a valid git repository at %s. Did you \"git init\"?\n", CLI.cwd)
	}

	tcfp := path.Join(CLI.cwd, timecardFile)
	tc, err := timecard.Load(g, tcfp)
	if err != nil {
		return err
	}
	return tc.Start()
}

func checkpointFunc(args []string) error {
	log.Printf("checkpoint: %#v\n", args)
	return nil
}

func endFunc(args []string) error {
	g, err := git.New(CLI.cwd)
	if err != nil {
		log.Fatalf("Error: Could not find a valid git repository at %s. Did you \"git init\"?\n", CLI.cwd)
	}

	tcfp := path.Join(CLI.cwd, timecardFile)
	tc, err := timecard.Load(g, tcfp)
	if err != nil {
		return err
	}
	return tc.End()
}

////////////////////////////////////////////////////////////////////////////////

var fnMap = map[string]cmdFn{
	"init":       initFunc,
	"start":      startFunc,
	"checkpoint": checkpointFunc,
	"end":        endFunc,
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	args := flag.Args()
	if CLI.version {
		log.Printf("%s\n", version)
	} else if CLI.help || len(args) == 0 {
		log.Printf("%s\n", usage)
	} else {
		cmd, args := strings.ToLower(args[0]), args[1:]
		if fn, ok := fnMap[cmd]; ok {
			if err := fn(args); err != nil {
				log.Fatalf("%s command failed: %s\n", cmd, err.Error())
			}
		} else {
			log.Fatalf("Unknown command! %s\n", usage)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetPrefix("")
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	flag.BoolVar(&CLI.version, "version", false, "print the version")
	flag.BoolVar(&CLI.version, "v", false, "print the version (short)")
	flag.BoolVar(&CLI.help, "help", false, "print help")
	flag.BoolVar(&CLI.help, "h", false, "print help (short)")
	flag.Parse()

	var err error
	CLI.cwd, err = os.Getwd()
	if err != nil {
		log.Fatalf("Unable to query current working dir, %s\n", err.Error())
	}
}
