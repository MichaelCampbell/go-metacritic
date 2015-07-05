package metacritic

import (
        "encoding/json"
        "strings"

        "github.com/PuerkitoBio/goquery"
        )

type PersonBasic struct {
  Name, Url, AverageMovieScore, AverageTVScore, DOB, Categories string
}


func search_person(url string) (string, error) {
  var person_results []PersonBasic
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

    var avg_movie_score, avg_tv_score string
    s.Find("li.result li.avg_career_score").Each(func(i int, res *goquery.Selection) {
      val := res.Find("span.label").First().Text()
      if strings.Contains(val, "Movie") {
        avg_movie_score = res.Find("span.data").First().Text()
      } else {
        avg_tv_score = res.Find("span.data").First().Text()
      }
    })

    person := PersonBasic{
      Name: strings.TrimSpace(s.Find("a").First().Text()),
      Url: url,
      AverageMovieScore: avg_movie_score,
      AverageTVScore: avg_tv_score,
      DOB: strings.TrimSpace(s.Find("li.release_date span.data").First().Text()),
      Categories: strings.TrimSpace(s.Find("li.categories span.data").First().Text()),
    }
    person_results = append(person_results, person)
  })

  res, err := json.Marshal(person_results)
  if err != nil {
    return "", nil
  }

  return string(res), nil
}
