build:
	sudo podman build . -t postgres-platzi-go

run:
	sudo podman run -p 5432:5432 -v /home/lelouch/data_docker/pgsql:/var/lib/postgresql/data --name postgres-platzi -d postgres-platzi-go

all: build run
