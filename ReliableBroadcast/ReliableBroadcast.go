/*
  Construido como parte da disciplina de Sistemas Distribuidos
  PUCRS - Escola Politecnica
  Professor: Fernando Dotti  (www.inf.pucrs.br/~fldotti)
  Algoritmo baseado no livro:
  Introduction to Reliable and Secure Distributed Programming
  Christian Cachin, Rachid Gerraoui, Luis Rodrigues

  Implementa��o do Eager Reliable Broadcast
  Semestre 2019/1
  Estudantes:  Vinicius Sesti e Gabriel Waengertner
*/

package ReliableBroadcast

import "fmt"
import "strings"
import BEB "../BestEffortBroadcast"

type ReliableBroadcast_Req_Message struct {
	Addresses []string
	Message   string
	Sender    string
}

type ReliableBroadcast_Ind_Message struct {
	From    string
	Sender  string
	Message string
}

type ReliableBroadcast_Module struct {
	Addresses []string
	Self      string

	Delivered           map[string]bool
	Req                 chan ReliableBroadcast_Req_Message
	Ind                 chan ReliableBroadcast_Ind_Message
	BestEffortBroadcast BEB.BestEffortBroadcast_Module
}

func (module ReliableBroadcast_Module) Init() {

	fmt.Println("Init RB!")
	module.BestEffortBroadcast = BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message, 1),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message, 1)}
	module.BestEffortBroadcast.Init(module.Self)
	module.Start()

}

func (module ReliableBroadcast_Module) Start() {
	go func() {
		for {
			select {
			case y := <-module.Req:
				module.Broadcast(y)
			case y := <-module.BestEffortBroadcast.Ind:
				module.Deliver(BEB2RB(y))
			}
		}
	}()
}

// func (module ReliableBroadcast_Module) Start() {
// 	go func() {
// 		for {
// 			y := <-module.Req
// 			module.Broadcast(y)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			y := <-module.BestEffortBroadcast.Ind
// 			module.Deliver(BEB2RB(y))
// 		}
// 	}()
// }

func (module ReliableBroadcast_Module) Broadcast(message ReliableBroadcast_Req_Message) {
	fmt.Println(module.Self + " --- RB: got message: " + message.Message)
	module.BestEffortBroadcast.Req <- RB2BEB(message)
	fmt.Println(module.Self + " --- RB: delivered message " + message.Message + " to BEB")
}

func (module ReliableBroadcast_Module) Deliver(message ReliableBroadcast_Ind_Message) {

	key := message.Sender + ";" + message.Message
	// fmt.Println(module.Self + " --- RB: Deliver: Received: " + message.Message + " originally by " + message.Sender + " from " + message.From)

	_, found := module.Delivered[key]
	if found {
		// fmt.Println(module.Self + " --- Message Already Received!")
		return
	}

	module.Delivered[key] = true
	module.Ind <- message

	module.BestEffortBroadcast.Req <- RB2BEB(message.Retransmit(module, message.Sender))
	// fmt.Println(module.Self + " --- RB: Delivered to CBTOB")
}

func (message ReliableBroadcast_Ind_Message) Retransmit(module ReliableBroadcast_Module, originalSender string) ReliableBroadcast_Req_Message {
	// fmt.Println("RETRANSMIT: mensagem antes de adicionar novo sender: " + message.Message)

	return ReliableBroadcast_Req_Message{
		Addresses: module.Addresses,
		Message:   message.Message,
		Sender:    originalSender}

}

func RB2BEB(message ReliableBroadcast_Req_Message) BEB.BestEffortBroadcast_Req_Message {

	return BEB.BestEffortBroadcast_Req_Message{
		Addresses: message.Addresses,
		Message:   message.Sender + ";" + message.Message}

}

func BEB2RB(message BEB.BestEffortBroadcast_Ind_Message) ReliableBroadcast_Ind_Message {

	parts := strings.Split(message.Message, ";")
	content := strings.Join(parts[1:], ";")
	sender := parts[0]

	return ReliableBroadcast_Ind_Message{
		From:    message.From,
		Sender:  sender,
		Message: content}

}

// func main() {

// 	if len(os.Args) < 2 {
// 		fmt.Println("Please specify at least one address:port!")
// 		return
// 	}

// 	addresses := os.Args[1:]
// 	fmt.Println(addresses)

// 	mod := ReliableBroadcast_Module{
// 		Req:       make(chan ReliableBroadcast_Req_Message),
// 		Ind:       make(chan ReliableBroadcast_Ind_Message),
// 		Delivered: make(map[string]bool),
// 		Addresses: addresses[1:],
// 		Self:      addresses[0]}
// 	mod.Init()

// 	msg := ReliableBroadcast_Req_Message{
// 		Addresses: addresses,
// 		Sender:    addresses[0],
// 		Message:   "BATATA!"}

// 	yy := make(chan string)
// 	mod.Req <- msg
// 	<-yy
// }
