# kwk - cli

<img src="https://travis-ci.org/kwk-super-snippets/cli.svg?branch=master" alt="travis build" />

kwk was built because I felt that there was no snippet manager which really solved my needs. I felt snippets should not be limited to code since, as a developer , there are so many little pieces of information that I needed to store. Also I felt it was very annoying to have to open a browser and find a link, I wanted to get the snippet quickly. 

Since the introduction of git by Linus Torvalds and the rise of nodejs into the mainstream the command line has become more friendly to the average developer who has benefited from what was reserved for the arcane world of hardcore linux developers. Thats why I felt that the time had come to build a snippet manager designed specifically for the command line but targeted at all developers.

## Key Features

1. A UI for the cli which actually makes your snippets discoverable and easy to maintain.
[image of pouch list] 
2. Snippet names are unique and smartly addressable to save you keystrokes.
[gif running without extension and pouch]
3. Executable snippets for virtually any language or data.
[show a list of code animals]
4. Organise your snippets by project or subject in pouches.
[show a list of pouches]
5. Search globally or just your snippets
[show search results]
6. Create 'apps' by composing multiple snippets, even in different languages.
[show the details of an app with multiple languages]
7. Create snippets on the commandline and edit with your favorite editor.
[show snippet open in web-storm]
8. Auto-updating
[show update notification]
9. Globally redundant
[list data centers]
10. Public 'open source' snippets and Private
[show gif of locking]

## Get Started

``` npm install -g kwkcli ```

or if you prefer to install without npm:

``` 
curl  https://s3.amazonaws.com/kwk-cli/latest/bin/kwk-linux-amd64.tar.gz -o
// check checksum
// unzip
// symlink to /usr/local/bin
```

### Creating a snippet
A snippet is just a UTF-8 string so creating one is easy with the new or edit command:

Create with the edit command:  
`kwk edit <name>.<ext>`

Create a snippet name first:  
`kwk new [<name>.<ext>] [snippet]` 
 
Create a snippet snippet first:  
`kwk new [snippet] [<name>[.<ext>]]`  

Create from a file:  
`cat <file> | kwk new [<name>.<ext>]`  

Multi-line create with a here-doc: 
```
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
Newbie tip: Heredocs are really easy to use simply type '<<' followed by any string you like to use as a terminator- in the example it is 'eof', type or paste in multiple lines after this and then finish up by typing the terminator 'eof' again.


### Navigating
Show your homepage, which includes your pouches and unsorted snippets:  
` kwk `  

Show contents of a pouch: 
 
` kwk <pouch name> `  

Show details of a snippet:
  
` kwk [view] <snip name>[.<ext>]`  

Search for snippets by term:  

` kwk find <term>`  

Expand the listing of snippets to show more details in list, settings can be changed by calling `kwk edit prefs`.  
` kwk expand [<pouch name>]`  

Suggest snippets when you can't remember the full name:  
` kwk <partial snippet/pouch name> RETURN`  

### Running a snippet

```
kwk [run|r] <name>[<.ext>]

e.g.  

kwk run killPort3000

``` 

By default most snippets will execute only when you use the 'run' command - the exception is web bookmarks which will run immediately. To modify which extensions need the run command edit the prefs file: `kwk edit prefs`. 

### Pouch management
`kwk <pouch name>`  
`kwk mkdir <pouch name>`  
`kwk rm <pouch name>`  
`kwk lock <pouch name>`  
`kwk unlock <pouch name>`  
`kwk mv <pouch name> <new pouch name>`
`kwk rm <pouch name>`
### Snippet management
`kwk view <name>[.ext]`
`kwk raw|cat <name>[.ext]`
`kwk edit <name>[.ext]`  
`kwk describe <name>[.ext] <description>`  
`kwk clone [username/][pouch/]<name>[.ext]`  
`kwk tag <name>[.ext] tag ...`  
`kwk patch <name>[.ext] <text to match> <replacement>`  
`kwk mv <pouch>/<name>.<ext> ... <new pouch>`  
`kwk mv <name>.<ext> <newname>[.<ext>]`  
`kwk rm <name>.<ext>`

### Social and Discovery
List other users' pouches and snippets:  

`kwk /<username/[pouch name]/[snippet name]`

List popular snippets:
```
kwk [/<username>] popular
```

Share a snippet with another user:  
```
kwk share <name>.<ext> <username|email>
```

List your snippets which have not been used recently:  
`kwk stale` 

### Account Management
`kwk signup`  
`kwk signin`  
`kwk signout`  
`kwk reset-password`  
`kwk change-password`  
`kwk change-email`  (coming soon)

### App Management
`kwk edit env`  
`kwk edit prefs`  
`kwk version`  
`kwk help <command>`  
`kwk update`  

### Flags
`-y` auto-yes  
`-n` naked  
`-c` covert  
`-d` debug  
`-a` list all


### Examples
The following examples can be found on kwk here: `kwk /examples`

Charting Weather  
Deleting multiple snippets  
News  
JsonFormatting
Multi-language

### FAQ



### Security

- Checksums
- TLS encryption

### Contributing

At the moment the cli is closed source, but please feel free to make suggestions here and I'll make an effort to solve any issues or consider suggestions for improvements as soon as I can.

### Built with GO

Initially I build the first version of kwk with just bash scripts but this soon evolved into a Go version because it felt like 100% the right runtime for the job - not least because setting up GRPC between client and server was a huge benefit over laboring through a standard ReST api.

### Road Map

kwk is in alpha release so expect things to break from time to time, however, releases are happening everyday, and with the auto-updater you should receive new features and fixes often. Lookout for enhancements to the UI and smarter responses from the api.
