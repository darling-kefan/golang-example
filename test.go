package main

import (
	"fmt"
	"models"
	pubservice "services/ztc/publish"
	_ "strconv"
	_ "reflect"
)

func main() {
	// publishs := models.Publishs().JoinPubsByIds([]int{22983,22969})
	publishs := models.Publishs().JoinPubsByIds([]int{22983})
	pub := pubservice.NewPub(publishs[0])
	fmt.Println(pub)

	//publishs = models.Publishs().PubsByIds([]string{"1000", "22969", "22974"})
	//fmt.Println(publishs)

	//fmt.Println("------------------------------------------------------")
	//publishs := models.Publishs().LastestPubversion(strconv.Itoa(22983), 2)
	//fmt.Println(publishs)

	//fmt.Println("------------------------------------------------------")

	//cards := models.Cards().CardsByIds([]string{"160", "161", "177"})
	//fmt.Println(cards)

	//fmt.Println("------------------------------------------------------")

	//advertisers := models.Advertisers().AdversByIds([]string{"1", "3", "4", "5", "6"}, []string{"id", "tvmid", "name", "mobileno", "email", "type", "subtype"})
	//fmt.Println(advertisers)

	//advQualifis := models.Advertisers().Qualifis([]string{"1", "3", "4", "5", "6"})
	//fmt.Println(advQualifis)
}
