---
date: 2016-03-08T21:07:13+01:00
title: Snippets
weight: 0
---

## Store anything
Snippets are often thought of as code but kwk can be used to store any small piece of information provided its no longer than 2000 *normal characters (about 100 lines), just pick a file extension which is suitable e.g. **.json** for some JSON or **.rb** for a ruby file. 

URLs are well suited to be stored as snippets, in fact we have some great features for URLs  including running with covert mode, meta data previews, auto url detection and by default all you need to do to open a url is type `kwk my-bookmark`.

> ** Use '.url' as an extension for URLs   
> *  Normal characters being UTF-8 non-combined characters e.g. é can be 2 characters.

## Creating Snippets
A snippet is just a UTF-8 string so creating one is easy with the ` new ` or ` edit ` command:

``` sh
kwk new "console.log('hello');" hello.js

# kwk new [snippet] [<name>[.<ext>]]
```  
Note: *If you add a file extension to the name the order of arguments above does not matter. Also you don't need to provide an extension it will be given a .txt extension by default and if no name is given a random animal name will be used.*

> Snippet names are lowercase (they will be lowered automatically)

### Create and edit in one step:  
``` sh
kwk edit get-weather.py

# kwk edit <name>.[<ext>]
```

### Create a snippet from a local file:  
``` sh
cat ./list-open-ports.sh | kwk new open-ports.sh

# cat <file> | kwk new [<name>.<ext>]
```  

### Multi-line create with a here-doc: 
``` sh
kwk new <name>.<ext> <<eof
> package main
> 
> import "fmt"
> 
> func main() {
> 	fmt.Println("Hello, world")
> }
> eof
```

{{< note title="Newbie Tip:" >}}
Heredocs are really easy to use.  Just type: '<<' followed by any string you like to use as a terminator (in the example above it is 'eof'). Type or paste in multiple lines after this and then finish up by typing the terminator 'eof' again.
{{< /note >}}


## Running Snippets
By default most snippets will execute only when you use the `run` command - the exception is web bookmarks which will run immediately. To modify which extensions need the run command edit the prefs file: `kwk edit prefs`. 
``` sh
kwk run killPort3000

# kwk [run|r] <name>[<.ext>] [args...]
``` 

> Note: that in most cases bash arguments will resolve starting from $0 instead of $1 since the env.yml file uses the -c option to run snippets. The reason can be found in the bash manual:

"If arguments remain after option processing, and neither the -c nor the -s option has been supplied, the first argument is assumed to be the name of a file  containing  shell
        commands."

## Viewing Snippets
``` sh
kwk hello.js
# or
kwk view hello.js

# kwk [view] <name>[.ext]
```
### Get the raw snippet
``` sh
kwk raw hello.js

# kwk raw|cat <name>[.ext]
```
## Editing Snippets
For each machine you use kwk on you can change which editors are used for which file extensions. Use `kwk edit env` to change the mappings and add new editors.
``` sh
kwk edit hello.js

# kwk edit <name>[.ext]
```  
*Note:*  When editing a snippet, kwk will watch for changes to a temporary file on the local disk, if your editor has **auto-save enabled** you may need to disable it or type `kwk edit <name>[.<ext>]` in the cli again after you've finished your changes so that your changes are actually saved, we are working on a more elegant fix for this.

### Patching
For the most part you would use the `edit` command, however, `patch` can be useful if you want to update a snippet programatically, such as setting a version number.
``` sh
kwk patch build.json "^(\d+\.)?(\d+\.)?(\*|\d+)$" "2.3.10"

# kwk patch <name>[.ext] <regex> <replacement> 
```

## Descriptions
If you are like most dev's then remembering what a snippet does can be difficult, add a short description for list views to be more helpful.
``` sh
kwk describe hello.js "A script which says hello"

# kwk describe <name>[.ext] <description>
```
## Cloning Snippets
``` sh
kwk clone /rjarmstrong/pretty-json.go slick-json.go

# kwk clone [username/][pouch/]<name>[.ext] [new name][.<ext>]
```  

## Tagging
``` sh
kwk tag hello.js javascript hello-world sample

# kwk tag <name>[.ext] tag ...
# kwk untag <name>[.ext] tag ...
```

## Moving Snippets
You can move multiple snippets between pouches provided **the full relative path is  given** (except snippets in the root pouch can be referenced by name alone).  
> Note: if the target pouch does not exist then you may end up renaming the file just like on unix systems.
``` sh
kwk mv hello.js hello.go hello

# Moves snippets from the root to pouch 'hello'


kwk mv playground/hello.rs playground/hello.php hello

# Moves snippets from pouch 'playground' to 'hello'

# kwk mv <pouch>/<name>.<ext> ... <new pouch>
```

## Renaming Snippets
``` sh
kwk mv hello.js hello-world

# kwk mv <name>[.<ext>] <newname>[.<ext>]
```  
> Note: If you don't provide an extension for the new name its possible that you may move the snippet into a pouch if one exists with that same name.

## Deleting Snippets
``` sh
kwk rm hello.js hello.go

# kwk rm <name>.<ext> ...
```


## Flags
<pre>
-y --auto-yes     Automatically enter 'yes' for prompts.  
-n --naked        When listing don't add styles   
-c --covert       When running a url, open it in covert mode  
-d --debug        Print output of debug  
-a --all          List all pouches including private  
-g --global       When searching include other peoples snippets
</pre>

