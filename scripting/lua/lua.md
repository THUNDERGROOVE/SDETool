## Lua

For lua we export a global table `sde`

| Function name 	| Args           	| Return                        	| Explanation                                           	|
|---------------	|----------------	|-------------------------------	|-------------------------------------------------------	|
| getVersions   	| None           	| A table of available versions 	| Lists all versions known by SDETool                   	|
| loadVersion   	| version string 	| None                          	| Sets the internal SDE object to the version specified 	|
| getTypeByID   	| TypeID int     	| An exported SDEType table     	|                                                       	|
| search            | search string     | A table of SDEType tables         | Uses SDETool's internal SDE.Search method                 |

If you would like more methods add an issue and I'll add it fairly quickly.