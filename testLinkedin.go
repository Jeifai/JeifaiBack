package main

import (
	"fmt"
	// "io/ioutil"
	// "log"
	// "net/http"
	"context"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func mainLinkedin() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res1 interface{}
	var res2 interface{}
	var res3 []byte
	var initialResponse string
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.linkedin.com/login/de?fromSignIn=true&trk=guest_homepage-basic_nav-header-signin"),
		chromedp.Sleep(5*time.Second),
		chromedp.EvaluateAsDevTools(`document.getElementById('username').value='massimo.perlomini@gmail.com';`, &res1),
		chromedp.Sleep(1*time.Second),
		chromedp.EvaluateAsDevTools(`document.getElementById('password').value='nxnxxmyy9891!linkedin';`, &res2),
		chromedp.Sleep(1*time.Second),
		chromedp.EvaluateAsDevTools(`document.getElementsByClassName("btn__primary--large from__button--floating mercado-button--primary")[0].click()`, &res3),
		chromedp.WaitVisible(".global-nav"),
		chromedp.Navigate("https://www.linkedin.com/company/nen-energia"),
		chromedp.WaitVisible(".company-hero-image"),
		chromedp.OuterHTML("html", &initialResponse),
	); err != nil {
		panic(err)
	}
	fmt.Println(res2)
	SaveResponseToFileWithFileNameLink(initialResponse, "linkedin.html")
}

func SaveResponseToFileWithFileNameLink(response string, filename string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	f, err := os.Create(dir + "/" + filename)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	f.WriteString(response)
}
