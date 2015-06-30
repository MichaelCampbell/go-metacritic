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
  r.HandleFunc("/search/{category}/{q}", SearchHandler)
  // r.HandleFunc("/find/{category}/{q}", FindHandler)
  http.ListenAndServe(":3000", r)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  action := args["category"]
  query := args["q"]

  result, err := metacritic.Search(action, query)

  if err != nil {
    fmt.Println(err)
  }

  // var game Game
  // var movie Movie
  // var movie_results []Movie
  var game_results []Game

  err = json.Unmarshal([]byte(result), &game_results)
  if err != nil {
    fmt.Println(err)
  }

  gme, err := json.Marshal(game_results)
  if err != nil {
    fmt.Println(err)
  }
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprint(w, string(gme))
}
