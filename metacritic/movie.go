package metacritic

import (
        "encoding/json"
        "strings"

        "github.com/PuerkitoBio/goquery"
        )

type MovieResult struct {
  Name, Url, Certificate, Runtime, ReleaseDate, Genres, UserRating, MetacriticRating string
}

func search_movie(url string) (string, error) {
  var movie_results []MovieResult

  doc, err := goquery.NewDocument(url)

  if err != nil {
    return "", err
  }

  doc.Find(".body .search_results li.result").Each(func(i int, s *goquery.Selection) {
    url, exists := s.Find("h3.product_title a").Attr("href")
    if !exists {
      url = "Not Available"
    } else {
      url = BASE_URL + url
    }
    m := MovieResult{
      Name: strings.TrimSpace(s.Find("a").First().Text()),
      Url: url,
      Certificate: strings.TrimSpace(s.Find("li.rating .data").Text()),
      Runtime: strings.TrimSpace(s.Find("li.runtime .data").Text()),
      ReleaseDate: strings.TrimSpace(s.Find("li.release_date .data").Text()),
      Genres: strings.TrimSpace(s.Find("li.genre .data").Text()),
      UserRating: strings.TrimSpace(s.Find("li.product_avguserscore .data").Text()),
      MetacriticRating: strings.TrimSpace(s.Find("span.metascore_w").Text()),
    }

    movie_results = append(movie_results, m)
  })

  res, err := json.Marshal(movie_results)

  if err != nil {
    return "", nil
  }

  return string(res), nil
}
