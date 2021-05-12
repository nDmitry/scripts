.SILENT:

build-freespace:
	GOOS=linux GOARCH=amd64 go build -o ./bin/freespace cmd/freespace/main.go

build-pgdumpdoc:
	GOOS=linux GOARCH=amd64 go build -o ./bin/pgdumpdoc cmd/pgdumpdoc/main.go
