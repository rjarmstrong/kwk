FROM golang:1.7.4

RUN apt-get update; apt-get install tree -y

RUN mkdir -p $GOPATH/src/bitbucket.com/sharingmachine/kwkcli && mkdir -p /builds

COPY . $GOPATH/src/bitbucket.com/sharingmachine/kwkcli/
WORKDIR $GOPATH/src/bitbucket.com/sharingmachine/kwkcli/

RUN ./build.sh

VOLUME /builds

CMD ["/bin/bash"]