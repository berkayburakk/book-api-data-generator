package clients

import (
	"bytes"
	constant "datageneratorbookapi/constants"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	faker2 "syreclabs.com/go/faker"

	"github.com/spf13/cobra"
)

func PostBookRequest(requestInfo RequestInfo, bookRequest BookCreateRequest) string {

	var bookUrl = getBookApiUrl(requestInfo.Environment)

	requestBody, err := json.Marshal(bookRequest)
	cobra.CheckErr(err)

	request, err := http.NewRequest(
		http.MethodPost,
		bookUrl+"/books",
		bytes.NewBuffer(requestBody),
	)
	cobra.CheckErr(err)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	cobra.CheckErr(err)

	response, err := http.DefaultClient.Do(request)
	cobra.CheckErr(err)

	body, err := ioutil.ReadAll(response.Body)
	cobra.CheckErr(err)

	bookResponse := BookResponse{}

	if response.StatusCode == http.StatusCreated {
		if err := json.Unmarshal(body, &bookResponse); err != nil {
			fmt.Println(bookResponse)
			log.Printf("Could not unmarshal response -%v", err)
		}
	} else {
		errorResponse := ErrorModel{}

		if err := json.Unmarshal(body, &errorResponse); err != nil {
			log.Printf("HttpResponse Code -%v", response.StatusCode)
		}
	}

	fmt.Printf("Book ID : %d \n", bookResponse.ID)
	fmt.Printf("Book Barcode : %s \n", bookResponse.Barcode)

	id := strconv.Itoa(bookResponse.ID)

	WriteAll(id)

	os.Exit(1)
	fmt.Print(bookResponse.Barcode)
	return bookResponse.Barcode
}

func GetBooksRequest(requestInfo RequestInfo) BookListResponse {

	form := url.Values{
		"environment": {requestInfo.Environment},
	}

	var bookUrl = getBookApiUrl(requestInfo.Environment)

	request, err := http.NewRequest(
		http.MethodGet,
		bookUrl+"/books",
		strings.NewReader(form.Encode()),
	)

	request.Header.Add("Accept", "application/json")
	cobra.CheckErr(err)

	response, err := http.DefaultClient.Do(request)
	cobra.CheckErr(err)

	body, err := ioutil.ReadAll(response.Body)
	cobra.CheckErr(err)

	bookListResponse := BookListResponse{}

	if response.StatusCode == http.StatusOK {
		if err := json.Unmarshal(body, &bookListResponse); err != nil {
			log.Printf("Could not unmarshal response -%v", err)
		}
	} else {
		errorResponse := ErrorModel{}

		if err := json.Unmarshal(body, &errorResponse); err != nil {
			log.Printf("HttpResponse Code -%v", response.StatusCode)
		}
	}

	var formattedResponse string
	for _, book := range bookListResponse.Books {
		formattedResponse += fmt.Sprintf("Author: %s\n", book.Author)
		formattedResponse += fmt.Sprintf("Barcode: %s\n", book.Barcode)
		formattedResponse += fmt.Sprintf("Book Name: %s\n", book.BookName)
		formattedResponse += fmt.Sprintf("Category: %s\n", book.Category)
		formattedResponse += "-----------------\n"
	}

	fmt.Print(formattedResponse)

	return bookListResponse

}

func getBookApiUrl(environment string) string {
	var url = constant.UrlList[environment]

	if len(url) == 0 {
		fmt.Println("Environment error")
		os.Exit(1)
	}

	return url
}

type RequestInfo struct {
	Environment string `json:"environment"`
}

type BookCreateRequest struct {
	Author   string `json:"author"`
	Barcode  string `json:"barcode"`
	BookName string `json:"bookName"`
	Category string `json:"category"`
}

type ErrorModel struct {
	Error string `json:"error"`
}

type Message struct {
	TYPE    string `json:"type"`
	CONTENT string `json:"content"`
}

type BookResponse struct {
	Author   string `json:"Author"`
	Barcode  string `json:"Barcode"`
	BookName string `json:"BookName"`
	Category string `json:"Category"`
	ID       int    `json:"id"`
}

type BookListResponse struct {
	Books []struct {
		BookName string `json:"BookName"`
		Category string `json:"Category"`
		Author   string `json:"Author"`
		Barcode  string `json:"Barcode"`
	} `json:"books"`
}

func CreateBook() BookCreateRequest {
	return GenerateBookData()
}

func GenerateBookData() BookCreateRequest {

	bookRequestModel, err := NewBookCreateRequestBuilder().
		BookName(faker2.Name().Title()).
		Author(faker2.Name().Name()).
		Barcode(faker2.Code().Ean8()).
		Category(faker2.Commerce().Department()).
		Build()

	cobra.CheckErr(err)

	return *bookRequestModel

}

type BookCreateRequestBuilder struct {
	bookCreateRequest *BookCreateRequest
}

func NewBookCreateRequestBuilder() *BookCreateRequestBuilder {
	bookCreateRequest := &BookCreateRequest{}
	b := &BookCreateRequestBuilder{bookCreateRequest: bookCreateRequest}
	return b
}

func (b *BookCreateRequestBuilder) Author(author string) *BookCreateRequestBuilder {
	b.bookCreateRequest.Author = author
	return b
}

func (b *BookCreateRequestBuilder) Barcode(barcode string) *BookCreateRequestBuilder {
	b.bookCreateRequest.Barcode = barcode
	return b
}

func (b *BookCreateRequestBuilder) BookName(bookName string) *BookCreateRequestBuilder {
	b.bookCreateRequest.BookName = bookName
	return b
}

func (b *BookCreateRequestBuilder) Category(category string) *BookCreateRequestBuilder {
	b.bookCreateRequest.Category = category
	return b
}

func (b *BookCreateRequestBuilder) Build() (*BookCreateRequest, error) {
	return b.bookCreateRequest, nil
}
