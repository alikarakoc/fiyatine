package parsers

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aligoren/fiyatine/internal/models"
)

type HepsiburadaParser struct {
	Content io.Reader
}

func (p HepsiburadaParser) parseServiceResponse() []models.ResponseModel {

	doc, err := goquery.NewDocumentFromReader(p.Content)

	if err != nil {
		//return nil, err
		log.Fatal(err)
	}

	var items []models.ResponseModel

	doc.Find(".productListContent-item div a").Each(func(i int, s *goquery.Selection) {
		productTitle, titleExist := s.Attr("title")
		url, _ := s.Attr("href")
		priceData := s.Find("div[data-test-id='price-current-price']").Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
			return !s.Is("span")
		}).Text()
		priceField, _ := strconv.ParseFloat(strings.Replace(priceData, ",", ".", 1), 64)

		splitUrl := strings.Split(url, "-")
		id := splitUrl[len(splitUrl)-1]

		if titleExist {
			items = append(items, models.ResponseModel{
				ID:         id,
				Title:      productTitle,
				Url:        url,
				Price:      fmt.Sprintf("₺%s", priceData),
				PriceField: priceField,
			})
		}
	})

	sort.Slice(items, func(i, j int) bool {
		return items[i].PriceField < items[j].PriceField
	})

	return items
}
