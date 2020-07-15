package main

import (
	"fmt"
	"github.com/hexvalid/midori-go/database"
	"github.com/hexvalid/midori-go/tormdr"
)

func main() {
	//fmt.Println(getnada.GenerateMail())
	db, _ := database.Open("midori-go.db")
	x, err := database.GetAllAccounts(db)
	fmt.Println(err)

	tormdrN, _ := tormdr.NewTorMDR(1, &tormdr.Config{
		TorMDRBinaryPath: "/usr/bin/tormdr",
		DataDirectory:    "/tmp/tormdr_data",
	})

	if err := tormdrN.Start(); err != nil {
		panic(err)
	}
	tormdrN.Start()

	x[3].OpenBrowser(tormdrN)
	x[3].Home()
	x[3].Roll()

	tormdrN.Stop()

	/*	a, _ := bot.GenerateNewAccount(0)

		tormdrN, _ := tormdr.NewTorMDR(1, &tormdr.Config{
			TorMDRBinaryPath: "/usr/bin/tormdr",
			DataDirectory:    "/tmp/tormdr_data",
		})

		if err := tormdrN.Start(); err != nil {
			panic(err)
		}
		tormdrN.Start()
		a.OpenBrowser(tormdrN)
		fmt.Println(a.Login(true))
		fmt.Println(a.Home())
		fmt.Println(a.VerifyEmail())
		time.Sleep(10 * time.Second)
		fmt.Println(a.Home())
		a.Roll()

		database.InsertAccount(db, &a)

		tormdrN.Stop()*/

	/*for {
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
	}*/

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

	fmt.Println(tormdrN.CheckIP())

	err = tormdrN.Stop()
	if err != nil {
		panic(err)
	}*/
}

/*You have an invalid email address attached to your account. Please change it to one that is valid (by clicking on the PROFILE button in the top bar) so that you can receive important emails from us in the future.
We recommend signing up to and using Google Mail instead of your current email provider.
To get specific details about the error, please click here.
If you think this is a mistake, please click here to validate your email address.*/
