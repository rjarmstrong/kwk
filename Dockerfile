FROM golang:1.8.1

RUN apt-get update; apt-get install tree -y; apt-get install zip -y

RUN curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip" \
    && unzip awscli-bundle.zip \
    && ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws

RUN mkdir -p $GOPATH/src/github.com/kwk-super-snippets/cli && mkdir -p /builds

WORKDIR $GOPATH/src/github.com/kwk-super-snippets/cli

VOLUME /builds

CMD ["/bin/bash"]