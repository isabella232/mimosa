FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y wget
RUN apt-get install -y golang
RUN wget https://apt.puppet.com/puppet-tools-release-bionic.deb
RUN dpkg -i puppet-tools-release-bionic.deb
RUN apt-get update
RUN apt-get install puppet-bolt

RUN mkdir /mimosa
COPY . /mimosa
WORKDIR /mimosa

CMD [ "go", "run", "bolt.go" ]
