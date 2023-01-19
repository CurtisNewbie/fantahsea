package fclient

import (
	"context"
	"testing"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/consul"
)

func TestGenFileTempTokens(t *testing.T) {
	common.LoadConfigFromFile("../app-conf-dev.yml")			
	consul.MustInitConsulClient()

	keys := []string{"ZZZ473737645015040965849", "ZZZ473734540804096965849"}
	tokenMap, err := GenFileTempTokens(context.Background(), keys)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokenMap) < 1 {
		t.Fatal("map is empty")
	}
	for _, k := range keys {
		if _, ok := tokenMap[k]; !ok {
			t.Fatalf("map doesn't have token for key '%s'", k)
		}
	}
	t.Logf("%+v", tokenMap)
}