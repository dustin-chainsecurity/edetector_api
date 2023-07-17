FROM        golang
RUN         mkdir -p /app
WORKDIR     /app
COPY        . .
RUN         go mod download
RUN         go build -o edetector_api
ENTRYPOINT  ["./edetector_api"]