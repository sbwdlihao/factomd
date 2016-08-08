// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wsapi

import (
	"encoding/hex"
	//"fmt"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/entryBlock"
	"github.com/FactomProject/factomd/common/entryCreditBlock"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

func HandleV2FactoidACK(state interfaces.IState, params interface{}) (interface{}, *primitives.JSONError) {
	ackReq := new(AckRequest)
	err := MapToObject(params, ackReq)
	if err != nil {
		return nil, NewInvalidParamsError()
	}

	if ackReq.TxID == "" && ackReq.FullTransaction == "" {
		return nil, NewInvalidParamsError()
	}

	txid := ackReq.TxID

	if txid == "" {
		b, err := hex.DecodeString(ackReq.FullTransaction)
		if err != nil {
			return nil, NewUnableToDecodeTransactionError()
		}
		tx := new(factoid.Transaction)
		err = tx.UnmarshalBinary(b)
		if err != nil {
			return nil, NewUnableToDecodeTransactionError()
		}
		txid = tx.GetSigHash().String()
	}

	txhash, err := primitives.NewShaHashFromStr(txid)
	if err != nil {
		return nil, NewInvalidParamsError()
	}

	status, h, txTime, blockTime, err := state.GetACKStatus(txhash)
	if err != nil {
		return nil, NewInternalError()
	}

	answer := new(FactoidTxStatus)
	answer.TxID = h.String()

	if txTime != nil {
		answer.TransactionDate = txTime.GetTimeMilli()
		if txTime.GetTimeMilli() > 0 {
			answer.TransactionDateString = txTime.String()
		}
	}
	if blockTime != nil {
		answer.BlockDate = blockTime.GetTimeMilli()
		if blockTime.GetTimeMilli() > 0 {
			answer.BlockDateString = blockTime.String()
		}
	}

	switch status {
	case constants.AckStatusInvalid:
		answer.Status = AckStatusInvalid
		break
	case constants.AckStatusUnknown:
		answer.Status = AckStatusUnknown
		break
	case constants.AckStatusNotConfirmed:
		answer.Status = AckStatusNotConfirmed
		break
	case constants.AckStatusACK:
		answer.Status = AckStatusACK
		break
	case constants.AckStatus1Minute:
		answer.Status = AckStatus1Minute
		break
	case constants.AckStatusDBlockConfirmed:
		answer.Status = AckStatusDBlockConfirmed
		break
	default:
		return nil, NewInternalError()
		break
	}

	return answer, nil
}

func HandleV2EntryACK(state interfaces.IState, params interface{}) (interface{}, *primitives.JSONError) {
	ackReq := new(AckRequest)
	err := MapToObject(params, ackReq)
	if err != nil {
		return nil, NewInvalidParamsError()
	}

	if ackReq.TxID == "" && ackReq.FullTransaction == "" {
		return nil, NewInvalidParamsError()
	}

	eTxID := ""
	ecTxID := ""

	if ackReq.TxID == "" {
		eTxID, ecTxID = DecodeTransactionToHashes(ackReq.FullTransaction)
		if ecTxID == "" && eTxID == "" {
			return nil, NewUnableToDecodeTransactionError()
		}
	}

	//TODO: fetch entries, ec TXs from state as well

	//We didn't receive a full transaction, but a transaction hash
	//We have to figure out which transaction hash we got
	if ecTxID == "" && eTxID == "" {
		h, err := primitives.NewShaHashFromStr(ackReq.TxID)
		if err != nil {
			return nil, NewInvalidParamsError()
		}
		entry, err := state.FetchEntryByHash(h)
		if err != nil {
			return nil, NewInternalError()
		}
		if entry != nil {
			eTxID = ackReq.TxID
		} else {
			ec, err := state.FetchECTransactionByHash(h)
			if err != nil {
				return nil, NewInternalError()
			}

			//fmt.Printf("ec - %v\n", ec)
			if ec != nil {
				ecTxID = ackReq.TxID
				eTxID = ec.GetEntryHash().String()
			}
		}
	}

	answer := new(EntryStatus)
	answer.CommitTxID = ecTxID
	answer.EntryHash = eTxID

	if answer.CommitTxID == "" && answer.EntryHash == "" {
		//We know nothing about the transaction, so we return unknown status
		answer.CommitData.Status = AckStatusUnknown
		answer.EntryData.Status = AckStatusUnknown
		return answer, nil
	}

	//Fetching the second part of the transaction pair
	if answer.EntryHash == "" {
		h, err := primitives.NewShaHashFromStr(answer.EntryHash)
		if err != nil {
			return nil, NewInvalidParamsError()
		}
		ec, err := state.FetchECTransactionByHash(h)
		if err != nil {
			return nil, NewInternalError()
		}
		if ec != nil {
			answer.EntryHash = ec.GetEntryHash().String()
		}
	}

	if answer.CommitTxID == "" {
		h, err := primitives.NewShaHashFromStr(answer.EntryHash)
		if err != nil {
			return nil, NewInvalidParamsError()
		}
		ec, err := state.FetchPaidFor(h)
		if err != nil {
			return nil, NewInternalError()
		}
		if ec != nil {
			answer.CommitTxID = ec.String()
		}
	}

	//Fetching statuses
	if answer.CommitTxID == "" {
		answer.CommitData.Status = AckStatusUnknown
	} else {
		h, err := primitives.NewShaHashFromStr(answer.CommitTxID)
		if err != nil {
			return nil, NewInvalidParamsError()
		}

		status, txid, txTime, blockTime, err := state.GetACKStatus(h)
		if err != nil {
			return nil, NewInternalError()
		}

		answer.CommitTxID = txid.String()

		if txTime != nil {
			answer.CommitData.TransactionDate = txTime.GetTimeMilli()
			if txTime.GetTimeMilli() > 0 {
				answer.CommitData.TransactionDateString = txTime.String()
			}
		}
		if blockTime != nil {
			answer.CommitData.BlockDate = blockTime.GetTimeMilli()
			if blockTime.GetTimeMilli() > 0 {
				answer.CommitData.BlockDateString = blockTime.String()
			}
		}

		switch status {
		case constants.AckStatusInvalid:
			answer.CommitData.Status = AckStatusInvalid
			break
		case constants.AckStatusUnknown:
			answer.CommitData.Status = AckStatusUnknown
			break
		case constants.AckStatusNotConfirmed:
			answer.CommitData.Status = AckStatusNotConfirmed
			break
		case constants.AckStatusACK:
			answer.CommitData.Status = AckStatusACK
			break
		case constants.AckStatus1Minute:
			answer.CommitData.Status = AckStatus1Minute
			break
		case constants.AckStatusDBlockConfirmed:
			answer.CommitData.Status = AckStatusDBlockConfirmed
			break
		default:
			return nil, NewInternalError()
			break
		}
	}

	if answer.EntryHash == "" {
		answer.EntryData.Status = AckStatusUnknown
	} else {
		h, err := primitives.NewShaHashFromStr(answer.EntryHash)
		if err != nil {
			return nil, NewInvalidParamsError()
		}

		status, txid, txTime, blockTime, err := state.GetACKStatus(h)
		if err != nil {
			return nil, NewInternalError()
		}

		answer.EntryHash = txid.String()

		if txTime != nil {
			answer.EntryData.TransactionDate = txTime.GetTimeMilli()
			if txTime.GetTimeMilli() > 0 {
				answer.EntryData.TransactionDateString = txTime.String()
			}
		}
		if blockTime != nil {
			answer.EntryData.BlockDate = blockTime.GetTimeMilli()
			if blockTime.GetTimeMilli() > 0 {
				answer.EntryData.BlockDateString = blockTime.String()
			}
		}

		switch status {
		case constants.AckStatusInvalid:
			answer.EntryData.Status = AckStatusInvalid
			break
		case constants.AckStatusUnknown:
			answer.EntryData.Status = AckStatusUnknown
			break
		case constants.AckStatusNotConfirmed:
			answer.EntryData.Status = AckStatusNotConfirmed
			break
		case constants.AckStatusACK:
			answer.EntryData.Status = AckStatusACK
			break
		case constants.AckStatus1Minute:
			answer.EntryData.Status = AckStatus1Minute
			break
		case constants.AckStatusDBlockConfirmed:
			answer.EntryData.Status = AckStatusDBlockConfirmed
			break
		default:
			return nil, NewInternalError()
			break
		}
	}

	return answer, nil
}

func DecodeTransactionToHashes(fullTransaction string) (eTxID string, ecTxID string) {
	//fmt.Printf("DecodeTransactionToHashes - %v\n", fullTransaction)
	b, err := hex.DecodeString(fullTransaction)
	if err != nil {
		return
	}

	cc := new(entryCreditBlock.CommitChain)
	rest, err := cc.UnmarshalBinaryData(b)
	if err != nil || len(rest) > 0 {
		//fmt.Printf("err - %v\n", err)
		ec := new(entryCreditBlock.CommitEntry)
		rest, err = ec.UnmarshalBinaryData(b)
		if err != nil || len(rest) > 0 {
			//fmt.Printf("err - %v\n", err)
			e := new(entryBlock.Entry)
			rest, err = e.UnmarshalBinaryData(b)
			if err != nil || len(rest) > 0 {
				//fmt.Printf("err - %v\n", err)
				return
			} else {
				//fmt.Println("e")
				eTxID = e.GetHash().String()
			}
		} else {
			//fmt.Println("ec")
			eTxID = ec.GetEntryHash().String()
			ecTxID = ec.GetHash().String()
		}
	} else {
		//fmt.Println("cc")
		eTxID = cc.GetEntryHash().String()
		ecTxID = cc.GetHash().String()
	}

	//fmt.Printf("eTxID - %v\n", eTxID)
	//fmt.Printf("ecTxID - %v\n", ecTxID)
	return
}

type AckRequest struct {
	TxID            string `json:"txid,omitempty"`
	FullTransaction string `json:"fulltransaction,omitempty"`
}

type FactoidTxStatus struct {
	TxID string `json:"txid"`
	GeneralTransactionData
}

type EntryStatus struct {
	CommitTxID string `json:"committxid"`
	EntryHash  string `json:"entryhash"`

	CommitData GeneralTransactionData `json:"commitdata"`
	EntryData  GeneralTransactionData `json:"entrydata"`

	ReserveTransactions          []ReserveInfo `json:"reserveinfo,omitempty"`
	ConflictingRevealEntryHashes []string      `json:"conflictingrevealentryhashes,omitempty"`
}

type ReserveInfo struct {
	TxID    string `json:"txid"`
	Timeout int64  `json:"timeout"` //Unix time
}

type GeneralTransactionData struct {
	TransactionDate       int64  `json:"transactiondate,omitempty"`       //Unix time
	TransactionDateString string `json:"transactiondatestring,omitempty"` //ISO8601 time
	BlockDate             int64  `json:"blockdate,omitempty"`             //Unix time
	BlockDateString       string `json:"blockdatestring,omitempty"`       //ISO8601 time

	Malleated *Malleated `json:"malleated,omitempty"`
	Status    string     `json:"status"`
}

type Malleated struct {
	MalleatedTxIDs []string `json:"malleatedtxids"`
}

const (
	AckStatusInvalid         = "Invalid"
	AckStatusUnknown         = "Unknown"
	AckStatusNotConfirmed    = "NotConfirmed"
	AckStatusACK             = "TransactionACK"
	AckStatus1Minute         = "1Minute"
	AckStatusDBlockConfirmed = "DBlockConfirmed"
)
