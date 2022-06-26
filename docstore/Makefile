
dev: docker-run-mongo-v3 docker-run-mongo-v4

docker-run-mongo-v3:
	@echo "Starting MongoDB v4 listening on 27017..."
	docker rm -f mongo4
	docker run -d --name mongo4  -p 27017:27017 mongo:4 
	@echo "...done. Run \"docker rm -f mongo4\" to clean up the container."

docker-run-mongo-v4:
	@echo "Starting MongoDB v3 listening on 27018..."
	docker rm -f mongo3 
	docker run -d --name mongo3  -p 27018:27017 mongo:3
	@echo "...done. Run \"docker rm -f mongo3\" to clean up the container."

clean:
	docker rm -f mongo3
	docker rm -f mongo4
