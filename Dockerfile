FROM golang:1.17.6-buster

ENV AKASH_VERSION="0.16.3"
ENV AKASH_NET="https://raw.githubusercontent.com/ovrclk/net/master/mainnet"
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -y curl
RUN apt-get install -y unzip
RUN apt-get install -y python
RUN curl https://raw.githubusercontent.com/ovrclk/akash/master/godownloader.sh | sh -s -- "v$AKASH_VERSION"
ENV AKASH_CHAIN_ID="akashnet-2"
ENV AKASH_NODE=http://akash-sentry01.skynetvalidators.com:26657
ENV AKASH_KEY_NAME=terraform
ENV AKASH_KEYRING_BACKEND=os

RUN