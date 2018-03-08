package consistency

import (
	"common"
	"fmt"
)

var ItemIDMap map[string]common.Item

type Cart struct {
	content map[string]uint32 //[itemid]count
}

var cart Cart

func initData() {

	CartContent := make(map[string]uint32)
	cart = Cart{content: CartContent}
	ItemIDMap = make(map[string]common.Item)
	ItemIDMap["3kd7a8d9lf"] = common.Item{"python", 100, "3kd7a8d9lf", 50}
	ItemIDMap["kfy3ksd8ks"] = common.Item{"golang", 30, "kfy3ksd8ks", 45}
	// here to add some initial item to our shop(the ItemIDMap)

}

// give ItemIDMap to client
func GetClientItemIDMap() []common.Item {
	items := []common.Item{}
	for _, v := range ItemIDMap {
		items = append(items, v)
	}
	return items
}

func GetItemIDMapFromCart() []common.CartItem {
	items := []common.CartItem{}
	for k, v := range cart.content {
		items = append(items, common.CartItem{ItemIDMap[k].Name, v, ItemIDMap[k].ID, ItemIDMap[k].Price})
	}
	return items
}

func AddItemToCartForClient(itemid string, num uint32) {
	//TODO: Recieve
	if item, exist := ItemIDMap[itemid]; exist {
		if count, ok := cart.content[itemid]; ok {
			cart.content[itemid] = count + num
		} else {
			cart.content[itemid] = num
		}
		fmt.Printf("add %d <%s> successfully, now you have %d\n", num, item.Name, cart.content[itemid])
	} else {
		fmt.Printf("no such item\n")
	}
	//TODO: send confirm
}

func RemoveItemFromCartForClient(itemid string, num uint32) {
	//TODO: Recieve
	if item, exist := ItemIDMap[itemid]; exist {
		if count, ok := cart.content[itemid]; ok {
			if count >= num {
				cart.content[itemid] = count - num
				fmt.Printf("remove %d <%s> successfully, now you have %d\n", num, item.Name, cart.content[itemid])
				if count == num {
					delete(cart.content, itemid)
				}
			} else {
				fmt.Printf("You have %d <%s> but want remove %d. Nothing happens.\n", cart.content[itemid], item.Name, num)
			}
		} else {
			fmt.Printf("You do not have such item in your cart.\n")
		}
	} else {
		fmt.Printf("no such item\n")
	}
	//TODO: send confirm
}

func ClearCartForServer() error {
	//TODO: Recieve
	ClearContent := make(map[string]uint32)
	cart.content = ClearContent
	fmt.Printf("Now you have nothing in your cart\n")
	return nil
	//TODO: send confirm
}

func checkoutForServer() {
	//TODO: recieve a req

	success := 1

	for itemid, count := range cart.content {
		if ItemIDMap[itemid].Volume < count {
			success = 0
			break
		}
	}

	if success == 1 { // I think only this is a red operation
		for itemid, count := range cart.content {
			tempitem := ItemIDMap[itemid]
			tempitem.Volume -= count
			ItemIDMap[itemid] = tempitem
		}
		ClearContent := make(map[string]uint32)
		cart.content = ClearContent
		//TODO:send OK
	} else {
		//TODO:send not enough
	}

}

func AddNewItem(item common.Item) {
	ItemIDMap[item.ID] = item
}
