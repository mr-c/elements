FROM antha/antha
RUN apt-get update && apt-get install -y \
    libglpk-dev \
 && rm -rf /var/lib/apt/lists/*
ADD . /go/src/github.com/antha-lang/elements
RUN go get github.com/antha-lang/elements/cmd/...
