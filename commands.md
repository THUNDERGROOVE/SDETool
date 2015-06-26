# Commands

lookup:
  -sde <string|filename>
    Explicitly loads this file instead of the newest SDE from the fileserver
  -t <int|TypeID>
    A typeID to look for
  -tn <string|TypeName>
    A TypeName to look for
  -td <string|mDisplayName>
    A mDisplayName to look for
  -attr
    Print the types attributes
search:
  -sde <string|filename>
    Explicitly loads this file instead of the newest SDE from the fileserver
  -t <string|TypeID, TypeName, mDisplayName>
    Lists all values found by calling search
  -attr
    Prints all of the attributes.  Only works when one value is returned by the search

```
// All global flags are booleans
```
global flags:
  -debug
    Prints debug information
  -nocolor
    Does not colorize output
    
