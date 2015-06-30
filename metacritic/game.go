package metacritic

import (
        "encoding/json"
        "strings"

        "github.com/PuerkitoBio/goquery"
        )

type GameBasic struct {
  Name, Url, Poster, Summary, ReleaseDate, Certificate, Genres, Publisher, Platform string
  CriticRating Rating
  UserRating Rating
}

type Game struct {
  GameBasic
  CriticReviews []CriticReview
  UserReviews []UserReview
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

func find_game(query string) (string, error) {
  var gme Game
  search_url := BASE_URL + "/search/game/" + query + "/results"
  url, err := first_result(search_url)
  if err != nil {
    return "", err
  }

  doc, err := goquery.NewDocument(url)
  if err != nil {
    return "", err
  }

  poster, exists := doc.Find(".product_data_summary .product_image img.product_image").First().Attr("src")
  if !exists {
    poster = "Not Available"
  }

  var genres string
  doc.Find(".summary_wrap .side_details .summary_details li.product_genre .data").Each(func(i int, s *goquery.Selection) {
    genres = genres + strings.TrimSpace(s.Text()) + ", "
    })
  genres = genres[0:len(genres)-2]

  crs := critic_reviews(url)
  urs := user_reviews(url)
  user_rating_count := strings.TrimSpace(doc.Find(".product_scores .side_details .score_summary span.count a").First().Text())
  user_rating_count = user_rating_count[0:len(user_rating_count)-7]
  gme = Game{
          GameBasic: GameBasic{
            Name: strings.TrimSpace(doc.Find(".content_head .product_title a span").First().Text()),
            Url: url,
            Poster: poster,
            Summary: strings.TrimSpace(doc.Find(".product_details ul.summary_details li.product_summary span.data span").Text()),
            Certificate: strings.TrimSpace(doc.Find(".summary_wrap .side_details .summary_details li.product_rating span.data").First().Text()),
            ReleaseDate: strings.TrimSpace(doc.Find(".product_data ul.summary_details li.release_data span.data").Text()),
            Genres: genres,
            Publisher: strings.TrimSpace(doc.Find(".product_data ul.summary_details li.publisher span.data a span").Text()),
            Platform: strings.TrimSpace(doc.Find(".content_head .product_title span.platform a span").Text()),
            UserRating: Rating{
              Average: strings.TrimSpace(doc.Find(".product_scores .side_details .score_summary a div").First().Text()),
              Count: strings.TrimSpace(user_rating_count),
            },
            CriticRating: Rating{
              Average: strings.TrimSpace(doc.Find(".product_scores .metascore_summary a span").First().Text()),
              Count: strings.TrimSpace(doc.Find(".product_scores .metascore_summary .summary span.count a span").Text()),
            },
          },
          CriticReviews: crs,
          UserReviews: urs,
        }
  res, err := json.Marshal(gme)
  if err != nil {
    return "", nil
  }

  return string(res), nil
}
