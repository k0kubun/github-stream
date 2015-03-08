# github-stream

GitHub Events API v3 client optimized for ETag header.

## What's this?

GitHub API v3 has [Events API](https://developer.github.com/v3/activity/events/).
This API's pagination is strictly limited. Thus it requires frequent polling.
This library polls the API and queues distinct events to its channel.

### Optimized for ETag header
The API requests can be optimized by "ETag" header.
This library requests APIs with given ETag header each time.

### Obey X-Poll-Interval
And the requests include "X-Poll-Interval".
This library's goroutine limits its API requests interval with the header.
GitHub API doc requires users to obey it.

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
    pp.Println(q)
  }
}
```

## License

MIT License
