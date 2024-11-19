.SILENT:

build-freespace:
	GOOS=linux GOARCH=amd64 go build -o ./bin/freespace cmd/freespace/main.go

build-tempmon:
	GOOS=linux GOARCH=amd64 go build -o ./bin/tempmon cmd/tempmon/main.go
build-tempmon-arm:
	GOOS=linux GOARCH=arm go build -o ./bin/tempmon cmd/tempmon/main.go

build-pgdumpdoc:
	GOOS=linux GOARCH=amd64 go build -o ./bin/pgdumpdoc cmd/pgdumpdoc/main.go

build-mysqldump:
	GOOS=linux GOARCH=amd64 go build -o ./bin/mysqldump cmd/mysqldump/main.go

build-dirbackup:
	GOOS=linux GOARCH=amd64 go build -o ./bin/dirbackup cmd/dirbackup/main.go

build-ping:
	GOOS=linux GOARCH=amd64 go build -o ./bin/ping cmd/ping/main.go