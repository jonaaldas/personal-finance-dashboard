package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	plaid "github.com/plaid/plaid-go/v40/plaid"
)

var (
	PLAID_CLIENT_ID                      = ""
	PLAID_SECRET                         = ""
	PLAID_ENV                            = ""
	PLAID_PRODUCTS                       = ""
	PLAID_COUNTRY_CODES                  = ""
	PLAID_REDIRECT_URI                   = ""
	SIGNAL_RULESET_KEY                   = ""
	APP_PORT                             = ""
	client              *plaid.APIClient = nil
)

var environments = map[string]plaid.Environment{
	"sandbox":    plaid.Sandbox,
	"production": plaid.Production,
}

func init() {
	// load env vars from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error when loading environment variables from .env file %w", err)
	}

	// set constants from env
	PLAID_CLIENT_ID = os.Getenv("PLAID_CLIENT_ID")
	PLAID_SECRET = os.Getenv("PLAID_SECRET")

	if PLAID_CLIENT_ID == "" || PLAID_SECRET == "" {
		log.Fatal("Error: PLAID_SECRET or PLAID_CLIENT_ID is not set. Did you copy .env.example to .env and fill it out?")
	}

	PLAID_ENV = os.Getenv("PLAID_ENV")
	PLAID_PRODUCTS = os.Getenv("PLAID_PRODUCTS")
	PLAID_COUNTRY_CODES = os.Getenv("PLAID_COUNTRY_CODES")
	PLAID_REDIRECT_URI = os.Getenv("PLAID_REDIRECT_URI")

	// set defaults
	if PLAID_PRODUCTS == "" {
		PLAID_PRODUCTS = "transactions, liabilities"
	}
	if PLAID_COUNTRY_CODES == "" {
		PLAID_COUNTRY_CODES = "US"
	}
	if PLAID_ENV == "" {
		PLAID_ENV = "sandbox"
	}

	if PLAID_CLIENT_ID == "" {
		log.Fatal("PLAID_CLIENT_ID is not set. Make sure to fill out the .env file")
	}
	if PLAID_SECRET == "" {
		log.Fatal("PLAID_SECRET is not set. Make sure to fill out the .env file")
	}

	// create Plaid client
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", PLAID_CLIENT_ID)
	configuration.AddDefaultHeader("PLAID-SECRET", PLAID_SECRET)
	configuration.UseEnvironment(environments[PLAID_ENV])
	client = plaid.NewAPIClient(configuration)
}

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"Hello": "World",
		})
	})

	app.Post("/api/set_access_token", func(c *fiber.Ctx) error {
		return getAccessToken(c)
	})

	app.Post("/api/create_link_token", func(c *fiber.Ctx) error {
		linkToken, err := linkTokenCreate(nil)

		if err != nil {
			fmt.Printf("Error creating link token1: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"link_token": linkToken,
		})
	})

	app.Get("/api/accounts", func(c *fiber.Ctx) error {
		return accounts(c)
	})

	app.Get("/api/liabilities", func(c *fiber.Ctx) error {
		return liabilities(c)
	})

	log.Fatal(app.Listen(":" + getPort()))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}

var accessToken string
var userToken string
var itemID string

var paymentID string

// The authorizationID is only relevant for the Transfer ACH product.
// We store the authorizationID in memory - in production, store it in a secure
// persistent data store
var authorizationID string
var accountID string

func getAccessToken(c *fiber.Ctx) error {
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

	return c.JSON(fiber.Map{
		"access_token": accessToken,
		"item_id":      itemID,
	})
}

// linkTokenCreate creates a link token using the specified parameters
func linkTokenCreate(
	paymentInitiation *plaid.LinkTokenCreateRequestPaymentInitiation,
) (string, error) {
	ctx := context.Background()

	// Institutions from all listed countries will be shown.
	countryCodes := convertCountryCodes(strings.Split(PLAID_COUNTRY_CODES, ","))
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

	products := convertProducts(strings.Split(PLAID_PRODUCTS, ","))
	fmt.Printf("Products being sent: %v\n", products)
	if paymentInitiation != nil {
		request.SetPaymentInitiation(*paymentInitiation)
		request.SetProducts([]plaid.Products{plaid.PRODUCTS_PAYMENT_INITIATION})
	} else {
		request.SetProducts(products)
	}

	if containsProduct(products, plaid.PRODUCTS_STATEMENTS) {
		statementConfig := plaid.NewLinkTokenCreateRequestStatements(
			time.Now().Local().Add(-30*24*time.Hour).Format("2006-01-02"),
			time.Now().Local().Format("2006-01-02"),
		)
		request.SetStatements(*statementConfig)
	}

	if containsProduct(products, plaid.PRODUCTS_CRA_BASE_REPORT) ||
		containsProduct(products, plaid.PRODUCTS_CRA_INCOME_INSIGHTS) ||
		containsProduct(products, plaid.PRODUCTS_CRA_PARTNER_INSIGHTS) {
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

func convertCountryCodes(countryCodeStrs []string) []plaid.CountryCode {
	countryCodes := []plaid.CountryCode{}

	for _, countryCodeStr := range countryCodeStrs {
		countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
	}

	return countryCodes
}

func convertProducts(productStrs []string) []plaid.Products {
	products := []plaid.Products{}
	productMap := map[string]plaid.Products{
		"transactions":         plaid.PRODUCTS_TRANSACTIONS,
		"liabilities":          plaid.PRODUCTS_LIABILITIES,
		"assets":               plaid.PRODUCTS_ASSETS,
		"auth":                 plaid.PRODUCTS_AUTH,
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

func containsProduct(products []plaid.Products, product plaid.Products) bool {
	for _, p := range products {
		if p == product {
			return true
		}
	}
	return false
}

func renderError(c *fiber.Ctx, originalErr error) {
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

func accounts(c *fiber.Ctx) error {
	ctx := context.Background()

	accountsGetResp, _, err := client.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		// *plaid.NewAccountsGetRequest(accessToken),
		*plaid.NewAccountsGetRequest("access-sandbox-41aa2e07-3853-4287-8a1f-8664465aaeb2"),
	).Execute()

	if err != nil {
		renderError(c, err)
		return err
	}

	c.JSON(fiber.Map{
		"accounts": accountsGetResp.GetAccounts(),
	})
	return nil
}

func liabilities(c *fiber.Ctx) error {
	ctx := context.Background()

	liabilitiesGetResp, _, err := client.PlaidApi.LiabilitiesGet(ctx).LiabilitiesGetRequest(
		// *plaid.NewAccountsGetRequest(accessToken),
		// access-sandbox-83f1ac58-dd31-4374-87f7-fd0c1e72eef4
		*plaid.NewLiabilitiesGetRequest("access-sandbox-83f1ac58-dd31-4374-87f7-fd0c1e72eef4"),
	).Execute()

	if err != nil {
		renderError(c, err)
		return err
	}

	c.JSON(fiber.Map{
		"liabilities": liabilitiesGetResp.GetLiabilities().Credit,
	})
	return nil
}
