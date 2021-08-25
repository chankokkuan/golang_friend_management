### Friend Management
# Install go dependency
```
go mod download
```

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

## To access the APIs which required JWT Token
Go to https://jwt.io/, choose HS256 Algorithm, input the payload with your claims and replace your jwt secret key with `your-256-bit-secret`

Claims Example:
```
{
	"user_id" "497f6eca-6276-4993-bfeb-53cbbbba6f08",
	"name": "John",
	"email": "john@example.com",
	"roles": ["User"]
}
```

After you get the token, add the token into the `Authorization: Bearer ${token}` in your http request header