package channelserver

import (
	"encoding/hex"
	ps "erupe-ce/common/pascalstring"
	"fmt"
	"strings"
	"time"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"

	"github.com/lib/pq"
	"github.com/sachaos/lottery"
	"go.uber.org/zap"
)

type ShopItem struct {
	ID            uint32 `db:"id"`
	ItemID        uint16 `db:"itemid"`
	Cost          uint32 `db:"cost"`
	Quantity      uint16 `db:"quantity"`
	MinHR         uint16 `db:"min_hr"`
	MinSR         uint16 `db:"min_sr"`
	MinGR         uint16 `db:"min_gr"`
	ReqStoreLevel uint16 `db:"req_store_level"`
	MaxQuantity   uint16 `db:"max_quantity"`
	CharQuantity  uint16 `db:"char_quantity"`
	RoadFloors    uint16 `db:"road_floors"`
	RoadFatalis   uint16 `db:"road_fatalis"`
	EnableWeeks   string `db:"enable_weeks"`
}

type Gacha struct {
	ID    uint32 `db:"id"`
	MinGR uint32 `db:"min_gr"`
	MinHR uint32 `db:"min_hr"`
	Name  string `db:"name"`
	Link1 string `db:"link1"`
	Link2 string `db:"link2"`
	Link3 string `db:"link3"`
	Icon  uint16 `db:"icon"`
	Type  uint16 `db:"type"`
	Hide  bool   `db:"hide"`
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func handleMsgMhfEnumerateShop(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateShop)
	// Generic Shop IDs
	// 0: basic item
	// 1: gatherables
	// 2: hr1-4 materials
	// 3: hr5-7 materials
	// 4: decos
	// 5: other item
	// 6: g mats
	// 7: limited item
	// 8: special item
	switch pkt.ShopType {
	case 1: // Running gachas
		var count uint16
		shopEntries, err := s.Server.db.Queryx("SELECT id, min_gr, min_hr, name, link1, link2, link3, icon, type, hide FROM gacha_shop ORDER BY id")
		if err != nil {
			doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
			return
		}
		resp := byteframe.NewByteFrame()
		resp.WriteUint32(0)
		var gacha Gacha
		for shopEntries.Next() {
			err = shopEntries.StructScan(&gacha)
			if err != nil {
				continue
			}
			resp.WriteUint32(gacha.ID)
			resp.WriteBytes(make([]byte, 16)) // Rank restriction
			resp.WriteUint32(gacha.MinGR)
			resp.WriteUint32(gacha.MinHR)
			resp.WriteUint32(0) // only 0 in known packet
			ps.Uint8(resp, gacha.Name, true)
			ps.Uint8(resp, gacha.Link1, false)
			ps.Uint8(resp, gacha.Link2, false)
			resp.WriteBool(gacha.Hide)
			ps.Uint8(resp, gacha.Link3, false)
			resp.WriteUint16(gacha.Icon)
			resp.WriteUint16(gacha.Type)
			count++
		}
		// nolint:errcheck
		resp.Seek(0, 0)
		resp.WriteUint16(count)
		resp.WriteUint16(count)
		doAckBufSucceed(s, pkt.AckHandle, resp.Data())
	case 2: // Actual gacha
		shopEntries, err := s.Server.db.Query("SELECT entryType, itemhash, currType, currNumber, currQuant, percentage, rarityIcon, rollsCount, itemCount, dailyLimit, itemType, itemId, quantity FROM gacha_shop_items WHERE shophash=$1 ORDER BY itemhash", pkt.ShopID)
		if err != nil {
			doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
			return
		}
		var entryType, currType, rarityIcon, rollsCount, itemCount, dailyLimit uint8
		var currQuant, currNumber, percentage uint16
		var itemhash uint32
		var itemType, itemId, quantity pq.Int64Array
		var count uint16
		resp := byteframe.NewByteFrame()
		resp.WriteUint32(pkt.ShopID)
		resp.WriteUint16(0) // total defs
		for shopEntries.Next() {
			err = shopEntries.Scan(&entryType, &itemhash, &currType, &currNumber, &currQuant, &percentage, &rarityIcon, &rollsCount, &itemCount, &dailyLimit, (*pq.Int64Array)(&itemType), (*pq.Int64Array)(&itemId), (*pq.Int64Array)(&quantity))
			if err != nil {
				panic(err)
			}
			resp.WriteUint8(entryType)
			resp.WriteUint32(itemhash)
			resp.WriteUint8(currType)
			resp.WriteUint16(0)          // unk, always 0 in existing packets
			resp.WriteUint16(currNumber) // it's either item ID or quantity for gacha coins
			resp.WriteUint16(currQuant)  // only for item ID
			resp.WriteUint16(percentage)
			resp.WriteUint8(rarityIcon)
			resp.WriteUint8(rollsCount)
			resp.WriteUint8(itemCount)
			resp.WriteUint16(0) // unk, always 0 in existing packets
			resp.WriteUint8(dailyLimit)
			resp.WriteUint8(0) // unk, always 0 in existing packets
			for i := 0; i < int(itemCount); i++ {
				resp.WriteUint16(uint16(itemType[i])) // unk, always 0 in existing packets
				resp.WriteUint16(uint16(itemId[i]))   // unk, always 0 in existing packets
				resp.WriteUint16(uint16(quantity[i])) // unk, always 0 in existing packets
			}
			count++
		}

		// nolint:errcheck // Error return value of `resp.Seek` is not checked
		resp.Seek(4, 0)
		resp.WriteUint16(count)
		doAckBufSucceed(s, pkt.AckHandle, resp.Data())
	case 4: // N Points, 0-6
		doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
	case 5: // GCP->Item, 0-6
		doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
	case 6: // Gacha coin->Item
		doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
	case 7: // Item->GCP
		data, _ := hex.DecodeString("000300033a9186fb000033860000000a000100000000000000000000000000000000097fdb1c0000067e0000000a0001000000000000000000000000000000001374db29000027c300000064000100000000000000000000000000000000")
		doAckBufSucceed(s, pkt.AckHandle, data)
	case 8: // Diva
		switch pkt.ShopID {
		case 0: // Normal exchange
			doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
		case 5: // GCP skills
			data, _ := hex.DecodeString("001f001f2c9365c1000000010000001e000a0000000000000000000a0000000000001979f1c2000000020000003c000a0000000000000000000a0000000000003e5197df000000030000003c000a0000000000000000000a000000000000219337c0000000040000001e000a0000000000000000000a00000000000009b24c9d000000140000001e000a0000000000000000000a0000000000001f1d496e000000150000001e000a0000000000000000000a0000000000003b918fcb000000160000003c000a0000000000000000000a0000000000000b7fd81c000000170000003c000a0000000000000000000a0000000000001374f239000000180000003c000a0000000000000000000a00000000000026950cba0000001c0000003c000a0000000000000000000a0000000000003797eae70000001d0000003c000a012b000000000000000a00000000000015758ad8000000050000003c00000000000000000000000a0000000000003c7035050000000600000050000a0000000000000001000a00000000000024f3b5560000000700000050000a0000000000000001000a00000000000000b600330000000800000050000a0000000000000001000a0000000000002efdce840000001900000050000a0000000000000001000a0000000000002d9365f10000001a00000050000a0000000000000001000a0000000000001979f3420000001f00000050000a012b000000000001000a0000000000003f5397cf0000002000000050000a012b000000000001000a000000000000319337c00000002100000050000a012b000000000001000a00000000000008b04cbd0000000900000064000a0000000000000002000a0000000000000b1d4b6e0000000a00000064000a0000000000000002000a0000000000003b918feb0000000b00000064000a0000000000000002000a0000000000001b7fd81c0000000c00000064000a0000000000000002000a0000000000001276f2290000000d00000064000a0000000000000002000a00000000000022950cba0000000e000000c8000a0000000000000002000a0000000000003697ead70000000f000001f4000a0000000000000003000a00000000000005758a5800000010000003e8000a0000000000000003000a0000000000003c7035250000001b000001f4000a0000000000010003000a00000000000034f3b5d60000001e00000064000a012b000000000003000a00000000000000b600030000002200000064000a0000000000010003000a000000000000")
			doAckBufSucceed(s, pkt.AckHandle, data)
		case 7: // Note exchange
			doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
		}
	case 10: // Item shop, 0-8
		_, week := time.Now().ISOWeek()
		season := fmt.Sprintf("%d", week%4)
		shopEntries, err := s.Server.db.Queryx(`SELECT id, itemid, cost, quantity, min_hr, min_sr, min_gr, req_store_level, max_quantity,
       		COALESCE((SELECT usedquantity FROM shop_item_state WHERE itemhash=nsi.id AND char_id=$3), 0) as char_quantity,
       		road_floors, road_fatalis, COALESCE(enable_weeks, '') as enable_weeks FROM normal_shop_items nsi WHERE shoptype=$1 AND shopid=$2
       		`, pkt.ShopType, pkt.ShopID, s.CharID)
		if err != nil {
			doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
			return
		}
		var count uint16
		resp := byteframe.NewByteFrame()
		resp.WriteBytes(make([]byte, 4))
		var shopItem ShopItem
		for shopEntries.Next() {
			err = shopEntries.StructScan(&shopItem)
			if err != nil {
				continue
			}

			if len(shopItem.EnableWeeks) > 0 && !contains(strings.Split(shopItem.EnableWeeks, ","), season) {
				continue
			}

			resp.WriteUint32(shopItem.ID)
			resp.WriteUint16(0)
			resp.WriteUint16(shopItem.ItemID)
			resp.WriteUint32(shopItem.Cost)
			resp.WriteUint16(shopItem.Quantity)
			resp.WriteUint16(shopItem.MinHR)
			resp.WriteUint16(shopItem.MinSR)
			resp.WriteUint16(shopItem.MinGR)
			resp.WriteUint16(shopItem.ReqStoreLevel)
			resp.WriteUint16(shopItem.MaxQuantity)

			if shopItem.CharQuantity > 0 {
				var itemWeek int
				err = s.Server.db.QueryRow("SELECT COALESCE(usedquantity,0), COALESCE(week,-1) FROM shop_item_state WHERE itemhash=$1 AND char_id=$2", shopItem.ID, s.CharID).Scan(&shopItem.CharQuantity, &itemWeek)
				if err != nil {
					resp.WriteUint16(0)
				} else if pkt.ShopID == 7 && itemWeek >= 0 && itemWeek != week {
					clearShopItemState(s, s.CharID, uint32(shopItem.ID))
					resp.WriteUint16(0)
				} else {
					resp.WriteUint16(shopItem.CharQuantity)
				}
			} else {
				resp.WriteUint16(shopItem.CharQuantity)
			}

			resp.WriteUint16(shopItem.RoadFloors)
			resp.WriteUint16(shopItem.RoadFatalis)
			count++
		}
		// nolint:errcheck
		resp.Seek(0, 0)
		resp.WriteUint16(count)
		resp.WriteUint16(count)
		doAckBufSucceed(s, pkt.AckHandle, resp.Data())
	}
}

func handleMsgMhfAcquireExchangeShop(s *Session, p mhfpacket.MHFPacket) {
	_, week := time.Now().ISOWeek()
	// writing out to an editable shop enumeration
	pkt := p.(*mhfpacket.MsgMhfAcquireExchangeShop)
	if pkt.DataSize == 10 {
		bf := byteframe.NewByteFrameFromBytes(pkt.RawDataPayload)
		_ = bf.ReadUint16() // unk, always 1 in examples
		itemHash := bf.ReadUint32()
		buyCount := bf.ReadUint32()
		_, err := s.Server.db.Exec(`INSERT INTO shop_item_state (char_id, itemhash, usedquantity, week)
  														 VALUES ($1,$2,$3,$4) ON CONFLICT (char_id, itemhash)
  														 DO UPDATE SET usedquantity = shop_item_state.usedquantity + $3
  														 WHERE EXCLUDED.char_id=$1 AND EXCLUDED.itemhash=$2`, s.CharID, itemHash, buyCount, week)
		if err != nil {
			s.logger.Fatal("Failed to update shop_item_state in db", zap.Error(err))
		}
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func clearShopItemState(s *Session, charId uint32, itemHash uint32) {
	_, err := s.Server.db.Exec(`DELETE FROM shop_item_state WHERE char_id=$1 AND itemhash=$2`, charId, itemHash)
	if err != nil {
		s.logger.Fatal("Failed to delete shop_item_state in db", zap.Error(err))
	}
}

func handleMsgMhfGetGachaPlayHistory(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetGachaPlayHistory)
	bf := byteframe.NewByteFrame()
	bf.WriteUint8(0)
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGetGachaPoint(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetGachaPoint)
	var fp, gp, gt uint32
	_ = s.Server.db.QueryRow("SELECT COALESCE(frontier_points, 0), COALESCE(gacha_prem, 0), COALESCE(gacha_trial,0) FROM characters WHERE id=$1", s.CharID).Scan(&fp, &gp, &gt)
	resp := byteframe.NewByteFrame()
	resp.WriteUint32(gp) // Real Gacha Points?
	resp.WriteUint32(gt) // Trial Gacha Point?
	resp.WriteUint32(fp) // Frontier Points?
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

type gachaItem struct {
	itemhash   uint32
	percentage uint16
	rarityIcon byte
	itemCount  byte
	itemType   pq.Int64Array
	itemId     pq.Int64Array
	quantity   pq.Int64Array
}

func (i gachaItem) Weight() int {
	return int(i.percentage)
}

func handleMsgMhfPlayNormalGacha(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfPlayNormalGacha)
	// needs to query db for input gacha and return a result or number of results
	// uint8 number of results
	// uint8 item type
	// uint16 item id
	// uint16 quantity

	var currType, rarityIcon, rollsCount, itemCount byte
	var currQuant, currNumber, percentage uint16
	var itemhash uint32
	var itemType, itemId, quantity pq.Int64Array
	var items []lottery.Weighter
	// get info for updating data and calculating costs
	err := s.Server.db.QueryRow("SELECT currType, currNumber, currQuant, rollsCount FROM gacha_shop_items WHERE shophash=$1 AND entryType=$2", pkt.GachaHash, pkt.RollType).Scan(&currType, &currNumber, &currQuant, &rollsCount)
	if err != nil {
		panic(err)
	}
	// get existing items in storage if any
	var data []byte
	_ = s.Server.db.QueryRow("SELECT gacha_items FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if len(data) == 0 {
		data = []byte{0x00}
	}
	// get gacha items and iterate through them for gacha roll
	shopEntries, err := s.Server.db.Query("SELECT itemhash, percentage, rarityIcon, itemCount, itemType, itemId, quantity FROM gacha_shop_items WHERE shophash=$1 AND entryType=100 ORDER BY itemhash", pkt.GachaHash)
	if err != nil {
		panic(err)
	}
	for shopEntries.Next() {
		err = shopEntries.Scan(&itemhash, &percentage, &rarityIcon, &itemCount, (*pq.Int64Array)(&itemType), (*pq.Int64Array)(&itemId), (*pq.Int64Array)(&quantity))
		if err != nil {
			panic(err)
		}
		items = append(items, &gachaItem{itemhash: itemhash, percentage: percentage, rarityIcon: rarityIcon, itemCount: itemCount, itemType: itemType, itemId: itemId, quantity: quantity})
	}
	// execute rolls, build response and update database
	results := byte(0)
	resp := byteframe.NewByteFrame()
	dbUpdate := byteframe.NewByteFrame()
	resp.WriteUint8(0) // results go here later
	l := lottery.NewDefaultLottery()
	for x := 0; x < int(rollsCount); x++ {
		ind := l.Draw(items)
		results += items[ind].(*gachaItem).itemCount
		for y := 0; y < int(items[ind].(*gachaItem).itemCount); y++ {
			// items in storage don't get rarity
			dbUpdate.WriteUint8(uint8(items[ind].(*gachaItem).itemType[y]))
			dbUpdate.WriteUint16(uint16(items[ind].(*gachaItem).itemId[y]))
			dbUpdate.WriteUint16(uint16(items[ind].(*gachaItem).quantity[y]))
			data = append(data, dbUpdate.Data()...)
			// nolint:errcheck // Error return value of `rollFrame.Seek` is not checked
			dbUpdate.Seek(0, 0)
			// response needs all item info and the rarity
			resp.WriteBytes(dbUpdate.Data())
			resp.WriteUint8(items[ind].(*gachaItem).rarityIcon)
		}
	}
	// nolint:errcheck
	resp.Seek(0, 0)
	resp.WriteUint8(results)
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())

	// add claimables to DB
	data[0] = data[0] + results
	_, err = s.Server.db.Exec("UPDATE characters SET gacha_items = $1 WHERE id = $2", data, s.CharID)
	if err != nil {
		s.logger.Fatal("Failed to update minidata in db", zap.Error(err))
	}
	// deduct gacha coins if relevant, items are handled fine by the standard savedata packet immediately afterwards
	if currType == 19 {
		_, err = s.Server.db.Exec("UPDATE characters SET gacha_trial = CASE WHEN (gacha_trial > $1) then gacha_trial - $1 else gacha_trial end, gacha_prem = CASE WHEN NOT (gacha_trial > $1) then gacha_prem - $1 else gacha_prem end WHERE id=$2", currNumber, s.CharID)
	}
	if err != nil {
		s.logger.Fatal("Failed to update gacha_items in db", zap.Error(err))
	}
}

func handleMsgMhfUseGachaPoint(s *Session, p mhfpacket.MHFPacket) {
	// should write to database when that's set up
	pkt := p.(*mhfpacket.MsgMhfUseGachaPoint)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfExchangeFpoint2Item(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfExchangeFpoint2Item)
	var balance uint32
	var itemValue, quantity int
	_ = s.Server.db.QueryRow("SELECT quantity, fpoints FROM fpoint_items WHERE id=$1", pkt.TradeID).Scan(&quantity, &itemValue)
	cost := (int(pkt.Quantity) * quantity) * itemValue
	// nolint:errcheck
	s.Server.db.QueryRow("UPDATE characters SET frontier_points=frontier_points::int - $1 WHERE id=$2 RETURNING frontier_points", cost, s.CharID).Scan(&balance)
	bf := byteframe.NewByteFrame()
	bf.WriteUint32(balance)
	doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfExchangeItem2Fpoint(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfExchangeItem2Fpoint)
	var balance uint32
	var itemValue, quantity int
	// nolint:errcheck
	s.Server.db.QueryRow("SELECT quantity, fpoints FROM fpoint_items WHERE id=$1", pkt.TradeID).Scan(&quantity, &itemValue)
	cost := (int(pkt.Quantity) / quantity) * itemValue
	// nolint:errcheck
	s.Server.db.QueryRow("UPDATE characters SET frontier_points=COALESCE(frontier_points::int + $1, $1) WHERE id=$2 RETURNING frontier_points", cost, s.CharID).Scan(&balance)
	bf := byteframe.NewByteFrame()
	bf.WriteUint32(balance)
	doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGetFpointExchangeList(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetFpointExchangeList)
	resp := byteframe.NewByteFrame()
	resp.WriteUint32(0)
	var buyables, sellables uint16
	var id uint32
	var itemType uint8
	var itemID, quantity, fPoints uint16

	buyRows, err := s.Server.db.Query("SELECT id,item_type,item_id,quantity,fpoints FROM fpoint_items WHERE trade_type=0")
	if err == nil {
		for buyRows.Next() {
			err = buyRows.Scan(&id, &itemType, &itemID, &quantity, &fPoints)
			if err != nil {
				continue
			}
			resp.WriteUint32(id)
			resp.WriteUint16(0)
			resp.WriteUint16(0)
			resp.WriteUint16(0)
			resp.WriteUint8(itemType)
			resp.WriteUint16(itemID)
			resp.WriteUint16(quantity)
			resp.WriteUint16(fPoints)
			buyables++
		}
	}

	sellRows, err := s.Server.db.Query("SELECT id,item_type,item_id,quantity,fpoints FROM fpoint_items WHERE trade_type=1")
	if err == nil {
		for sellRows.Next() {
			err = sellRows.Scan(&id, &itemType, &itemID, &quantity, &fPoints)
			if err != nil {
				continue
			}
			resp.WriteUint32(id)
			resp.WriteUint16(0)
			resp.WriteUint16(0)
			resp.WriteUint16(0)
			resp.WriteUint8(itemType)
			resp.WriteUint16(itemID)
			resp.WriteUint16(quantity)
			resp.WriteUint16(fPoints)
			sellables++
		}
	}
	// nolint:errcheck
	resp.Seek(0, 0)
	resp.WriteUint16(buyables)
	resp.WriteUint16(sellables)

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfPlayStepupGacha(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfPlayStepupGacha)
	results := byte(0)
	stepResults := byte(0)
	resp := byteframe.NewByteFrame()
	rollFrame := byteframe.NewByteFrame()
	stepFrame := byteframe.NewByteFrame()
	stepData := []byte{}
	var currType, rarityIcon, rollsCount, itemCount byte
	var currQuant, currNumber, percentage uint16
	var itemhash uint32
	var itemType, itemId, quantity pq.Int64Array
	var items []lottery.Weighter
	// get info for updating data and calculating costs
	err := s.Server.db.QueryRow("SELECT currType, currNumber, currQuant, rollsCount, itemCount, itemType, itemId, quantity FROM gacha_shop_items WHERE shophash=$1 AND entryType=$2", pkt.GachaHash, pkt.RollType).Scan(&currType, &currNumber, &currQuant, &rollsCount, &itemCount, (*pq.Int64Array)(&itemType), (*pq.Int64Array)(&itemId), (*pq.Int64Array)(&quantity))
	if err != nil {
		panic(err)
	}
	// get existing items in storage if any
	var data []byte
	_ = s.Server.db.QueryRow("SELECT gacha_items FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if len(data) == 0 {
		data = []byte{0x00}
	}
	// roll definition includes items with step up gachas that are appended last
	for x := 0; x < int(itemCount); x++ {
		stepFrame.WriteUint8(uint8(itemType[x]))
		stepFrame.WriteUint16(uint16(itemId[x]))
		stepFrame.WriteUint16(uint16(quantity[x]))
		stepData = append(stepData, stepFrame.Data()...)
		stepFrame.WriteUint8(0) // rarity still defined
		stepResults++
	}
	// get gacha items and iterate through them for gacha roll
	shopEntries, err := s.Server.db.Query("SELECT itemhash, percentage, rarityIcon, itemCount, itemType, itemId, quantity FROM gacha_shop_items WHERE shophash=$1 AND entryType=100 ORDER BY itemhash", pkt.GachaHash)
	if err != nil {
		panic(err)
	}
	for shopEntries.Next() {
		err = shopEntries.Scan(&itemhash, &percentage, &rarityIcon, &itemCount, (*pq.Int64Array)(&itemType), (*pq.Int64Array)(&itemId), (*pq.Int64Array)(&quantity))
		if err != nil {
			panic(err)
		}
		items = append(items, &gachaItem{itemhash: itemhash, percentage: percentage, rarityIcon: rarityIcon, itemCount: itemCount, itemType: itemType, itemId: itemId, quantity: quantity})
	}
	// execute rolls, build response and update database
	resp.WriteUint16(0) // results count goes here later
	l := lottery.NewDefaultLottery()
	for x := 0; x < int(rollsCount); x++ {
		ind := l.Draw(items)
		results += items[ind].(*gachaItem).itemCount
		for y := 0; y < int(items[ind].(*gachaItem).itemCount); y++ {
			// items in storage don't get rarity
			rollFrame.WriteUint8(uint8(items[ind].(*gachaItem).itemType[y]))
			rollFrame.WriteUint16(uint16(items[ind].(*gachaItem).itemId[y]))
			rollFrame.WriteUint16(uint16(items[ind].(*gachaItem).quantity[y]))
			data = append(data, rollFrame.Data()...)
			// nolint:errcheck // Error return value of `rollFrame.Seek` is not checked
			rollFrame.Seek(0, 0)
			// response needs all item info and the rarity
			resp.WriteBytes(rollFrame.Data())
			resp.WriteUint8(items[ind].(*gachaItem).rarityIcon)
		}
	}
	resp.WriteBytes(stepFrame.Data())
	// nolint:errcheck
	resp.Seek(0, 0)
	resp.WriteUint8(results + stepResults)
	resp.WriteUint8(results)
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())

	// add claimables to DB
	data = append(data, stepData...)
	data[0] = data[0] + results + stepResults
	_, err = s.Server.db.Exec("UPDATE characters SET gacha_items = $1 WHERE id = $2", data, s.CharID)
	if err != nil {
		s.logger.Fatal("Failed to update gacha_items in db", zap.Error(err))
	}
	// deduct gacha coins if relevant, items are handled fine by the standard savedata packet immediately afterwards
	// reduce real if trial don't cover cost
	if currType == 19 {
		_, err = s.Server.db.Exec(`UPDATE characters
			SET gacha_trial = CASE WHEN (gacha_trial > $1) then gacha_trial - $1 else gacha_trial end,
			gacha_prem = CASE WHEN NOT (gacha_trial > $1) then gacha_prem - $1 else gacha_prem end
			WHERE id=$2`, currNumber, s.CharID)
	}
	if err != nil {
		s.logger.Fatal("Failed to update gacha_items in db", zap.Error(err))
	}
	// update step progression
	_, err = s.Server.db.Exec("UPDATE stepup_state SET step_progression = $1 WHERE char_id = $2", pkt.RollType+1, s.CharID)
	if err != nil {
		s.logger.Fatal("Failed to update step_progression in db", zap.Error(err))
	}

}

func handleMsgMhfReceiveGachaItem(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfReceiveGachaItem)
	// persistent for claimable items on cat
	var data []byte
	err := s.Server.db.QueryRow("SELECT COALESCE(gacha_items, $2) FROM characters WHERE id = $1", s.CharID, []byte{0x00}).Scan(&data)
	if err != nil {
		panic("Failed to get gacha_items")
	}
	// limit of 36 items are returned
	if data[0] > 36 {
		outData := make([]byte, 181)
		copy(outData, data[0:181])
		outData[0] = byte(36)
		saveData := append(data[:1], data[181:]...)
		saveData[0] = saveData[0] - 36
		doAckBufSucceed(s, pkt.AckHandle, outData)
		if pkt.Unk0 != 0x2401 {
			_, err := s.Server.db.Exec("UPDATE characters SET gacha_items = $2 WHERE id = $1", s.CharID, saveData)
			if err != nil {
				s.logger.Fatal("Failed to update gacha_items in db", zap.Error(err))
			}
		}
	} else {
		doAckBufSucceed(s, pkt.AckHandle, data)
		if pkt.Unk0 != 0x2401 {
			_, err := s.Server.db.Exec("UPDATE characters SET gacha_items = null WHERE id = $1", s.CharID)
			if err != nil {
				s.logger.Fatal("Failed to update gacha_items in db", zap.Error(err))
			}
		}
	}
}

func handleMsgMhfGetStepupStatus(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetStepupStatus)
	// get the reset time from db
	var step_progression int
	var step_time time.Time
	err := s.Server.db.QueryRow(`SELECT COALESCE(step_progression, 0), COALESCE(step_time, $1) FROM stepup_state WHERE char_id = $2 AND shophash = $3`, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), s.CharID, pkt.GachaHash).Scan(&step_progression, &step_time)
	if err != nil {
		s.logger.Fatal("Failed to Select coalesce in db", zap.Error(err))
	}

	// calculate next midday
	var t = time.Now().In(time.FixedZone("UTC+9", 9*60*60))
	year, month, day := t.Date()
	midday := time.Date(year, month, day, 12, 0, 0, 0, t.Location())
	if t.After(midday) {
		midday = midday.Add(24 * time.Hour)
	}
	// after midday or not set
	if t.After(step_time) {
		step_progression = 0
	}
	_, err = s.Server.db.Exec(`INSERT INTO stepup_state (shophash, step_progression, step_time, char_id)
		VALUES ($1,$2,$3,$4) ON CONFLICT (shophash, char_id)
		DO UPDATE SET step_progression=$2, step_time=$3
		WHERE EXCLUDED.char_id=$4 AND EXCLUDED.shophash=$1`, pkt.GachaHash, step_progression, midday, s.CharID)
	if err != nil {
		s.logger.Fatal("Failed to update platedata savedata in db", zap.Error(err))
	}
	resp := byteframe.NewByteFrame()
	resp.WriteUint8(uint8(step_progression))
	resp.WriteUint32(uint32(time.Now().In(time.FixedZone("UTC+9", 9*60*60)).Unix()))
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfPlayFreeGacha(s *Session, p mhfpacket.MHFPacket) {
	// not sure this is used anywhere, free gachas use the MSG_MHF_PLAY_NORMAL_GACHA method in captures
}

func handleMsgMhfGetBoxGachaInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetBoxGachaInfo)
	count := 0
	var used_itemhash pq.Int64Array
	// pull array of used values
	// single sized respone with 0x00 is a valid with no items present
	_ = s.Server.db.QueryRow("SELECT used_itemhash FROM lucky_box_state WHERE shophash=$1 AND char_id=$2", pkt.GachaHash, s.CharID).Scan((*pq.Int64Array)(&used_itemhash))
	resp := byteframe.NewByteFrame()
	resp.WriteUint8(0)
	for ind := range used_itemhash {
		resp.WriteUint32(uint32(used_itemhash[ind]))
		resp.WriteUint8(1)
		count++
	}
	// nolint:errcheck
	resp.Seek(0, 0)
	resp.WriteUint8(uint8(count))
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfPlayBoxGacha(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfPlayBoxGacha)
	// needs to query db for input gacha and return a result or number of results
	// uint8 number of results
	// uint8 item type
	// uint16 item id
	// uint16 quantity

	var currType, rarityIcon, rollsCount, itemCount byte
	var currQuant, currNumber, percentage uint16
	var itemhash uint32
	var itemType, itemId, quantity, usedItemHash pq.Int64Array
	var items []lottery.Weighter
	// get info for updating data and calculating costs
	err := s.Server.db.QueryRow("SELECT currType, currNumber, currQuant, rollsCount FROM gacha_shop_items WHERE shophash=$1 AND entryType=$2", pkt.GachaHash, pkt.RollType).Scan(&currType, &currNumber, &currQuant, &rollsCount)
	if err != nil {
		panic(err)
	}
	// get existing items in storage if any
	var data []byte
	_ = s.Server.db.QueryRow("SELECT gacha_items FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if len(data) == 0 {
		data = []byte{0x00}
	}
	// get gacha items and iterate through them for gacha roll
	shopEntries, err := s.Server.db.Query(`SELECT itemhash, percentage, rarityIcon, itemCount, itemType, itemId, quantity
		FROM gacha_shop_items
		WHERE shophash=$1 AND entryType=100
		ORDER BY itemhash
		EXCEPT ALL SELECT itemhash, percentage, rarityIcon, itemCount, itemType, itemId, quantity
		FROM gacha_shop_items gsi JOIN lucky_box_state lbs ON gsi.itemhash = ANY(lbs.used_itemhash)
		WHERE lbs.char_id=$2`, pkt.GachaHash, s.CharID)
	if err != nil {
		panic(err)
	}
	for shopEntries.Next() {
		err = shopEntries.Scan(&itemhash, &percentage, &rarityIcon, &itemCount, (*pq.Int64Array)(&itemType), (*pq.Int64Array)(&itemId), (*pq.Int64Array)(&quantity))
		if err != nil {
			panic(err)
		}
		items = append(items, &gachaItem{itemhash: itemhash, percentage: percentage, rarityIcon: rarityIcon, itemCount: itemCount, itemType: itemType, itemId: itemId, quantity: quantity})
	}
	// execute rolls, build response and update database
	results := byte(0)
	resp := byteframe.NewByteFrame()
	dbUpdate := byteframe.NewByteFrame()
	resp.WriteUint8(0) // results go here later
	l := lottery.NewDefaultLottery()
	for x := 0; x < int(rollsCount); x++ {
		ind := l.Draw(items)
		results += items[ind].(*gachaItem).itemCount
		for y := 0; y < int(items[ind].(*gachaItem).itemCount); y++ {
			// items in storage don't get rarity
			dbUpdate.WriteUint8(uint8(items[ind].(*gachaItem).itemType[y]))
			dbUpdate.WriteUint16(uint16(items[ind].(*gachaItem).itemId[y]))
			dbUpdate.WriteUint16(uint16(items[ind].(*gachaItem).quantity[y]))
			data = append(data, dbUpdate.Data()...)
			// nolint:errcheck // Error return value of `dbUpdate.Seek` is not checked
			dbUpdate.Seek(0, 0)
			// response needs all item info and the rarity
			resp.WriteBytes(dbUpdate.Data())
			resp.WriteUint8(items[ind].(*gachaItem).rarityIcon)

			usedItemHash = append(usedItemHash, int64(items[ind].(*gachaItem).itemhash))
		}
		// remove rolled
		items = append(items[:ind], items[ind+1:]...)
	}
	// nolint:errcheck
	resp.Seek(0, 0)
	resp.WriteUint8(results)
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())

	// add claimables to DB
	data[0] = data[0] + results
	_, err = s.Server.db.Exec("UPDATE characters SET gacha_items = $1 WHERE id = $2", data, s.CharID)
	if err != nil {
		s.logger.Fatal("Failed to update gacha_items in db", zap.Error(err))
	}
	// update lucky_box_state
	_, err = s.Server.db.Exec(`INSERT INTO lucky_box_state (char_id, shophash, used_itemhash)
		VALUES ($1,$2,$3) ON CONFLICT (char_id, shophash)
		DO UPDATE SET used_itemhash = COALESCE(lucky_box_state.used_itemhash::int[] || $3::int[], $3::int[])
		WHERE EXCLUDED.char_id=$1 AND EXCLUDED.shophash=$2`, s.CharID, pkt.GachaHash, usedItemHash)
	if err != nil {
		s.logger.Fatal("Failed to update lucky box state in db", zap.Error(err))
	}
	// deduct gacha coins if relevant, items are handled fine by the standard savedata packet immediately afterwards
	if currType == 19 {
		_, err = s.Server.db.Exec(`UPDATE characters
			SET gacha_trial = CASE WHEN (gacha_trial > $1) then gacha_trial - $1 else gacha_trial end, gacha_prem = CASE WHEN NOT (gacha_trial > $1) then gacha_prem - $1 else gacha_prem end
			WHERE id=$2`, currNumber, s.CharID)
	}
	if err != nil {
		s.logger.Fatal("Failed to update gacha_trial in db", zap.Error(err))
	}
}

func handleMsgMhfResetBoxGachaInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfResetBoxGachaInfo)
	_, err := s.Server.db.Exec("DELETE FROM lucky_box_state WHERE shophash=$1 AND char_id=$2", pkt.GachaHash, s.CharID)
	if err != nil {
		panic(err)
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}
