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
  r.HandleFunc("/api/{mode}/{category}/{q}", ServeHandler)
  http.ListenAndServe(":3000", r)
}

func ServeHandler(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  action := args["category"]
  query := args["q"]
  mode := args["mode"]
  action = strings.ToLower(action)

  var result string
  var err error

  if (mode == "search") {
    result, err = metacritic.Search(action, query)
  } else if (mode == "find") {
    result, err = metacritic.Find(action, query)
  } else {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintln(w, "Invalid Operation. Avaialble search, movie")
  }

  if err != nil {
    fmt.Fprintln(w, err)
  }

  var res []byte
  var movie_results []MovieBasic
  var game_results []GameBasic
  var movie Movie
  var game Game

  if (mode == "search") {
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
  } else {
    if (action == "game") {
      err = json.Unmarshal([]byte(result), &game)
      if err != nil {
       fmt.Println(err)
      }

      res, err = json.Marshal(game)
    } else {
      err = json.Unmarshal([]byte(result), &movie)
      if err != nil {
        fmt.Println(err)
      }

      res, err = json.Marshal(movie)
    }
  }

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprint(w, string(res))
}

