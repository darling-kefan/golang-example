package main

import (
	"fmt"
	"models"
	"strconv"
	_ "reflect"
)

func main() {
	//publishs := models.Publishs().GetJoinPubsByIds([]string{"1000", "22969"})
	//fmt.Println(publishs)

	//publishs = models.Publishs().GetPubsByIds([]string{"1000", "22969", "22974"})
	//fmt.Println(publishs)

	//fmt.Println("------------------------------------------------------")

	//cards := models.Cards().GetCardsByIds([]string{"160", "161", "177"})
	//fmt.Println(cards)

	//fmt.Println("------------------------------------------------------")

	//advertisers := models.Advertisers().GetAdversByIds([]string{"1", "3", "4", "5", "6"}, []string{"id", "tvmid", "name", "mobileno", "email", "type", "subtype"})
	//fmt.Println(advertisers)

	//advQualifis := models.Advertisers().GetQualifis([]string{"1", "3", "4", "5", "6"})
	//fmt.Println(advQualifis)

	fmt.Println("------------------------------------------------------")
	publishs := models.Publishs().LastestPubversion(strconv.Itoa(22983), 2)
	fmt.Println(publishs)
}
