# ahl

A simple CLI tool to highlight patterns in text streams using ANSI colors. Useful for log analysis, debugging, and making terminal output more readable.

## Features
- Highlight multiple regex patterns with different colors
- Supports a wide range of ANSI colors
- Easy to use in pipelines and scripts

## Usage

Pipe text into `ahl` and specify patterns:

```sh
echo "error: something failed" | ahl --pattern 'error=red' --pattern 'failed=yellow'
# or using short flags
ahl -p 'error=red' -p 'failed=yellow'
```

Or use a single pattern (defaults to red):

```sh
echo "warning: disk space low" | ahl 'warning'
```

### Stripping existing color codes

If your input already contains ANSI color codes and you want to remove them before highlighting, use the `--cleanup` flag (or `-c`):

```sh
echo -e "\033[31merror\033[0m: something failed" | ahl --cleanup --pattern 'error=red'
# or
ahl -c -p 'error=red'
```

This will strip all existing color codes before applying your own highlighting.

## Installation

### Build from source

You need Go 1.18 or newer installed. Then run:

```sh
git clone https://github.com/yourusername/ahl.git
cd ahl
go build -o ahl
```

This will produce an `ahl` binary in the current directory.

## Supported Colors

- black
- red
- green
- yellow
- blue
- magenta
- cyan
- white
- gray
- brightred
- brightgreen
- brightyellow
- brightblue
- brightmagenta
- brightcyan
- brightwhite

## Pattern Syntax

Patterns are specified as `REGEX=COLOR`. Example:

```sh
ahl --pattern 'foo=green' --pattern 'bar=yellow'
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

See [LICENSE](LICENSE)

## Running Tests

To run the test suite:

```sh
go test
``` 