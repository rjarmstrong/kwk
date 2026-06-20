---
date: 2016-03-08T21:07:13+01:00
title: Get started
weight: 0
---

## Install
``` sh
npm install -g kwkcli

# then

kwk signup
```

or if you prefer to install without npm:
``` sh
curl  https://s3.amazonaws.com/kwk-cli/latest/bin/kwk-linux-amd64.tar.gz -o
// check checksum
// unzip
// symlink to /usr/local/bin

kwk signup
```

## Hello world

``` sh
kwk new "echo 'Hello, world!'" hello.sh

# Creates a snippet

kwk run hello

# Executes it and prints: Hello, world!
```
