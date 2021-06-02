package main

import (
	"fmt"

	"gtrade/entity"
	"gtrade/portfolio"
	"gtrade/terminal"
)

func main() {
	var order entity.Order

	t := terminal.Get("USD")

	order = t.NewOrder("ZYNE", 4.23, 10, "buy")
	t.ProcessOrder(order)
	//
	order = t.NewOrder("RIG", 3.42, 10, "buy")
	t.ProcessOrder(order)

	order = t.NewOrder("FOLD", 9.92, 10, "buy")
	t.ProcessOrder(order)

	order = t.NewOrder("AAPL", 21.9, 10, "sell")
	t.ProcessOrder(order)

	p := portfolio.Get("USD")

	fmt.Println(p)
}