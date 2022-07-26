module gobol

go 1.13

require (
	github.com/codenotary/immudb v1.3.0
	gobol/src/cmd v0.0.0
	google.golang.org/grpc v1.39.0
)

// https://stackoverflow.com/questions/53682247/how-to-point-go-module-dependency-in-go-mod-to-a-latest-commit-in-a-repo
replace gobol/src/cmd => ./src/cmd
replace github.com/codenotary/immudb => github.com/alexbezu/immudb v1.2.3-0.20220523142458-1f3f6469b8b2
