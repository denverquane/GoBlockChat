package wallet

import (
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
	"strconv"
	"crypto/x509"
	"encoding/pem"
	"crypto/rsa"
)

type Wallet struct {
	lastProcessedBlock int
	addresses          []transaction.PersonalAddress
	balance            float64
	ChannelRecords     map[string]ChannelRecord //map of all channels we (might) have access to
}

func MakeNewWallet() Wallet {
	address := transaction.GenerateNewPersonalAddress()
	fmt.Println(address.Address)
	return Wallet{-1, []transaction.PersonalAddress{address}, 0.0, make(map[string]ChannelRecord, 0)}
}

func (wallet Wallet) GetAddress() transaction.PersonalAddress {
	return wallet.addresses[0]
}

func (wallet Wallet) GetBalances() string {
	str := "Wallet Balance: " + strconv.FormatFloat(wallet.balance, 'f', -1, 64) + " REP"
	for name, _ := range wallet.ChannelRecords {
		str += "\n                " + "1 " + name
	}

	return str
}

//TODO This is for testing!!! Don't rely on this!
func (wallet Wallet) getOriginInfo() transaction.OriginInfo {
	return transaction.AddressToOriginInfo(wallet.addresses[0])
}

func (wallet Wallet) MakeTransaction(quantity float64, currency string, dest transaction.Base64Address) transaction.SignedTransaction {
	unsigned := transaction.SignedTransaction{wallet.getOriginInfo(), dest, quantity, currency,
		"Sending!", nil, nil}
	return unsigned.SignMessage(&wallet.addresses[0].PrivateKey)
}

func (wallet *Wallet) UpdateBalances(blockchain blockchain.BlockChain) {
	for _, addr := range wallet.addresses {
		wallet.balance += blockchain.GetAddrBalanceFromInclusiveIndex(wallet.lastProcessedBlock+1, addr.Address, "REP")
	}

	newCurrencies := wallet.getNewTokenRecords(blockchain)
	//for _, v := range newCurrencies {
	//	fmt.Println(v.channelPublic)
	//	signed := v.MakeSendMyPublicTransaction()
	//	fmt.Println(signed.ToString())
	//}
	wallet.ChannelRecords = mergeChannelMaps(wallet.ChannelRecords, newCurrencies)

	wallet.lastProcessedBlock = int(blockchain.GetNewestBlock().Index)
	fmt.Println(wallet.GetBalances())
}

func (wallet Wallet) getNewTokenRecords(chain blockchain.BlockChain) map[string]ChannelRecord {
	currencies := make(map[string]ChannelRecord, 0)
	for i, block := range chain.Blocks {
		if i > wallet.lastProcessedBlock {
			for _, tx := range block.Transactions {
				if tx.SignedTrans.Currency != "REP" { //don't even bother with the rep ones we should've already processed
					if _, recordExists := wallet.ChannelRecords[tx.SignedTrans.Currency]; !recordExists {
						//ensure we don't already know about this currency
						for _, addr := range wallet.addresses {
							if tx.SignedTrans.DestAddr == addr.Address { //we received a transaction
								fmt.Println("Received " + tx.SignedTrans.Currency)
								if _, ok := currencies[tx.SignedTrans.Currency]; ok {
									fmt.Println("RECEIVED TOKEN FOR A CURRENCY I ALREADY HAVE!")
								} else {
									record := GenerateNewChannelRecord(tx.SignedTrans.Currency,
										tx.SignedTrans.Origin.Address, addr)
									//make a new channel record


									block, _ := pem.Decode([]byte(tx.SignedTrans.Payload))
									if block == nil {
										fmt.Println("failed to parse PEM block containing the key")
									}
									pubkey, _ := x509.ParsePKIXPublicKey(block.Bytes)
									fmt.Println("Pubkey: ")
									fmt.Println(pubkey)
									record.channelPublic = *pubkey.(*rsa.PublicKey)

									//switch pub := pubkey.(type) {
									//	case *rsa.PublicKey:
									//		record.channelPublic = *pub
									//		fmt.Println("\nHEYRecive\n")
									//		fmt.Println(record.channelPublic)
									//}

									record.status = ReceivedTokenAndChannelPub
									currencies[tx.SignedTrans.Currency] = record
								}
							}
						}
					}
				}
			}
		}
	}
	return currencies
}

func mergeChannelMaps(map1 map[string]ChannelRecord, map2 map[string]ChannelRecord) map[string]ChannelRecord {
	map3 := map1
	for name2, val2 := range map2 {
		map3[name2] = val2
	}
	return map3
}
