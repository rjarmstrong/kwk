---
date: 2016-03-08T21:07:13+01:00
title: Creating Apps
weight: 0
---

## Embedding Snippets

kwk 'run' commands can be embedded in snippets and so this opens up a whole world of opportunity to create interesting apps. Snippet references can also be dynamically passed to apps e.g.` /richard/apps/countdown.sh 3 'kwk tmap'` will countdown and run tmap after 3 seconds.

Caveats:
1. Apps only allow 3 levels.
2. When referencing a snippet it should be an absolute name e.g. `kwk run /richard/examples/hello.js`.
3. Apps should have the shell .sh extension, although it is of course possible to invoke the shell from your favorite programming language.
4. There are several restrictions in .sh files e.g. you can't use the 'rm' command.
5. .sh file are run in pseudo-safe mode but you can use 'eval' to dynamically run other kwk commands (see /kwk/apps/countdown.sh )

