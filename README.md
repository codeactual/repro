# Steps

- Clone this repository to a location outside `GOPATH`.
- Run the `repro` bash script. It should clone `x/tools`, run `repro.go` which consumes `10` packages in `x/tools`, then display the memory profile from `go tool pprof`.
