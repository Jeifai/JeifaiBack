
package main
import	(
	"fmt"
	"github.com/gocolly/colly/v2"
)

func main() {
	
	start_url := "https://it.wikipedia.org/wiki/Maximilian_Nisi"
	start_url_2 := "https://it.wikipedia.org/wiki/Gaetano_Quagliariello"

	c := colly.NewCollector()

	c.OnHTML("p", func(e *colly.HTMLElement) {
	
		paragraph := e.Text

		fmt.Println(paragraph)

	})

	sArray := [2]string{start_url, start_url_2}

	for i := 0; i < len(sArray) ; i++{

		fmt.Println(sArray[i])

		c.Visit(sArray[i])

	}

}
