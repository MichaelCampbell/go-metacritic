package metacritic

import (
        "strings"
        "errors"
        )

const BASE_URL = "http://www.metacritic.com"
var valid_types = []string{"movie", "games", "albums", "tv", "person", "trailers", "companies"}

func Search(kind, query string) (string, error) {
  kind = strings.ToLower(kind)
  valid := is_valid_type(kind, valid_types)
  if valid != true {
    return "", errors.New("Invalid Type")
  }

  url := BASE_URL + "/search/" + kind + "/" + query + "/results"
  result, err := search_movie(url)

  if err != nil {
    return "", err
  }

  return result, nil
}

func is_valid_type(kind string, arr []string) bool {
  for _,  el := range arr {
    if kind == el {
      return true
    }
  }
  return false
}
