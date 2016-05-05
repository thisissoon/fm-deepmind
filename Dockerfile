# Ubuntu just works
FROM alpine:3.3
MAINTAINER SOON_ <dorks@thisissoon.com>

## Environment Variables
ENV GOPATH /deepmind
ENV GOBIN /usr/local/bin
ENV PATH $PATH:$GOPATH/bin

# OS Dependencies
RUN apk update && apk add go \
        tzdata \
        go-tools \
        git \
        ca-certificates \
        make \
        bash \
        wget \
    && rm -rf /var/cache/apk/*

# Set working Directory
WORKDIR /deepmind

# GPM (Go Package Manager)
RUN wget https://raw.githubusercontent.com/pote/gpm/v1.4.0/bin/gpm \
        && chmod +x gpm \
        && mv gpm /usr/local/bin

# Set our final working dir to be where the source code lives
WORKDIR /deepmind/src/github.com/thisissoon/fm-deepmind

# Copy source code into the deepmind src directory so Go can build the package
COPY ./ /deepmind/src/github.com/thisissoon/fm-deepmind
RUN gpm install

# Install the go package
RUN go install main.go

# Set the default entrypoint to be deepmind
ENTRYPOINT main
