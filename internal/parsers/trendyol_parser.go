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

type TrendyolParser struct {
	Content io.Reader
}

func (p TrendyolParser) parseServiceResponse() []models.ResponseModel {

	doc, err := goquery.NewDocumentFromReader(p.Content)

	if err != nil {
		//return nil, err
		log.Fatal(err)
	}

	var items []models.ResponseModel

	doc.Find(".p-card-wrppr .p-card-chldrn-cntnr a").Each(func(i int, s *goquery.Selection) {
		productTitle, titleExist := s.Find(".prdct-desc-cntnr-ttl").Attr("title")
		url, _ := s.Attr("href")
		priceData := strings.Replace(s.Find(".prc-box-dscntd").Text(), " TL", "", 1)
		priceField, _ := strconv.ParseFloat(strings.Replace(priceData, ",", ".", 1), 64)

		if titleExist {
			items = append(items, models.ResponseModel{
				Title:      productTitle,
				Url:        fmt.Sprintf("https://www.trendyol.com%s", url),
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
