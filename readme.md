### Friend Management
# Run MongoDB
```
docker run --rm -it --name mongodb \
	-e MONGO_INITDB_ROOT_USERNAME=root \
	-e MONGO_INITDB_ROOT_PASSWORD=password \
	-p 27017:27017 \
	mongo:4
```

# Config
```
# copy and update config.yaml accordingly 
cp config.sample.yaml config.yaml
```

# Run Server
```
go run ./cmd/httpd/main.go
```