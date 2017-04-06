FROM golang:1.7

SHELL ["/bin/bash", "-c"]

RUN apt-get -q update \
    && apt-get -q install -y \
        git \
        build-essential \
        libtool \
        pkg-config \
        autotools-dev \
        autoconf \
        automake \
        cmake \
        uuid-dev \
        libpcre3-dev \
        valgrind

ENV LIBSODIUM_VERSION=1.0.12 \
    LIBZMQ_VERSION=v4.2.2 \
    CZMQ_VERSION=v4.0.2 \
    ZYRE_VERSION=v2.0.0

RUN declare -A _deps=( \
        ["jedisct1/libsodium"]=${LIBSODIUM_VERSION} \
        ["zeromq/libzmq"]=${LIBZMQ_VERSION} \
        ["zeromq/czmq"]=${CZMQ_VERSION} \
        ["zeromq/zyre"]=${ZYRE_VERSION} \
    ) \
    && for repo in "${!_deps[@]}"; do git clone --depth=1 --branch="${_deps[$repo]}" "https://github.com/$repo.git" "/tmp/$repo"; done \
    && cd /tmp/jedisct1/libsodium && ./autogen.sh && ./configure && make install && ldconfig && cd - \
    && cd /tmp/zeromq/libzmq && ./autogen.sh && ./configure --with-libsodium && make install && ldconfig && cd - \
    && cd /tmp/zeromq/czmq && ./autogen.sh && ./configure && make install && ldconfig && cd - \
    && cd /tmp/zeromq/zyre && ./autogen.sh && ./configure && make install && ldconfig && cd - \
    && rm -rf /tmp/*
