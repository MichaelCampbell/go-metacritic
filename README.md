##Go-Metacritic

go-metacritic is a go library for accessing data from metacritic.com using golang.

To get started,

`go get https://github.com/avinoth/go-metacritic`

and import it into the code

`import "github.com/avinoth/go-metacritic/metacritic"`

####Available APIs
* Search(category, query)
* Find(category, query)

```go
import "github.com/avinoth/go-metacritic/metacritic"

type MovieBasic struct {
  Name, Url, Poster, Certificate, Runtime, ReleaseDate, Genres string
}

movie_results := []MovieBasic

//Search API takes two parameters and returns a string.
result, _ = metacritic.Search("movie", "fight")
json.Unmarshal([]byte(result), &movie_results)
res, _ = json.Marshal(movie_results)

fmt.Print(string(res))
```




#####Available categories
* Game
* Movie
* Album
* Person
