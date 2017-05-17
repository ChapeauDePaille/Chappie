package main

import (
	"github.com/tealeg/xlsx"
)

//ReadXlsx : Lire dans un fichier excel
func ReadXlsx() (a []string) {
	var flag bool //Permet de ne pas prendre le titre de la colonne dans la liste d'adresses (on "saute" la 1ère itération)
	excelFileName := "/home/anton/Documents/Projets/Go/src/github.com/ChapeauDePaille/Chappie/ListeIP.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				if flag == true {
					text, _ := cell.String()
					//fmt.Printf("%s\n", text)
					a = append(a, text)
				}
				flag = true
			}
		}
	}
	return a
}
