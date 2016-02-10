
Getting started
---------------
 - Copy the contents of `../../source` except `simulator.con` into the `driver` folder.
 - Copy `../../source/simulator.con` into this folder.
 - Compile `simelev.a` in the `driver` folder: `dmd sim_backend.d timer_event.d -lib -ofsimelev`
 - Run `rdmd driver/sim_frontend.d` from this folder.
 - Run `go run main.go` from this folder.
 - The simulator frontend should show a light on the command button on the bottom floor, and a print count of 1.
 
 
Things to note
--------------
 - CGO is weird. The comment above `import "C"` in `driver/io.go` is not a comment, and there must be no additional newline between the "comment" and the `import "C"` line.
 - `${SRCDIR}` is a Go 1.5 feature.
