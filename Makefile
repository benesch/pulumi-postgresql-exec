VERSION ?= $(patsubst v%,%,$(shell git describe))

bin/pulumi-sdkgen-postgresql-exec: cmd/pulumi-sdkgen-postgresql-exec/*.go
	go build -o bin/pulumi-sdkgen-postgresql-exec ./cmd/pulumi-sdkgen-postgresql-exec

python-sdk: bin/pulumi-sdkgen-postgresql-exec
	rm -rf sdk
	bin/pulumi-sdkgen-postgresql-exec $(VERSION)
	cp README.md sdk/python/
	cd sdk/python/ && \
		sed -i.bak -e "s/\$${VERSION}/$(VERSION)/g" -e "s/\$${PLUGIN_VERSION}/$(VERSION)/g" setup.py && \
		rm setup.py.bak
