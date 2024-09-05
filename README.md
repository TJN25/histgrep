# HistGrep

Improved search through command history with multiple terms, excludes and formatting options that can be selected for each search.

## What is it?

HistGrep is a terminal based command line tool for searching through history files or other logs, offering a quick way to find what you want
and providing a number of easy formatting options for each search.

## Installation

### Build (requires Go 1.21+)

```bash
git clone https://github.com/TJN25/histgrep.git
cd histgrep/src/histgrep
go build
```

## Basic search

Run with `histgrep s -i input_file.txt search terms here` or `cat input_file.txt | histgrep s search terms here`.
Output can be redirected to a file with `-o` flag.

## Configuration

Histgrep allows for a wide range of configuration options. When running `histgrep s` it will search for two configuration files in `$HOME` and `$XDG_CONFIG_HOME`.
To use your own custom configuration files, add `HISTGREP_CONFIG_PATH` to your environment.

HistGrep expects to find `defaults.json` and `formats.json` in the config directory.

### formats.json
`formats.json` contains a list of search and output formats.
Json object begins with the name of the format. This can be specified at runtime with `histgrep s -n NAME_OF_FORMAT` or `histgrep s --name NAME_OF_FORMAT`
- Inputs
    - keys is a list of names to label each segment of each line. These can be any ascii string, but it's recommended to use a short name that is descriptive.
    - separators is a list of strings that will be used to separate each segment of a line. The first key will occur before the first separator, the second key will occur after the second separator, etc.
- Outputs
    - keys is a list of the names used in the inputs indicating which keys to keep and in which order to show them.
    - separators work the same as inputs, but it is possible to include an extra separator that will be appended on the end of the line.
- Color
    - Specify the color of each key. This can also contain conditionals that will change the color of the key if strings match.
    - Specify the color of the separators.
- Excludes: removes lines containing certain strings.
    -   Specify the key to search for a term within and whether the term is at the start, end or anywhere in the line.
```
{
    "simple":{
        "Input":{
            "keys":["date","time","directory","command"],
            "separators":["."," ",": "]
        },
        "Output":{
            "keys":["command","directory","date"],
            "separators":[" # from ", " :: "]
        },
        "Color":{
            "command":{"default":"green","commit":"red"},
            "directory":{"default":"grey"},
            "date":{"default":"grey"},
            "SEPARATOR":{"default":"grey"}},
        "Excludes":{
            "command": {
                "starts_with" : ["cd","clear","ls","ll","pwd","less","more","cat","echo","exit"]
            },
            "directory": {
                "contains" : ["EXAMPLE", "REMINDER"],
                "ends_with" : ["FOO", "BAR"]
            }
        }
    }
}

###

```

### defaults.json
Provide the name of the format you wish to use as a default. If this is missing then no formatting will be done.
```
{
 "Name": "simple"
}
```


[MIT](https://github.com/charmbracelet/glow/raw/master/LICENSE)

