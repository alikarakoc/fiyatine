package services

import (
	"log"
	"net/http"
	"net/url"

	"github.com/aligoren/fiyatine/internal/models"
	"github.com/aligoren/fiyatine/internal/parsers"
	"github.com/aligoren/fiyatine/internal/utils"
)

type N11Service struct {
	SearchParams models.ProductSearchModel
}

func (service N11Service) buildUrl() string {
	requestUrl := url.URL{
		Scheme: "https",
		Host:   "n11.com",
		Path:   "arama",
	}

	query := requestUrl.Query()

	query.Add("q", service.SearchParams.ProductName)

	requestUrl.RawQuery = query.Encode()

	return requestUrl.String()
}

func (service N11Service) searchProduct() []models.ResponseModel {

	baseUrl := service.buildUrl()

	httpClient := utils.HttpClient{
		Method:  http.MethodGet,
		BaseUrl: baseUrl,
		Header: map[string]string{
			"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"referer":    "https://www.n11.com/",
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
		},
		Body: nil,
	}

	response, err := httpClient.MakeGet()
	if err != nil {
		log.Printf("error: %v", err)
	}

	parser := parsers.BaseParser{
		ParserService: parsers.N11Parser{
			Content: response.Body,
		},
	}

	return parser.Parse()

}
