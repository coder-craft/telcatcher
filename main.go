package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"./opesun/goquery"
)

var config = &Config{
	Url:   "https://zhengzhou.anjuke.com/sale/p%v",
	Count: 50,
}

type Config struct {
	Url   string
	Count int
}

func LoadConfig() bool {
	buff, err := ioutil.ReadFile("Config.json")
	if err != nil {
		fmt.Println("Load config error:", err)
	} else {
		err = json.Unmarshal(buff, config)
		if err != nil {
			fmt.Println("Unmarshal config file error:", err)
		}
		return true
	}
	return false
}
func main() {
	if LoadConfig() == false {
		return
	}
	file := Contacts{}
	for i := 1; i < config.Count; i++ {
		baseUrl := fmt.Sprintf(config.Url, i)
		page, err := goquery.ParseUrl(baseUrl)
		if err != nil {
			fmt.Printf("Get base url %v error.\n", baseUrl, err)
			continue
		}
		title := page.Find(".house-title a")
		//fmt.Println(title)
		for i := 0; i < title.Length(); i++ {
			time.Sleep(time.Second)
			detailUrl := title.Eq(i).Attr("href")
			fmt.Println(detailUrl)
			detail, detailErr := goquery.ParseUrl(detailUrl)
			if detailErr != nil {
				fmt.Println("Detail err:", detailErr)
			} else {
				contact := Contact{
					Starred:   Starred{InnerText: "0"},
					Account:   Account{Value: "0"},
					GroupList: GroupList{Group: []Group{{InnerText: "二手房中介"}}},
				}
				name := detail.Find(".brokercard-name").Text()
				fmt.Println("Name:", name)
				contact.CName = Name{InnerText: name}
				tel := detail.Find(".broker-mobile").Text()
				strTel := ""
				for i := 0; i < len(tel); i++ {
					if tel[i] >= '0' && tel[i] <= '9' {
						strTel = strTel + string(tel[i])
					}
				}
				fmt.Println("Tel:", strTel)
				contact.PhoneList = PhoneList{
					Phone: []Phone{Phone{Type: "2", InnerText: strTel}},
				}
				coDetail := detail.Find(".broker-company")
				coDetail = coDetail.Find("p")
				orl := Organization{
					Type: "1",
				}
				or := ""
				co := ""
				for i := 0; i < coDetail.Length(); i++ {
					node := coDetail.Eq(i)
					if i == 0 {
						or = node.Find("a").Text()
						orl.InnerText = or
					}
					if i == 1 {
						co = node.Find("a").Text()
						orl.Detail = append(orl.Detail, Detail{
							Company: []Company{
								Company{InnerText: co},
							},
						})
					}
				}
				fmt.Printf("Name:%v,Tel:%v,Or:%v,Co:%v\n", name, strTel, or, co)
				contact.OrganizationList.Organization = append(contact.OrganizationList.Organization, orl)
				file.Contact = append(file.Contact, contact)
				buff, err := xml.Marshal(&file)
				if err != nil {
					fmt.Println(err)
				}
				ioutil.WriteFile("out.xml", buff, os.FileMode(0))
			}
		}
	}
	fmt.Println("Catch complete.")
}

type Contacts struct {
	Contact []Contact `xml:"Contact"`
}
type Contact struct {
	CName            Name             `xml:"Name"`
	Starred          Starred          `xml:"Starred"`
	PhoneList        PhoneList        `xml:"PhoneList"`
	OrganizationList OrganizationList `xml:"OrganizationList"`
	Account          Account          `xml:Account`
	GroupList        GroupList        `xml:"GroupList"`
}
type Name struct {
	InnerText string `xml:",innerxml"`
}
type Starred struct {
	InnerText string `xml:",innerxml"`
}
type PhoneList struct {
	Phone []Phone
}
type Phone struct {
	Type      string `xml:"Type,attr"`
	InnerText string `xml:",innerxml"`
}
type OrganizationList struct {
	Organization []Organization `xml:"Organization"`
}
type Organization struct {
	Type      string   `xml:"Type,attr"`
	InnerText string   `xml:",innerxml"`
	Detail    []Detail `xml:"Detail"`
}
type Detail struct {
	Company []Company `xml:"Company"`
}
type Company struct {
	InnerText string `xml:",innerxml"`
}
type Account struct {
	Value string `xml:"value,attr"`
	Name  Name   `xml:Name`
	Type  Type   `xml:Type`
}
type Type struct {
}
type GroupList struct {
	Group []Group `xml:"GroupName"`
}
type Group struct {
	InnerText string `xml:",innerxml"`
}
