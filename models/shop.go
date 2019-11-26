package models

type Shop struct {
	Name       string
	Company    string
	URL        string
	Сategories []Category
}

type Сategories struct {
	Сategories []Category `xml:"category"`
}

type Category struct {
	Id       string `xml:"id,attr"`
	Category string `xml:",chardata"`
}
