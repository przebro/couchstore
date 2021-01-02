dockerbuid:
	docker build -t databazaar/couchdbdrv -f ./docker/Dockerfile .
start:
	docker run -d -l couchbzr1 -p5300:5984 -p6300:6984 -v ${PWD}/docker/etc:/opt/couchdb/etc/local.d databazaar/couchdbdrv
stop:
	docker rm -f $$( docker ps -qaf "label=couchbzr1")
tests:
	go test ./... -covermode=count --coverprofile='coverage.out'
	go tool cover -html coverage.out 