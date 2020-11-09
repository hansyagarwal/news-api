# news-api

## Requirements
- Postman
- Go ver 1.4 or later
- Text editor (VSCode recommended)

## What operations it can do?
- Create an Article
- Get an Article by it's id
- List/Get all the Articles
- Search for an artice by it's title,subtitle or content

## How to run?
import/download all package and libraries using `go get`
``` 
go run main.go
```
> port: 3000

open postman and create new request

### for creating an article, 
switch from GET to POST and type `localhost:3000/articles`
select `body` and then `raw` and type
```
{
	"Title": "news article",
	"Subtitle": "subtitle of new article",
	"Content": "content of new article"
}
```
click send

### for other operations,
switch to GET and type:
`localhost:3000/articles` to list all the articles

`localhost:3000/articles/<id>` to get the article by its id (try `5fa927a02950189579ca32f4` as id)

`localhost:3000/articles/search?q=<search term>` to get the article by its title/subtitle/content (if the title is more than 1 word use `%20` for space, eg: `/articles/search?q=sup%20sub`)
