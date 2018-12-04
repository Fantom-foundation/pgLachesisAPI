//**********************************************************************************
//	TODO:   (Depending on what dat we are storing)
// 		- Update blocks for when additional transactions are added
//		- Update accounts with additional transactions
//		- Update transactions with new entries for a block
//		- Update accounttransactions with new transactions for an account
//
//		Delete entries?
//
//
//**********************************************************************************
//
//	TODO to cater for explorer:
//		- Summary block : Create, update, read
//		- Add "From", "To" and amount to Transaction table
//		- Add Ether transaction history  (14 days?)
//		- Add "block reward" and "mined by" to block table
//
//		-	Login?
//
//**********************************************************************************

package pgLachesis

import (
	"database/sql" 
	"fmt"
	"math/rand"
	"strconv"
//	"time"

	_ "github.com/lib/pq"
)

//  If this package is open source then these details need to be hidden from public view
const (
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  password = "P0stGr3sSU"
  dbname   = "Lachesis"
)

var	Ppgsql *sql.DB


//**********************************************************************************
//  init function
//**********************************************************************************
func init() {
	fmt.Println("init innit? ")
	db := ConnectPostgres()

	Ppgsql = db

	fmt.Println("init done", db) 
}

type AccountPG struct {
	Account 	string
	Address 	string          
	Balance 	float32
	DateTime 	string
}		

type AccountStruct struct {
	Account 		string `json:"account`
	Address 		string `json:"address`
	Balance 		float32 `json:"balance`
	AccountDateTime string `json:"accountdatetime`
}

// Account represents an Ethereum account located at a specific location defined
// by the optional URL field.
//type Account struct {
//	Address common.Address `json:"address"` // Ethereum account address derived from the key
//	URL     URL            `json:"url"`     // Optional resource locator within a backend
//}
//**********************************************************************************
//  Write Accounts
//**********************************************************************************
func WriteAccounts(account []byte, address string) error {

	fmt.Println("WriteAccounts in: ", account, address)
	var apg  AccountPG

	apg.Account = string(account)
	apg.Address = address

	q := fmt.Sprintf("INSERT INTO public.accounts(account, address, balance, account_datetime) VALUES ($1, $2, 0, NOW());")	

	_, err := Ppgsql.Exec(q, apg.Account, apg.Address)

	if err != nil {
		fmt.Println("fail to write account ", account, "error: ", err)
	}

	return err

}


//**********************************************************************************
//  Read Accounts
//**********************************************************************************
func ReadAccounts(account []byte) ([]AccountPG, error) {

	fmt.Println("ReadAccounts in: ", string(account))

	var err error
	var apg  AccountPG
	var apgs  []AccountPG

	q := ` SELECT account, address from accounts WHERE account = $1;`

	rows, err := Ppgsql.Query(q, account)

	if err != nil {
		fmt.Println("Error reading accounts : ", err)
		return apgs, err
	} 
	defer rows.Close()


	for rows.Next() {
		err = rows.Scan(&apg.Account, &apg.Address)
		if err != nil {
			fmt.Println("Error reading accounts : ", err)
		} 

		apgs = append(apgs, apg)
	}	
	
	return apgs, err
}


//**********************************************************************************
//  Read Accounts
//**********************************************************************************
func ReadAccountsBalance(account string, address string) (AccountPG, error) {

	fmt.Println("ReadAccountsBalance in: ", string(account))

	var err error
	var apg  AccountPG
	var apgs  []AccountPG

	q := fmt.Sprintf("SELECT account, address, balance from accounts " + 
					" WHERE account = '%s' and address = '%s';", account, address)

	rows, err := Ppgsql.Query(q)

	if err != nil {
		fmt.Println("Error reading accounts : ", err)
		return apg, err
	} 
	defer rows.Close()

	fmt.Println("ReadAccountsBalance rows: ", rows)
// fix data then remove the loop, or alternatively leave loop to use for fault detection.
	for rows.Next() {
		err = rows.Scan(&apg.Account, &apg.Address, &apg.Balance)
		if err != nil {
			fmt.Println("Error reading accounts : ", err)
		} 

		// !!!
		break

		apgs = append(apgs, apg)
	}	
	
	return apg, err
}

//**********************************************************************************
//  Update Accounts Balance
//**********************************************************************************
func UpdateAccounts(account string, address string, balance float32) error {

	fmt.Println("UpdateAccounts in: ", string(account))

	var err error
	var apg  AccountPG

	//  1st we need to retrieve the balance to update. If its updated to less then zero
	// then we have a problem

	q := ` SELECT account, address, balance from accounts WHERE account = $1 and address = $2;`

	row := Ppgsql.QueryRow(q, account, address)

	err = row.Scan(&apg.Account, &apg.Address, &apg.Balance)
	if err != nil {
		fmt.Println("Error reading accounts : ", err)
	} 

	if apg.Balance + balance < 0 {
		fmt.Println("cant have negative balance", )
	} else {

		q := ` UPDATE account, set balance = balance + $1 WHERE account = $2 and address = $3;`

		_, err := Ppgsql.Exec(q, account, address)

		if err != nil {
			fmt.Println("Error updating accounts with new balance : ", err)
			return err
		} 
	}
		
	return err
}



//**********************************************************************************
//  Write Account Transactions
//**********************************************************************************
func WriteAccountTrans(account []byte, transaction []byte) error {

	fmt.Println("WriteAccountTrans in: ", string(account), string(transaction))
		
	q := fmt.Sprintf("INSERT INTO public.accounttransactions(account, transaction, at_datetime) VALUES ($1, $2, NOW());")	

	_, err := Ppgsql.Exec(q, string(account), string(transaction))
	
	if err != nil {
		fmt.Println("fail to write accounttransactions ", string(account) , string(transaction), "error: ", err)
	}

	return err
}

//**********************************************************************************
//  Read Account Transactions
//**********************************************************************************
func ReadAccountTrans(account []byte) ([][]byte, error) {

	fmt.Println("ReadAccountTrans in: ", string(account))

	var trans [][]byte
	var tran  string
	var err error

	q := ` SELECT transaction from accounttransactions WHERE account = $1;`

	rows, err := Ppgsql.Query(q, string(account))

	if err != nil {
		fmt.Println("Error reading accounts : ", err)
		return trans, err
	} 
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&tran)
		if err != nil {
			fmt.Println("Error looping through rows", err)
		}
		trans = append(trans, []byte(tran))
	}	

	return trans, err

}


//	fields := map[string]interface{}{
//		"blockHash":         blockHash,
//		"blockNumber":       hexutil.Uint64(blockNumber),
//		"transactionHash":   hash,
//		"transactionIndex":  hexutil.Uint64(index),
//		"from":              from,
//		"to":                tx.To(),
//		"gasUsed":           hexutil.Uint64(receipt.GasUsed),
//		"cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed),
//		"contractAddress":   nil,
//		"logs":              receipt.Logs,
//		"logsBloom":         receipt.Bloom,
//	}

type TranFields struct {
	BlockHash 				string
	BlockNumber 			uint64
	TransactionHash 		string
	TransactionIndex        uint64
	From 					string
	To 						string	
	GasUsed                 uint64
	CumulativeGasUsed       uint64
	ContractAddress         string
	Logs                    interface{} //  receipt.Logs,
	LogsBloom          		interface{} //  receipt.Bloom,
	}

//**********************************************************************************
// Write Transactions
//**********************************************************************************
func WriteTransactions(t  TranFields) error {

	fmt.Println("WriteTransactions in: ", t.BlockHash)

	var err error

	q := `INSERT INTO transactions(
 				transaction,
 				tx_from,
 				tx_to,
				GasUsed,   
 				transaction_datetime) 
 			VALUES ($1, $2, $3, $4, $R5, $6 NOW());	`

	
	_, err = Ppgsql.Exec(q, t.TransactionHash, t.From, t.To, t.GasUsed)
	if err != nil {
		fmt.Println("fail to write transaction ", t.TransactionHash, "error: ", err)
	}

	return err
}



//**********************************************************************************
// Write BlockTransactions
//**********************************************************************************
func WriteBlockTransactions(Transactions  [][]byte) (string, error) {

	fmt.Println("WriteTransactions in: ", string(Transactions[0][:]))

	transactionblockid := rand.Int63()
	var err error


	q := `INSERT INTO blocktransactions(transaction, transactionblockid, transaction_datetime) VALUES ($1, $2, NOW());	`

	
	for i := 0; i < len(Transactions); i++ {
		
		_, err := Ppgsql.Exec(q, string(Transactions[i]), transactionblockid)
		if err != nil {
			fmt.Println("fail to write transaction ", string(Transactions[i]), "error: ", err)
			break
		}
	}

	return strconv.FormatInt(transactionblockid, 10), err
}


type TransactionDetail struct {
	Transaction 			string
	Tx_from 				string
	Tx_to 					string
	Tx_value 				float32
	TransactionBlockIndex 	int
	Transaction_DateTime 	string
}


//**********************************************************************************
// Write Transactions
//**********************************************************************************
// get single most recent transaction for block or for in general
//**********************************************************************************
func ReadLatestTransaction(transactionid string) (TransactionDetail, error) {

	fmt.Println("ReadLatestTransaction in: ")

	var pbrtran TransactionDetail
	var q string

 	if transactionid == "0" { 
 		// get the latest 		
 		q = fmt.Sprintf("SELECT transaction, tx_from, tx_to, tx_value, transactionblockid, transaction_datetime  " + 
 						" FROM transactions WHERE  transaction = ( SELECT MAX(transaction) FROM transactions);")
 	}	else {
 		// get specific
 		q = fmt.Sprintf("SELECT transaction, tx_from, tx_to, tx_value, transactionblockid, transaction_datetime  " + 
 						" FROM transactions WHERE transaction = %d ;", transactionid)
 	}

	fmt.Println("q: ", q)

	row := Ppgsql.QueryRow(q)

	err := row.Scan(&pbrtran.Transaction, 
						&pbrtran.Tx_from,
						&pbrtran.Tx_to,
						&pbrtran.Tx_value,
						&pbrtran.TransactionBlockIndex, 
						&pbrtran.Transaction_DateTime)

	if err != nil {
		fmt.Println("Error getting data from block:", err)
	} else {
		fmt.Println("Tran: ", pbrtran)
	}
			
	fmt.Println("pbrtran:", pbrtran)	

	return pbrtran, err
}


//**********************************************************************************
//  TODO - Need to understand what we tryign to achieve with Pagination
//
//	- Consider for an alternative function: 
//  func ReadListTransactions(transactionblockid string, TranStart int, TranPerPage int) ([]string, error) {
//
//**********************************************************************************
func ReadListTransactions(block int, transactionid string, page int) ([]TransactionDetail, error) {

	fmt.Println("ReadListTransactions in: ", block, transactionid, page)

	var pbrtran TransactionDetail
	var pbrtrans []TransactionDetail

	var q string

	OrderBY 	:= "DESC"
	Sign 		:= ">="
	andTran		:= ""
	andBlock	:= ""



 	if block == 0 { 	
	 		// do nothing
		if transactionid != "0" {
			andTran 	= fmt.Sprintf(" AND transaction %s '%s' ", Sign, transactionid)
			andBlock   	= fmt.Sprintf(" AND transactionblockid = %d ", block)
		}
 	} else {	
 		if page == 0 {
 			// do nothing
		} else {
	 		OrderBY = "ASC"
			Sign 	= "<="
		}
		if transactionid != "0" {
			andTran 	= fmt.Sprintf(" AND transaction %s '%s' ", Sign, transactionid)
			andBlock   	= fmt.Sprintf(" AND transactionblockid = %d ", block)
		}
 	}			

 	q = fmt.Sprintf("SELECT transaction, tx_from, tx_to, tx_value, transactionblockid, transaction_datetime  " + 
 						" FROM transactions WHERE transaction <> '' %s %s ORDER BY transactionblockid %s limit 20;", 
 							andTran, andBlock, OrderBY)


 	fmt.Println("q = ", q)

	rows, err := Ppgsql.Query(q)
	if err != nil {
		fmt.Println("Error reading transactions db: ", err)
		return pbrtrans, err
	} 

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&pbrtran.Transaction, 
						&pbrtran.Tx_from,
						&pbrtran.Tx_to,
						&pbrtran.Tx_value,
						&pbrtran.TransactionBlockIndex, 
						&pbrtran.Transaction_DateTime)
		if err != nil {
			fmt.Println("Error looping through rows of transactions:", err)
		} else {
			fmt.Println("Block: ", pbrtran)

		}
		pbrtrans = append(pbrtrans, pbrtran)
	}	
	
	fmt.Println("pbrtrans:", pbrtrans)	

	return pbrtrans, err
}

//**********************************************************************************
//  Update Transactions for a specific block
// um...?  Unless we know we are only getting a complete block
//**********************************************************************************
//func UpdateTransactions(transactionblockid string, Transactions  [][]byte) (AccountPG, error) {

//}



// StateHash as taken from block.go:
//StateHash is the hash of the current state of transactions, if you have one
//node talking to an app, and another set of nodes talking to inmem, the
//stateHash will be different
//statehash should be ignored for validator checking
//   .. therefore StateHash will be ignored for now
type BlockBody struct {
	Index         int64
	RoundReceived int64
	StateHash     []byte
	FrameHash     []byte
	Transactions  [][]byte
}

type PGBlockBody struct {
	Index         		string   
	RoundReceived 		string
	StateHash     		string
	FrameHash     		string
	TransactionsBlockID string		
	TransactionsBlockCnt int 		//   <<<<---- RoundReceived and TransactionsBlockCnt should be equal
	Block_Reward 		int
	Block_Datetime 		string

}

//**********************************************************************************
// Write Block, includes the writing of all transaction to the transaction table
//**********************************************************************************
func WriteBlock(block BlockBody) error {

	fmt.Println("WriteBlock in: [0][:]", string(block.Transactions[0][:]))
	fmt.Println("WriteBlock in: ", string(block.Transactions[:][0]))

	indexStr := strconv.FormatInt(block.Index, 10)
	rrStr := strconv.FormatInt(block.RoundReceived, 10)

	var pbblock PGBlockBody
	
	pbblock.Index 			= indexStr
	pbblock.RoundReceived 	= rrStr

	pbblock.StateHash 		= string(block.StateHash)
	pbblock.FrameHash 		= string(block.FrameHash)

	transactionBlockID, err :=  WriteBlockTransactions(block.Transactions)

	if err != nil {
		fmt.Println("Error: ", err)
	} else {

		q := `INSERT INTO blocks(
								blockIndex, 
								framehash, 
								block_miner,
								transactionblockid, 
								transactionblockcount, 
								block_datetime) 
				VALUES ($1, $2, $3, $4, $5, NOW());	`


	   // check if can use pbblock.RoundReceived or if we need the count returned from WriteTransactions
	 	_, err = Ppgsql.Exec(q, pbblock.Index, pbblock.FrameHash, "", 
	 							transactionBlockID, pbblock.RoundReceived )  

	if err != nil {
			fmt.Println("fail to write blocks : error: ", err)
		}
	}
	return err
}
			

//**********************************************************************************
// Reads th eblock table only, and a separate call must be made to retrieve all
// relevant transaction for this block
//**********************************************************************************

type BlockStruct struct {
	Index         int64
	RoundReceived int64
	FrameHash     string
	Block_Datetime string
}

func ReadBlocks(block int, page int) ( []BlockStruct, error) {

	fmt.Println("ReadBlocks in: ", block, page)

	var pbrblock BlockStruct
	var pbrblocks []BlockStruct


	OrderBY := "DESC"
	Sign 	:= ">="
	where   := ""

 	if block == 0 { 	
 		// do nothing
 	} else {	
 		if page == 0 {
 			// do nothing
		} else {
	 		OrderBY = "ASC"
			Sign 	= "<="
		}

		where = fmt.Sprintf("WHERE blockIndex %s %d", Sign, block)
 	}			

 	q := fmt.Sprintf("SELECT  blockIndex, transactionblockcount, framehash, block_datetime " + 
 						" FROM blocks %s ORDER BY blockIndex %s limit 20;", where, OrderBY)


	fmt.Println("q: ", q)

	rows, err := Ppgsql.Query(q)

	if err != nil {
		fmt.Println("Error reading blocks: ", err)
		return pbrblocks, err
	} 
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&pbrblock.Index, 
						&pbrblock.RoundReceived,
						&pbrblock.FrameHash, 
						&pbrblock.Block_Datetime)
		if err != nil {
			fmt.Println("Error looping through rows of block:", err)
		} else {
			fmt.Println("Block: ", pbrblock)

		}
		pbrblocks = append(pbrblocks, pbrblock)
	}	
	
	fmt.Println("pbrblocks:", pbrblocks)	

	return pbrblocks, err
}


//**********************************************************************************
// Retrieves details on a single block
// If blockindex supplied is 0 then returns the most recent block
// otherwise will return details for specified block
//**********************************************************************************
func ReadBlock(block int) ( BlockStruct, error) {

	fmt.Println("ReadBlocks in: ", block)

	var pbrblock BlockStruct
	var q string

 	if block == 0 { 
 		// get the latest 		
 		q = fmt.Sprintf("SELECT blockIndex, transactionblockcount, framehash, block_datetime " + 
 						" FROM blocks WHERE  blockIndex = ( SELECT MAX(blockIndex) FROM blocks);")
 	}	else {
 		q = fmt.Sprintf("SELECT  blockIndex, transactionblockcount, framehash, block_datetime " + 
 						" FROM blocks WHERE blockIndex = %d ;", block)
 	}

	fmt.Println("q: ", q)

	row := Ppgsql.QueryRow(q)

	err := row.Scan(&pbrblock.Index, 
					&pbrblock.RoundReceived,
					&pbrblock.FrameHash, 
					&pbrblock.Block_Datetime)

	if err != nil {
		fmt.Println("Error getting data from block:", err)
	} else {
		fmt.Println("Block: ", pbrblock)
	}
			
	fmt.Println("pbrblock:", pbrblock)	

	return pbrblock, err
}


//**********************************************************************************
// Update Summary fields
//**********************************************************************************
type SummaryStruct struct {
	Market_cap     		float64
	BTC_ETH_amount 		float64
	BTC_ETH_at     		float64
	BTC_ETH_movement 	float64
 	Lastblockno 		int64
 	Hashrate 			float64
 	Transactions 		float64
 	Network_difficulty 	float32
}

func UpdateSummaryBlock(bi int64, hr float64, t float64, nd float32)  error {

	fmt.Println("UpdateSummaryBlock in: ", bi)

	q := ` UPDATE summary, 
			set lastblockno = $1, 
			set hashrate = $2, 
			set transactions = $3, 
			set network_difficulty = $4, 
			set lastUpdate = NOW();`

	_, err := Ppgsql.Exec(q, bi, hr, t, nd)

	return err
}


func UpdateSummaryBTC(btc_amount float64, btc_at float64, btc_mvmnt float64)  error {

	fmt.Println("UpdateSummaryBTC in: ", btc_amount)

	q := ` UPDATE summary, 
			SET BTC_ETH_amount = $1, 
				BTC_ETH_at REAL = $2,
 				BTC_ETH_movement = $3;`

	_, err := Ppgsql.Exec(q, btc_amount, btc_at, btc_mvmnt)

	return err
}


func UpdateSummaryMarketcap(mc float64)  error {

	fmt.Println("UpdateSummaryMarketcap in: ", mc)

	q := ` UPDATE summary, set Market_cap = $1;`

	_, err := Ppgsql.Exec(q, mc)

	return err
}


//**********************************************************************************
// Fetch summary details
//**********************************************************************************

func ReadSummary() (SummaryStruct, error) {

	fmt.Println("ReadSummary in: ")

	var pbsum SummaryStruct

 	q := fmt.Sprintf("SELECT ...  FROM summary ;")
 	
	fmt.Println("q: ", q)

	row := Ppgsql.QueryRow(q)

	err := row.Scan(&pbsum.Market_cap, 
					&pbsum.BTC_ETH_amount,
					&pbsum.BTC_ETH_at, 
					&pbsum.BTC_ETH_movement, 
					&pbsum.Lastblockno,
					&pbsum.Hashrate, 
					&pbsum.Transactions, 
					&pbsum.Network_difficulty)

	if err != nil {
		fmt.Println("Error getting data from block:", err)
	} else {
		fmt.Println("Block: ", pbsum)
	}
			
	fmt.Println("pbsum:", pbsum)	

	return pbsum, err
}


//*****************************************************************************************
//	Test if tables exist, if not, then create them
//*****************************************************************************************
func testTablesExist() error {

	var triedAgain bool = false

retryAccount:	
	//  test account table exists
	err := HelloAccounts()
	if err != nil && triedAgain == false {
		fmt.Println("accounts table doesnt exist")
		fmt.Println("Creating accounts table")

		err = CreateAccounts()

		if err != nil {
			if err.Error() == `pq: relation "accounts" already exists` {
				fmt.Println("accounts table: ", err)
			} else {
				fmt.Println("Create accounts error: ", err)
			}
		} else {
			if triedAgain == false {

				triedAgain = true

				goto retryAccount

			} else {
				fmt.Println("Problem creating summary table - Alert admin ")
			}
		}
	} else {
		fmt.Println("accounts exists")
	}


	triedAgain = false

retryTransaction:	
	//  test transaction table exists
	err = HelloTransaction()
	if err != nil {
		fmt.Println("transaction table doesnt exist")
		fmt.Println("Creating transaction table")

		err = CreateTransaction()

		if err != nil {
			if err.Error() == `pq: relation "transaction" already exists` {
				fmt.Println("transaction table: ", err)
			} else {
				fmt.Println("Create transaction error: ", err)
			}
		} else {
			if triedAgain == false {

				triedAgain = true

				goto retryTransaction

			} else {
				fmt.Println("Problem creating transaction table - Alert admin ")
			}
		}
	} else {
		fmt.Println("transaction exists")
	}


	triedAgain = false

retryBlockTransaction:	
	//  test blocktransaction table exists
	err = HelloBlockTransaction()
	if err != nil {
		fmt.Println("blocktransaction table doesnt exist")
		fmt.Println("Creating blocktransaction table")

		err = CreateBlockTransaction()

		if err != nil {
			if err.Error() == `pq: relation "blocktransaction" already exists` {
				fmt.Println("blocktransaction table: ", err)
			} else {
				fmt.Println("Create blocktransaction error: ", err)
			}
		} else {
			if triedAgain == false {

				triedAgain = true

				goto retryBlockTransaction

			} else {
				fmt.Println("Problem creating blocktransaction table - Alert admin ")
			}
		}
	} else {
		fmt.Println("transaction exists")
	}


	triedAgain = false

retryactran:	
	//  test transaction table exists
	err = HelloAccountTrans()
	if err != nil {
		fmt.Println("accounttrans table doesnt exist")
		fmt.Println("Creating accounttrans table")

		err = CreateAccountTrans()

		if err != nil {
			if err.Error() == `pq: relation "accounttrans" already exists` {
				fmt.Println("accounttrans table: ", err)
			} else {
				fmt.Println("Create accounttrans error: ", err)
			}
		} else {
			if triedAgain == false {

				triedAgain = true

				goto retryactran

			} else {
				fmt.Println("Problem creating summary table - Alert admin ")
			}
		}
	} else {
		fmt.Println("accounttrans table exists")
	}


	triedAgain = false

retryBlock:	
	//  test transaction table exists
	err = HelloBlock()
	if err != nil {
		fmt.Println("blocks table doesnt exist")
		fmt.Println("Creating blocks table")

		err = CreateBlock()

		if err != nil {
			if err.Error() == `pq: relation "blocks" already exists` {
				fmt.Println("blocks table: ", err)
			} else {
				fmt.Println("Create blocks error: ", err)
			}
		} else {
			if triedAgain == false {

				triedAgain = true

				goto retryBlock

			} else {
				fmt.Println("Problem creating summary table - Alert admin ")
			}
		}
	} else {
		fmt.Println("blocks table exists")
	}


	triedAgain = false

retrySummary:	
	//  test transaction table exists
	err = HelloSummary()
	if err != nil {
		fmt.Println("blocks table doesnt exist")
		fmt.Println("Creating blocks table")

		err = CreateSummary()

		if err != nil {
			if err.Error() == `pq: relation "summary" already exists` {
				fmt.Println("summary table: ", err)
			} else {
				fmt.Println("Create summary error: ", err)
			}
		} else {
			if triedAgain == false {

				triedAgain = true

				goto retrySummary

			} else {
				fmt.Println("Problem creating summary table - Alert admin ")
			}
		}
	} else {
		fmt.Println("summary table exists")
	}


	return nil
}



//*****************************************************************************************
//	Hello? can we read from these tables
//*****************************************************************************************

func HelloAccounts() error {

	_, err := Ppgsql.Query("SELECT account from accounts LIMIT 1")

	return err
}

func HelloTransaction() error {
	_, err := Ppgsql.Query("SELECT transaction from transactions LIMIT 1")

	return err
	
}

func HelloBlockTransaction() error {
	_, err := Ppgsql.Query("SELECT transaction from blocktransactions LIMIT 1")

	return err
	
}

func HelloAccountTrans() error {
	_, err := Ppgsql.Query("SELECT account from accounttransactions LIMIT 1")

	return err
	
}

func HelloBlock() error {
	_, err := Ppgsql.Query("SELECT blockIndex from blocks LIMIT 1")

	return err
	
}

func HelloSummary() error {
	_, err := Ppgsql.Query("SELECT market_cap from summary LIMIT 1")

	return err
	
}


//*****************************************************************************************
//	Lets create our tables!!!
//*****************************************************************************************
//*****************************************************************************************
//	Lets create our DB????
//*****************************************************************************************
func CreateLachesisDB() error {
	// can we do this with PostGres????
	return nil
}

//*****************************************************************************************
//	Lets create our tables!!!
//*****************************************************************************************
func CreateAccounts() error {

	q := `	CREATE TABLE accounts (
 				account VARCHAR (70) UNIQUE,
 				address VARCHAR (70),
 				balance REAL,
 				account_datetime TIMESTAMP
			);`
	
	_, err := Ppgsql.Exec(q)
	
	return err
}

//*****************************************************************************************
//	Lets create our tables!!!
//*****************************************************************************************
func CreateTransaction() error {
	
	q := `	CREATE TABLE transactions (
 				transaction 		VARCHAR (70) UNIQUE,
 				tx_from 			VARCHAR (70),
 				tx_to 				VARCHAR (70),
 				tx_value 			REAL,
				Gas               	BIGINT,
				GasUsed           	BIGINT,  
				GasPrice          	BIGINT,   
 				transaction_datetime TIMESTAMP
			);`

  
	
	_, err := Ppgsql.Exec(q)
	
	return err
}

//*****************************************************************************************
//	Lets create our tables!!!
//*****************************************************************************************
func CreateBlockTransaction() error {
	
	q := `	CREATE TABLE blocktransactions (
 				transaction VARCHAR (70) UNIQUE,
 				transactionblockid VARCHAR (70),
 				transaction_datetime TIMESTAMP
			);`
  
	_, err := Ppgsql.Exec(q)
	
	return err
}

//*****************************************************************************************
//	Lets create our tables!!!
//*****************************************************************************************
func CreateAccountTrans() error {
	
	q := `	CREATE TABLE accounttransactions (
 				account VARCHAR (70),
 				transaction VARCHAR (70),
 				at_DateTime TIMESTAMP
			);`	

	_, err := Ppgsql.Exec(q)
	
	return err
}

//*****************************************************************************************
//	Lets create our tables!!!
//*****************************************************************************************
func CreateBlock() error {
	
	q := `	CREATE TABLE blocks (
 				blockIndex INTEGER UNIQUE,
 				framehash VARCHAR (70),
 				block_miner VARCHAR (70),
 				transactionblockid VARCHAR (70),
 				transactionblockcount INTEGER,
 				block_reward REAL,
 				block_datetime TIMESTAMP 
			);`
	
	_, err := Ppgsql.Exec(q)
	
	return err
}


//*****************************************************************************************
//	Lets create our tables!!!
//*****************************************************************************************
func CreateSummary() error {
	
	q := `	CREATE TABLE summary (
 				market_cap REAL,
 				BTC_ETH_amount REAL,
 				BTC_ETH_at REAL,
 				BTC_ETH_movement REAL,
 				lastblockno INTEGER,
 				hashrate REAL,
 				transactions REAL,
 				network_difficulty REAL,
 				lastUpdate TIMESTAMP
			);`
	
	_, err := Ppgsql.Exec(q)
	
	return err
}


//*****************************************************************************************
//	Drop tables
//*****************************************************************************************

//*****************************************************************************************
//	Drop all tables!!!
//*****************************************************************************************
func DropAllTables() error {
	err := DropAccounts()
	fmt.Println("err", err)
	err = DropTransaction()
	fmt.Println("err", err)
	err = DropBlockTransaction()
	fmt.Println("err", err)
	err = DropAccountTrans()
	fmt.Println("err", err)
	err = DropBlock()
	fmt.Println("err", err)
	err = DropSummary()
	fmt.Println("err", err)

	return err
}


//*****************************************************************************************
//	Drop accounts table
//*****************************************************************************************
func DropAccounts() error {

	fmt.Println("DropAccounts")

	_, err := Ppgsql.Exec("DROP TABLE accounts")

	return err
}

//*****************************************************************************************
//	Drop transactions table
//*****************************************************************************************
func DropBlockTransaction() error {
	_, err := Ppgsql.Exec("DROP TABLE transactions")

	return err
	
}
//*****************************************************************************************
//	Drop transactions table
//*****************************************************************************************
func DropTransaction() error {
	_, err := Ppgsql.Exec("DROP TABLE transactions")

	return err
	
}

//*****************************************************************************************
//	Drop accounttrans table
//*****************************************************************************************
func DropAccountTrans() error {
	_, err := Ppgsql.Exec("DROP TABLE accounttransactions")

	return err
	
}

//*****************************************************************************************
//	Drop blocks table
//*****************************************************************************************
func DropBlock() error {
	_, err := Ppgsql.Exec("DROP TABLE blocks")

	return err
	
}

//*****************************************************************************************
//	Drop summary table
//*****************************************************************************************
func DropSummary() error {
	_, err := Ppgsql.Exec("DROP TABLE summary")

	return err
	
}



//*****************************************************************************************
//	Connect to PostGtres instance
//*****************************************************************************************
func ConnectPostgres() *sql.DB {

	fmt.Println("ConnectPostgres innit?")
  	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    						"password=%s dbname=%s sslmode=disable",
    						host, port, user, password, dbname)

  	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Excuse me kind person, do you have postgres loaded..?")
  		panic(err)
	} 

	Ppgsql = db
	
	err = db.Ping()
	if err != nil {
		fmt.Println("db.Ping unsuccessful", err)
  		panic(err)
	} else {
		fmt.Println("db.Ping successful", err)
	}

	err = Ppgsql.Ping()
	if err != nil {
		fmt.Println("Ppgsql.Ping unsuccessful")
  		panic(err)
	} else {
		fmt.Println("Ppgsql.Ping successful")
	}

	// Create tables if don't exist  -->  TODO: Chat to Andre to see if needed    
	err = testTablesExist()
	if err != nil {
		fmt.Println("Problem with accessing postgres tables")
		fmt.Println("Do you have postgres loaded?")
	}

	fmt.Println("ConnectPostgres done")

	return db
}
