# FROM is the go image we'll be building
# alpine is a lightweight go image made for docker
FROM golang:1.19-alpine

# WORKDIR is the directory in our VM that will hold
# all of our files in the docker image
WORKDIR /app

# copy the mod and sum so we can see what dependencies we need
COPY go.mod ./
COPY go.sum ./

# RUN does what it says on the tin
# download dependencies
RUN go mod download

# the docker docs recommend specifically copying which file types
# you need for the image, like COPY *.go ./, but for this example
# I don't think there's any need to get so granular
COPY . ./

# build the go image
RUN go build -o /stream/stream.exe cmd/stream/main.go 

# I think this is a sort of internal port forwarding
EXPOSE 8080

# Honestly I don't really understand this syntax
CMD [ "/stream/stream.exe" ]