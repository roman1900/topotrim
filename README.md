# topotrim

A simple tool to trim a topoJSON file based on values of a particular property and output the result to a new file

## Usage:
```
topotrim -m property -s match -i source.json -o output.json
```
## Where:
- -m is the objects geometry property to match on. (defaults to POA_CODE16 because that is what i use it for :D)
- -s is the match criteria string (defaults to 7 because again that is what I'm using it for)
- -i is the input json file
- -o is the output json file (defaults to output.json)


## ToDos:
- The matching property is string only. Modify this for multiple types
- The match criteria is begins with only. Modify this to be regexy
- The arcs in the original file are just written directly to the new file. Modify this so only the arcs required from the matched geometries are written to the ouput
