FROM golang:1.7.4

RUN mkdir -p $GOPATH/src/bitbucket.com/sharingmachine/kwkcli && mkdir -p /builds

COPY . $GOPATH/src/bitbucket.com/sharingmachine/kwkcli/
WORKDIR $GOPATH/src/bitbucket.com/sharingmachine/kwkcli/

RUN ./build.sh

CMD ["/bin/bash"]