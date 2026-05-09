package main

import (
	"fmt"

	"github.com/Glootea/gettextgocodegen/examples/translations"
)

func main() {
	en, _ := translations.New(translations.LocaleEnUS)
	ru, _ := translations.New(translations.LocaleRuRU)

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

	fmt.Println("\n=== Conflict Test ===")
	fmt.Println("HelloWorld:", en.HelloWorld())
	fmt.Println("HelloWorld:", ru.HelloWorld())
}
