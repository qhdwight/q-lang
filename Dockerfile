FROM ubuntu:latest

RUN apt-get -y update && apt-get -y upgrade && apt-get install -y build-essential

RUN mkdir -p /qlang/
COPY q-linux program.qq /qlang/

WORKDIR /qlang
