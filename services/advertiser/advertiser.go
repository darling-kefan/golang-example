package advertiser

import (
	"strconv"
	
	"models"
)

type advertiser struct {
	id int
}

func New(adverid int) *advertiser {
	return &advertiser{id: adverid}
}

func (a *advertiser) HasSyncCondition() (ret bool) {
	qualifis := models.Advertisers().Qualifis([]string{strconv.Itoa(a.id)})
	// 是否为新广告主
	isNew := false
	// 审核通过的资质
	quatsMap := make(map[string]bool, 10)
	for _, qualifi := range qualifis {
		if qualifi["qualifis_on_type"] == "2" {
			switch qualifi["qualifis_type"] {
			case "1", "2", "7" :
				isNew = true
			}
		}
		if qualifi["qualifis_status"] == "2" {
			quatsMap[qualifi["qualifis_type"]] = true
		}
	}
	if isNew {
		_, ok1 := quatsMap["1"];
		_, ok2 := quatsMap["2"];
		_, ok3 := quatsMap["7"];
		if ok1 && ok2 && ok3 {
			ret = true
		}
	} else {
		_, ok1 := quatsMap["1"];
		_, ok2 := quatsMap["2"];
		if ok1 && ok2 {
			ret = true
		}
	}
	return
}
