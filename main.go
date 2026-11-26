package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/carlos7ags/folio/zugferd"
)

func main() {
	issueDate := time.Now().UTC()
	number := fmt.Sprintf("%s-%s", issueDate.Format("2006"), "4711")
	data := &zugferd.Invoice{
		Number:    number,
		IssueDate: issueDate,
		Currency:  "EUR",
		Seller: zugferd.Party{
			Name:  "Jens Schumacher",
			VATID: "DE123456789",
		},
		Buyer: zugferd.Party{
			Name:  "Max Mustermann",
			VATID: "DE987654321",
		},
		Lines: []zugferd.LineItem{
			{
				ID:              "1",
				Description:     "Produkt A",
				Quantity:        "10",
				UnitCode:        "C62",
				UnitPrice:       zugferd.NewAmount(1, 23),
				LineTotal:       zugferd.NewAmount(12, 30),
				TaxCategoryCode: "S",
				TaxRatePercent:  "19.00",
			},
		},
		TaxTotals: []zugferd.TaxBreakdown{
			{
				RatePercent:   "19.00",
				CategoryCode:  "S",
				TaxableAmount: zugferd.NewAmount(12, 30),
				TaxAmount:     zugferd.NewAmount(2, 34),
			},
		},
		Totals: zugferd.MonetarySummation{
			LineTotal:        zugferd.NewAmount(12, 30),
			TaxBasisTotal:    zugferd.NewAmount(12, 30),
			TaxTotal:         zugferd.NewAmount(2, 34),
			GrandTotal:       zugferd.NewAmount(14, 64),
			DuePayableAmount: zugferd.NewAmount(14, 64),
		},
		PaymentTerms: "-",
	}
	doc, err := BuildInvoice(data)
	if err != nil {
		log.Panic(err)
	}
	fileName := fmt.Sprintf("rechnung-%s.pdf", data.Number)
	err = doc.Save(fileName)
	if err != nil {
		log.Panic(fmt.Errorf("save '%s' failed: %w", fileName, err))
	}
	log.Printf("created '%s'", fileName)
	os.Exit(0)
}
