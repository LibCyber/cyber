package country

import (
	"fmt"
	"strings"
)

var countryMap = map[string]string{
	"中国":    "CN",
	"美国":    "US",
	"日本":    "JP",
	"英国":    "GB",
	"法国":    "FR",
	"德国":    "DE",
	"意大利":   "IT",
	"西班牙":   "ES",
	"韩国":    "KR",
	"加拿大":   "CA",
	"澳大利亚":  "AU",
	"俄罗斯":   "RU",
	"印度":    "IN",
	"瑞士":    "CH",
	"荷兰":    "NL",
	"瑞典":    "SE",
	"巴西":    "BR",
	"波兰":    "PL",
	"土耳其":   "TR",
	"墨西哥":   "MX",
	"丹麦":    "DK",
	"泰国":    "TH",
	"挪威":    "NO",
	"奥地利":   "AT",
	"印度尼西亚": "ID",
	"比利时":   "BE",
	"伊朗":    "IR",
	"捷克":    "CZ",
	"南非":    "ZA",
	"芬兰":    "FI",
	"以色列":   "IL",
	"爱尔兰":   "IE",
	"新加坡":   "SG",
	"葡萄牙":   "PT",
	"希腊":    "GR",
	"新西兰":   "NZ",
	"阿根廷":   "AR",
	"匈牙利":   "HU",
	"沙特阿拉伯": "SA",
	"埃及":    "EG",
	"马来西亚":  "MY",
	"乌克兰":   "UA",
	"智利":    "CL",
	"菲律宾":   "PH",
	"卡塔尔":   "QA",
	"克罗地亚":  "HR",
	"阿联酋":   "AE",
	"卢森堡":   "LU",
	"巴基斯坦":  "PK",
	"罗马尼亚":  "RO",
	"哥伦比亚":  "CO",
	"越南":    "VN",
	"印尼":    "ID",
	"冰岛":    "IS",

	"香港": "HK",
	"台湾": "TW",
	"澳门": "MO",

	"深港": "HK",
	"沪日": "JP",
	"深台": "TW",
	"沪美": "US",
	"深新": "SG",

	// 添加更多的国家和相应的ISO 3166-1 Alpha-2代码...
}

func ParseCountry(input string) (string, error) {
	for country, code := range countryMap {
		if strings.Contains(input, country) {
			return code, nil
		}
	}
	return "", fmt.Errorf("no matching country found in input")
}
