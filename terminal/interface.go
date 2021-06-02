package terminal


type Broker interface {
	NewOrder()
	GetOrders()
}

type Trader interface {
	getPrice()
}