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

func AccountsTable(accounts []*meroxa.Account, currentAccount string, hideHeaders bool) string {
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
				{Align: simpletable.AlignCenter, Text: checkCurrent(p.Name, p.UUID, currentAccount)},
			}

			table.Body.Cells = append(table.Body.Cells, r)
		}
		table.SetStyle(simpletable.StyleCompact)
		return table.String()
	}
	return ""
}

func checkCurrent(name, uuid, currentUUID string) string {
	if uuid == currentUUID {
		return name + " (current)"
	}
	return name
}

func PrintAccountsTable(accounts []*meroxa.Account, currentAccount string, hideHeaders bool) {
	fmt.Println(AccountsTable(accounts, currentAccount, hideHeaders))
}
