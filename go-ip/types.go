package goip

type IpCountry struct {
	IsoCode       string `json:"iso_code"` //iso 编码 https://zh.m.wikipedia.org/zh/ISO_3166-1
	Name          string `json:"name"`
	Continent     string `json:"continent"`      //洲
	ContinentCode string `json:"continent_code"` //洲编码
	Isp           string `json:"isp"`            //isp
}

type IP2Locationrecord struct {
	Country_short      string
	Country_long       string
	Region             string
	City               string
	Isp                string
	Latitude           float32
	Longitude          float32
	Domain             string
	Zipcode            string
	Timezone           string
	Netspeed           string
	Iddcode            string
	Areacode           string
	Weatherstationcode string
	Weatherstationname string
	Mcc                string
	Mnc                string
	Mobilebrand        string
	Elevation          float32
	Usagetype          string
	Addresstype        string
	Category           string
	District           string
	Asn                string
	As                 string
}
