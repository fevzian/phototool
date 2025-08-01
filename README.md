# Project photosort

## Description
Let's organize photos according dates. Provides various functions such as sort, exif, etc.

## Prerequisites
exif-tool is required to be installed inside oprating system.

For debian based linux: sudo apt install exiftool
For windows based os:   refer to https://exiftool.org/install.html 

## Sorting command
phototool sort --src_dir <source_dir_path> --dest_dir <destination_dir_path>

## All metadata for a spacific image
phototool exif --file <file_path>


## Building

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
