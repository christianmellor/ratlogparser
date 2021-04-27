# Ratlog Parser 🐀

Ratlog parser is a go package that parses [ratlog](https://github.com/ratlog/ratlog-spec) formatted logs.

## Usage

```go 
import "github.com/christianmellor/ratlogparser"

// Parse logs from Stdin
r := os.Stdin
// Create a new parsing instance
p := ratlogparser.SimpleParser{}
// Give it somewhere to place the entries
entries := ratlogparser.NewEntryReaderWriter(nil)
err := p.Parse(r, entries)
if err != nil {
	log.Fatal(err)
}

for {
    // Read each entry
    e, err := entries.Entry()
    if err != nil {
        // err == io.EOF signals end of input
        fmt.Println(err)
        break;
    }
    fmt.Println(e.Message)
}
```
