# syntax=docker/dockerfile:1
FROM golang:1.18-alpine as build_ua_contracts_cli

RUN apk update && apk upgrade && apk add --no-cache git openssh make

#Pass the content of the private key into the container
ARG SSH_PRIVATE_KEY
RUN mkdir /root/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa

#Github requires a private key with strict permission settings
RUN chmod 600 /root/.ssh/id_rsa

#Add Github to known hosts
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

WORKDIR /Code
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"
RUN git clone git@github.com:canonical/ua-contracts.git

WORKDIR /Code/ua-contracts
RUN make install


FROM golang:1.18-alpine AS base

COPY --from=build_ua_contracts_cli /go/bin/contract /go/bin/contract

WORKDIR /Code/ua-qa-go

COPY go.mod go.sum main.go cmd ./
RUN go mod download

COPY main.go main.go 
COPY cmd cmd
COPY state state

RUN go build -o ua-qa-go main.go


FROM alpine:latest

ENV PATH=$PATH:/go/bin
ENV CONTRACTS_URL="https://contracts.staging.canonical.com/"

COPY --from=build_ua_contracts_cli /go/bin/contract /go/bin/contract
COPY --from=base /Code/ua-qa-go/ua-qa-go /go/bin/ua-qa-go

CMD [ "ua-qa-go" ]
