module soloos/sdfs

go 1.12

require (
	github.com/google/flatbuffers v1.11.0
	github.com/stretchr/testify v1.3.0
	soloos/common v0.0.0
	soloos/sdbone v0.0.0
	soloos/soloboat v0.0.0
)

replace (
	soloos/common v0.0.0 => /soloos/common
	soloos/sdbone v0.0.0 => /soloos/sdbone
	soloos/soloboat v0.0.0 => /soloos/soloboat
)
