package main

import (
	"github.com/k0kubun/github-stream"
	"github.com/k0kubun/pp"
	"os"
)

func main() {
	token := getEnv("TOKEN", "")
	s := stream.NewStream(token)
	defer s.Stop()

	for {
		pp.Println(<-s.Events)
	}
}

func getEnv(key string, def string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return def
	}

	return v
}
