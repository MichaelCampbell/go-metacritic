package metacritic

import (
        "strings"
        "errors"
        )

const BASE_URL = "http://www.metacritic.com"
var valid_types = []string{"movie", "game", "album"}

func Search(kind, query string) (string, error) {
  var err error
  var result string
  valid := is_valid_type(&kind, valid_types)
  if !valid {
    return "", errors.New("Invalid Type")
  }

  url := BASE_URL + "/search/" + kind + "/" + query + "/results"
  switch kind {
  case "movie":
    result, err = search_movie(url)
  case "game":
    result, err = search_game(url)
  case "album":
    result, err = search_album(url)
  }
  if err != nil {
    return "", err
  }

  return result, nil
}

func Find(kind, query string) (string, error) {
  var err error
  var result string
  valid := is_valid_type(&kind, valid_types)
  if !valid {
    return "", errors.New("Invalid Type")
  }

  switch kind {
  case "movie":
    result, err = find_movie(query)
  case "game":
    result, err = find_game(query)
  case "album":
    result, err = find_album(query)
  }
  if err != nil {
    return "", err
  }

  return result, nil
}

func is_valid_type(kind *string, arr []string) bool {
  *kind = strings.ToLower(*kind)
  for _,  el := range arr {
    if *kind == el {
      return true
    }
  }
  return false
}
