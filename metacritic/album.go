package metacritic

import (
        "encoding/json"
        "strings"

        "github.com/PuerkitoBio/goquery"
        )

type AlbumBasic struct {
  Name, Url, Summary, ReleaseDate, Genres string
  CriticRating Rating
}


func search_album(url string) (string, error) {
  var album_results []AlbumBasic
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
    album := AlbumBasic{
      Name: strings.TrimSpace(s.Find("a").First().Text()),
      Url: url,
      Summary: strings.TrimSpace(s.Find("p.basic_stat").First().Text()),
      ReleaseDate: strings.TrimSpace(s.Find("li.release_date span.data").First().Text()),
      CriticRating: Rating{
        Average: strings.TrimSpace(s.Find("div.main_stats span").First().Text()),
      },
    }
    album_results = append(album_results, album)
  })

  res, err := json.Marshal(album_results)
  if err != nil {
    return "", nil
  }

  return string(res), nil
}
