.DEFAULT: default

default: all

all: server-docker-cp

ensure-bin:
	mkdir -p ./bin

clean:
	rm -vrf ./bin


server-build:
	docker-compose -f ./dockerfiles/docker-compose.yml up --build server-build

server-run:
	docker-compose -f ./dockerfiles/docker-compose.yml up --build server-run

server-docker-cp: server-build ensure-bin
	docker cp sample-app-server-build:/usr/local/sample-app/bin/sample-app-server ./bin

server-run-local: server-docker-cp
	./bin/sample-app-server --log-disable-file

