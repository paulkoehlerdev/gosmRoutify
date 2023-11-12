all:
	# sync version with ../go.mod
	go install -v github.com/golang/protobuf/protoc-gen-go@v1.5.3

	rm *.go || true
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go_opt=Mfileformat.proto=github.com/paulkoehlerdev/pkg/libraries/osmpbf/osmproto \
		--go_opt=Mosmformat.proto=github.com/paulkoehlerdev/pkg/libraries/osmpbf/osmproto \
		*.proto