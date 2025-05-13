# HistGrep: Enhanced Command History Search 
Easily search your command history using multiple terms, exclusions, and formatting
options that can be customized for each query.

## What is it?
HistGrep is a terminal-based command-line tool for searching through history files 
or other logs. It offers improved search capabilities, allowing you to search using
multiple terms, exclude specific terms, and apply formatting options for each search. This provides
a quick and efficient way to find the information you need.

## Installation

Binaries for linux and macOS are provided in the `linux` and `macos` subdirectories.

### Build (requires Go 1.21+)

```bash
git clone https://github.com/TJN25/histgrep.git
cd histgrep/src/histgrep
go build
```

## Basic search

Run with `histgrep s -i input_file.txt foo bar baz` or `cat input_file.txt | histgrep s foo bar baz`.
You can redirect the output to a file using the -o flag.

## Options

**Colors**
Turn colors off with the `-f` or `--no-color` flag.

**Case Sensitivity**
Use case-sensitive search with the `-c` or `--case-sensitive` flag.

**Pager**
Enable paging with the `-p` or `--pager` flag.

**Line Numbering**
Include line numbers with the `-n` or `--numbered` flag.

**Exclude Terms**
Exclude specific terms with the `-x` or `--exclude` flag followed by the terms to exclude in quotes e.g. `-x "exclude_term1 exclude_term2"`.

## Configuration

Histgrep allows for a wide range of configuration options. When running `histgrep s`, it will search for two configuration files in `$HOME` and `$XDG_CONFIG_HOME`. To use your own custom configuration files, add `HISTGREP_CONFIG_PATH` to your environment.
Histgrep expects to find `histgrep.toml` and `formats.json` in the config directory.

Histgrep now supports a `histgrep.toml` file where default flags can be set. It can also be provided with the location of log files to be used as the search files when none are provided (`histgrep s my search terms` will search all files in `~/.logs/` matching the file pattern).
These changes consolidate the configuration into a single TOML file, making it easier for users to manage their settings. The `defaults.json` file can be removed, and users should be instructed to update their configurations accordingly when upgrading to this new version.

### histgrep.toml
```
[default_logs]
directory = "~/.logs/"
file_pattern = "{SHELL}-history-{YYYY}-{MM}-{DD}.log"

[search]
case_sensitive = false

[display]
color_enabled = true
pager_enabled = false
```

### formats.json
The `formats.json` file contains a list of search and output formats. Each JSON object
begins with the name of the format. This can be specified at runtime with 
`histgrep s -n NAME_OF_FORMAT` or `histgrep s --name NAME_OF_FORMAT`.

Inputs:
   -	`keys` is a list of names to label each segment of each line. These can be any ASCII string, but it's recommended to use a short, descriptive name.
   -	`separators` is a list of strings that will be used to separate each segment of a line. The first key will occur before the first separator, the second key will occur after the second separator, and so on.
Outputs:
   -	`keys` is a list of the names used in the inputs, indicating which keys to keep and in which order to show them.
   -	`separators` work the same as inputs, but it is possible to include an extra separator that will be appended to the end of the line.
Color:
   -	Specify the color of each key. This can also contain conditionals that will change the color of the key if certain strings match.
   -	Specify the color of the separators.
Excludes:
   -	Removes lines containing certain strings.
   -	Specify the key to search for a term within and whether the term is at the start, end, or anywhere in the line.
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

## New Features
- Automatic log file selection: If no input file is specified and stdin is empty, HistGrep will automatically use log files matching the pattern specified in the TOML config.
- Live search in pager mode: When using the pager, you can press / to search or ? to exclude terms. The search updates in real-time as you type.
- Navigation in pager mode: Use vim-like motions (j, k, g, G) or arrow keys to navigate through the results.

This updated README includes information about:
1. The new TOML configuration file and its structure
2. The new command-line flags for case sensitivity, pager mode, and line numbers
3. The automatic log file selection feature
4. The live search and navigation features in pager mode

The existing information about installation, basic usage, and the JSON configuration files has been preserved for backwards compatibility.
