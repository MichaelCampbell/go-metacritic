package metacritic

import (
        "encoding/json"
        "strings"
        "strconv"
        "errors"

        "github.com/PuerkitoBio/goquery"
        )

type MovieBasic struct {
  Name, Url, Poster, Summary, Certificate, Runtime, ReleaseDate, Genres string
  UserRating Rating
  CriticRating Rating
}

type Rating struct {
  Average, Count string
}

type CriticReview struct {
  Score, Source, Author, Summary, Url string
}

type UserReview struct {
  Username, ProfileUrl, Score, ReviewDate, Review, Like, Dislike string
}

type Movie struct{
  MovieBasic
  CriticReviews []CriticReview
  UserReviews []UserReview
}

func search_movie(url string) (string, error) {
  var movie_results []MovieBasic
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
    movie := MovieBasic{
      Name: strings.TrimSpace(s.Find("a").First().Text()),
      Url: url,
      Certificate: strings.TrimSpace(s.Find("li.rating .data").Text()),
      Runtime: strings.TrimSpace(s.Find("li.runtime .data").Text()),
      ReleaseDate: strings.TrimSpace(s.Find("li.release_date .data").Text()),
      Genres: strings.TrimSpace(s.Find("li.genre .data").Text()),
      UserRating: Rating{
        Average: strings.TrimSpace(s.Find("li.product_avguserscore .data").Text()),
      },
      CriticRating: Rating{
        Average: strings.TrimSpace(s.Find("span.metascore_w").Text()),
      },
    }
    movie_results = append(movie_results, movie)
  })

  res, err := json.Marshal(movie_results)
  if err != nil {
    return "", nil
  }

  return string(res), nil
}

func find_movie(query string) (string, error) {
  var mov Movie
  search_url := BASE_URL + "/search/movie/" + query + "/results"
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

  var cert, runtime string
  doc.Find(".summary_wrap .side_details .summary_details li.product_rating").Each(func(i int, s *goquery.Selection) {
    if s.Find("span.label").Text() == "Rating:" {
      cert = s.Find("span.data").Text()
    } else {
      runtime = s.Find("span.data").Text()
    }
  })

  crs := critic_reviews(url)
  urs := user_reviews(url)
  user_rating_count := strings.TrimSpace(doc.Find(".product_scores .side_details .score_summary span.count a").First().Text())
  user_rating_count = user_rating_count[0:len(user_rating_count)-7]
  mov = Movie{
          MovieBasic: MovieBasic{
            Name: strings.TrimSpace(doc.Find(".content_head .product_title a span").Text()),
            Url: url,
            Poster: poster,
            Summary: strings.TrimSpace(doc.Find(".product_details ul.summary_details li.product_summary span.data").Text()),
            Certificate: strings.TrimSpace(cert),
            Runtime: strings.TrimSpace(runtime),
            ReleaseDate: strings.TrimSpace(doc.Find(".product_data ul.summary_details li.release_data span.data").Text()),
            Genres: genres,
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
    cr := CriticReview{
      Score: strings.TrimSpace(s.Find(".review_content .review_grade .metascore_w").Text()),
      Source: strings.TrimSpace(s.Find(".review_critic .source a").Text()),
      Author: strings.TrimSpace(s.Find(".review_critic .author a").Text()),
      Summary: strings.TrimSpace(s.Find(".review_body").Text()),
      Url: s.Find(".review_content ul.review_actions li.full_review a").AttrOr("href", "Not Available"),
    }
    critic_reviews = append(critic_reviews, cr)
  })

  return critic_reviews
}

func user_reviews(url string) []UserReview{
  var user_reviews []UserReview
  url = url + "/user-reviews"
  doc, err := goquery.NewDocument(url)
  if err != nil {
    return user_reviews
  }

  pages := doc.Find("div.pages li.page").Nodes

  for i := 0; i <= len(pages); i++ {
    nxturl := url + "?page=" + strconv.Itoa(i)
    page, err := goquery.NewDocument(nxturl)
    if err != nil {
      return user_reviews
    }

    page.Find("ol.user_reviews li.user_review").Each(func(j int, s *goquery.Selection) {
      votes, _ := strconv.Atoi(s.Find(".review_content ul.review_actions .thumb_count span.total_ups").Text())
      t_votes, _ := strconv.Atoi(s.Find(".review_content ul.review_actions .thumb_count span.total_thumbs").Text())
      ur := UserReview{
        Username: strings.TrimSpace(s.Find(".review_content .review_critic .name a").Text()),
        ProfileUrl: strings.TrimSpace(s.Find(".review_content .review_critic .name a").AttrOr("href", "Not Available")),
        Score: strings.TrimSpace(s.Find(".review_content .review_grade div.user").Text()),
        ReviewDate: strings.TrimSpace(s.Find(".review_content .review_critic .date").Text()),
        Review: strings.TrimSpace(s.Find(".review_content .review_body span.blurb_expanded").Text()),
        Like: strconv.Itoa(votes),
        Dislike: strconv.Itoa(t_votes - votes),
      }
      user_reviews = append(user_reviews, ur)
    })
  }

  return user_reviews
}


func first_result(url string) (string, error) {
  doc, err := goquery.NewDocument(url)
  if err != nil {
    return "", err
  }

  url, exists := doc.Find(".body .search_results li.result h3.product_title a").First().Attr("href")
  if !exists {
    return "", errors.New("Movie Not Found")
  } else {
    return BASE_URL + url, nil
  }
}
