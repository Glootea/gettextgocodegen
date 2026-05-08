package main

import (
	"fmt"

	"github.com/gettextcodegen/examples/translations"
)

func main() {
	en, _ := translations.New("en_US")
	ru, _ := translations.New("ru_RU")

	fmt.Println("=== English ===")
	fmt.Println(en.HelloWorld())
	fmt.Println(en.WelcomeToOurApp())
	fmt.Println(en.HiMyNameIs("John"))
	fmt.Println(en.ButtonsSubmit())
	fmt.Println(en.GetItemsCount(1))
	fmt.Println(en.GetItemsCount(5))
	fmt.Println(en.CartTotal("$100"))

	fmt.Println("\n=== Russian ===")
	fmt.Println(ru.HelloWorld())
	fmt.Println(ru.WelcomeToOurApp())
	fmt.Println(ru.HiMyNameIs("Джон"))
	fmt.Println(ru.ButtonsSubmit())
	fmt.Println(ru.GetItemsCount(1))
	fmt.Println(ru.GetItemsCount(5))
	fmt.Println(ru.CartTotal("100₽"))
}
