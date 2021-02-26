package recurly_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/fubotv/go-recurly"
	"github.com/google/go-cmp/cmp"
)

func TestGiftCards_Encoding(t *testing.T) {
	t.Parallel()

	moment, _ := time.Parse(recurly.DateTimeFormat, "2014-01-01T07:00:00Z")
	tests := []struct {
		v        recurly.GiftCard
		expected string
	}{
		{
			v: recurly.GiftCard{
				XMLName: xml.Name{Local: "gift_card"},
				ID:      2003020297591186183,
			},
			expected: MustCompactString(`
				<gift_card>
					<id>2003020297591186183</id>
				</gift_card>
			`),
		},
		{
			v: recurly.GiftCard{
				XMLName:  xml.Name{Local: "gift_card"},
				ID:       2003020297591186183,
				Delivery: &recurly.Delivery{},
			},
			expected: MustCompactString(`
				<gift_card>
					<id>2003020297591186183</id>
					<delivery></delivery>
				</gift_card>
			`),
		},
		{
			v: recurly.GiftCard{
				XMLName:           xml.Name{Local: "gift_card"},
				ID:                2003020297591186183,
				RedemptionCode:    "518822D87268C142",
				BalanceInCents:    2999,
				ProductCode:       "gift_card",
				UnitAmountInCents: 2999,
				Currency:          "USD",
				CreatedAt:         recurly.NewTime(moment),
				UpdatedAt:         recurly.NewTime(moment),
				DeliveredAt:       recurly.NewTime(moment),
				RedeemedAt:        recurly.NewTime(moment),
				CanceledAt:        recurly.NewTime(moment),
			},
			expected: MustCompactString(`
				<gift_card>
					<id>2003020297591186183</id>
					<redemption_code>518822D87268C142</redemption_code>
					<balance_in_cents>2999</balance_in_cents>
					<product_code>gift_card</product_code>
					<unit_amount_in_cents>2999</unit_amount_in_cents>
					<currency>USD</currency>
					<created_at>2014-01-01T07:00:00Z</created_at>
					<updated_at>2014-01-01T07:00:00Z</updated_at>
					<delivered_at>2014-01-01T07:00:00Z</delivered_at>
					<redeemed_at>2014-01-01T07:00:00Z</redeemed_at>
					<canceled_at>2014-01-01T07:00:00Z</canceled_at>
				</gift_card>
			`),
		},
	}

	for i, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("Encode/%d", i), func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)
			if err := xml.NewEncoder(buf).Encode(tt.v); err != nil {
				t.Fatal(err)
			} else if diff := cmp.Diff(buf.String(), tt.expected); diff != "" {
				t.Fatal(diff)
			}
		})

		t.Run(fmt.Sprintf("Decode/%d", i), func(t *testing.T) {
			t.Parallel()

			var g recurly.GiftCard
			if err := xml.Unmarshal([]byte(tt.expected), &g); err != nil {
				t.Fatal(err)
			} else if diff := cmp.Diff(tt.v, g); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestGiftCards_List(t *testing.T) {
	client, s := recurly.NewTestServer()
	defer s.Close()

	invocations := 0
	s.HandleFunc("GET", "/v2/gift_cards", func(w http.ResponseWriter, r *http.Request) {
		invocations++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(MustOpenFile("gift_cards.xml"))
	}, t)

	pager := client.GiftCards.List(nil)
	for pager.Next() {
		var giftCards []recurly.GiftCard
		if err := pager.Fetch(context.Background(), &giftCards); err != nil {
			t.Fatal(err)
		} else if !s.Invoked {
			t.Fatal("expected to be invoked")
		} else if diff := cmp.Diff(giftCards, NewTestGiftCards()); diff != "" {
			t.Fatal(diff)
		}
	}
}

func NewTestGiftCards() []recurly.GiftCard {
	return []recurly.GiftCard{
		{
			XMLName:           xml.Name{Local: "gift_card"},
			ID:                2003020297591186183,
			RedemptionCode:    "518822D87268C142",
			BalanceInCents:    500,
			ProductCode:       "gift_card",
			UnitAmountInCents: 1000,
			Currency:          "USD",
			Delivery: &recurly.Delivery{
				XMLName:      xml.Name{Local: "delivery"},
				Method:       "post",
				EmailAddress: "john@example.com",
				FirstName:    "John",
				LastName:     "Smith",
				Address: &recurly.Address{
					XMLName: xml.Name{Local: "address"},
					Address: "123 B St.",
					City:    "San Francisco",
					State:   "CA",
					Zip:     "94110",
					Country: "USA",
				},
				GifterName:      "Sally",
				PersonalMessage: "\n                Hi John, Happy Birthday! I hope you have a great day! Love, Sally",
			},
			CreatedAt:  recurly.NewTime(MustParseTime("2016-07-26T15:23:46Z")),
			UpdatedAt:  recurly.NewTime(MustParseTime("2016-07-29T04:25:39Z")),
			RedeemedAt: recurly.NewTime(MustParseTime("2016-07-29T04:25:38Z")),
		},
		{
			XMLName:           xml.Name{Local: "gift_card"},
			ID:                1988596186827727838,
			RedemptionCode:    "3E687AE878D37EBD",
			ProductCode:       "gift_card",
			UnitAmountInCents: 1000,
			Currency:          "USD",
			Delivery: &recurly.Delivery{
				XMLName:         xml.Name{Local: "delivery"},
				Method:          "email",
				EmailAddress:    "jill@example.com",
				Address:         &recurly.Address{XMLName: xml.Name{Local: "address"}},
				FirstName:       "Jill",
				LastName:        "Wilson",
				PersonalMessage: "\n                Happy Holidays!",
			},
			CreatedAt:   recurly.NewTime(MustParseTime("2016-12-14T15:23:46Z")),
			UpdatedAt:   recurly.NewTime(MustParseTime("2016-12-14T15:23:46Z")),
			DeliveredAt: recurly.NewTime(MustParseTime("2016-12-14T15:23:46Z")),
		},
	}
}