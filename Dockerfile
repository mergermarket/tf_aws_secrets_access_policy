FROM golang

ARG TERRAFORM_VERSION=0.11.11

RUN apt-get update && apt-get install -y curl unzip && \
    curl -O "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip" && \
    unzip -d /usr/local/bin/ "terraform_${TERRAFORM_VERSION}_linux_amd64.zip"

COPY go.mod go.sum /var/cache/

RUN cd /var/cache && \
    go mod download
