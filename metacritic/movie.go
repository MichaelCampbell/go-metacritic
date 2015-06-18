package metacritic

import (
        "encoding/json"
        "strings"

        "github.com/PuerkitoBio/goquery"
        )

type Basic struct {
  Name, Url, Poster, Summary, Certificate, Runtime, ReleaseDate, Genres, UserRating, MetacriticRating string
}

type CriticReview struct {
  Score, Source, Author, Summary, Url string
}

type Movie struct{
  Basic
  CriticReviews []CriticReview
}

func search_movie(url string) (string, error) {
  var movie_results []Basic
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
    m := Basic{
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

func find_movie(url string) (string, error) {
  var mov Movie

  doc, err := goquery.NewDocument(url)
  if err != nil {
    return "", err
  }

  url, exists := doc.Find(".content_head .product_title a").Attr("href")
  if !exists {
    url = "Not Available"
  } else {
    url = BASE_URL + url
  }

  poster, exists := doc.Find(".product_data_summary .product_image img.product_image").First().Attr("src")
  if !exists {
    poster = "Not Available"
  }

  crs := critic_reviews(url)
  var cert, runtime string
  doc.Find(".summary_wrap .side_details .summary_details li.product_rating").Each(func(i int, s *goquery.Selection) {
    if s.Find("span.label").Text() == "Rating:" {
      cert = s.Find("span.data").Text()
    } else {
      runtime = s.Find("span.data").Text()
    }
  })

  var genres string
  doc.Find(".summary_wrap .side_details .summary_details li.product_genre .data").Each(func(i int, s *goquery.Selection) {
    genres = genres + ", " + strings.TrimSpace(s.Text())
    })
  mov = Movie{
          Basic: Basic{
            Name: strings.TrimSpace(doc.Find(".content_head .product_title a span").Text()),
            Url: url,
            Poster: poster,
            Summary: strings.TrimSpace(doc.Find(".product_details ul.summary_details li.product_summary span.data").Text()),
            Certificate: strings.TrimSpace(cert),
            Runtime: strings.TrimSpace(runtime),
            ReleaseDate: strings.TrimSpace(doc.Find(".product_data ul.summary_details li.release_data span.data").Text()),
            Genres: genres,
            UserRating: strings.TrimSpace(doc.Find(".product_scores .side_details .score_summary a div").First().Text()),
            MetacriticRating: strings.TrimSpace(doc.Find(".product_scores .metascore_summary a span").First().Text()),
          },
          CriticReviews: crs,
        }

  res, err := json.Marshal(mov)
  if err != nil {
    return "", nil
  }
  return string(res), nil
}

func critic_reviews(url string) []CriticReview{
  var critic_reviews []CriticReview
  url = url + "/critic-reviews"
  doc, err := goquery.NewDocument(url)
  if err != nil {
    return critic_reviews
  }

  doc.Find(".product_reviews ol.critic_reviews li.critic_review").Each(func(i int, s *goquery.Selection) {
    url, exists := s.Find(".review_content ul.review_actions li.full_review a").Attr("href")
    if !exists {
      url = "Not Available"
    } else {
      url = BASE_URL + url
    }
    cr := CriticReview{
      Score: strings.TrimSpace(s.Find(".review_content .review_grade .metascore_w").Text()),
      Source: strings.TrimSpace(s.Find(".review_critic .source a").Text()),
      Author: strings.TrimSpace(s.Find(".review_critic .author a").Text()),
      Summary: strings.TrimSpace(s.Find(".review_body").Text()),
      Url: url,
    }

    critic_reviews = append(critic_reviews, cr)
  })
  return critic_reviews
}
