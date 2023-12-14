# nats-streaming
NATS_IMAGE_NAME=nats-streaming
NATS_CONTAINER_NAME=nats-streaming-trivial
# postgres
DB_IMAGE_NAME=postgres-trivial
DB_CONTAINER_NAME=postgres-trivial


build:
	docker build -t $(DB_IMAGE_NAME) .
	docker pull $(NATS_IMAGE_NAME)

run:
	docker run -d -p 5432:5432 --name $(DB_CONTAINER_NAME) $(DB_IMAGE_NAME)
	docker run -d -p 4222:4222 -p 8222:8222 -e STAN_CLUSTER_ID=L0-cluster-id -e STAN_CLIENT_ID=L0-client-id -e STAN_SUBJECT=L0-subject --name $(NATS_CONTAINER_NAME) $(NATS_IMAGE_NAME)

start:
	docker start $(NATS_CONTAINER_NAME)
	docker start $(DB_CONTAINER_NAME)

stop:
	docker stop $(DB_CONTAINER_NAME)
	docker stop $(NATS_CONTAINER_NAME)

clean: stop
	docker rm $(DB_CONTAINER_NAME) $(NATS_CONTAINER_NAME)
	docker rmi $(DB_CONTAINER_NAME)
	docker rmi $(NATS_CONTAINER_NAME)


