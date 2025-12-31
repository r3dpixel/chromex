# chromex

A minimal Go wrapper around [chromedp](https://github.com/chromedp/chromedp) for Chrome automation.

## Installation

```bash
go get github.com/r3dpixel/chromex
```

## Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/chromedp/chromedp"
    "github.com/r3dpixel/chromex"
)

func main() {
    title, err := chromex.RunChrome(chromex.ChromeConfig{
        Timeout: 30 * time.Second,
    }, func(ctx context.Context) (string, error) {
        var title string
        err := chromedp.Run(ctx,
            chromedp.Navigate("https://example.com"),
            chromedp.Title(&title),
        )
        return title, err
    })

    if err != nil {
        panic(err)
    }
    fmt.Println(title)
}
```

## Configuration

| Option       | Description               | Default          |
|--------------|---------------------------|------------------|
| `ChromePath` | Path to Chrome executable | Auto-detected    |
| `Timeout`    | Operation timeout         | 60s              |
| `Flags`      | Custom chromedp flags     | `DefaultFlags()` |