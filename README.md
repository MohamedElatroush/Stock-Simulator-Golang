# Stock Simulator-Golang
 + A proof of concept project, creating a stock buy/sell simulator using Golang programming language.
 + A user can Sign Up/Log in, check balance, buy, sell a stock. (using the help of SQL DB)
 + When a user creates an account he gets awarded with a random value in his wallet ($10 - $25)
 + The user has to track the value of the stock presented on Localhost/8080 in order to buy/sell at the correct time
 + A user can't buy the stock if he/she already owns it
 + A user can't sell the stock if he/she doesn't own it
 + A client-server TCP approach have been used to send credentials and commands from the client to the backend server that deals with the DB (transactions and authorization)

 # Libraries Used
 + HTTP
 + time
 + SQL
