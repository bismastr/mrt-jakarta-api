package main

import (
	"fmt"

	colly "github.com/gocolly/colly/v2"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	c.OnHTML("[data-stasiun]", func(e *colly.HTMLElement) {
		station := e.Attr("data-stasiun")
		fmt.Println("Stasiun:", station)
		e.ForEach(".col-12.col-xl-6", func(i int, h *colly.HTMLElement) {
			direction := h.ChildTexts("span")
			fmt.Println(direction)
		})

	})

	c.Visit("https://jakartamrt.co.id/id/jadwal-keberangkatan-mrt?dari=null")
}
