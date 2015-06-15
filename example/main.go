package main

import (
        "fmt"
        "github.com/avinoth/go-metacritic/metacritic"
        "encoding/json"
        )

type Movie struct {
  Name, Url, Certificate, Runtime, ReleaseDate, Genres, UserRating, MetacriticRating string
}

func main() {
  result, err := metacritic.Search("Movie", "fight")
  if err != nil {
    fmt.Println(err)
  }

  var mov []Movie
  err = json.Unmarshal([]byte(result), &mov)

  if err != nil {
    fmt.Println(err)
  }

  fmt.Printf("|Name\t Url \t UserRating \tMetacriticRating\n")
  for _, el := range mov {
    fmt.Printf("|%-50s \t | %-100s \t |%s \t|%s \t| \n", el.Name, el.Url, el.UserRating, el.MetacriticRating)
  }
}
