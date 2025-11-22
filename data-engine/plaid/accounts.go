package plaid

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	plaid "github.com/plaid/plaid-go/v40/plaid"
)

func Accounts(c *fiber.Ctx) error {
	ctx := context.Background()
	fmt.Println("accessToken: " + accessToken)
	accountsGetResp, _, err := client.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		// *plaid.NewAccountsGetRequest(accessToken),
		// access-sandbox-3b63e2a9-1d22-4139-afd0-dda1b3bd07e8 -> USAA
		*plaid.NewAccountsGetRequest(accessToken),
	).Execute()

	if err != nil {
		RenderError(c, err)
		return err
	}

	SaveAccount(&accountsGetResp)

	c.JSON(fiber.Map{
		"accounts": accountsGetResp.GetAccounts(),
	})
	return nil
}
