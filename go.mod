module gobol

go 1.15

require (
	github.com/codenotary/immudb v1.3.1
	github.com/alexbezu/gobol/cmd v0.0.0
	google.golang.org/grpc v1.39.0
)

replace github.com/alexbezu/gobol/cmd => ./src/cmd
replace github.com/alexbezu/gobol => ./

replace github.com/codenotary/immudb => github.com/alexbezu/immudb v1.2.3-0.20220728165515-6ceb4f57a4af
