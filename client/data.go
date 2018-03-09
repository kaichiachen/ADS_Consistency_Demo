// This file is NOT USED

package main

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

}

// get ItemIDMap from server
func RefreshItemIDMap() {
	//TODO
}

func AddItemToCartForClient(itemid string, num uint32) {

	if item, exist := ItemIDMap[itemid]; exist {
		//TODO:send a req to server
		//TODO:get a confirm from server
		if count, ok := cart.content[itemid]; ok {
			cart.content[itemid] = count + num
		} else {
			cart.content[itemid] = num
		}
		fmt.Printf("add %d <%s> successfully, now you have %d\n", num, item.Name, cart.content[itemid])
	} else {
		fmt.Printf("no such item\n")
	}

}

func RemoveItemFromCartForClient(itemid string, num uint32) {
	if item, exist := ItemIDMap[itemid]; exist {
		if count, ok := cart.content[itemid]; ok {
			if count >= num {
				//TODO:send a req to server
				//TODO:get a confirm from server
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
}

func ClearCart() error {
	//TODO:send a req to server
	//TODO:get a confirm from server
	ClearContent := make(map[string]uint32)
	cart.content = ClearContent
	fmt.Printf("Now you have nothing in your cart\n")
	return nil
}

func checkout() {
	//send a req to server

	//if return success
	RefreshItemIDMap()
	var total uint32 = 0
	for itemid, count := range cart.content {
		total = total + (count)*ItemIDMap[itemid].Price
	}

	fmt.Printf("You should pay $%d for your items\n", total)
	err := ClearCart()

	if err != nil {
		fmt.Printf("Something wrong, checkout failed\n")
	}

}

func TestForClientData() {
	ItemIDMap["a"] = common.Item{"XYJ", 0, "a", 13}

	ListItemToUser()

	AddItemToCartForClient("a", 4)
	RemoveItemFromCartForClient("a", 5)
	RemoveItemFromCartForClient("a", 2)
	AddItemToCartForClient("b", 2)

	ListCartToUser()

	checkout()
}

func AddNewItem(item common.Item) {
	//TODO: send req to server
	RefreshItemIDMap()
}

func ListItemToUser() {
	RefreshItemIDMap()
	fmt.Printf("Here are some goods:\n")
	for _, item := range ItemIDMap {
		fmt.Printf("<%s> : %d for sale, $%d each\n", item.Name, item.Volume, item.Price)
	}
}

func ListCartToUser() {
	fmt.Printf("Here are the goods in your cart:\n")
	for itemid, count := range cart.content {
		fmt.Printf("<%s> : %d in cart\n", ItemIDMap[itemid].Name, count)
	}
}
