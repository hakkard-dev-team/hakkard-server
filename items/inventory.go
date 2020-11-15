package items

type Inventory struct {
	items []Item `json:"items"`
}

func (I *Inventory) NewInventory() *Inventory {

	return &Inventory
}

func (I *Inventory) Add(i Item) {
	append(i, I.items)
}
