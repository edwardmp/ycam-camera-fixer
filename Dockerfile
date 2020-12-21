FROM golang:1.14

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/edwardmp/ycam-camera-fixer

# Copy everything from the current directory to the PWD
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Run the executable
CMD ["ycam-camera-fixer"]