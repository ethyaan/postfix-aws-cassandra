build:
	cd src/ && GOOS=linux GOARCH=amd64 go build -o postfix-aws-cassandra main.go

build-arm:
	GOOS=darwin GOARCH=arm64 go build -o postfix-aws-cassandra src/main.go

build-docker:
	docker --debug build -t postfix-aws-cassandra . --progress=plain --platform linux/amd64

extract-rpm: build-docker
	@docker create --name extract-container postfix-aws-cassandra
	@docker cp extract-container:/home/builder/rpmbuild/RPMS/x86_64/postfix-aws-cassandra-1.0.0-1.amzn2023.x86_64.rpm .
	@docker rm extract-container