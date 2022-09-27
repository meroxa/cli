package display

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

func AccountTable(account *meroxa.Account) string {
	mainTable := simpletable.New()
	mainTable.Body.Cells = [][]*simpletable.Cell{
		{
			{Align: simpletable.AlignRight, Text: "UUID:"},
			{Text: account.UUID},
		},
		{
			{Align: simpletable.AlignRight, Text: "Name:"},
			{Text: account.Name},
		},
	}

	mainTable.SetStyle(simpletable.StyleCompact)

	return mainTable.String()
}

func PrintAccountTable(account *meroxa.Account) {
	fmt.Println(AccountTable(account))
}

func AccountsTable(accounts []*meroxa.Account, hideHeaders bool) string {
	if len(accounts) != 0 {
		table := simpletable.New()

		if !hideHeaders {
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "UUID"},
					{Align: simpletable.AlignCenter, Text: "NAME"},
				},
			}
		}

		for _, p := range accounts {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: p.UUID},
				{Align: simpletable.AlignCenter, Text: p.Name},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		return table.String()
	}
	return ""
}

func PrintAccountsTable(accounts []*meroxa.Account, hideHeaders bool) {
	fmt.Println(AccountsTable(accounts, hideHeaders))
}
