# Ubuntu just works
FROM alpine:3.2
MAINTAINER SOON_ <dorks@thisissoon.com>

## Environment Variables
ENV GOPATH /deepmind
ENV GOBIN /usr/local/bin
ENV PATH $PATH:$GOPATH/bin

# OS Dependencies
RUN echo 'http://dl-4.alpinelinux.org/alpine/edge/main' >> /etc/apk/repositories \
    && apk update && apk add go go-tools git build-base ca-certificates make bash && rm -rf /var/cache/apk/*

# Set working Directory
WORKDIR /deepmind

# GPM (Go Package Manager)
RUN git clone https://github.com/pote/gpm.git \
    && cd gpm \
    && git checkout v1.3.2 \
    && ./configure \
    && make install

# Set our final working dir to be where the source code lives
WORKDIR /deepmind/src/github.com/thisissoon/fm-deepmind

# Copy source code into the deepmind src directory so Go can build the package
COPY ./ /deepmind/src/github.com/thisissoon/fm-deepmind
RUN gpm install

# Install the go package
RUN go install main.go

# Set the default entrypoint to be deepmind
ENTRYPOINT main
