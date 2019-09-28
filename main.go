package main

import (
	"chaos-stack-tesco/database"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Fetch struct {
	Uk struct {
		Ghs struct {
			Products struct {
				InputQuery  string `json:"input_query"`
				OutputQuery string `json:"output_query"`
				Filters     struct {
				} `json:"filters"`
				QueryPhase string `json:"queryPhase"`
				Totals     struct {
					All   int `json:"all"`
					New   int `json:"new"`
					Offer int `json:"offer"`
				} `json:"totals"`
				Config  string `json:"config"`
				Results []struct {
					Image                    string   `json:"image"`
					SuperDepartment          string   `json:"superDepartment"`
					Tpnb                     int      `json:"tpnb"`
					ContentsMeasureType      string   `json:"ContentsMeasureType"`
					Name                     string   `json:"name"`
					UnitOfSale               int      `json:"UnitOfSale"`
					AverageSellingUnitWeight float64  `json:"AverageSellingUnitWeight"`
					Description              []string `json:"description"`
					UnitQuantity             string   `json:"UnitQuantity"`
					ID                       int      `json:"id"`
					ContentsQuantity         int      `json:"ContentsQuantity"`
					Department               string   `json:"department"`
					Price                    float64  `json:"price"`
					Unitprice                float64  `json:"unitprice"`
				} `json:"results"`
				Suggestions []interface{} `json:"suggestions"`
			} `json:"products"`
		} `json:"ghs"`
	} `json:"uk"`
}

func main() {
	err := database.Init()
	if err != nil {
		fmt.Println("Error connecting to the database at startup")
		panic(err)
	}
	ApiResults, err := FetchTescoAPI("drink", 1, 5)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(ApiResults.Uk.Ghs.Products.Results)
	results := ApiResults.Uk.Ghs.Products.Results
	for _, v := range results {
		fmt.Println(v.Name)
	}
	r := gin.Default()
	r.Use(gin.Recovery())
	r.POST("/products/rate", func(c *gin.Context) {
		product := database.Product{}
		product.Name = c.PostForm("name")
		product.Rating, err = strconv.Atoi(c.PostForm("rating"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "rating is not a number"})
		}

		//todo: elllenorizni hogy a Tesco api-ban benne van-e a neve alapjan
		fetchResult, err := FetchTescoAPI(product.Name, 0, 1)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "Tesco Api returned a bad request"})
		}
		for _, v := range fetchResult.Uk.Ghs.Products.Results {
			if v.Name == "" {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "no results"})
			}
			if v.Name == product.Name {
				//if we find it in tesco api, then it's good, and we can rate it
				err = database.SaveProduct(product)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "failed to save to database"})
				}
			}
		}
	})
	r.GET("/products/", func(c *gin.Context) {
		name := c.Query("name")
		rating := c.Query("rating")
		//todo: lekerdezni az adatbazisbol nev es rating alapjan

	})
	log.Fatal(r.Run())

}

func FetchTescoAPI(query string, offset int, limit int) (Fetch, error) {
	ofs := strconv.Itoa(offset)
	l := strconv.Itoa(limit)
	token := "c6a9390ece40410dbdc5c3587eb78c3a"
	url := "https://dev.tescolabs.com/grocery/products/?query=" + query + "&offset=" + ofs + "&limit=" + l
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Fetch{}, err
	}

	req.Header.Add("Ocp-Apim-Subscription-Key", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Fetch{}, err
	}
	//fmt.Println(res)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Fetch{}, err
	}

	//fmt.Println(res)
	var grocery Fetch
	err = json.Unmarshal(body, &grocery)
	//fmt.Println(string(body))
	return grocery, nil
}
