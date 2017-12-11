FROM golang:1.9.2

# Install Git repository
RUN apt update
RUN apt-get install git

# Build application
RUN go get github.com/ryanbradynd05/go-tmdb
RUN go get github.com/araddon/dateparse
RUN go get github.com/davecgh/go-spew/spew
RUN go get github.com/ferhatelmas/levenshtein

# Copy R2D2 into container
COPY ./ /go/src/github.com/bnmcg/r2d2/

WORKDIR /go/src/github.com/bnmcg/r2d2
RUN go install

# Run the application
WORKDIR /

CMD ["/go/bin/r2d2", "/var/lib/r2d2-input", "/var/lib/r2d2-output/"]