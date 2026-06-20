---
date: 2016-03-08T21:07:13+01:00
title: Organise with Pouches
weight: 0
---

## Stopping the mess
Pouches make categorizing snippets easy and provide a way to make a group of snippets private. Pouches have at most 24 snippets. Its important to note that a pouch is different from a folder in that snippet names still need to be unique in your entire snippet collection, with folders you can have files with the same name in different folders - not so in kwk. This decision allows us to:

1. Reduce entropy but avoiding duplicate code and other types of snippets to accumulate.
2. Use fewer keystrokes to address a snippet, thus even though we use pouch in an absolute path e.g.
 /richard/samples/hello.js - 'samples' is not usually used for addressing purposes, but rather cognitive categorization.*

> As with snippets pouch names are all lower case
> * Pouch names are actually required when moving snippets - [might need to revisit this]

### Create a pouch
``` sh
kwk mkdir big-project
# kwk mkdir <pouch name>
```  

### View a pouch
``` sh
kwk big-project
# kwk <pouch name>
```
 
### Delete a pouch
``` sh
kwk rm <pouch name>
```

### Make a pouch private
``` sh
kwk lock <pouch name>
```
### Make a pouch public 
``` sh
kwk unlock <pouch name>
```

### Rename a pouch
``` sh
kwk mv <pouch name> <new pouch name>
```

### Delete a pouch
``` sh
kwk rm <pouch name>
```