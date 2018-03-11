package consistency

import (
	"common"
	"fmt"

	"bytes"
	"encoding/binary"
	"encoding/json"
	"math"
)

var ItemIDMap map[string]common.Item

type Cart struct {
	content map[string]uint32 //[itemid]count
}

var cart Cart
var Comuport int // port of this server

func ComuportInit(serverport int) {
	Comuport = serverport
}

func UnMarshalCart(d []byte) uint32 {
       bs := bytes.NewBuffer(d)

       var rednum uint32
       binary.Read(bytes.NewBuffer(bs.Next(4)), binary.LittleEndian, &rednum)

/*
       bs = bytes.NewBuffer(bs.Next(math.MaxInt32))
       e := gob.NewDecoder(bs)
       var tmpmap map[string]common.Item
       err := e.Decode(&tmpmap)
       if err != nil {
                fmt.Println(`failed gob Decode`, err);
       } else {
                fmt.Println(`gob Decode success!`);
       }

       // Ta da! It is a map!
       fmt.Printf("%#v\n", tmpmap)
       ItemIDMap = tmpmap
*/
       bs = bytes.NewBuffer(bs.Next(math.MaxInt32))
       json.Unmarshal(bs.Bytes(), &ItemIDMap)
       fmt.Printf("%#v\n", ItemIDMap)

       return rednum
}

func ReplyStatusIsNew() {
       mes := NewMessage(MESSAGE_STATUS_IS_NEW)
       res := &bytes.Buffer{}

       mes.Data = res.Bytes()
       //fmt.Println("Ready to send MESSAGE_STATUS_IS_NEW")
       //Core.Network.StartupMessageQueue <- *mes
       Core.Network.StartupReplyQueue <- *mes
       //Core.Network.BroadcastQueue <- *mes
}

func ReplyCurStatus() {
        mes := NewMessage(MESSAGE_START_UPDATE_REPLY)
        res := &bytes.Buffer{}

        bs := &bytes.Buffer{}
        binary.Write(bs, binary.LittleEndian, RedNum)
        res.Write(bs.Bytes())

/*
        bs = &bytes.Buffer{}
        e := gob.NewEncoder(bs)
        err := e.Encode(&ItemIDMap)
        if err != nil {
                fmt.Println(`failed gob Encode`, err);
        } else {
                fmt.Println(`gob Encode success!`);
        }
        res.Write(bs.Bytes())
*/
        //emp := make(map[string]interface{})
        //emp = ItemIDMap
        empData, err := json.Marshal(ItemIDMap)
        if err != nil {
                fmt.Println(err.Error())
                return
        }

        jsonStr := string(empData)
        //fmt.Println("Current map is", jsonStr)
        res.Write([]byte(jsonStr))

        mes.Data = res.Bytes()
        //fmt.Println("Ready to send MESSAGE_START_UPDATE_REPLY")
        //Core.Network.BroadcastQueue <- *mes
        Core.Network.StartupReplyQueue <- *mes
        //Core.Network.StartupMessageQueue <- *mes
}

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

func AddItemToCartForClient(itemid string, num uint32) OP_RESULT {
	if _, exist := ItemIDMap[itemid]; exist {
		if count, ok := cart.content[itemid]; ok {
			cart.content[itemid] = count + num
		} else {
			cart.content[itemid] = num
		}
		//fmt.Printf("add %d <%s> successfully, now %d in cart\n", num, item.Name, cart.content[itemid])
		return OPERATION_SUCCESS
	} else {
		//fmt.Printf("no such item\n")
		return OPERATION_FAIL
	}
}

func RemoveItemFromCartForClient(itemid string, num uint32) OP_RESULT {
	if _, exist := ItemIDMap[itemid]; exist {
		if count, ok := cart.content[itemid]; ok {
			if count >= num {
				cart.content[itemid] = count - num
				// fmt.Printf("remove %d <%s> successfully, now %d in cart\n", num, item.Name, cart.content[itemid])
				if count == num {
					delete(cart.content, itemid)
				}
				return OPERATION_SUCCESS

			} else {
				// fmt.Printf("You have %d <%s> but want remove %d. Nothing happens.\n", cart.content[itemid], item.Name, num)
			}
		} else {
			// fmt.Printf("No such item in your cart.\n")
			return OPERATION_SUCCESS
		}
	} else {
		// fmt.Printf("no such item\n")
	}
	return OPERATION_FAIL
}

func ClearCartForServer() OP_RESULT {
	//TODO: Recieve
	ClearContent := make(map[string]uint32)
	cart.content = ClearContent
	// fmt.Printf("Now you have nothing in your cart\n")
	return OPERATION_SUCCESS
	//TODO: send confirm
}

func CheckItemVolume() OP_RESULT {
	for itemid, count := range cart.content {
		if ItemIDMap[itemid].Volume < count {
			return false
		}
	}
	return true
}

func ArchiveCartItems() []byte {
	bs := []byte{}
	for itemid, _ := range cart.content {
		tempitem := ItemIDMap[itemid]
		item := common.Item{tempitem.Name, cart.content[itemid], tempitem.ID, tempitem.Price}
		bytes, _ := item.MarshalBinary()
		bs = append(bs, bytes...)
	}
	return bs
}

func CheckoutForServer(op Operation) OP_RESULT {
	itemCount := int(op.PayloadLength) / 118
	bs := op.Payload
	items := []common.Item{}
	for i := 1; i <= itemCount; i++ {
		item := common.Item{}
		item.UnMarshalBinary(bs[118*(i-1) : 118*i])
		items = append(items, item)
	}

	for i := 0; i < len(items); i++ {
		tempitem := ItemIDMap[items[i].ID]
		tempitem.Volume -= items[i].Volume
		ItemIDMap[items[i].ID] = tempitem
	}

	ClearContent := make(map[string]uint32)
	cart.content = ClearContent

	return OPERATION_SUCCESS

}

func AddNewItem(item common.Item) OP_RESULT {
	ItemIDMap[item.ID] = item
	return OPERATION_SUCCESS
}
