module github.com/edwardofclt/newrelic-go-agent/v3/integrations/nrmysql

// 1.10 is the Go version in mysql's go.mod
go 1.10

require (
	// v1.5.0 is the first mysql version to support gomod
	github.com/go-sql-driver/mysql v1.5.0
	// v3.3.0 includes the new location of ParseQuery
	github.com/edwardofclt/newrelic-go-agent/v3 v3.3.0
)
