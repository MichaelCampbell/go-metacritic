package metacritic

import (
        "strings"
        "errors"
        )

const BASE_URL = "http://www.metacritic.com"
var valid_types = []string{"movie", "games", "albums", "tv", "person", "trailers", "companies"}

func Search(kind, query string) (string, error) {
  valid := is_valid_type(&kind, valid_types)
  if !valid {
    return "", errors.New("Invalid Type")
  }

  url := BASE_URL + "/search/" + kind + "/" + query + "/results"
  result, err := search_movie(url)

  if err != nil {
    return "", err
  }

  return result, nil
}

func Find(kind, query string) (string, error) {
  valid := is_valid_type(&kind, valid_types)
  if !valid {
    return "", errors.New("Invalid Type")
  }

  url := BASE_URL + "/movie/" + query
  result, err := find_movie(url)
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
