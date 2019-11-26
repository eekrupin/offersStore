package modules

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/eekrupin/offersStore/db"
	"github.com/eekrupin/offersStore/models"
	"github.com/gocql/gocql"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

var ErrNoMoreOffer = errors.New("offers: no more rows")
var MaxWorkers = 1
var BatchSize = 250

func FileToCassandra(xmlfile string) (offers int, err error) {
	startTime := time.Now()

	xmlStream, err := os.Open(xmlfile)
	if err != nil {
		return 0, errors.New("failed to open XML file: " + err.Error())
	}
	defer xmlStream.Close()

	offers, err = XmlStreamToCassandra(xmlStream)

	log.Printf("Processed offers: %d, elapsed seconds: %f, file: %s", offers, time.Since(startTime).Seconds(), xmlfile)

	return
}

func HttpFileToCassandra(url string) (offers int, err error) {
	startTime := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return 0, errors.New("failed to open url file: " + err.Error())
	}
	defer resp.Body.Close()

	offers, err = XmlStreamToCassandra(resp.Body)

	log.Printf("Processed offers: %d, elapsed seconds: %f, url: %s", offers, time.Since(startTime).Seconds(), url)

	return
}

func XmlStreamToCassandra(xmlStream io.Reader) (offers int, err error) {

	decoder := xml.NewDecoder(xmlStream)

	date, err := parseDate(decoder, "yml_catalog", "date")
	if err != nil {
		return
	}
	_ = date

	shop, err := parseShop(decoder)
	if err != nil {
		return
	}
	//_ = shop

	//var offers []*models.Offer

	session, _ := db.DB.CreateSession()
	defer session.Close()

	c := make(chan int, MaxWorkers)
	var wg sync.WaitGroup
	rowsCount := 0
	var batch *gocql.Batch
Loop:
	for {
		offer, err := parseOffer(decoder)
		if err != nil && err != ErrNoMoreOffer {
			return offers, err
		} else if err == ErrNoMoreOffer {
			if rowsCount > 0 {
				wg.Add(1)
				go executeBatch(session, batch, c, &wg)
				rowsCount = 0
			}
			break Loop
		}
		//offers = append(offers, offer)
		if offer.Id == nil || *offer.Id == "" {
			continue
		}

		if rowsCount == 0 {
			batch = session.NewBatch(gocql.UnloggedBatch)
		}
		query := offerQueryInsert(shop, offer, session)
		values := query.Values()
		batch.Query(query.Statement(), values...)
		rowsCount++

		if rowsCount%BatchSize == 0 {
			wg.Add(1)
			go executeBatch(session, batch, c, &wg)
			rowsCount = 0
		}

		//wg.Add(1)
		//go saveOffer(shop, offer, session, c, &wg)

		offers++
	}

	wg.Wait()
	runtime.GC()
	debug.FreeOSMemory()

	return
}

func executeBatch(session *gocql.Session, batch *gocql.Batch, c chan int, wg *sync.WaitGroup) {
	c <- 1
	err := session.ExecuteBatch(batch)
	<-c
	wg.Done()
	if err != nil {
		log.Println(err)
	}
}

func saveOffer(shop *models.Shop, offer *models.Offer, session *gocql.Session, c chan int, wg *sync.WaitGroup) {

	c <- 1

	query := offerQueryInsert(shop, offer, session)

	err := query.Exec()
	<-c
	wg.Done()
	if err != nil {
		log.Printf("Name %s, id %d, %s", shop.Name, offer.Id, err)
	}
	offer = nil
}

func offerQueryInsert(shop *models.Shop, offer *models.Offer, session *gocql.Session) *gocql.Query {
	var Сategories []string
	for _, item := range shop.Сategories {
		Сategories = append(Сategories, item.Category)
	}
	var ConditionType *string
	var ConditionReason []string
	if offer.Condition != nil {
		ConditionType = offer.Condition.Type
		for _, reason := range offer.Condition.Reason {
			ConditionReason = append(ConditionReason, *reason)
		}
	}
	var DeliveryOptions []string
	if offer.DeliveryOptions != nil {
		for _, item := range *offer.DeliveryOptions {
			DeliveryOptions = append(DeliveryOptions, fmt.Sprintf("{Cost:%s, OrderBefore:%s, Cost:%s}", RepresentStringPointerJson(item.Days), RepresentStringPointerJson(item.OrderBefore), RepresentStringPointerJson(item.Cost)))
		}
	}
	var ShipmentOptions []string
	if offer.ShipmentOptions != nil {
		for _, item := range *offer.ShipmentOptions {
			ShipmentOptions = append(ShipmentOptions, fmt.Sprintf("{Cost:%s, OrderBefore:%s}", RepresentStringPointerJson(item.Days), RepresentStringPointerJson(item.OrderBefore)))
		}
	}
	var AgeUnit, Age *string
	if offer.Age != nil {
		AgeUnit = offer.Age.Unit
		Age = offer.Age.Age
	}
	var Param []string
	if offer.Param != nil {
		for _, item := range *offer.Param {
			Param = append(Param, fmt.Sprintf("{Name:%s, Param:%s, Unit:%s}", RepresentStringPointerJson(item.Name), RepresentStringPointerJson(item.Param), RepresentStringPointerJson(item.Unit)))
		}
	}
	query := session.Query(`INSERT INTO offers (
						shopname,
						shopcompany,
						shopurl,
						shopcategories,
						id,
			
						type,
						available,
						bid,
						cbid,
						fee,
			
						group_id,
						url,
						price,
						oldprice,
						currencyid,
			
						categoryid,
						picture,
						store,
						pickup,
						delivery,
			
						name,
						vendor,
						vendorcode,
						description,
						sales_notes,
			
						manufacturer_warranty,
						country_of_origin,
						barcode,
						typeprefix,
						model,
			
						adult,
						expiry,
						weight,
						dimensions,
						conditiontype,
			
						conditionreason,
						deliveryoptions,
						shipmentoptions,
						ageunit,
						age,
			
						param
			) VALUES ( 	?, ?, ?, ?, ?,
						?, ?, ?, ?, ?,
						?, ?, ?, ?, ?,
						?, ?, ?, ?, ?,
						?, ?, ?, ?, ?,
						?, ?, ?, ?, ?,
						?, ?, ?, ?, ?,
						?, ?, ?, ?, ?
						, ?)`,
		shop.Name, shop.Company, shop.URL, Сategories, offer.Id,
		offer.Type, offer.Available, offer.Bid, offer.Cbid, offer.Fee,
		offer.Group_id, offer.Url, offer.Price, offer.Oldprice, offer.CurrencyId,
		offer.CategoryId, offer.Picture, offer.Store, offer.Pickup, offer.Delivery,
		offer.Name, offer.Vendor, offer.VendorCode, offer.Description, offer.Sales_notes,
		offer.Manufacturer_warranty, offer.Country_of_origin, offer.Barcode, offer.TypePrefix, offer.Model,
		offer.Adult, offer.Expiry, offer.Weight, offer.Dimensions, ConditionType,
		ConditionReason, DeliveryOptions, ShipmentOptions, AgeUnit, Age,
		Param)
	return query
}

func parseShop(decoder *xml.Decoder) (*models.Shop, error) {
	shop := new(models.Shop)

Loop:
	for {
		t, _ := decoder.Token()
		if t == nil {
			return nil, errors.New("can't find next element")
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "name" {
				decodeElementPanic(decoder, &shop.Name, &se)
			} else if se.Name.Local == "company" {
				decodeElementPanic(decoder, &shop.Company, &se)
			} else if se.Name.Local == "url" {
				decodeElementPanic(decoder, &shop.URL, &se)
			} else if se.Name.Local == "categories" {
				var cats models.Сategories
				decodeElementPanic(decoder, &cats, &se)
				shop.Сategories = cats.Сategories
			} else if se.Name.Local == "offers" {
				break Loop
			}
		}
	}

	return shop, nil
}

func parseOffer(decoder *xml.Decoder) (*models.Offer, error) {
	offer := new(models.Offer)

	//Loop:
	for {
		t, _ := decoder.Token()
		if t == nil {
			return nil, ErrNoMoreOffer
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "offer" {
				decodeElementPanic(decoder, &offer, &se)
				return offer, nil
			}
		case xml.EndElement:
			if se.Name.Local == "offers" {
				return nil, ErrNoMoreOffer
			}
		}
	}

	return nil, ErrNoMoreOffer
}

func decodeElementPanic(decoder *xml.Decoder, v interface{}, start *xml.StartElement) {
	err := decoder.DecodeElement(v, start)
	if err != nil {
		errText := fmt.Sprint("categories", err)
		panic(errText)
	}
}

func parseDate(decoder *xml.Decoder, nodeName string, attrName string) (string, error) {
	value := ""

Loop:
	for {
		t, _ := decoder.Token()
		if t == nil {
			return "", errors.New("can't find start position")
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == nodeName {
				value = attribute(se.Attr, attrName)
				if value == "" {
					return "", errors.New(attrName + " is empty in " + nodeName)
				}
				break Loop
			} else {
				return "", errors.New("can't find start position " + nodeName)
			}
		}
	}

	return value, nil
}

func attribute(attrs []xml.Attr, name string) (value string) {
	for _, attr := range attrs {
		if attr.Name.Local == name {
			value = attr.Value
			return
		}
	}
	return
}
