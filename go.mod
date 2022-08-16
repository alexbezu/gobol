module github.com/alexbezu/gobol

go 1.15

require (
	github.com/alexbezu/gobol/cmd v0.0.0-20220806084713-011edb07e99f
	github.com/codenotary/immudb v1.3.1
	github.com/go-redis/redis/v9 v9.0.0-beta.2
	google.golang.org/grpc v1.39.0
)

replace github.com/alexbezu/gobol/cmd => ./cmd

replace github.com/codenotary/immudb => github.com/alexbezu/immudb v1.2.3-0.20220728165515-6ceb4f57a4af
