package main

import (
	"github.com/hexvalid/midori-go/getnada"
	"sync"
	"time"
)

func main() {
	//fmt.Println(getnada.GenerateMail())
	for {
		x, _ := getnada.GetInbox("test@getnada.com")
		var wg sync.WaitGroup
		for _, mail := range x {
			wg.Add(1)
			go func() {
				mail.Delete()
				wg.Done()
			}()
			time.Sleep(10 * time.Millisecond)
		}
		wg.Wait()
	}
	/*tormdrN, _ := tormdr.NewTorMDR(1, &tormdr.Config{
		TorMDRBinaryPath:    "/usr/bin/tormdr",
		DataDirectory:       "/tmp/tormdr_data",
		KeepalivePeriod:     60,
		UseSocks5Proxy:      true,
		Socks5ProxyAddress:  "40.70.243.118:32416",
		Socks5ProxyUserName: "e4cf6e290c0cf8ae8fb91fcf818e1e40",
		Socks5ProxyPassword: "a565ab1f3802afbf4d07c1674069d813",
	})

	if err := tormdrN.Start(); err != nil {
		panic(err)
	}

	en, err := tormdr.FindExitNode(nil, 25000, true, true, true)

	_ = tormdrN.SetExitNode(en)

	if err != nil {
		panic(err)
	}

	fmt.Println(tormdrN.TestIP())

	err = tormdrN.Stop()
	if err != nil {
		panic(err)
	}*/
}
