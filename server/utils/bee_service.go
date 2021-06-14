package utils

import (
	"encoding/json"
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	// "github.com/panjf2000/ants/v2" # goroutine pool, TODO
)

type GetNonceParams struct {
	Address string `json:"address"`
}

type Peer struct {
	Peer         string       `json:"peer"`
	LastSent     interface{}  `json:"lastsent"`
	LastReceived LastReceived `json:"lastreceived"`
}

type LastReceived struct {
	Beneficiary string `json:"beneficiary"`
	Chequebook  string `json:"chequebook"`
	Payout      int    `json:"payout"`
}

type CashInfo struct {
	Total    int
	UnCashed int
}

const BZZ_DECIMAL = 100000000000000000
const MIN_AMOUT = 5

type NodeState struct {
	Address       string
	CashCount     int
	TotalUnCashed float64
	EthBalance    float64
	BzzBalance    float64
	PeerCount     int
	Version       string
}

func GetBeeNodeStateInConcurrently(beeNodess []model.BeeNodes) {
	var wg sync.WaitGroup
	for index, node := range beeNodess {
		wg.Add(1)
		go func(index int, node model.BeeNodes) {
			defer wg.Done()
			nodeState := GetBeeNodeState(node.Ip, strconv.Itoa(node.DebugPort))
			node.BzzBalance = nodeState.BzzBalance
			node.EthBalance = nodeState.EthBalance
			node.PeerCount = nodeState.PeerCount
			node.UncashedCount = nodeState.CashCount
			node.UncashedAmount = nodeState.TotalUnCashed
			node.WalletAddress = nodeState.Address
			node.Version = nodeState.Version
			node.UpdatedAt = time.Now()
			beeNodess[index] = node
			// fmt.Printf("index %v, %v, %v\n", index, beeNodess[index], nodeState)
		}(index, node)
	}
	wg.Wait()
}

func GetBeeNodeState(ip string, port string) (nodeState *NodeState) {
	nodeState = new(NodeState)
	address := getAddress(ip, port)
	if address == "" {
		return
	}
	peers := getPeers(ip, port)
	if peers == nil {
		return
	}
	peerMap := make(map[string]Peer)
	for _, peer := range peers {
		peerMap[peer.Peer] = peer
	}
	cashMap := getAllCash(ip, port, peerMap)
	len := len(cashMap)
	total_amount := 0
	total_uncashed := 0
	for _, cash := range cashMap {
		// fmt.Printf("%v,%v\n", peer, cash)
		total_amount += cash.Total
		total_uncashed += cash.UnCashed
	}
	nodeState.Address = address
	nodeState.CashCount = len
	nodeState.TotalUnCashed = float64(total_uncashed) / BZZ_DECIMAL
	nodeState.EthBalance, _ = GetEthBalance(address).Float64()
	nodeState.BzzBalance, _ = GetTokenBalance(BZZ_CONTRACT, address, 16).Float64()
	nodeState.PeerCount = getPeerLength(ip, port)
	nodeState.Version = healthCheck(ip, port)

	// fmt.Printf("node state : %+v\n", nodeState)
	return
}

// return version
func healthCheck(ip string, port string) string {
	msgMap := QueryBee(ip, port, "health")
	if result, ok := msgMap.(map[string]interface{})["version"]; ok {
		if result == nil {
			return ""
		}
		return result.(string)
	}
	return ""
}

func getAddress(ip string, port string) string {
	msgMap := QueryBee(ip, port, "addresses")
	if msgMap == nil {
		return ""
	}
	return msgMap.(map[string]interface{})["ethereum"].(string)
}

func getPeers(ip string, port string) []Peer {
	peers := make([]Peer, 0)

	msgMap := QueryBee(ip, port, "chequebook/cheque").(map[string]interface{})
	if msgMap == nil {
		return nil
	}
	byteData, _ := json.Marshal(msgMap["lastcheques"])
	err := json.Unmarshal(byteData, &peers)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return peers
}

func getAllCash(ip string, port string, peerMap map[string]Peer) map[string]CashInfo {
	cashMap := make(map[string]CashInfo)
	for key, peer := range peerMap {
		lastCashedPayout := getLastCashedPayout(ip, port, key)
		left := peer.LastReceived.Payout - lastCashedPayout
		if left > 0 {
			cashMap[key] = CashInfo{peer.LastReceived.Payout, left}
		}
	}
	return cashMap
}

func getLastCashedPayout(ip string, port string, peer string) int {
	msgMap := QueryBee(ip, port, "chequebook/cashout/"+peer)
	if msgMap == nil {
		return 0
	}
	// log.Printf("%v, %+v\n", peer, msgMap)
	if result, ok := msgMap.(map[string]interface{})["result"]; ok {
		if result == nil {
			return 0
		}
		if lastPayout, ok := result.(map[string]interface{})["lastPayout"]; ok {
			return int(lastPayout.(float64))
		}
		return 0
	}
	return 0
}

func getPeerLength(ip string, port string) int {
	msgMap := QueryBee(ip, port, "peers").(map[string]interface{})
	if msgMap == nil {
		return 0
	}
	return len(msgMap["peers"].([]interface{}))
}

func QueryBee(ip string, port string, path string) interface{} {
	url := "http://" + ip + ":" + port + "/" + path
	client := http.Client{
		Timeout: time.Second * 3, // Timeout after 2 seconds
	}
	// global.GVA_LOG.Info("QueryBee: " + url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		global.GVA_LOG.Info(err.Error())
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Info(err.Error())
		return nil
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)

	var msgMapTemplate interface{}
	err = json.Unmarshal([]byte(body), &msgMapTemplate)
	if err != nil {
		global.GVA_LOG.Info(err.Error())
		return nil
	}
	return msgMapTemplate
}

func PostBee(ip string, port string, path string, b io.Reader, headers map[string]string) interface{} {
	url := "http://" + ip + ":" + port + "/" + path
	client := http.Client{
		Timeout: time.Second * 10, // Timeout after 2 seconds
	}
	req, err := http.NewRequest(http.MethodPost, url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		global.GVA_LOG.Info(err.Error())
		return nil
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Info(err.Error())
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)

	var msgMapTemplate interface{}
	err = json.Unmarshal([]byte(body), &msgMapTemplate)
	if err != nil {
		global.GVA_LOG.Info(err.Error())
		return nil
	}
	return msgMapTemplate
}
