module github.com/deref/fsw

go 1.16

require (
	github.com/deref/pier v0.0.0-20210620044641-0f71544154e7
	github.com/fsnotify/fsevents v0.1.1
	github.com/google/uuid v1.3.0
)

// https://github.com/deref/httpie-go/tree/pier
replace github.com/nojima/httpie-go => github.com/deref/httpie-go v0.7.1-0.20210620034715-00ad2c785a86
