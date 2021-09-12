package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type Block struct {
	Transaction string
	PrevPointer *Block
}

func RemoveIndex(s []net.Conn, index int) []net.Conn {
	return append(s[:index], s[index+1:]...)
}
func input(s string) {
	time.Sleep(time.Second * 10)
	s = "go"
}

var partition bool = false

func routine() {
	conn, err := net.Dial("tcp", "localhost:"+"2022")
	if err != nil {
		fmt.Println("no server is Listening at port # ", "2022")
	} else {
		var recvdBlock Block
		dec := gob.NewDecoder(conn)
		err = dec.Decode(&recvdBlock)
		if err != nil {
			fmt.Println("Error, Cant find message")
		}
	}
	//conn.SetDeadline(time.Time{})
	heartResponse1 := make([]byte, 1024)
	// var client_msg []byte// = make([]byte, 1024)
	_, err = conn.Read(heartResponse1[0:])
	// assci := [9]byte{112,97,114,116,105,116,105,111,110}

	fmt.Println(" msg is --------->", string(heartResponse1))
	partition = true

}

func handleConnection(c net.Conn, s string) {

	log.Println("A client new has connected", c.RemoteAddr())
	block1 := &Block{"successfully connected from server " + s, nil}
	block2 := &Block{s, block1}

	gobEncoder := gob.NewEncoder(c)
	err := gobEncoder.Encode(block2)
	if err != nil {
		log.Println(err)
	}
}

func main() {

	if os.Args[1] == "client" {

		var port string = ":" + os.Args[2]

		ln, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Service Started at Port Number  ", port, "\n")
		for {
			conn, err := ln.Accept()
			//time.Sleep(10 * time.Second)
			fmt.Print("Enter the message: ") //Print function is used to display output in same line
			var first string
			fmt.Scanln(&first)
			//heartbeat := "message   is   from the client"
			conn.Write([]byte(first))
			if err != nil {
				log.Println(err)
				continue
			}
			go handleConnection(conn, first)

		}

	} else { // ch

		var port string = ":" + os.Args[1]

		var IndCon = 0
		var mynum int = len(os.Args) - 2

		var total int = 4 // enter the no of nodes here
		rand.Seed(time.Now().UnixNano())
		min := 100
		max := 300
		reply := make([]byte, 1024)
		var loop int = 0
		var counter int = total - 2
		excluded := make([]int, total-1)
		//connectionIndex := make()
		//excluded[mynum] = 1
		//p := "abc"
		winnerNum := -1
		var count int = 0
		check := false
		winner := -1
		connections := make([]net.Conn, total-1)
		loopCheck := false
		ack := "received"

		log.Println("Current Node Number : ", mynum)
		var size int = len(os.Args)
		for i := 2; i < size; i++ {
			conn, err := net.Dial("tcp", "localhost:"+os.Args[i])
			if err != nil {
				fmt.Println("no server is Listening at port # ", os.Args[i])
			} else {
				var recvdBlock Block
				dec := gob.NewDecoder(conn)
				err = dec.Decode(&recvdBlock)

				connections[IndCon] = conn
				IndCon++
				/*if mynum == IndCon { // skipping for self connection index
					IndCon++
				}
				connections[IndCon] = conn
				IndCon++
				*/
				if err != nil {
					fmt.Println("Error, Cant find message")

				}
			}
			counter = counter - 1
		}
		ln, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Service Started at Port Number  ", port, "\n")

		for num := 0; num <= counter; num++ {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go handleConnection(conn, port)
			connections[IndCon] = conn
			IndCon++
			/*
				if mynum == IndCon {
					IndCon++
				} // storing connections for ahead notes
				connections[IndCon] = conn
				IndCon++*/
		}
		fmt.Println("Cluster completed ")
		total = total - 1

		for l := 0; l < 5; l++ { // again elections
			if count >= total {
				break
			}

			println("Another election started ")
			loop = 0
			check = false
			winner = -1
			randtimer := rand.Intn(max-min+1) + min
			for i := 0; i < randtimer; i++ {
				loopCheck = false
				for j := 0; j < total; j++ {
					if excluded[loop] != 1 {
						loopCheck = true
					}
				}

				if loopCheck {
					// disabling block of read
					connections[loop].SetReadDeadline(time.Now().Add(time.Millisecond))
					_, err = connections[loop].Read(reply[0:])
					if err != nil {
					} else {
						println("Message Received: = ", string(reply), " From Node: ", loop)
						check = true
						winner = loop // this node's timer was ended the first
						break
					}
				}
				loop++
				if loop == total {
					loop = 0
				}

			}

			vote := "yes"
			if check == true { // loser nodes sending vote to winner
				fmt.Println("Writing vote")
				connections[winner].Write([]byte(vote))
			}

			if check == false {
				fmt.Println(" Node timed out: ", mynum)
				message := "stop"

				go routine()

				for i := 0; i < total; i++ {
					loopCheck = false
					for j := 0; j < total; j++ {
						if excluded[i] != 1 {
							loopCheck = true
						}

					}

					if loopCheck {
						_, err = connections[i].Write([]byte(message))
						if err != nil {
							println("Write to server failed:", err.Error())
							os.Exit(1)
						} else {
							println("Message sent to Node Number: ", i)
						}
					}
				}
				for i := 0; i < total; i++ {
					loopCheck = false
					for j := 0; j < total; j++ {
						if excluded[i] != 1 {
							loopCheck = true
						}

					}
					if loopCheck { // disabling the non block nature
						connections[i].SetDeadline(time.Time{})
					}
				}

				response := make([]byte, 1024)
				mycount := 0

				if check == false {
					for i := 0; i > -1; i++ {
						if i == total {
							i = 0
						}
						loopCheck = false
						for j := 0; j < total; j++ {
							if excluded[i] != 1 {
								loopCheck = true
							}

						}

						if loopCheck {
							_, err = connections[i].Read(response[0:])
							if err != nil {
							} else {
								println("Message Received : ", string(response), " From Node: ", i)

								mycount++
							}
							if mycount >= (total - count) {
								i = -2
								break
							}
						}
					}

				}
			}
			// Now setting up heartbeat
			loop = 0
			if check == false { // leader sending heartbeat
				for {
					loopCheck = false
					for j := 0; j < total; j++ {
						if excluded[loop] != 1 {
							loopCheck = true
						}

					}
					if loopCheck {
						heartResponse := make([]byte, 1024)
						heartbeat := "tooktook"
						t := strconv.Itoa(mynum)
						heartbeat = heartbeat + t
						connections[loop].Write([]byte(heartbeat))
						fmt.Println("Heartbeat sent")
						_, err = connections[loop].Read(heartResponse[0:])
						println("Response Received: ", string(heartResponse), " From Node: ", loop)
						if partition {
							connections = connections[:1]
							total = 1
							partition = false
						}

					}
					loop++
					if loop == total {
						loop = 0
					}
				}
			}
			rand.Seed(time.Now().UnixNano())
			randtimer2 := rand.Intn(max-min+1) + min
			if check == true {
				for i := 0; i < randtimer2; i++ {
					loopCheck = false
					for j := 0; j < total; j++ {
						//	println("Value of loop is: ", loop)
						if excluded[loop] != 1 {
							loopCheck = true
						}
					}

					if loopCheck {
						connections[loop].SetReadDeadline(time.Now().Add(time.Millisecond))
						temp := make([]byte, 1024)
						reply = temp
						_, err = connections[loop].Read(reply[0:])
						if err != nil {
							// println("read from client failed:", err.Error())
						} else {

							println("Heartbeat Received: ", string(reply), " From Node index: ", loop)
							winner = loop
							strRep := string(reply)
							ss := strRep[8]
							winnerNum, err = strconv.Atoi(string(ss))
							//winnerNum = string(reply)
							fmt.Println("winner num is: ", winnerNum)
							i = 0
							connections[winner].Write([]byte(ack))
							if partition {
								if winnerNum != 0 {
									if mynum != 0 {
										connections = RemoveIndex(connections, winner) // you also have to remove the node connected with leader
										connections = RemoveIndex(connections, 0)
										total = total - 2 // total = total - 2 after removing second node as well
									} else {

										connections = connections[winner:winner]
										total = 1

									}
								} else {
									if mynum != 1 {
										total = total - 2
										connections = RemoveIndex(connections, winner)
										connections = RemoveIndex(connections, 1)

									} else {
										connections = connections[winner:winner]
										total = 1
									}
								}
								loop = 0
								partition = false

								break
							}

						}
					}

					loop++
					if loop == total {
						loop = 0
					}

				}
				/*	if check == false {
						connections = connections[:1]
						total = 1
					} else {
						if winnerNum != 0 {
							if mynum != 0 {
								connections = RemoveIndex(connections, winner) // you also have to remove the node connected with leader
								total = total - 1                              // total = total - 2 after removing second node as well
							} else {

								connections = connections[winner:winner]
								total = 1
							}
						} else {
							if mynum != 1 {
								total = total - 1
								connections = RemoveIndex(connections, winner)
							} else {
								connections = connections[winner:winner]
								total = 1
							}
						}
					}*/
				//	excluded[winner] = 1
				count++

			}

		}
	}
}
