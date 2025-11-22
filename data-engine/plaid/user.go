package plaid

import (
	"context"

	"github.com/gofiber/fiber/v2"
	plaid "github.com/plaid/plaid-go/v40/plaid"
)

func CreatePlaidUser(c *fiber.Ctx) {
	ctx := context.Background()

	city := "New Brunswick"
	region := "NJ"
	street := "19 S Ward St"
	postalCode := "08901"
	country := "US"

	addressData := plaid.AddressData{
		City:       *plaid.NewNullableString(&city),
		Region:     *plaid.NewNullableString(&region),
		Street:     street,
		PostalCode: *plaid.NewNullableString(&postalCode),
		Country:    *plaid.NewNullableString(&country),
	}
	DateOfBirth := "1995-03-28"
	emails := []string{"jonaaldas@gmail.com"}
	phoneNumbers := []string{"+17324852784"}

	firstName := "Jonathan"
	lastName := "Aldas"
	ssnLast4 := "1234"

	consumerReportUserIdentity := plaid.ConsumerReportUserIdentity{
		FirstName:      firstName,
		LastName:       lastName,
		PhoneNumbers:   phoneNumbers,
		Emails:         emails,
		SsnLast4:       *plaid.NewNullableString(&ssnLast4),
		DateOfBirth:    *plaid.NewNullableString(&DateOfBirth),
		PrimaryAddress: addressData,
	}

	request := plaid.NewUserCreateRequest(
		"c0e2c4ee-b763-4af5-cfe9-46a46bce883d",
	)

	request.SetConsumerReportUserIdentity(consumerReportUserIdentity)

	response, _, err := client.PlaidApi.UserCreate(ctx).UserCreateRequest(*request).Execute()
	if err != nil {
		return
	}
	_ = response
}
