package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var cred = ""

func main() {
	con, err := net.Dial("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(con)

	for {
		// Waiting for the client request
		fmt.Println("Would you like to sign up or log in? (1: Sign Up, 2: Log In): ")

		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			for clientRequest != "1" && clientRequest != "2" {
				fmt.Println("Would you like to sign up or log in? (1: Sign Up, 2: Log In): ")
				clientRequest, err = clientReader.ReadString('\n')
				clientRequest = strings.TrimSpace(clientRequest)
			}
			if clientRequest == "1" { // 1 to Sign Up and 2 log in
				cred = SignUp()
			} else if clientRequest == "2" {
				cred = login()
			}
			if _, err = con.Write([]byte(clientRequest + "," + cred + "\n")); err != nil {
				log.Printf("failed to send the client request: %v\n", err)
			}

		case io.EOF:
			log.Println("client closed the connection")
			return
		default:
			log.Printf("client error: %v\n", err)
			return
		}

		// Waiting for the server response
		serverResponse, err := serverReader.ReadString('\n')

		switch err {
		case nil:
			server_response := strings.TrimSpace(serverResponse)
			if server_response == "Logged in!" {
				fmt.Println("Would you like to [1:buy stock], [2:sell stock], [3:check balance],[4:check stock price], [5:Sign out]")
				var user_choice string
				for user_choice != "1" || user_choice != "2" || user_choice != "3" || user_choice != "4" {
					fmt.Scanln(&user_choice)
					if user_choice == "1" { // The user wants to buy a stock (conditions he has to have enough funds to buy a stock, if he doesn't have enough the purchase gets declined, flag have_stock in DB has to be 0)
						fmt.Println("Buy Stock")
						for {
							fmt.Println("Press 1 to confirm!")
							new_request, err := clientReader.ReadString('\n')
							switch err {
							case nil:
								new_request := strings.TrimSpace(new_request)
								if new_request != "1" {
									break
								}
								if _, err = con.Write([]byte(new_request + "\n")); err != nil {
									log.Printf("failed to send the client request: %v\n", err)
								}
								serverResponse, err := serverReader.ReadString('\n')

								switch err {
								case nil:
									server_response := strings.TrimSpace(serverResponse)
									fmt.Println(server_response)
								case io.EOF:
									log.Println("server closed the connection")
									return
								default:
									log.Printf("server error: %v\n", err)
									return
								}

							case io.EOF:
								log.Println("server closed the connection")
								return
							default:
								log.Printf("server error: %v\n", err)
								return
							}
							break
						}
					} else if user_choice == "2" { // The user wants to sell a stock, he can sell at anytime only if have_stock flag in the DB is set to 1)
						fmt.Println("Sell Stock")
						for {
							fmt.Println("Press 2 to confirm!")
							new_request, err := clientReader.ReadString('\n')
							switch err {
							case nil:
								new_request := strings.TrimSpace(new_request)
								if new_request != "2" {
									break
								}
								if _, err = con.Write([]byte(new_request + "\n")); err != nil {
									log.Printf("failed to send the client request: %v\n", err)
								}
								serverResponse, err := serverReader.ReadString('\n')

								switch err {
								case nil:
									server_response := strings.TrimSpace(serverResponse)
									fmt.Println(server_response)
								case io.EOF:
									log.Println("server closed the connection")
									return
								default:
									log.Printf("server error: %v\n", err)
									return
								}

							case io.EOF:
								log.Println("server closed the connection")
								return
							default:
								log.Printf("server error: %v\n", err)
								return
							}
							break
						}
					} else if user_choice == "3" {
						fmt.Println("Check Balance")
						for {
							fmt.Println("Press 3 to confirm!")
							new_request, err := clientReader.ReadString('\n')
							switch err {
							case nil:
								new_request := strings.TrimSpace(new_request)
								if new_request != "3" {
									break
								}
								if _, err = con.Write([]byte(new_request + "\n")); err != nil {
									log.Printf("failed to send the client request: %v\n", err)
								}
								serverResponse, err := serverReader.ReadString('\n')

								switch err {
								case nil:
									server_response := strings.TrimSpace(serverResponse)
									fmt.Println(server_response)

								case io.EOF:
									log.Println("server closed the connection")
									return
								default:
									log.Printf("server error: %v\n", err)
									return
								}

							case io.EOF:
								log.Println("server closed the connection")
								return
							default:
								log.Printf("server error: %v\n", err)
								return
							}
							break
						}
					} else if user_choice == "4" {
						fmt.Println("Check Stock price")
						for {
							fmt.Println("Press 4 to confirm!")
							new_request, err := clientReader.ReadString('\n')
							switch err {
							case nil:
								new_request := strings.TrimSpace(new_request)
								if new_request != "4" {
									break
								}
								if _, err = con.Write([]byte(new_request + "\n")); err != nil {
									log.Printf("failed to send the client request: %v\n", err)
								}
								serverResponse, err := serverReader.ReadString('\n')

								switch err {
								case nil:
									server_response := strings.TrimSpace(serverResponse)
									fmt.Println(server_response)

								case io.EOF:
									log.Println("server closed the connection")
									return
								default:
									log.Printf("server error: %v\n", err)
									return
								}

							case io.EOF:
								log.Println("server closed the connection")
								return
							default:
								log.Printf("server error: %v\n", err)
								return
							}
							break

						}

					} else if user_choice == "5" {
						for {
							fmt.Println("Press 5 to confirm!")
							new_request, err := clientReader.ReadString('\n')
							switch err {
							case nil:
								new_request := strings.TrimSpace(new_request)
								if new_request != "5" {
									break
								}
								if _, err = con.Write([]byte(new_request + "\n")); err != nil {
									log.Printf("failed to send the client request: %v\n", err)
								}
								serverResponse, err := serverReader.ReadString('\n')

								switch err {
								case nil:
									server_response := strings.TrimSpace(serverResponse)
									fmt.Println(server_response)

								case io.EOF:
									log.Println("server closed the connection")
									return
								default:
									log.Printf("server error: %v\n", err)
									return
								}

							case io.EOF:
								log.Println("server closed the connection")
								return
							default:
								log.Printf("server error: %v\n", err)
								return
							}
							break

						}
						break
					}
					fmt.Println("Would you like to [1:buy stock], [2:sell stock], [3:check balance],[4:check stock price], [5:Sign out]")
				}

			}
		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}
		cred = ""
	}
}

func SignUp() string { // Sign Up function
	fmt.Println(".::YOU ARE IN SIGN UP::.")

	fmt.Println("Username: ")

	var username string
	fmt.Scanln(&username)

	fmt.Println("Password (needs to start with a capital letter, you can't add commas in the password): ")
	var password string
	fmt.Scanln(&password)

	fmt.Println("Your username is: ", username)
	fmt.Println("Your password is:", password)

	up := concat(username, password)
	return up
}

func login() string { // Log in function
	fmt.Println("You are in log in function!")
	fmt.Println("Username: ")
	var username string
	fmt.Scanln(&username)

	fmt.Println("Password: ")
	var password string
	fmt.Scanln(&password)
	log := concat(username, password)
	return log
}

func concat(username string, password string) string {
	var user_pass = username + "," + password
	return user_pass
}
