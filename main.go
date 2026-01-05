package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	errorCodes []ErrorCode
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("WiiLink WFC Error API has encountered an error! Reason: %v\n", err)
	}
}

func lookupCode(c *gin.Context) {
	errorCode := c.Query("code")
	if errorCode == "" {
		c.JSON(http.StatusBadRequest, [1]ErrorResponse{{
			Error:    errorCode,
			Found:    0,
			InfoList: []ErrorInfo{},
		}})
		return
	}
	errorCodeRunes := []rune(errorCode)

	var matchingCodes []ErrorCode

	for _, codeInfo := range errorCodes {
		matches, err := regexp.MatchString(codeInfo.Regex, errorCode)
		checkError(err)

		if matches {
			// Replace placeholder characters
			cardRunes := []rune(codeInfo.Card)
			for index, character := range cardRunes {
				if character == '?' {
					cardRunes[index] = errorCodeRunes[index]
				}
			}
			codeInfo.Card = string(cardRunes)

			matchingCodes = append(matchingCodes, codeInfo)
		}
	}

	if len(matchingCodes) == 0 {
		c.JSON(http.StatusOK, [1]ErrorResponse{{
			Error:    errorCode,
			Found:    0,
			InfoList: []ErrorInfo{},
		}})
		return
	}

	var codeInfoList []ErrorInfo

	for index, codeInfo := range matchingCodes {
		// Last item found may have a full description
		if index == len(matchingCodes)-1 {
			if len(codeInfo.Description) != 0 {
				codeInfoList = append(codeInfoList, ErrorInfo{
					Type: "Error",
					Name: codeInfo.Card,
					Info: strings.Join(codeInfo.Description, "\n"),
				})
			} else {
				codeInfoList = append(codeInfoList, ErrorInfo{
					Type: "Error",
					Name: codeInfo.Card,
					Info: codeInfo.Comment,
				})
			}

			break
		}

		var itemType string
		switch index {
		case 0:
			itemType = "Class"
		case 1:
			itemType = "Section"
		default:
			itemType = "Group"
		}

		codeInfoList = append(codeInfoList, ErrorInfo{
			Type: itemType,
			Name: codeInfo.Card,
			Info: codeInfo.Comment,
		})
	}

	c.JSON(http.StatusOK, [1]ErrorResponse{{
		Error:    errorCode,
		Found:    1,
		InfoList: codeInfoList,
	}})
}

func main() {
	// Load the config
	rawConfig, err := os.ReadFile("./config.xml")
	checkError(err)

	config := &Config{}
	err = xml.Unmarshal(rawConfig, config)
	checkError(err)

	rawCodes, err := os.ReadFile("./error_codes.json")
	checkError(err)

	err = json.Unmarshal(rawCodes, &errorCodes)
	checkError(err)

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "https://wfc.wiilink.ca")
	})
	r.GET("/error", lookupCode)

	log.Fatal(r.Run(config.APIAddress))
}
