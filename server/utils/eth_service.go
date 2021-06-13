package utils

// https://goethereumbook.org/smart-contract-read-erc20/
import (
	"context"
	"gin-vue-admin/global"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const RPC_URL = "http://159.138.149.40:8888"
const BZZ_CONTRACT = "0x2ac3c1d3e24b45c6c310534bc2dd84b5ed576335"

func GetEthBalance(address string) *big.Float {
	fbal := new(big.Float)
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		global.GVA_LOG.Info(err.Error())
		return fbal
	}
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		global.GVA_LOG.Info(err.Error())
		return fbal
	}
	fbal.SetString(balance.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(18))))
	return value
}

func GetTokenBalance(contract string, address string, decimals int) *big.Float {
	fbal := new(big.Float)
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		global.GVA_LOG.Info(err.Error())
		return fbal
	}

	tokenAddress := common.HexToAddress(contract)
	instance, err := NewToken(tokenAddress, client)
	if err != nil {
		global.GVA_LOG.Info(err.Error())
		return fbal
	}

	accountAddress := common.HexToAddress(address)
	bal, err := instance.BalanceOf(&bind.CallOpts{}, accountAddress)
	if err != nil {
		global.GVA_LOG.Info(err.Error())
		return fbal
	}
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	return value
}

/*
func main() {
	fmt.Println(GetEthBalance("0xba05240EbF7257b04eFa859345d2b945820fD7CB"))
	fmt.Println(GetTokenBalance(BZZ_CONTRACT ,"0xba05240EbF7257b04eFa859345d2b945820fD7CB", 16))
        }
*/
