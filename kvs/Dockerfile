#Stage 1: Compile the binary in a containerized Golang environment 
FROM golang:1.15 as build

# Copy the source files from the host
COPY . /go/src/github.com/cloud-native-go/kvs

# Set the working directory to the same place we copied the code
WORKDIR /go/src/github.com/cloud-native-go/kvs

# Download the dependency code
RUN go get github.com/gorilla/mux
RUN go get github.com/lib/pq

# Build the Binary!
RUN CGO_ENABLED=0 GOOS=linux go build -o kvs

#Stage 2: Build the Key-Value store image proper

# We use a "scratch" image, which contains no distribution files.
FROM scratch as image

# Copy the existing binary from host
COPY  --from=build /go/src/github.com/cloud-native-go/kvs .

# Tell Docker we'll be using port 8080
EXPOSE 8080

# Tell Docker to execute this command on a "docker run"
CMD ["/kvs"]