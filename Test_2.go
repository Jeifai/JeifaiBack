
package main
import	(
	"fmt"
	"github.com/gocolly/colly/v2"
)

func main() {
	
	start_url := "https://it.wikipedia.org/wiki/Maximilian_Nisi"

	c := colly.NewCollector()

	c.OnHTML("html", func(e *colly.HTMLElement) {
	
		paragraphs := e.ChildTexts("p")

		for i := 0; i < len(paragraphs) ; i++{

			fmt.Println(i, paragraphs[i])

		}
	})

	c.Visit(start_url) 

}
