package plaid

import (
	"context"

	"github.com/gofiber/fiber/v2"
	plaid "github.com/plaid/plaid-go/v40/plaid"
)

func Liabilities(c *fiber.Ctx) error {
	ctx := context.Background()

	liabilitiesGetResp, _, err := client.PlaidApi.LiabilitiesGet(ctx).LiabilitiesGetRequest(
		// *plaid.NewAccountsGetRequest(accessToken),
		// access-sandbox-83f1ac58-dd31-4374-87f7-fd0c1e72eef4
		*plaid.NewLiabilitiesGetRequest(AMEX_TOKEN),
	).Execute()

	if err != nil {
		RenderError(c, err)
		return err
	}

	c.JSON(fiber.Map{
		"liabilities": liabilitiesGetResp.GetLiabilities().Credit,
	})
	return nil
}
