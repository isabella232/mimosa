FROM golang:buster AS builder

ENV NONINTERACTIVE=true
ENV GOOS=linux

# eneables go modules
ENV GO111MODULE=on

RUN apt-get update && apt-get -y install wget apt-utils
RUN wget https://apt.puppet.com/puppet-tools-release-buster.deb
RUN dpkg -i puppet-tools-release-buster.deb
RUN apt-get update && apt-get -y install puppet-bolt

WORKDIR  /go/src/github.com/puppetlabs/mimosa/iot/bolt
COPY . .
RUN go mod download
RUN go build -o bolt_iot

FROM gcr.io/reusabolt/reusabolt-preview

RUN rm -rf /root/reusabolt/Boltdir/*
COPY Boltdir /root/reusabolt/Boltdir
RUN bolt puppetfile install
RUN bolt task show

WORKDIR /app
COPY wrapper_script.sh /app
COPY roots.pem /app
COPY key.pem /app
RUN chmod +x wrapper_script.sh
COPY --from=builder /go/src/github.com/puppetlabs/mimosa/iot/bolt/bolt_iot /app
RUN chmod +x bolt_iot
RUN pwd
RUN ls
CMD ["/app/wrapper_script.sh"]
