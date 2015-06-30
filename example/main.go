package main

import (
        "fmt"
        "github.com/avinoth/go-metacritic/metacritic"
        "encoding/json"
        "github.com/gorilla/mux"
        "net/http"
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
  r := mux.NewRouter()
  r.HandleFunc("/game/{q}", GameHandler)
  http.ListenAndServe(":3000", r)
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  // action := args["action"]
  query := args["q"]

  result, err := metacritic.Find("game", query)

  if err != nil {
    fmt.Println(err)
  }

  var game Game
  err = json.Unmarshal([]byte(result), &game)
  if err != nil {
    fmt.Println(err)
  }

  gme, err := json.Marshal(game)
  if err != nil {
    fmt.Println(err)
  }
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprint(w, string(gme))
}
