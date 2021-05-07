package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name string `json:"username"`
}

var x = 0

var c_price = 0

func Stock_price() {

	min := 10
	max := 100

	for {
		rand.Seed(time.Now().UnixNano())
		c_price = rand.Intn(max-min) + min
		time.Sleep(5 * time.Second)
	}
}

var temp = template.Must(template.ParseFiles("index.html")) // Get the HTML file

func handler(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println("Error = ", err)
	}
	err = templ.Execute(w, c_price)
	if err != nil {
		fmt.Println("Error = ", err)
	}

}
func main() {

	go Stock_price()
	go server()
	fmt.Println("Server is running!")
	http.HandleFunc("/", handler) // http://127.0.0.1:8080/Go
	http.ListenAndServe(":8080", nil)

}

func server() {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// If you want, you can increment a counter here and inject to handleClientRequest below as client identifier
		go handleClientRequest(con)
	}
}

func handleClientRequest(con net.Conn) {

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/user")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	fmt.Println("Successfully connected to Database")

	defer con.Close()

	clientReader := bufio.NewReader(con)

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')
		x = 0
		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			fmt.Println("Recieved:", clientRequest)
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			} else {
				split := strings.Split(clientRequest, ",")

				if split[0] == "1" {
					fmt.Println("Signing up!")
					username := split[1]
					password := split[2]

					fmt.Println("Your username registerd is:", username)
					fmt.Println("Your password is:", password)

					random := random_wallet()
					fmt.Println("You have $", random, "in your wallet")
					/////// Adding new user to database ////////
					query := "INSERT INTO user_info VALUES" + "('" + username + "'" + ",'" + password + "'," + "'" + strconv.Itoa(random) + "'," + "'0'" + ");"
					// fmt.Println(query)
					insert, err := db.Query(query)

					if err != nil {
						panic(err.Error())
					}
					defer insert.Close()
					////////////////////////////////////////////

					if _, err = con.Write([]byte("Welcome, You've succefully signed up to concepts stock simulator!\n")); err != nil {
						log.Printf("failed to respond to client: %v\n", err)
					}

				}
				if split[0] == "2" {
					login_(split[1], split[2], clientReader, con, db)
					x = -1
				}

				if x == -1 {
					fmt.Println("Logging out")
					if _, err = con.Write([]byte("Server: Logged out")); err != nil {
						log.Printf("failed to respond to client: %v\n", err)
					}
					break
				}

			}

		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}

		// Responding to the client request
		if _, err = con.Write([]byte("\n")); err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}

	}
}

func random_wallet() int { // A function that returns a random number between $5 and $12
	min := 10
	max := 25
	rand.Seed(time.Now().UnixNano())
	amt := rand.Intn(max-min) + min
	return amt
}

func login_(split_1 string, split_2 string, clientReader *bufio.Reader, con net.Conn, db *sql.DB) {
	fmt.Println("Logging in!")
	var user string
	username := split_1
	password := split_2
	results, err := db.Query("SELECT username FROM user_info WHERE username='" + username + "' AND password='" + password + "'")
	for results.Next() {

		err = results.Scan(&user)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("\n Hello, " + user)
	}

	if len(user) <= 0 {
		fmt.Println("\nFaied to log in! (username or password is incorrect")
		return

	} else {
		fmt.Println("You've successfully logged in!")
		if _, err = con.Write([]byte("Logged in!\n")); err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}

	}
	for {
		new_request, err := clientReader.ReadString('\n')
		switch err {
		case nil:
			new_request := strings.TrimSpace(new_request)
			fmt.Println("Recieved:", new_request)
			if new_request == "1" {
				fmt.Println("Buy Stock")
				// Check user balance
				var balance string
				balance_query := "SELECT amount FROM user_info WHERE username='" + username + "';"
				err = db.QueryRow(balance_query).Scan(&balance)
				if err != nil {
					if err == sql.ErrNoRows {
						// there were no rows, but otherwise no error occurred
					} else {
						log.Fatal(err)
					}
				}
				bal, err := strconv.Atoi(balance)
				if err != nil {
					fmt.Println(err)
				}
				// check stock price
				current_stock := c_price
				// if stock price > user balance, then he can't buy
				if current_stock <= bal { // if stock price is less than or equal balance
					var have_stock string
					have_stock = "SELECT have_stock FROM user_info WHERE username='" + username + "';"

					err = db.QueryRow(have_stock).Scan(&have_stock)
					if err != nil {
						if err == sql.ErrNoRows {
							// there were no rows, but otherwise no error occurred
						} else {
							log.Fatal(err)
						}
					}

					if have_stock == "0" {
						old_bal := bal
						bal = bal - current_stock

						fmt.Println(bal)
						update_query := "UPDATE user_info SET amount='" + strconv.Itoa(bal) + "' WHERE username='" + username + "';"
						update, err := db.Query(update_query)
						if err != nil {
							panic(err.Error())
						}
						defer update.Close()
						update_query = "UPDATE user_info SET have_stock='1' WHERE username='" + username + "';"
						update, err = db.Query(update_query)
						if err != nil {
							panic(err.Error())
						}
						defer update.Close()
						if _, err = con.Write([]byte("Balance before: $" + strconv.Itoa(old_bal) + " - Stock bought for: $" + strconv.Itoa(current_stock) + " - Current balance: $" + strconv.Itoa(bal) + " ::You now own concepts stock:: \n")); err != nil {
							log.Printf("failed to respond to client: %v\n", err)
						}
						if _, err = con.Write([]byte("")); err != nil {
							log.Printf("failed to respond to client: %v\n", err)
						}
					} else {
						if _, err = con.Write([]byte("Server: You already own the stock!!\n")); err != nil {
							log.Printf("failed to respond to client: %v\n", err)
						}
						break
					}
				} else {
					if _, err = con.Write([]byte("Server: You don't have sufficient funds\n")); err != nil {
						log.Printf("failed to respond to client: %v\n", err)
					}
					break
				}
				// if the user already owns the stock, then he can't buy

			} else if new_request == "2" {
				fmt.Println("Sell stock")

				var have_stock string
				have_stock = "SELECT have_stock FROM user_info WHERE username='" + username + "';"

				err = db.QueryRow(have_stock).Scan(&have_stock)
				if err != nil {
					if err == sql.ErrNoRows {
						// there were no rows, but otherwise no error occurred
					} else {
						log.Fatal(err)
					}
				}

				if have_stock == "0" {
					if _, err = con.Write([]byte("Server: You don't own the stock\n")); err != nil {
						log.Printf("failed to respond to client: %v\n", err)
					}
					break
				} else {
					var current_balance string
					current_balance = "SELECT amount FROM user_info WHERE username='" + username + "';"
					err = db.QueryRow(current_balance).Scan(&current_balance)
					if err != nil {
						if err == sql.ErrNoRows {
							// there were no rows, but otherwise no error occurred
						} else {
							log.Fatal(err)
						}
					}
					current_price := c_price
					fmt.Println("Current stock price: ", current_price)
					bal, err := strconv.Atoi(current_balance)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("You currently have: $", bal)
					new_bal := bal + current_price
					if err != nil {
						fmt.Println(err)
					}
					// UPDATE Customers SET ContactName='Juan'WHERE Country='Mexico';
					update_query := "UPDATE user_info SET amount='" + strconv.Itoa(new_bal) + "' WHERE username='" + username + "';"
					update, err := db.Query(update_query)
					if err != nil {
						panic(err.Error())
					}
					defer update.Close()

					update2 := "UPDATE user_info SET have_stock='0' WHERE username='" + username + "';"
					update, err = db.Query(update2)
					if err != nil {
						panic(err.Error())
					}
					defer update.Close()

					if _, err = con.Write([]byte("Your balance before: $" + strconv.Itoa(bal) + " - Current Stock price: $" + strconv.Itoa(current_price) + " - You now have: $" + strconv.Itoa(new_bal) + "\n")); err != nil {
						log.Printf("failed to respond to client: %v\n", err)
					}
				}

			} else if new_request == "3" {
				var balance string
				balance_query := "SELECT amount FROM user_info WHERE username='" + username + "';"
				err = db.QueryRow(balance_query).Scan(&balance)
				if err != nil {
					if err == sql.ErrNoRows {
						// there were no rows, but otherwise no error occurred
					} else {
						log.Fatal(err)
					}
				}
				if _, err = con.Write([]byte("Current balance: $" + balance + "\n")); err != nil {
					log.Printf("failed to respond to client: %v\n", err)
				}
			} else if new_request == "4" {
				if _, err = con.Write([]byte("Current price: $" + strconv.Itoa(c_price) + "\n")); err != nil {
					log.Printf("failed to respond to client: %v\n", err)
				}
			} else if new_request == "5" {
				return
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}

	}
}
