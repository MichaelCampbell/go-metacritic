package main

import (
        "fmt"
        "github.com/avinoth/go-metacritic/metacritic"
        )

func main() {
  result, err := metacritic.Search("Movie", "fight")
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(result)
}
