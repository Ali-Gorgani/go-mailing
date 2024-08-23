package main

import (
	"go-mailing/cmd/goMailing"
)

func main() {
	// kn := kavenegar.New("566346566C492F4E682F703148745135506948744F43515234586247586C4373755A7445746E704D7766383D")

	// output, err := kn.Send(kavenegar.SendInputParams{
	// 	Receptor: []string{"09114418131"},
	// 	Message:  "Hello, World!",
	// 	Sender:   "10008663",
	// })

	// output, err := kn.Select([]int32{643239283, 643239284})

	// output, err := kn.SelectOutbox(kavenegar.SelectOutboxInputParams{
	// 	StartDate: 1724299528,
	// })

	// output, err := kn.LatestOutBox(kavenegar.LatestOutboxInputParams{})

	// output, err := kn.CountOutbox(kavenegar.CountOutboxInputParams{
	// 	StartDate: 1724132490,
	// })

	// output, err := kn.Cancel([]int64{647491105})

	// output, err := kn.Receive(kavenegar.ReceiveInputParams{
	// 	LineNumber: "10008663",
	// 	IsRead: 0,
	// })

	// output, err := kn.Info()

	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(output)

	err := goMailing.StartServer()
	if err != nil {
		panic(err)
	}
}
