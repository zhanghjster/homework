FROM ubuntu:14.04

RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
        git \
        curl \
        wget \
        telnet \
        ca-certificates \
        software-properties-common \
        python3-software-properties \
        python-software-properties \
	&& rm -rf /var/lib/apt/lists/*

RUN add-apt-repository ppa:jonathonf/vim 

RUN apt-get update && apt-get install -y vim

ENV GOLANG_VERSION 1.10.1

RUN git clone --depth=1 https://github.com/amix/vimrc.git ~/.vim_runtime && \
    sh ~/.vim_runtime/install_awesome_vimrc.sh

RUN wget -O go.tgz https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go.tgz && \
    rm go.tgz

ENV GOPATH /go

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR $GOPATH

