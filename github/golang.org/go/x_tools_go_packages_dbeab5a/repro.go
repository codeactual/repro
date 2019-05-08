// repro.go loads -pkgmax number of packages in -pkgdir found by `go list`
// in order to allow generate CPU/memory profiles of x/tools/go/packages.Load.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strings"

	"golang.org/x/tools/go/packages"
)

var pkgdir = flag.String("pkgdir", "", "directory from which to load all packages found recursively")
var pkgmax = flag.Int("pkgmax", 0, "maximum number of packages to load")

// From: https://golang.org/pkg/runtime/pprof/#hdr-Profiling_a_Go_program
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	// From: https://golang.org/pkg/runtime/pprof/#hdr-Profiling_a_Go_program
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", "./...")
	cmd.Dir = *pkgdir

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(stdoutStderr))
		log.Fatalf("'go list' failed: %+v\n", err)
	}

	for n, dir := range strings.Split(strings.TrimSpace(string(stdoutStderr)), "\n") {
		fmt.Printf("loading [%s]\n", dir)

		cfg := packages.Config{
			Dir:   *pkgdir,
			Mode:  packages.LoadSyntax, // align with observed case
			Tests: true,                // align with observed case
		}

		_, loadErr := packages.Load(&cfg, dir)
		if loadErr != nil {
			log.Fatalf("packages.Load of [%s] failed: %+v\n", dir, loadErr)
		}

		if *pkgmax > 0 && n+1 == *pkgmax {
			fmt.Printf("stopping after -pkgmax [%d]\n", *pkgmax)
			break
		}
	}

	// From: https://golang.org/pkg/runtime/pprof/#hdr-Profiling_a_Go_program
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
