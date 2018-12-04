package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

//	pgl "github.com/Fantom-foundation/pgLachesisAPI/pgLachesis"
	pgl "github.com/BrianActon/pgLachesisAPI/pgLachesis"
)


func main() {
	 

  	muxin := mux.NewRouter()
// account  	
    muxin.HandleFunc("/createAccount", createAccount).Methods("POST") 			// done     // EVM
    muxin.HandleFunc("/updateAccount", updateAccount).Methods("POST")           // Not done // EVM
    muxin.HandleFunc("/getBalance", getBalance).Methods("GET") 					// done

    // -- Handled by NodeJS team?
//    muxin.HandleFunc("/login", login).Methods("POST")
//    muxin.HandleFunc("/logout", logout).Methods("GET")  
//    muxin.HandleFunc("/forgotpw", forgotpw).Methods("POST")
//    muxin.HandleFunc("/changepw", changepw).Methods("POST")

// blocks and transactions
    muxin.HandleFunc("/writeBlock", writeBlock).Methods("POST")					// Lachesis - poset.go
    muxin.HandleFunc("/fetchBlocks", fetchBlocks).Methods("GET")                // done
    muxin.HandleFunc("/fetchSingleBlock", fetchSingleBlock).Methods("GET")      // done

    muxin.HandleFunc("/fetchTransList", fetchTransList).Methods("GET")  		// done
    muxin.HandleFunc("/fetchSingleTrans", fetchSingleTrans).Methods("GET")   	// done

// summary
    muxin.HandleFunc("/writeSummary", writeSummary).Methods("POST")				// done
    muxin.HandleFunc("/fetchSummary", fetchSummary).Methods("GET")				// done

    muxin.HandleFunc("/", nothing)                 

  	server := &http.Server{
    	Addr:    "127.0.0.1:8105",
    	Handler: muxin,
  	}

  // 3. ListenAndServe()

  	log.Fatal(server.ListenAndServe())

}



//**********************************************************************************
// Someone, somewhere isnt sending correct format
// Does it require a reply?  eg "not found"
//**********************************************************************************
func nothing(w http.ResponseWriter , r *http.Request) {
	fmt.Println("err... incorrect IP:port/folder")
}


//**********************************************************************************
// Create account - POST
//**********************************************************************************
func createAccount(w http.ResponseWriter , r *http.Request) {

	fmt.Println("createAccount")

    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)

    fmt.Println("body: ", string(body))

    var as 	pgl.AccountStruct

    err := json.Unmarshal(body, &as)
    if err != nil {
    	fmt.Println("err:", err)
    }

    fmt.Println("as: ", as.Account, as.Address)

    err = pgl.WriteAccounts([]byte(as.Account), as.Address) 
    if err != nil {
    	fmt.Println("err:", err)
    	w.WriteHeader(http.StatusInternalServerError)
    	w.Write([]byte(err.Error()))
    } else {
    	w.WriteHeader(http.StatusOK)
    	w.Write([]byte("200 OK"))
    }
}

// Add address 
//**********************************************************************************
// Update account - POST
// Update balance ?
//**********************************************************************************
func updateAccount(w http.ResponseWriter , r *http.Request) {

	fmt.Println("updateAccount")

    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)

    fmt.Println("body: ", string(body))

    var as 	pgl.AccountStruct

    err := json.Unmarshal(body, &as)
    if err != nil {
    	fmt.Println("err:", err)
    }

    fmt.Println("as: ", as.Account, as.Address)

    err = pgl.UpdateAccounts(as.Account, as.Address, as.Balance) 
    if err != nil {
    	fmt.Println("err:", err)
    	w.WriteHeader(http.StatusInternalServerError)
    	w.Write([]byte(err.Error()))
    } else {
    	w.WriteHeader(http.StatusOK)
    	w.Write([]byte("200 OK"))
    }
}


/*
for NodeJS team???

//**********************************************************************************
// login to  account  -- POST
//**********************************************************************************
func login(w http.ResponseWriter , r *http.Request) {

    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)
}


//**********************************************************************************
// logout of  account POST or GET
//**********************************************************************************
func logout(w http.ResponseWriter , r *http.Request) {

    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)
    fmt.Println(string(body))

}


//**********************************************************************************
// ChangePW POST
//**********************************************************************************
func changepw(w http.ResponseWriter , r *http.Request) {

    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)
    fmt.Println(string(body))

}


//**********************************************************************************
// ChangePW POST
//**********************************************************************************
func changepw(w http.ResponseWriter , r *http.Request) {

    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)
    fmt.Println(string(body))

}
*/


//**********************************************************************************
// Get account balance GET
//**********************************************************************************
func getBalance(w http.ResponseWriter , r *http.Request) {
	fmt.Println("pgLachesisAPI - getBalance")
	queryValues := r.URL.Query()
	if len(queryValues) < 1 {
		fmt.Println("Invalid query")
	}

	Account    :=  queryValues["account"][0]
	Address    :=  queryValues["address"][0]

	fmt.Println("pgLachesisAPI - getBalance: Account ", Account, " Address", Address)

	AccountPG, err := pgl.ReadAccountsBalance(Account, Address)

	if err != nil {
		fmt.Println("ReadAccounts unsuccessful", err)
	} 

	fmt.Println("pgLachesisAPI - getBalance:", AccountPG)
	byt, err := json.Marshal(AccountPG)

	fmt.Println("pgLachesisAPI - getBalance:", string(byt))

	w.Header().Set("Content-Type", "application/json")

	w.Write(byt)

}



//**********************************************************************************
// Write block POST
//**********************************************************************************
func writeBlock(w http.ResponseWriter , r *http.Request) {

	fmt.Println("writeBlock")

    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)

    fmt.Println("body: ", string(body))

	var bb pgl.BlockBody

    err := json.Unmarshal(body, &bb)
    if err != nil {
    	fmt.Println("err:", err)
    }

    fmt.Println("bb: ", bb.Index, bb.StateHash)

    err = pgl.WriteBlock(bb) 
    if err != nil {
    	fmt.Println("err:", err)
    	w.WriteHeader(http.StatusInternalServerError)
    	w.Write([]byte(err.Error()))
    } else {
    	w.WriteHeader(http.StatusOK)
    	w.Write([]byte("200 OK"))
    }
}




//**********************************************************************************
// Fetch list of blocks   GET
// variables currentblock=xxx & page=n
// currentblock is reference point for pagination
// page=0 for most recent (currentblock will be ignored)
// page=1 foritems news than current block
// page=2 for older items 
//**********************************************************************************
func fetchBlocks(w http.ResponseWriter , r *http.Request) {
	fmt.Println("pgLachesisAPI - fetchBlocks")
	queryValues := r.URL.Query()
	if len(queryValues) < 1 {
		fmt.Println("Invalid query")
	}

	BlockIndex    	:=  queryValues["blockindex"][0]
	Page    		:=  queryValues["page"][0]

	BlockIndexInt, err 	:= strconv.Atoi(BlockIndex)
	if err != nil {
		BlockIndexInt = 0
	}

	PageInt, err		:=  strconv.Atoi(Page)
	if err != nil {
		PageInt = 0
	}

	fmt.Println("pgLachesisAPI - fetchBlocks: BlockIndex ", BlockIndexInt, " Page", PageInt)

	BlockPG, err := pgl.ReadBlocks(BlockIndexInt, PageInt)

	if err != nil {
		fmt.Println("ReadBlocks unsuccessful", err)
	} 

	fmt.Println("pgLachesisAPI - fetchBlocks:", BlockPG)
	byt, err := json.Marshal(BlockPG)

	fmt.Println("pgLachesisAPI - fetchBlocks:", string(byt))

	w.Header().Set("Content-Type", "application/json")

	w.Write(byt)


}


//**********************************************************************************
// Fetch details of most recent block  - GET
//**********************************************************************************
func fetchSingleBlock(w http.ResponseWriter , r *http.Request) {
    
	fmt.Println("pgLachesisAPI - fetchSingleBlock")

	queryValues := r.URL.Query()
	if len(queryValues) < 1 {
		fmt.Println("Invalid query")
	}

	BlockIndex    	:=  queryValues["blockindex"][0]

	BlockIndexInt, err 	:= strconv.Atoi(BlockIndex)
	if err != nil {
		BlockIndexInt = 0
	}

	BlockPG, err := pgl.ReadBlock(BlockIndexInt)

	if err != nil {
		fmt.Println("ReadBlock unsuccessful", err)
	} 

	fmt.Println("pgLachesisAPI - fetchBlocks:", BlockPG)
	byt, err := json.Marshal(BlockPG)

	fmt.Println("pgLachesisAPI - fetchBlocks:", string(byt))

	w.Header().Set("Content-Type", "application/json")

	w.Write(byt)

}


//**********************************************************************************
// Fetch list of Transaction
// variables currenttrans=xxx & page=n
// currenttrans is reference point for pagination
// page=0 for most recent (currenttrans will be ignored)
// page=1 foritems news than current block
// page=2 for older items 
//**********************************************************************************
func fetchTransList(w http.ResponseWriter , r *http.Request) {

	fmt.Println("pgLachesisAPI - fetchTransList")
	queryValues := r.URL.Query()
	if len(queryValues) < 1 {
		fmt.Println("Invalid query")
	}

	BlockIndex    	:=  queryValues["blockindex"][0]
	TransactionID  	:=  queryValues["transactionid"][0]
	Page    		:=  queryValues["page"][0]

	BlockIndexInt, err	:=  strconv.Atoi(BlockIndex)
	if err != nil {
		BlockIndexInt = 0
	}
	PageInt, err		:=  strconv.Atoi(Page)
	if err != nil {
		PageInt = 0
	}

	fmt.Println("pgLachesisAPI - fetchTransList: TransactionID ", TransactionID, 
												" Page", PageInt, 
												" BlockIndex", BlockIndexInt)

	TransListPG, err := pgl.ReadListTransactions(BlockIndexInt, TransactionID, PageInt)

	if err != nil {
		fmt.Println("ReadListTransactions unsuccessful", err)
	} 

	fmt.Println("pgLachesisAPI - fetchTransList:", TransListPG)
	byt, err := json.Marshal(TransListPG)

	fmt.Println("pgLachesisAPI - fetchTransList:", string(byt))

	w.Header().Set("Content-Type", "application/json")

	w.Write(byt)
}


//**********************************************************************************
// Fetch details of most recent Transaction
//**********************************************************************************
func fetchSingleTrans(w http.ResponseWriter , r *http.Request) {
  
	fmt.Println("pgLachesisAPI - fetchsingleTrans")

	queryValues := r.URL.Query()
	if len(queryValues) < 1 {
		fmt.Println("Invalid query")
	}

	TransactionID    	:=  queryValues["transactionid"][0]

	TransactionIDPG, err := pgl.ReadLatestTransaction(TransactionID)

	if err != nil {
		fmt.Println("ReadTransaction unsuccessful", err)
	} 

	fmt.Println("pgLachesisAPI - fetchsingleTrans:", TransactionIDPG)
	byt, err := json.Marshal(TransactionIDPG)

	fmt.Println("pgLachesisAPI - fetchsingleTrans:", string(byt))

	w.Header().Set("Content-Type", "application/json")

	w.Write(byt)
}



//**********************************************************************************
// Fetch details of most recent Transaction
//**********************************************************************************
func writeSummary(w http.ResponseWriter , r *http.Request) {
    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)
    fmt.Println(string(body))



}

//**********************************************************************************
// Fetch details of most recent Transaction
//**********************************************************************************
func fetchSummary(w http.ResponseWriter , r *http.Request) {
 
	fmt.Println("pgLachesisAPI - fetchSummary")

	queryValues := r.URL.Query()
	if len(queryValues) < 1 {
		fmt.Println("Invalid query")
	}

	TransactionID    	:=  queryValues["transactionid"][0]

	TransactionIDPG, err := pgl.ReadLatestTransaction(TransactionID)

	if err != nil {
		fmt.Println("ReadTransaction unsuccessful", err)
	} 

	fmt.Println("pgLachesisAPI - fetchSummary:", TransactionIDPG)
	byt, err := json.Marshal(TransactionIDPG)

	fmt.Println("pgLachesisAPI - fetchSummary:", string(byt))

	w.Header().Set("Content-Type", "application/json")

	w.Write(byt)

}