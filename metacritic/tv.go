package metacritic

import (
        "encoding/json"
        "strings"

        "github.com/PuerkitoBio/goquery"
        )

type TvBasic struct {
  Name, Url, Summary, ReleaseDate, Genres string
  CriticRating Rating
}


func search_tv(url string) (string, error) {
  var tv_results []TvBasic
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
    tv := TvBasic{
      Name: strings.TrimSpace(s.Find("a").First().Text()),
      Url: url,
      Summary: strings.TrimSpace(s.Find("p.basic_stat").First().Text()),
      ReleaseDate: strings.TrimSpace(s.Find("li.release_date span.data").First().Text()),
      Genres: strings.TrimSpace(s.Find("li.genre span.data").First().Text()),
      CriticRating: Rating{
        Average: strings.TrimSpace(s.Find("div.main_stats span").First().Text()),
      },
    }
    tv_results = append(tv_results, tv)
  })

  res, err := json.Marshal(tv_results)
  if err != nil {
    return "", nil
  }

  return string(res), nil
}
