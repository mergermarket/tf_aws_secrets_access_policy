FROM golang

ARG TERRAFORM_VERSION=0.11.7

RUN apt-get update && apt-get install -y curl unzip && \
    curl -O https://releases.hashicorp.com/terraform/0.11.7/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip -d /usr/local/bin/ terraform_${TERRAFORM_VERSION}_linux_amd64.zip

RUN go get -v github.com/hashicorp/terraform/terraform

COPY . .
