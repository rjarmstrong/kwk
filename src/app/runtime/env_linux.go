package runtime


const defaultEnv = `kwkenv: "1"
editors:
#  Specify one app for each file type to edit.
#  sh: [vim]
#  go: [emacs]
#  py: [nano]
#  url: [vim]
  default: ["vim"]
apps:
  vim: ["vi", "$FULL_NAME"]
  emacs: ["emacs", "$FULL_NAME" ]
  nano: ["nano", "$FULL_NAME" ]
  default: ["vi", "$FULL_NAME"]
runners:
  jl: ["julia", "-e", "$SNIP"]
  sh: ["/bin/bash", "-c", "$SNIP"]
  url: ["firefox", "--new-tab", "$SNIP"]
  url-covert: ["firefox", "--private-window", "$SNIP"]
  js: ["node", "-e", "$SNIP"] #nodejs
  py: ["python", "-c", "$SNIP"] #python
  php: ["php", "-r", "$SNIP"] #php
  scpt: ["osascript", "-e", "$SNIP"] #applescript
  applescript: ["osascript", "-e", "$SNIP"] #applescript
  rb: ["ruby", "-e", "$SNIP"] #ruby
  pl: ["perl", "-E", "$SNIP" ] #perl
  exs: ["elixir", "-e", "$SNIP"] # elixir
  java:
    compile: ["javac", "$FULL_NAME"]
    run: ["java", "$CLASS_NAME"]
  scala:
    compile: ["scalac", "-d", "$DIR", "$FULL_NAME"]
    run: ["scala", "$NAME"]
  go: #golang
    run: ["go", "run", "$FULL_NAME"]
  rs: #rust
    compile: ["rustc", "-o", "$NAME", "$FULL_NAME"]
    run: ["$NAME"]
  cpp: # c++
    compile: ["g++", "$FULL_NAME", "-o", "$FULL_NAME.out" ]
    run: ["$FULL_NAME.out"]
  path: ["echo", "$SNIP" ]
  xml: ["echo", "$SNIP"]
  json: ["echo", "$SNIP"]
  yml: ["echo", "$SNIP"]
  default: ["echo", "$SNIP"]
security: #https://gist.github.com/pmarreck/5388643
  encrypt: []
  decrypt: []
  sign: []
  verify: []`