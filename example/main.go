package main

import (
        "fmt"
        "encoding/json"
        "strings"
        "net/http"

        "github.com/avinoth/go-metacritic/metacritic"
        "github.com/gorilla/mux"
        )

type MovieBasic struct {
  Name, Url, Poster, Certificate, Runtime, ReleaseDate, Genres string
  UserRating Rating
  CriticRating Rating
}

type Movie struct {
  MovieBasic
  CriticReviews []CriticReview
  UserReviews []UserReview
}

type GameBasic struct {
  Name, Url, Summary, ReleaseDate, Certificate, Publisher, Platform string
  CriticRating Rating
}

type Game struct {
  GameBasic
  UserRating Rating
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
  action = strings.ToLower(action)
  result, err := metacritic.Search(action, query)

  if err != nil {
    fmt.Println(err)
  }

  var res []byte
  var movie_results []MovieBasic
  var game_results []GameBasic

  if (action == "game") {
    err = json.Unmarshal([]byte(result), &game_results)
    if err != nil {
     fmt.Println(err)
    }

    res, err = json.Marshal(game_results)
  } else {
    err = json.Unmarshal([]byte(result), &movie_results)
    if err != nil {
      fmt.Println(err)
    }

    res, err = json.Marshal(movie_results)
  }

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprint(w, string(res))
}
