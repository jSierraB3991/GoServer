build:
	sudo podman build . -t platzi-rest-ws-app
run:
	sudo podman run -d -p 5050:5050 --name platzi-rest-ws platzi-rest-ws-app
all: build run
