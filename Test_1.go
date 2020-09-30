package main
import	(
	"fmt"
	"github.com/gocolly/colly/v2"
)

func main() {
	
	start_url := "https://it.wikipedia.org/wiki/Maximilian_Nisi"

	c := colly.NewCollector()

	c.OnHTML("html", func(e *colly.HTMLElement) {


		//.ChildText accetta solo stringhe come parametri!!
		title := e.ChildText("title")
		fmt.Println(title)

		//.firstHeading il . si riferisce alla classe
		titleByClass := e.ChildText(".firstHeading")

		//#firstHeading il # si riferisce al id
		titleById := e.ChildText("#firstHeading")
		
		fmt.Println(titleByClass)
		fmt.Println(titleById)
		paragraphs := e.ChildTexts("p")
		description := paragraphs[1]
		fmt.Println(description)
		fmt.Println(len(paragraphs))
		image := e.ChildAttr(".image", "href")
		fmt.Println("https://it.wikipedia.org/" + image)




	})

	c.Visit(start_url)
}
