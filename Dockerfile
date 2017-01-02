FROM golang:1.7.4

RUN apt-get update; apt-get install tree -y; apt-get install zip -y

RUN mkdir -p $GOPATH/src/bitbucket.com/sharingmachine/kwkcli && mkdir -p /builds

RUN curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip" \
    && unzip awscli-bundle.zip \
    && ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws

COPY . $GOPATH/src/bitbucket.com/sharingmachine/kwkcli/
WORKDIR $GOPATH/src/bitbucket.com/sharingmachine/kwkcli/

RUN ./build.sh

VOLUME /builds

CMD ["/bin/bash"]