package network

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"encoding/json"
	"io"
	"chatProgram/blockchain"
	"reflect"
	"bufio"
	"os"
	"strings"
)

var globalChain* blockchain.BlockChain
var Peers []string

func MakeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/postTransaction", handleWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/chainUpdate", handleChainUpdate).Methods("POST")
	return muxRouter
}

func FetchOrMakeChain(router http.Handler) {

	// TODO discover existing chain here, using router

	users := make([]blockchain.UserPassPair, 1)
	users[0] = blockchain.UserPassPair{"user", "pass"}
	chain := blockchain.MakeInitialChain(users)
	globalChain = &chain
}

func handleGetBlockchain(w http.ResponseWriter, _ *http.Request) {
	data, err := json.MarshalIndent(globalChain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("GET")
	w.Header().Set(	"Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	io.WriteString(w, string(data))
}

func handleChainUpdate(w http.ResponseWriter, r *http.Request) {
	var m blockchain.BlockChain

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	if m.IsValid() {
		if m.Len() > globalChain.Len() {
			if blockchain.AreChainsSameBranch(m, *globalChain) {
				globalChain = &m
				fmt.Println("Valid blockchain supplied! Replaced with: ")
				spew.Dump(globalChain)
				//BroadcastToAllPeers(Peers, *globalChain)
			} else {
				fmt.Println("Chains are of different branches! Not replacing mine")
			}
		} else if globalChain.Len() == 1 && m.Len() == 1 {
			globalChain = &m
			fmt.Println("Both chains are 1 length; replacing mine!")
			//spew.Dump(Chain)
		}else {
			fmt.Println("Chains are the same or lesser length; not replacing anything")
		}
	} else {
		fmt.Println("Invalid blockchain supplied; not replacing anything")
	}

	data, err := json.MarshalIndent(globalChain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(data))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m blockchain.AuthTransaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	if m.TransactionType == "" || m.Message == "" {
		respondWithJSON(w, r, http.StatusBadRequest, "Supply transaction in this format: " + blockchain.GetTransactionFormat())
		return
	}

	oldBlock := globalChain.GetNewestBlock()
	newBlock := blockchain.GenerateBlock(oldBlock, m)

	if reflect.DeepEqual(oldBlock,newBlock) {
		respondWithJSON(w, r, http.StatusUnauthorized, "User is not authorized to post")
		return
	}
	fmt.Println("New block:\n" + newBlock.ToString())

	if blockchain.IsBlockSequenceValid(newBlock, oldBlock) {
		globalChain.Blocks = append(globalChain.Blocks, newBlock)
		//Block = blockchain.CheckLongerChain(newBlock, Block)
		fmt.Println("Successfully added: {" + m.RemovePassword().ToString() + "} to the chain")
		BroadcastToAllPeers(Peers, *globalChain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, _ *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Header().Set(	"Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Write(response)
}

func getPeersFromInput(){
	var done = false
	for !done {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter URL to broadcast to, WITH port (ex: 127.0.0.1:8090), or \"quit\" if you're done: ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", 1)
		text = strings.Replace(text, " ", "", 1)
		fmt.Println("Entered: \"" + text + "\"")
		if text == "quit" {
			done = true
		} else if text == "" {
			fmt.Println("Empty string supplied")
		} else {
			BroadcastChain("http://" + text + "/chainUpdate", *globalChain)
			Peers = append(Peers, "http://" + text + "/chainUpdate")
		}
	}
}