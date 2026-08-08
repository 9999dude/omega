.PHONY: build test healthy bad load-healthy load-cpu clean

build:
	./scripts/build.sh

test:
	go test ./healthy-service/... ./bad-service/... ./loadgen/...

healthy: build
	./bin/healthy-service

bad: build
	./bin/bad-service

load-healthy: build
	./bin/loadgen -url 'http://127.0.0.1:8080/work?size=2000' -duration 15s -concurrency 8

load-cpu: build
	./bin/loadgen -url 'http://127.0.0.1:8081/cpu?ms=100' -duration 15s -concurrency 8

clean:
	rm -rf ./bin ./profiles/*.pb.gz ./profiles/*.trace ./profiles/*.svg ./profiles/*.pdf ./profiles/*.png
