# Ratlog Parser 🐀

Ratlog parser is a go package that parses [ratlog](https://github.com/ratlog/ratlog-spec) formatted logs.

## Usage

```go 
import "github.com/christianmellor/ratlogparser"

// Parse logs from Stdin
r := os.Stdin
p := ratlogparser.SimpleParser{}
entries := ratlogparser.NewEntryReaderWriter(nil)
err := p.Parse(r, entries)
if err != nil {
	log.Fatal(err)
}

for {
    e, err := entries.Entry()
    if err != nil {
        break;
    }
    fmt.Println(e)
}

// if err == io.EOF signals end of input
fmt.Println(err)
```
