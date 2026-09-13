package payment

import (
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/webhook"
)

type StripeService struct {
	config *types.StripeConfig
}

func NewStripeService(sysConfig *types.SystemConfig) *StripeService {
	return &StripeService{config: &sysConfig.Payment.Stripe}
}

func (s *StripeService) UpdateConfig(config *types.StripeConfig) {
	s.config = config
}

func (s *StripeService) Pay(params PayRequest) (string, error) {
	if s.config == nil || !s.config.Enabled {
		return "", fmt.Errorf("stripe service disabled")
	}
	if s.config.SecretKey == "" {
		return "", fmt.Errorf("stripe secret key is empty")
	}

	stripe.Key = s.config.SecretKey
	currency := strings.ToLower(s.config.Currency)
	if currency == "" {
		currency = "usd"
	}

	amount := parseMoneyToMinorUnit(params.TotalFee, currency)
	if amount <= 0 {
		return "", fmt.Errorf("invalid stripe amount: %s", params.TotalFee)
	}

	successURL := strings.TrimRight(params.ReturnURL, "/")
	if successURL == "" {
		successURL = strings.TrimRight(s.config.Domain, "/")
	}
	cancelURL := successURL
	sessionParams := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(params.Subject),
					},
					UnitAmount: stripe.Int64(amount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successURL + "?session_id={CHECKOUT_SESSION_ID}&order_no=" + params.OutTradeNo),
		CancelURL:  stripe.String(cancelURL),
		Metadata: map[string]string{
			"order_no": params.OutTradeNo,
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"order_no": params.OutTradeNo,
			},
		},
	}

	sess, err := session.New(sessionParams)
	if err != nil {
		return "", fmt.Errorf("error with create stripe checkout session: %w", err)
	}

	return sess.URL, nil
}

func (s *StripeService) Query(outTradeNo string) (OrderInfo, error) {
	return OrderInfo{}, fmt.Errorf("stripe orders are verified by webhook")
}

func (s *StripeService) TradeVerify(request *http.Request) (OrderInfo, error) {
	if s.config == nil || !s.config.Enabled {
		return OrderInfo{}, fmt.Errorf("stripe service disabled")
	}
	payload, err := io.ReadAll(request.Body)
	if err != nil {
		return OrderInfo{}, fmt.Errorf("error with read stripe webhook body: %w", err)
	}
	body, err := webhook.ConstructEvent(payload, request.Header.Get("Stripe-Signature"), s.config.WebhookKey)
	if err != nil {
		return OrderInfo{}, fmt.Errorf("error with verify stripe webhook: %w", err)
	}

	switch body.Type {
	case "checkout.session.completed":
		var sess stripe.CheckoutSession
		if err := json.Unmarshal(body.Data.Raw, &sess); err != nil {
			return OrderInfo{}, fmt.Errorf("error with parse stripe session: %w", err)
		}
		if sess.Metadata["order_no"] == "" {
			return OrderInfo{}, fmt.Errorf("missing order number in stripe metadata")
		}
		payTime := time.Unix(body.Created, 0).Format("2006-01-02 15:04:05")
		tradeID := ""
		if sess.PaymentIntent != nil {
			tradeID = sess.PaymentIntent.ID
		}
		return OrderInfo{
			Status:     Success,
			OutTradeNo: sess.Metadata["order_no"],
			TradeId:    strings.TrimSpace(tradeID),
			Amount:     fmt.Sprintf("%.2f", float64(sess.AmountTotal)/100),
			PayTime:    payTime,
		}, nil
	default:
		return OrderInfo{}, fmt.Errorf("unsupported stripe event: %s", body.Type)
	}
}

func parseMoneyToMinorUnit(amount string, currency string) int64 {
	if amount == "" {
		return 0
	}
	if strings.EqualFold(currency, "jpy") {
		return int64(math.Round(parseFloat(amount)))
	}
	return int64(math.Round(parseFloat(amount) * 100))
}

func parseFloat(s string) float64 {
	var v float64
	_, _ = fmt.Sscanf(s, "%f", &v)
	return v
}

var _ PayService = (*StripeService)(nil)
