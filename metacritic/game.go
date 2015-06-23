package metacritic

import (
        "encoding/json"
        "strings"

        "github.com/PuerkitoBio/goquery"
        )

type GameBasic struct {
  Name, Url, Summary, ReleaseDate, Certificate, Publisher, Platform string
  CriticRating Rating
  UserRating Rating
}

type Game struct {
  GameBasic
}

func search_game(url string) (string, error) {
  var game_results []GameBasic
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
    game := GameBasic{
      Name: strings.TrimSpace(s.Find("a").First().Text()),
      Url: url,
      Summary: strings.TrimSpace(s.Find("p.basic_stat").First().Text()),
      ReleaseDate: strings.TrimSpace(s.Find("li.release_date span.data").First().Text()),
      Certificate: strings.TrimSpace(s.Find("li.maturity_rating span.data").First().Text()),
      CriticRating: Rating{
        Average: strings.TrimSpace(s.Find("div.main_stats span").First().Text()),
      },
      Publisher: strings.TrimSpace(s.Find("li.publisher span.data").First().Text()),
      Platform: strings.TrimSpace(s.Find("span.platform").First().Text()),
    }
    game_results = append(game_results, game)
  })

  res, err := json.Marshal(game_results)
  if err != nil {
    return "", nil
  }

  return string(res), nil
}
