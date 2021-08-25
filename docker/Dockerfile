# Test and build
FROM golang as builder

ARG USER=user
ARG WORK_DIR=/app

# Allows cross-compilation for the compiled Go app
ENV CGO_ENABLED 0

# update os and install required packages
RUN apt-get update -qq && apt-get upgrade -y

# Create directories for the user to avoid running as root
RUN groupadd ${USER} && useradd -m -g ${USER} -l ${USER}
RUN mkdir -p ${WORK_DIR} && chown -R ${USER}:${USER} ${WORK_DIR}

# set WORKDIR and USER
WORKDIR ${WORK_DIR}
USER ${USER}

# Copy the local package files to the container's workspace
COPY . ${WORK_DIR}

# verify the modules
RUN go mod verify

# Run tests
RUN go test -cover ./...

# Build app, produces ./app
RUN go build -o app ./cmd/httpd/main.go

# ----------------

# Produce docker image to run app
FROM alpine

ARG USER=versa_user
ARG WORK_DIR=/app
ARG PORT=8080

ENV GIN_MODE=release
ENV PORT=${PORT}

# Create directories for the user to avoid running as root
# https://github.com/mhart/alpine-node/issues/48
RUN addgroup -S ${USER} && adduser -S ${USER} -G ${USER}

# set WORKDIR
WORKDIR ${WORK_DIR}

# Copy the built file from the previous stage to the container's workspace
COPY --chown=0:0 --from=builder ${WORK_DIR}/app ${WORK_DIR}

# Install CA certificates to prevent x509: certificate signed by unknown authority errors
RUN apk update && \
    apk add ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

USER ${USER}
EXPOSE ${PORT}

CMD ["./app"]
