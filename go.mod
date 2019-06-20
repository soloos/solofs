module soloos/sdfs

go 1.12

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/google/flatbuffers v1.11.0
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/stretchr/testify v1.3.0
	google.golang.org/appengine v1.6.1 // indirect
	soloos/common v0.0.0
	soloos/sdbone v0.0.0
)

replace (
	soloos/common v0.0.0 => /soloos/common
	soloos/sdbone v0.0.0 => /soloos/sdbone
)
