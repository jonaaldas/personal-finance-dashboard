package plaid

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	plaid "github.com/plaid/plaid-go/v40/plaid"
)

func Investments(c *fiber.Ctx) error {
	ctx := context.Background()

	request := plaid.NewInvestmentsTransactionsGetRequest("access-sandbox-1770c1ff-9e41-4fb0-ab5e-30207357c06c", "2025-01-01", "2025-01-31")

	investmentsResp, _, err := client.PlaidApi.InvestmentsTransactionsGet(ctx).InvestmentsTransactionsGetRequest(*request).Execute()

	if err != nil {
		fmt.Printf("Error getting investments: %v\n", err)
		RenderError(c, err)
		return err
	}

	c.JSON(fiber.Map{
		"investments": investmentsResp.GetInvestmentTransactions(),
	})

	return nil
}
