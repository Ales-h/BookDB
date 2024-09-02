package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
    "os"
    "regexp"
	//"strconv"
	//"time"
	"github.com/labstack/echo/v4"
	//"github.com/labstack/echo/v4/middleware"
)

var API_KEY string

var logger echo.Logger

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
}

type Volume struct {
	Title         string   `json:"title"`
	Subtitle      string   `json:"subtitle,omitempty"`
	Authors       []string `json:"authors,omitempty"`
	YearOfRelease string   `json:"publishedDate,omitempty"`
	PageCount     int      `json:"pageCount"`
	ImageLinks    struct {
		SmallThumbnail string `json:"smallThumbnail,omitempty"`
	} `json:"imageLinks,omitempty"` //delete omitempty TODO
}

func GetVolumeData(url string) (Volume, error) {
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		return Volume{}, err
	}

	var Results struct {
		VolumeInfo Volume `json:"volumeInfo"`
	}

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&Results)
	if err != nil {
		return Volume{}, err
	}

	debug := fmt.Sprintf("%v", Results.VolumeInfo)
	logger.Info(debug)
	return Results.VolumeInfo, nil
}

func BooksSearch(query string) ([]Volume, error) {
	const baseURL = "https://www.googleapis.com/books/v1/volumes"

	url := fmt.Sprintf("%s?q=%s&printType=books&key=%s", baseURL, query, API_KEY)
	fmt.Print(url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Items []struct {
			VolumeInfo Volume `json:"volumeInfo"`
		} `json:"items"`
	}

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&result)
	if err != nil {
		fmt.Println("error1")
		return []Volume{}, err
	}

	volumes := make([]Volume, len(result.Items))

	for i, item := range result.Items {
		volumes[i] = item.VolumeInfo
	}

	return volumes, nil

}

func main() {
    API_KEY = os.Getenv(API_KEY)

	e := echo.New()
	logger = e.Logger
	e.Renderer = NewTemplate()
	e.Static("/public", "public")
	e.Debug = true
    re, err := regexp.Compile("\\s+")

    if err != nil {
        panic(fmt.Sprintf("Failed to compile Regex: %v", err))
    }
	e.GET("/", func(c echo.Context) error {
		fmt.Printf("HELLO")
		return c.Render(http.StatusOK, "index", nil)
	})

	e.POST("/search", func(c echo.Context) error {
		query := c.FormValue("search")
        if query == "" { // if search bar is empty clear the Volumes
            return c.Render(200, "results", []Volume{})
        }
        q := re.ReplaceAllString(query, "+")
		volumes, err := BooksSearch(q)
		if err != nil {
			fmt.Printf(err.Error())
			return c.String(http.StatusInternalServerError, "Error fetching books")
		}
		fmt.Printf("%v", volumes)

		return c.Render(200, "results", volumes)

	})

	e.Logger.Fatal(e.Start(":8080"))

}
