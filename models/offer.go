package models

type Offer struct {
	Id                    *string          `xml:"id,attr"`
	Type                  *string          `xml:"type,attr"`
	Available             *string          `xml:"available,attr"`
	Bid                   *string          `xml:"bid,attr"`
	Cbid                  *string          `xml:"cbid,attr"`
	Fee                   *string          `xml:"fee,attr"`
	Group_id              *int             `xml:"group_id,attr"`
	Url                   *string          `xml:"url"`
	Price                 *int             `xml:"price"`
	Oldprice              *int             `xml:"oldprice"`
	CurrencyId            *string          `xml:"currencyId"`
	CategoryId            *int             `xml:"categoryId"`
	Picture               *string          `xml:"picture"`
	Store                 *bool            `xml:"store"`
	Pickup                *bool            `xml:"pickup"`
	Delivery              *bool            `xml:"delivery"`
	Name                  *string          `xml:"name"`
	Vendor                *string          `xml:"vendor"`
	VendorCode            *string          `xml:"vendorCode"`
	Description           *string          `xml:"description"`
	Sales_notes           *string          `xml:"sales_notes"`
	Manufacturer_warranty *bool            `xml:"manufacturer_warranty"`
	Country_of_origin     *string          `xml:"country_of_origin"`
	Barcode               *string          `xml:"barcode"`
	TypePrefix            *string          `xml:"typePrefix"`
	Model                 *string          `xml:"model"`
	Adult                 *bool            `xml:"adult"`
	Expiry                *string          `xml:"expiry"`
	Weight                *string          `xml:"weight"`
	Dimensions            *string          `xml:"dimensions"`
	Condition             *Condition       `xml:"condition"`
	DeliveryOptions       *DeliveryOptions `xml:"delivery-options>option"`
	ShipmentOptions       *ShipmentOptions `xml:"shipment-options>option"`
	Age                   *Age             `xml:"age"`
	Param                 *Params          `xml:"param"`
}

type DeliveryOptions []struct {
	Cost        *string `xml:"cost,attr"`
	Days        *string `xml:"days,attr"`
	OrderBefore *string `xml:"order-before,attr"`
}

type ShipmentOptions []struct {
	Days        *string `xml:"days,attr"`
	OrderBefore *string `xml:"order-before,attr"`
}

type Params []struct {
	Name  *string `xml:"name,attr"`
	Unit  *string `xml:"unit,attr"`
	Param *string `xml:",chardata"`
}

type Condition struct {
	Type   *string   `xml:"type,attr"`
	Reason []*string `xml:"reason"`
}

type Age struct {
	Unit *string `xml:"unit,attr"`
	Age  *string `xml:",chardata"`
}
