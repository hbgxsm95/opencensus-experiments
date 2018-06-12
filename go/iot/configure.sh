#!/bin/bash
if [ "$#" -ne 3 ]; then
  echo "Illegal number of parameters"
  exit -1
fi

cd $GOPATH

go get -u github.com/hybridgroup/gobot

go get -u contrib.go.opencensus.io/exporter/stackdriver

go get -u go get -u go.opencensus.io

sudo apt-get install sshpass

sudo apt-get install gcc-arm-linux-gnueabihf

cd $GOPATH:/github.com/src/census-ecosystem/oepncensus-experiments/go/iot/

chmod u+x ./run.sh

sshpass -p $3 scp ./configurePi.sh $1@$2:~/
sshpass -p $3 ssh $1@$2 "chmod u+x ~/configurePi.sh"
sshpass -p $3 ssh $1@$2 ~/configurePi.sh
