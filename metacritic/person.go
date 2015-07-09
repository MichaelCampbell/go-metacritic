package metacritic

import (
        "encoding/json"
        "strings"

        "github.com/PuerkitoBio/goquery"
        )

type PersonBasic struct {
  Name, Url, DOB, Categories string
  AverageMovieScore Rating
  AverageTVScore Rating
}

type Person struct {
  PersonBasic
  Biography string
  MovieScores Distribution
  TVScores Distribution
  MovieCredits []CreditInfo
  TVCredits []CreditInfo
}

type Distribution struct {
  Positive, Mixed, Negative string
}

type CreditInfo struct {
  CriticRating, Name, Url, Year, Credit, UserRating string
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
      AverageMovieScore: Rating{
        Average: avg_movie_score,
      },
      AverageTVScore: Rating{
        Average: avg_tv_score,
      },
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

func find_person(query string) (string, error) {
  var person Person
  search_url := BASE_URL + "/search/person/" + query + "/results"
  url, err := first_result(search_url)
  if err != nil {
    return "", err
  }

  doc, err := goquery.NewDocument(url)
  if err != nil {
    return "", err
  }

  var positive, mixed, negative string
  doc.Find("#main_content ol.score_counts li.score_count").Each(func(i int, s *goquery.Selection) {
    if s.Find("span.label").Text() == "Positive:" {
      positive = s.Find("span.count").Text()
    } else if s.Find("span.label").Text() == "Mixed:" {
      mixed = s.Find("span.count").Text()
    } else if s.Find("span.label").Text() == "Negative:" {
      negative = s.Find("span.count").Text()
    }
  })

  movie_credits := credits(url)
  person = Person{
          PersonBasic: PersonBasic{
            Name: strings.TrimSpace(doc.Find("#main_content .person_title").First().Text()),
            Url: url,
            AverageMovieScore: Rating{
              Average: strings.TrimSpace(doc.Find("#main_content .reviews_total span.count a").Text()),
            },
          },
          MovieScores: Distribution{
            Positive: positive,
            Mixed: mixed,
            Negative: negative,
          },
          Biography: strings.TrimSpace(doc.Find("#main_content .bio span.blurb_expanded").Text()),
          MovieCredits: movie_credits,
        }
  res, err := json.Marshal(person)
  if err != nil {
    return "", nil
  }

  return string(res), nil
}

func credits(url string) []CreditInfo {
  var credit_summary []CreditInfo
  doc, err := goquery.NewDocument(url)
  if err != nil {
    return credit_summary
  }

  doc.Find("table.person_credits tbody tr").Each(func(i int, s *goquery.Selection) {
    cr := CreditInfo{
      CriticRating: strings.TrimSpace(s.Find(".title span").Text()),
      Name: strings.TrimSpace(s.Find(".title a").Text()),
      Url: BASE_URL + strings.TrimSpace(s.Find(".title a").AttrOr("href", "NA")),
      Year: strings.TrimSpace(s.Find(".year").Text()),
      Credit: strings.TrimSpace(s.Find(".role").Text()),
      UserRating: strings.TrimSpace(s.Find(".score span").Text()),
    }
    credit_summary = append(credit_summary, cr)
  })

  return credit_summary
}

