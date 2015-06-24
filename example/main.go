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

type Game struct {
  Name, Url, Summary, ReleaseDate, Certificate, Publisher, Platform string
  CriticRating Rating
  CriticReviews []CriticReview
  UserReviews []UserReview
}

type Rating struct {
  Average string
  Count string
}

type CriticReview struct {
  Score, Source, Author, Summary, Url string
}

type UserReview struct {
  Username, ProfileUrl, Score, ReviewDate, Review, Like, Dislike string
}

func main() {
  result, err := metacritic.Search("Game", "fight")
  if err != nil {
    fmt.Println(err)
  }

  var game_results []Game
  err = json.Unmarshal([]byte(result), &game_results)

  if err != nil {
    fmt.Println(err)
  }

  for _, game := range game_results {
    fmt.Println(game.Name)
  }
}
