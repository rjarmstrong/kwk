# kwk - the snippet manager for developers

kwk was built because I felt that there was no snippet manager which really solved my needs. I felt snippets should not be limited to code since, as a developer , there are so many little pieces of information that I needed to store. Also I felt it was very annoying to have to open a browser and find a link, I wanted to get the snippet quickly. 

Since the introduction of git by Linus Torvalds and the rise of nodejs into the mainstream the command line has become more friendly to the average developer who has benefited from what was reserved for the arcane world of hardcore linux developers. Thats why I felt that the time had come to build a snippet manager designed specifically for the command line but targeted at all developers.

## Key Features

1. A UI for the cli which actually makes your snippets discoverable and easy to maintain.
2. Snippet names are unique and smartly addressable to save your keystrokes.
3. Executable snippets for virtually any language!
4. Organise your snippets by project or subject in pouches.
5. Search globally or just your snippets
6. Create 'apps' by composing multiple snippets, even in different languages.
7. Create snippets on the commandline and edit with your favorite editor.
8. Auto-updating
9. Globally redundant
10. Public 'open source' snippets and Private

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
`kwk new [<name>.<ext>] [snippet]`   
`kwk new [snippet] [<name>[.<ext>]]`  
`cat <file> | kwk new [<name>.<ext>]`  
### Navigating
` kwk `  
` kwk <pouch name> `  
` kwk <snip name>[.<ext>]`  
` kwk find <term>`  
` kwk expand [<pouch name>]`  
` kwk <partial snippet/pouch name> RETURN`
### Running a snippet
`kwk run <name>[<.ext>]`  
### Pouch management
`kwk <pouch name>`  
`kwk mkdir <pouch name>`  
`kwk rm <pouch name>`  
`kwk lock <pouch name>`  
`kwk unlock <pouch name>`  
`kwk mv <pouch name> <new pouch name>`
`kwk rm <pouch name>`
### Snippet management
`kwk edit <name>[.ext]`  
`kwk describe <name>[.ext]`  
`kwk clone [username/][pouch/]<name>[.ext]`  
`kwk tag <name>[.ext] tag ...`  
`kwk patch <name>[.ext] <text to match> <replacement>`  
`kwk mv <pouch>/<name>.<ext> ... <new pouch>`
`kwk mv <name>.<ext> <newname>[.<ext>]`  
`kwk rm <name>.<ext>`

### Social and Discovery
`kwk share <name>.<ext>`  
`kwk popular`  
`kwk stale` 

### Account Management
`kwk signup`  
`kwk signin`  
`kwk signout`  
`kwk reset-password`  
`kwk change-password`  
`kwk change-email`  (coming soon)

### App Management
`kwk version`  
`kwk help <command>`  
`kwk update`  

### Full Docs


### Examples


### FAQ


### Security


### Contributing


### Built with GO


Road Map
--