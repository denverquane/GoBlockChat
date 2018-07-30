package blockchain

import (
	"github.com/denverquane/GoBlockShare/common"
	"log"
	"time"
	"math/rand"
)

func (chain *BlockChain) AddMockTransactions() {
	torr, _ := common.MakeTorrentFileFromFile(1000, "README.md", "readme")
	for i := 0; i < 3; i++ {
		address := common.GenerateNewPersonalAddress()
		trans := common.PublishTorrentTrans{torr}
		origin := common.AddressToOriginInfo(address)
		btt := common.SignableTransaction{origin, trans, common.PUBLISH_TORRENT, nil, nil, ""}
		signed := btt.SignAndSetTxID(&address.PrivateKey)
		log.Println("Gonna broadcast " + signed.TxID + " to blockchains")
		worked, err := chain.AddTransaction(signed, "test addr")
		if !worked {
			log.Println(err.Error())
		}

		for chain.IsProcessing() {
			time.Sleep(100)
		}

		address2 := common.GenerateNewPersonalAddress()
		res := rand.Intn(2) == 0
		res1 := rand.Intn(2) == 0
		res2 := rand.Intn(2) == 0

		trans2 := common.TorrentRepTrans{signed.TxID, signed.Transaction.(common.PublishTorrentTrans).Torrent.TotalHash,
			common.RepMessage{res, res1, res2}}
		origin2 := common.AddressToOriginInfo(address2)
		btt2 := common.SignableTransaction{origin2, trans2, common.TORRENT_REP, nil, nil, ""}
		signed2 := btt2.SignAndSetTxID(&address2.PrivateKey)
		log.Println("Gonna broadcast " + signed2.TxID + " to blockchains")
		worked, err = chain.AddTransaction(signed2, "test addr")
		if !worked {
			log.Println(err.Error())
		}

		for chain.IsProcessing() {
			time.Sleep(100)
		}

		address3 := common.GenerateNewPersonalAddress()
		address4 := common.GenerateNewPersonalAddress()
		trans3 := common.SharedLayerTrans{torr.LayerHashKeys[0], address4.Address}
		origin3 := common.AddressToOriginInfo(address3)
		btt3 := common.SignableTransaction{origin3, trans3, common.SHARED_LAYER, nil, nil, ""}
		signed3 := btt3.SignAndSetTxID(&address3.PrivateKey)
		log.Println("Gonna broadcast " + signed3.TxID + " to blockchains")
		worked, err = chain.AddTransaction(signed3, "test addr")
		if !worked {
			log.Println(err.Error())
		}

		for chain.IsProcessing() {
			time.Sleep(100)
		}


		trans4 := common.LayerRepTrans{signed3.TxID, torr.LayerHashKeys[0],
			true, true}
		origin4 := common.AddressToOriginInfo(address4)
		btt4 := common.SignableTransaction{origin4, trans4, common.LAYER_REP, nil, nil, ""}
		signed4 := btt4.SignAndSetTxID(&address4.PrivateKey)
		log.Println("Gonna broadcast " + signed4.TxID + " to blockchains")
		worked, err = chain.AddTransaction(signed4, "test addr")
		if !worked {
			log.Println(err.Error())
		}
	}


}
