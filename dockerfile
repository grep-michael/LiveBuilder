FROM ghcr.io/goreleaser/goreleaser-cross:latest

WORKDIR /workspace
COPY . .

ENV CGO_ENABLED=1
ENV GOOS=darwin
ENV GOARCH=amd64
ENV CC=o64-clang
ENV CXX=o64-clang++

# Add comprehensive warning suppression
ENV CGO_CFLAGS="-w -Wno-error -Wno-format-nonliteral -Wno-unused-command-line-argument -Wno-format-security"
ENV CGO_CXXFLAGS="-w -Wno-error -Wno-format-nonliteral -Wno-unused-command-line-argument -Wno-format-security"
ENV CGO_LDFLAGS="-w"

RUN go mod download

# Try building with maximum error suppression
RUN go build -ldflags="-s -w" -o LiveBuilder 2>/dev/null || \
    CGO_CFLAGS="-w -Wno-everything" CGO_CXXFLAGS="-w -Wno-everything" go build -ldflags="-s -w" -o LiveBuilder
