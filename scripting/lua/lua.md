## Lua

For lua we export a global table `sde`

| Function name 	| Args           	| Return                        	| Explanation                                           	| First Implemented 	|
|---------------	|----------------	|-------------------------------	|-------------------------------------------------------	|------------       	|
| getVersions   	| None           	| A table of available versions 	| Lists all versions known by SDETool                   	| master-0.1        	|
| loadVersion   	| version string 	| None                          	| Sets the internal SDE object to the version specified 	| master-0.1        	|
| getTypeByID   	| TypeID int     	| An exported SDEType table     	|                                                       	| master-0.1        	|
| search        	| search string 	| A table of SDEType tables     	| Uses SDETool's internal SDE.Search method             	| dev-0.1	        	|

If you would like more methods add an issue and I'll add it fairly quickly.