package main

import (
        "fmt"
        "github.com/avinoth/go-metacritic/metacritic"
        "encoding/json"
        )

type Movie struct {
  Name, Url, Poster, Certificate, Runtime, ReleaseDate, Genres string
  UserRating Rating
  CriticRating Rating
  CriticReviews []CriticReview
}

type Rating struct {
  Average string
  Count string
}

type CriticReview struct {
  Score, Source, Author, Summary, Url string
}

func main() {
  result, err := metacritic.Find("Movie", "fight-club")
  if err != nil {
    fmt.Println(err)
  }

  var mov Movie
  err = json.Unmarshal([]byte(result), &mov)

  if err != nil {
    fmt.Println(err)
  }

  fmt.Printf("Users rated the movie %s %s out of %s. \n", mov.Name, mov.UserRating.Average, mov.UserRating.Count)
  fmt.Printf("Critics rated the movie %s %s out of %s. \n", mov.Name, mov.CriticRating.Average, mov.CriticRating.Count)

  // fmt.Printf("|Name \t|UserRating \t|MetacriticRating |\n")
  // fmt.Printf("|%s \t|%s \t|%s |\n", mov.Name, mov.UserRating.Average, mov.MetacriticRating.Average)
  // for _, el := range mov {
  //   fmt.Printf("|%-50s \t |%s \t|%s \t| \n", el.Name, el.UserRating, el.MetacriticRating)
  // }
}
