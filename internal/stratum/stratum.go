package stratum

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"gitlab.com/afa7789/luxor_challenge/domain"
	"gitlab.com/afa7789/luxor_challenge/internal/middleware"
	"gitlab.com/afa7789/luxor_challenge/internal/postgres"
)

// methods names here
const (
	authorize        = "mining.authorize"
	subscribe        = "mining.subscribe"
	newworker        = "mining.new_worker"
	extranonce2_size = 4 // constant of length of extranonce that the miner will use
)

type Manager struct {
	mid *middleware.Middleware
	db  *postgres.PGClient
}

type Request struct {
	db   *postgres.PGClient
	mid  *middleware.Middleware
	conn net.Conn               // net connection
	ser  domain.StratumEntryRPC //  entry to stratum
}

// creates a Stratum manager to use receiver function
func NewManager(mid *middleware.Middleware, db *postgres.PGClient) *Manager {
	return &Manager{
		db:  db,
		mid: mid,
	}
}

// return a request that holds the functions for each method request from stratum
func newRequest(db *postgres.PGClient, mid *middleware.Middleware, conn net.Conn, ser domain.StratumEntryRPC) *Request {
	return &Request{
		db:   db,
		mid:  mid,
		conn: conn,
		ser:  ser,
	}
}

// handles the connection by scanning the input and passing it forward
func (s *Manager) HandleConn(conn net.Conn) {
	defer conn.Close()
	// scanner to read buffer
	scn := bufio.NewScanner(conn)
	for scn.Scan() {

		// scan bytes
		bytes := scn.Bytes()
		var ser domain.StratumEntryRPC
		err := json.Unmarshal(bytes, &ser)

		if err != nil {
			log.Printf("error : %s ", err.Error())
		}

		// set a new request
		req := newRequest(s.db, s.mid, conn, ser)
		// methodCaller will check the method receive and stipulate the correct answer inside
		// if any error happens will be using the default Stratum ERror response
		if err := req.methodCaller(); err != nil {
			req.handleError(err)
		}
	}
	if scn.Err() != nil {
		log.Print(scn.Err())
	}
}

//handleError response
func (r *Request) handleError(err *domain.StratumError) {
	// example :  {"id":3,"result":false,"error":[20,"invalid job",null]}
	resp := fmt.Sprintf(`{"id":%d,"result":false,"error":[%d,"%s",null]`, r.ser.ID, err.ID, err.ErrorMessage)
	// write at buffer
	r.conn.Write([]byte(resp + "\n"))
	log.Printf(err.ErrorMessage)
	// add in the DB that it was a failure
	errorr := r.db.InsertEntry(r.ser, false, r.conn.LocalAddr().String())
	if errorr != nil {
		log.Printf(errorr.Error())
	}
}

//methodCaller defines the function to be called and answer the connection
func (r *Request) methodCaller() *domain.StratumError {
	response := ""
	var err *domain.StratumError
	switch r.ser.Method {
	case newworker:
		response, err = r.newworker()
		if err != nil {
			return err
		}
	case authorize:
		// {"params": ["slush.miner1", "password"], "id": 1, "method": "mining.authorize"}
		response, err = r.authorize()
		// this could have been done outside switch to use the middleware independently of the case
		// but since we do not have that many cases I choose to stick with this to progress on the coding
		// just the subscribe
		if err != nil {
			return err
		}
	case subscribe:
		// {"id": 1, "method": "mining.subscribe", "params": []}
		response, err = r.subscribe()
		if err != nil {
			return err
		}
	default:
		return &domain.StratumError{
			ID:           20,
			ErrorMessage: "Other/Unknown, method not known",
		}
	}

	// Answer in the buffer
	r.conn.Write([]byte(response + "\n"))

	// add in the DB that it was a success
	errorr := r.db.InsertEntry(r.ser, true, r.conn.LocalAddr().String())
	if errorr != nil {
		log.Printf(errorr.Error())
	}
	return nil
}

// authorize request authorization to the middleware
func (r *Request) authorize() (string, *domain.StratumError) {
	err := r.mid.Authorize(r.conn, r.ser.ID)
	if err != nil {
		log.Printf("error in authorizing, id - %d : %s", r.ser.ID, err)
		return "", &domain.StratumError{
			ID:           24,
			ErrorMessage: "Other/Unknown",
			Traceback:    nil,
		}
	}
	return fmt.Sprintf(`{"id":%d,"result":true,err:null}`, r.ser.ID), nil
}

// authorize request authorization to the middleware
func (r *Request) newworker() (string, *domain.StratumError) {
	i, err := r.db.InsertMiner()
	if err != nil {
		log.Printf("error in creating a new worker : %s", err)
		return "", &domain.StratumError{
			ID:           20,
			ErrorMessage: "Other/Unknown",
			Traceback:    nil,
		}
	}
	return fmt.Sprintf(`{"id":%d,"result":true,err:null}`, *i), nil
}

// subscribe function does the subscribtion of the miner to request a "mining job"
func (r *Request) subscribe() (string, *domain.StratumError) {
	if err := r.isAuthorized(); err != nil {
		return "", err
	}

	// to make a hex encoded that's unique I'll use the timestamp and the id.
	src := []byte(fmt.Sprintf("%d_%d",
		r.ser.ID,
		time.Now().Unix(),
	))

	// hex encoding, a package taht does the hexadecimal encoding.
	hexed := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(hexed, src)

	// first_tuple , we should send here the subscribtion id
	t1, err := r.insertSubscription("mining.set_difficulty")
	if err != nil {
		return "", err
	}

	// second_tuple , we should send here the subscribtion id
	t2, err := r.insertSubscription("mining.notify")
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(`[[%s,%s],"%s",%d`, t1, t2, string(hexed), extranonce2_size)

	response := fmt.Sprintf(`{"id":%d,"result":%s,err:null`, r.ser.ID, result)
	return response, nil
}

// insert subscrition at the db
func (r *Request) insertSubscription(not_name string) (string, *domain.StratumError) {
	// creates the uuid that's a unique 16 byte hex number
	u := uuid.New()
	// convert it to byte[]
	bytea := u[:]

	// insert it to the DB
	err := r.db.InsertSubscription(
		context.Background(),
		not_name,
		r.ser.ID,
		bytea,
	)
	//error handling and log
	if err != nil {
		log.Printf("error inserting subscription at subscribe method:%s", err)
		return "", &domain.StratumError{
			ID:           20,
			ErrorMessage: "internal server error",
			Traceback:    nil,
		}
	}

	// mount the response to be used in the conn
	string_response := fmt.Sprintf(`["%s", "%s"]`, not_name, hex.EncodeToString(bytea))

	return string_response, nil
}

// calls the middlware
func (r *Request) isAuthorized() *domain.StratumError {
	// check with middleware if it is authorized
	if r.mid.IsAuthorized(r.ser.ID) {
		return nil
	}
	return &domain.StratumError{
		ID:           24,
		ErrorMessage: "Unauthorized worker",
		Traceback:    nil,
	}
}
