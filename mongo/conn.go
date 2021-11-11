package mongo

import (
	"fmt"
	"os"
	"sync"

	"github.com/OhYee/blotter/utils/initial"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _clientOptions *options.ClientOptions = nil
var _clientOptionsOnce = sync.Once{}

func initClientOptions() {
	_clientOptionsOnce.Do(func() {
		fmt.Println("Initial mongodb")
		addr := os.Getenv("mongoURI")
		if addr == "" {
			addr = "127.0.0.1:27017"
		}
		_clientOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", addr))
		fmt.Println("Initial mongodb finished")
	})
}

func getClientOptions() *options.ClientOptions {
	if _clientOptions == nil {
		initClientOptions()
	}
	return _clientOptions

}

func init() {
	initial.Register(initClientOptions)
}
