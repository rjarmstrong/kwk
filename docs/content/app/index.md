---
date: 2016-03-08T21:07:13+01:00
title: App Commands & Settings
weight: 0
---

## Your Environment

kwk relies on your local environment since it needs to have the appropriate runtimes installed (e.g. node js, ruby etc) to execute snippets and it needs editors to edit the snippets.
 
 Normally you don't need to install anything since kwk has been purposefully designed to work *around your environment*, that being said you may want to run some snippets in languages you don't normally use so you'll need to go to their website to grab the installers. 
 
## Tweaking Editors and Runtimes
For each machine you have kwk installed on you can configure the environment to behave in certain ways. e.g. For a build server you may want to edit snippets with <a href="http://www.vim.org/download.php" target="_blank">Vim</a>, for your work computer <a href="https://www.jetbrains.com/webstorm/" target="_blank">Webstorm</a> and your home computer <a href="https://www.sublimetext.com/">Sublime</a>. But that's not all! You can also specify certain file types be edited with certain editors. Explore the `env` file for details:

### Snippet Editors

``` sh 
kwk edit env

>>>

kwkenv: "1"
editors:
#  Specify one app for each file type to edit.
  sh: [vim]
  go: [gogland]
  py: [vscode]
  url: [textedit]
  default: ["textedit"]
apps:
  webstorm: ["open", "-a", "webstorm", "$DIR"]
  textedit: ["open", "-e", "$FULL_NAME"]
  vscode: ["open", "-a", "Visual Studio Code", "$DIR"]
  vim: ["vi", "$FULL_NAME" ]
  emacs: ["emacs", "$FULL_NAME" ]
  nano: ["nano", "$FULL_NAME" ]
  default: ["open", "-t", "$FULL_NAME"]
  ...
```

### Snippet Runners

kwk automatically creates configurations for Linux, Mac and Windows - however you can edit these if you need to run a non-default command for any particular file extension.

``` sh
runners:
  sh: ["/bin/bash", "-c", "$SNIP"]
  url: ["open", "$SNIP"]
  url-covert: ["/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", "$SNIP"]
  js: ["node", "-e", "$SNIP"] #nodejs
  py: ["python", "-c", "$SNIP"] #python
  php: ["php", "-r", "$SNIP"] #php
  scpt: ["osascript", "-e", "$SNIP"] #applescript
  applescript: ["osascript", "-e", "$SNIP"] #applescript
  rb: ["ruby", "-e", "$SNIP"] #ruby
  pl: ["perl", "-E", "$SNIP" ] #perl
...
```

## Preferences

Preferences in kwk are local in that a default preference file is installed on the host computer which you can edit. For any new installation on kwk, whether reinstalling on the same machine or a new one, the last edited preferences will be installed.

``` sh
kwk edit prefs

>>>

kwkprefs: v1
preferences:
  covert: false     # By default always open bookmarks in covert mode
  listall: true     # List all pouches in the root including locked
  debug: false      # Print detailed command output
  regulateupdates: true # When false will get updates as soon as they are available
  slimlines: 3      # When listing the number of snippet lines to show
  expandedlines: 15 # When expanding the number of snippet lines to show
  alwaysexpandlists: false # Always expand lines when listing
  ...
```

### On screen help

```sh 
kwk help [<command>]
```  

### Force update kwk and check version
``` sh
kwk update

kwk version 
```  

