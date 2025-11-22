package plaid

import (
	"context"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	plaid "github.com/plaid/plaid-go/v40/plaid"
)

func Transactions(c *fiber.Ctx) error {
	ctx := context.Background()

	// Set cursor to empty to receive all historical updates
	var cursor *string

	// New transaction updates since "cursor"
	var added []plaid.Transaction
	var modified []plaid.Transaction
	var removed []plaid.RemovedTransaction // Removed transaction ids
	hasMore := true
	// Iterate through each page of new transaction updates for item
	for hasMore {
		// request := plaid.NewTransactionsSyncRequest(accessToken)
		request := plaid.NewTransactionsSyncRequest(AMEX_TOKEN)
		if cursor != nil {
			request.SetCursor(*cursor)
		}
		resp, _, err := client.PlaidApi.TransactionsSync(
			ctx,
		).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			RenderError(c, err)
			return nil
		}

		// Update cursor to the next cursor
		nextCursor := resp.GetNextCursor()
		cursor = &nextCursor

		// If no transactions are available yet, wait and poll the endpoint.
		// Normally, we would listen for a webhook, but the Quickstart doesn't
		// support webhooks. For a webhook example, see
		// https://github.com/plaid/tutorial-resources or
		// https://github.com/plaid/pattern

		if *cursor == "" {
			time.Sleep(2 * time.Second)
			continue
		}

		// Add this page of results
		added = append(added, resp.GetAdded()...)
		modified = append(modified, resp.GetModified()...)
		removed = append(removed, resp.GetRemoved()...)
		hasMore = resp.GetHasMore()
	}

	sort.Slice(added, func(i, j int) bool {
		return added[i].GetDate() < added[j].GetDate()
	})

	start := len(added) - 9
	if start < 0 {
		start = 0
	}
	latestTransactions := added[start:]

	c.JSON(fiber.Map{
		"latest_transactions": latestTransactions,
	})
	return nil
}
