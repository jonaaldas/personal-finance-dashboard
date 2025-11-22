package plaid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	plaid "github.com/plaid/plaid-go/v40/plaid"
)

func GetAccessToken(c *fiber.Ctx) error {
	var body struct {
		PublicToken string `json:"publicToken"`
	}

	if err := c.BodyParser(&body); err != nil {
		fmt.Printf("Error parsing request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.PublicToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "publicToken is required",
		})
	}

	ctx := context.Background()

	exchangePublicTokenResp, _, err := client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
		*plaid.NewItemPublicTokenExchangeRequest(body.PublicToken),
	).Execute()
	if err != nil {
		fmt.Printf("Error exchanging public token: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	accessToken = exchangePublicTokenResp.GetAccessToken()
	itemID = exchangePublicTokenResp.GetItemId()

	fmt.Println("public token: " + body.PublicToken)
	fmt.Println("access token: " + accessToken)
	fmt.Println("item ID: " + itemID)

	Accounts(c)

	return c.JSON(fiber.Map{
		"access_token": accessToken,
		"item_id":      itemID,
	})
}

// linkTokenCreate creates a link token using the specified parameters
func LinkTokenCreate(
	paymentInitiation *plaid.LinkTokenCreateRequestPaymentInitiation,
) (string, error) {
	ctx := context.Background()

	// Institutions from all listed countries will be shown.
	countryCodes := ConvertCountryCodes(strings.Split(PLAID_COUNTRY_CODES, ","))
	fmt.Print(countryCodes)
	redirectURI := PLAID_REDIRECT_URI

	// This should correspond to a unique id for the current user.
	// Typically, this will be a user ID number from your application.
	// Personally identifiable information, such as an email address or phone number, should not be used here.
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: time.Now().String(),
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Plaid Quickstart",
		"en",
		countryCodes,
	)
	request.SetUser(user)

	products := ConvertProducts(strings.Split(PLAID_PRODUCTS, ","))
	fmt.Printf("Products being sent: %v\n", products)
	if paymentInitiation != nil {
		request.SetPaymentInitiation(*paymentInitiation)
		request.SetProducts([]plaid.Products{plaid.PRODUCTS_PAYMENT_INITIATION})
	} else {
		request.SetProducts(products)
	}

	if ContainsProduct(products, plaid.PRODUCTS_STATEMENTS) {
		statementConfig := plaid.NewLinkTokenCreateRequestStatements(
			time.Now().Local().Add(-30*24*time.Hour).Format("2006-01-02"),
			time.Now().Local().Format("2006-01-02"),
		)
		request.SetStatements(*statementConfig)
	}

	if ContainsProduct(products, plaid.PRODUCTS_CRA_BASE_REPORT) ||
		ContainsProduct(products, plaid.PRODUCTS_CRA_INCOME_INSIGHTS) ||
		ContainsProduct(products, plaid.PRODUCTS_CRA_PARTNER_INSIGHTS) {
		request.SetUserToken(userToken)
		request.SetConsumerReportPermissiblePurpose(plaid.CONSUMERREPORTPERMISSIBLEPURPOSE_ACCOUNT_REVIEW_CREDIT)
		request.SetCraOptions(*plaid.NewLinkTokenCreateRequestCraOptions(60))
	}

	if redirectURI != "" {
		request.SetRedirectUri(redirectURI)
	}

	linkTokenCreateResp, httpResp, err := client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

	if err != nil {
		fmt.Printf("Error creating link token2: %v\n", err)
		if httpResp != nil {
			fmt.Printf("HTTP Status: %d\n", httpResp.StatusCode)
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					fmt.Printf("Response Body: %s\n", string(bodyBytes))
				}
			}
		}
		return "", err
	}

	return linkTokenCreateResp.GetLinkToken(), nil
}

func ConvertCountryCodes(countryCodeStrs []string) []plaid.CountryCode {
	countryCodes := []plaid.CountryCode{}

	for _, countryCodeStr := range countryCodeStrs {
		countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
	}

	return countryCodes
}

func ConvertProducts(productStrs []string) []plaid.Products {
	products := []plaid.Products{}
	productMap := map[string]plaid.Products{
		"transactions":         plaid.PRODUCTS_TRANSACTIONS,
		"liabilities":          plaid.PRODUCTS_LIABILITIES,
		"assets":               plaid.PRODUCTS_ASSETS,
		"auth":                 plaid.PRODUCTS_AUTH,
		"balances":             plaid.PRODUCTS_BALANCE,
		"identity":             plaid.PRODUCTS_IDENTITY,
		"investments":          plaid.PRODUCTS_INVESTMENTS,
		"payment_initiation":   plaid.PRODUCTS_PAYMENT_INITIATION,
		"statements":           plaid.PRODUCTS_STATEMENTS,
		"cra_base_report":      plaid.PRODUCTS_CRA_BASE_REPORT,
		"cra_income_insights":  plaid.PRODUCTS_CRA_INCOME_INSIGHTS,
		"cra_partner_insights": plaid.PRODUCTS_CRA_PARTNER_INSIGHTS,
	}

	for _, productStr := range productStrs {
		productStr = strings.TrimSpace(productStr)
		productStr = strings.ToLower(productStr)
		if product, ok := productMap[productStr]; ok {
			products = append(products, product)
		} else {
			fmt.Printf("Warning: Unknown product '%s', attempting direct conversion\n", productStr)
			products = append(products, plaid.Products(productStr))
		}
	}

	return products
}

func ContainsProduct(products []plaid.Products, product plaid.Products) bool {
	for _, p := range products {
		if p == product {
			return true
		}
	}
	return false
}

func RenderError(c *fiber.Ctx, originalErr error) {
	if plaidError, err := plaid.ToPlaidError(originalErr); err == nil {
		// Return 200 and allow the front end to render the error.

		c.JSON(fiber.Map{
			"error": plaidError,
		})
		return
	}

	c.JSON(fiber.Map{
		"error": originalErr.Error(),
	})
}

func SaveAccount(accountsGetResp *plaid.AccountsGetResponse) {
	jsonConvertedRes, err := json.Marshal(accountsGetResp)
	if err != nil {
		log.Fatalf("Failed to Serialize to JSON from native Go struct type: %v", err)
	}

	// here http.Post method expects body as 'io.Redear' which should implement Read() method.
	// So, bytes package will take care of that.
	url := fmt.Sprintf("%s/api/save", CONVEX_URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonConvertedRes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// Set the Content-Type header to indicate that we are sending JSON
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Define a struct to match the response format
	type SaveResponse struct {
		Message string `json:"message"`
	}

	// Parse the JSON response
	var saveResp SaveResponse
	if err := json.Unmarshal(body, &saveResp); err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		fmt.Printf("Response body: %s\n", string(body))
		return
	}
	fmt.Printf("Account saved successfully: %s\n", saveResp.Message)
}

type AccessTokenResponse struct {
	AccessTokens []string `json:"access_tokens"`
}

func GetAllAccessTokens() (AccessTokenResponse, error) {
	url := fmt.Sprintf("%s/api/access_token", CONVEX_URL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return AccessTokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	fmt.Println(err)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return AccessTokenResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return AccessTokenResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return AccessTokenResponse{}, err
	}

	var accessTokenResponse AccessTokenResponse
	if err := json.Unmarshal(body, &accessTokenResponse); err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		fmt.Printf("Response body: %s\n", string(body))
		return AccessTokenResponse{}, err
	}

	return accessTokenResponse, nil

}
