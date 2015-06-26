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
  result, err := metacritic.Find("Game", "fight")
  if err != nil {
    fmt.Println(err)
  }

  var game Game
  err = json.Unmarshal([]byte(result), &game)

  if err != nil {
    fmt.Println(err)
  }

  fmt.Println(game.Name)
}
