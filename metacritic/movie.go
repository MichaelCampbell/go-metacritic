package metacritic

import (
        "encoding/json"
        "log"

        "github.com/PuerkitoBio/goquery"
        )

type MovieResult struct {
  name, url, certificate, runtime, releaseDate, genres, userRating, metacriticRating string
}

func search_movie(url string) (string, error) {
  var movie_results []MovieResult

  doc, err := goquery.NewDocument(url)

  if err != nil {
    return "", err
  }

  doc.Find(".body .search_results li").Each(func(i int, s *goquery.Selection) {
    movie_results[i] = MovieResult{
      name: s.Find("a").First().Text(),
      // url: BASE_URL + s.Find("a").Attr("href"),
      certificate: s.Find("li.rating .data").Text(),
      runtime: s.Find("li.runtime .data").Text(),
      releaseDate: s.Find("li.release_date .data").Text(),
      genres: s.Find("li.genre .data").Text(),
      userRating: s.Find("li.product_avguserscore .data").Text(),
      metacriticRating: s.Find("span.metascore_w").Text(),
    }
  })

  res, err := json.Marshal(movie_results)
  if err != nil {
    log.Fatal(err)
  }

  return string(res), nil

}
