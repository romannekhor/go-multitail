# go-multitail

## What is this?

This is just a very primitive clone of `multitail`.

## Usage example

```
./go-multitail \
-cmd "bash read_host_X_logs.sh" -l "host_X" -color yellow \
-cmd "bash read_host_Y_logs.sh" -l "host_Y" -color green
```

## TODO
* Save output history in memory
* Navigate Up/Down through the history
* Search in history
* Highlight substrings in output