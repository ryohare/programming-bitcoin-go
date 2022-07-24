FROM circleci/golang as builder

WORKDIR /build

COPY . .

# build the package
RUN go get -v ./...
RUN go build main.go
