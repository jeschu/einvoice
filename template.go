package main

import (
	_ "embed" // Blank-Import für die Direktive
	"fmt"
	"html/template"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/carlos7ags/folio/document"
	"github.com/carlos7ags/folio/html"
	"github.com/carlos7ags/folio/layout"
	"github.com/carlos7ags/folio/zugferd"
)

//go:embed invoice-template.html
var templateContent string

var (
	tmpl     *template.Template
	fontPath string
)

func init() {
	var err error
	if tmpl, err = template.New("einvoice").Parse(templateContent); err != nil {
		log.Panic(err)
	}
	fontPath = findSystemFont()
}

func fillInvoiceTemplate(data *zugferd.Invoice) (string, error) {
	d := make(map[string]any)
	d["fontPath"] = fontPath
	d["invoice"] = data
	var b strings.Builder
	err := tmpl.Execute(&b, d)
	return b.String(), err
}

func findSystemFont() string {
	var candidates []string
	switch runtime.GOOS {
	case "darwin":
		candidates = []string{
			"/System/Library/Fonts/Supplemental/Verdana.ttf",
			"/System/Library/Fonts/Supplemental/Arial.ttf",
			"/System/Library/Fonts/Supplemental/Georgia.ttf",
		}
	case "linux":
		candidates = []string{
			"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
			"/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf",
			"/usr/share/fonts/dejavu/DejaVuSans.ttf",
		}
	case "windows":
		candidates = []string{
			`C:\Windows\Fonts\verdana.ttf`,
			`C:\Windows\Fonts\arial.ttf`,
		}
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	log.Panic(fmt.Errorf("no suitable system font found for PDF/A embedding"))
	return ""
}

func BuildInvoice(data *zugferd.Invoice) (*document.Document, error) {
	invoiceHtml, err := fillInvoiceTemplate(data)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("==== 8< invoiceHtml 8< ====")
	fmt.Println(invoiceHtml)
	fmt.Println("---- >8 invoiceHtml >8 ----")

	elems, err := html.Convert(invoiceHtml, &html.Options{AllowAbsolutePaths: true})
	if err != nil {
		return nil, fmt.Errorf("html: %w", err)
	}
	doc := document.NewDocument(document.PageSizeA4)
	doc.SetMargins(layout.Margins{Top: 40, Right: 40, Bottom: 40, Left: 40})
	doc.Info.Title = fmt.Sprintf("Rechnung %s", data.Number)
	doc.Info.Author = "Jens Schumacher"
	for _, e := range elems {
		doc.Add(e)
	}

	err = data.Attach(doc, zugferd.ProfileBasic)
	return doc, err
}
