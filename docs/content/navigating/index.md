---
date: 2016-03-08T21:07:13+01:00
title: Finding Snippets
weight: 0
---

## Saving Keystrokes
kwk has been designed specifically to save keystrokes. That's why when typing snippet names we often try to guess which snippet you intended to use.
 
``` sh
# Whew!
kwk view /cjones/apps/hello.js

# Better
kwk view hello.js

# Even Better
kwk view hello

# Best!
kwk hello

# Use Suggest if you can't remember the name
kwk hel


# kwk [view] [/<username>/][pouch/]<name>[.<ext>]
``` 
> Note: The `view` command is optional. You can configure the default command for different file extensions by typing `kwk edit env`

## Canonical Snippet names
> "If it ends with a slash, it's a directory. If not, it's a file."
<a href="http://stackoverflow.com/questions/980255/should-a-directory-path-variable-end-with-a-trailing-slash" target="_blank">Dave Sherohman</a>

The kwk version is: "If it ends with a slash, it's a username or pouch. If not, it's a snippet."

## Listing made easy
 
No need to keep typing 'ls' to list things, we actually removed that after realising how many times we typed it, so now you can list all your pouches or snippets with zero or one words, either `kwk` or `kwk <pouch name>` is enough.

## Show your Homepage 
Includes your pouches and unsorted snippets:  
``` sh 
kwk 
```  

## Look inside a Pouch
 
``` sh
kwk apps

# kwk <pouch name> 
```

## Searching  

``` sh
kwk find worl

# kwk find [-g] <term>
# -g/--global = search everybody's public snippets.
```

## Expanding Lists
To show more details in list, number of lines to expand can be changed by calling `kwk edit prefs`.  
``` sh 
kwk expand apps

# kwk expand [<pouch name>]
```

## Suggesting 
``` sh 
kwk hell #RETURN
> hello.js|
```  

## Enchilada
**Only use this if you really need to list all your snippets**. Might be slower and lots of screen use.
``` sh 
kwk enchilada
```  

## Special lists
To list snippets that have been run recently:
```bash
kwk @today
kwk @week
kwk @month
kwk @old
```

## Community

If you use a forward slash `/` in front of anything kwk will interpret it as a Username, thus you can browse other peoples public snippets.

``` sh
kwk /alice/hawaii-trip
# Lists alice's hawaii pouch

kwk /alice/hawaii-trip/budget.txt
# Shows the details of the vacation budget

kwk /alice/hawaii-trip/hibiscus.url
# Opens the url for Hibiscus hotel

kwk run /alice/hawaii-trip/local-news.py
# Runs alice's news feed for Hawaii

# kwk /<username/[pouch name]/[snippet name]
```

## Popular Snippets
Show Alice's most recently used snippets

``` sh
kwk /alice pop

# kwk [/<username>] popular|pop
```

## Sharing Snippets

At the moment snippets are only shared by email.

``` sh
kwk share hello.js alice

# kwk share <name>.<ext> <email>
```

## Stale Snippets 

To help you clean up use this command to show snippets which haven't been used much or recently.

``` sh
kwk stale
```