ARG GO_VERSION

FROM golang:$GO_VERSION
RUN apt-get -y update && apt-get -y install \
      libiw-dev \
      bison     \
      curl      \
      file      \
      git       \
      jq        \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
RUN bash -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /bin

WORKDIR /go/src/
COPY ./ ./
