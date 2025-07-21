package envtest

import (
	"context"
	"fmt"
	"log"

	"github.com/blocto/solana-go-sdk/client"
)

func test() {
	c := client.NewClient("https://solana-devnet.g.alchemy.com/v2/M0quawxer68-8mXsa8UBw1y4Dtmp-qx3")
	// If you would like to customize the http client used to make the
	// requests you could do something like this
	// c := client.New(rpc.WithEndpoint(rpc.MainnetRPCEndpoint),rpc.WithHTTPClient(customHTTPClient))

	resp, err := c.GetVersion(context.TODO())
	if err != nil {
		log.Fatalf("failed to version info, err: %v", err)
	}

	fmt.Println("version", resp.SolanaCore)
	// version 2.3.2
}
