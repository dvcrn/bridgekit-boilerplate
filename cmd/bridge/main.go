package main

import (
	"context"
	"github.com/dvcrn/bridgekit-boilerplate/internal"
	"github.com/dvcrn/bridgekit/pkg"
)

func main() {
	br := pkg.NewBridgeKit(
		"MyBridge",
		"sh-mybridge",
		"",
		"Integration",
		"1.0",
		&internal.Config{},
		internal.ExampleConfig,
	)
	connector := internal.NewBridgeConnector(br)
	br.StartBridgeConnector(context.Background(), connector)
}
