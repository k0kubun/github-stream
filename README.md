# github-stream

GitHub Events API v3 client optimized for ETag header.

## What's this?

GitHub API v3 has [Events API](https://developer.github.com/v3/activity/events/).
This library polls the API and queues distinct events to its channel.

### Optimized for ETag header
This API's pagination is strictly limited. Thus it requires frequent polling.

The API requests can be optimized by "ETag" header.
This library requests APIs with given ETag header each time.

## Usage

```go
import (
  "github.com/k0kubun/github-stream"
  "github.com/k0kubun/pp"
)

func main() {
  s := stream.NewStream("[access token]")
  defer s.Stop()

  for ev := <-s.Events {
    pp.Println(ev)
  }
}
```

## API document

https://godoc.org/github.com/k0kubun/github-stream

## License

MIT License
