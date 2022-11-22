package triplegenerator

import (
	"errors"

	cs "github.com/acompany-develop/QuickMPC-BTS/src/BeaverTripleService/ConfigStore"
	logger "github.com/acompany-develop/QuickMPC-BTS/src/BeaverTripleService/Log"
	ts "github.com/acompany-develop/QuickMPC-BTS/src/BeaverTripleService/TripleStore"
	utils "github.com/acompany-develop/QuickMPC-BTS/src/BeaverTripleService/Utils"
)

var Db *ts.SafeTripleStore
var tripleRandMax = int64(1000)
var tripleRandMin = int64(-1000)
var sharizeRandMax = int64(1 << 60)
var sharizeRandMin = int64(-1 << 60)

func init() {
	Db = ts.GetInstance()
}

func sharize(data int64, size uint32) ([]int64, error) {
	shares, err := utils.GetRandInt64Slice(uint64(size-1), sharizeRandMin, sharizeRandMax)
	if err != nil {
		errText := "乱数取得に失敗"
		logger.Error(errText)
		return nil, errors.New(errText)
	}

	sum := int64(0)
	for _, x := range shares {
		sum += x
	}

	shares = append(shares, data-sum)
	return shares, nil
}

func GenerateTriples(amount uint32) (map[uint32]([]*ts.Triple), error) {
	ret := make(map[uint32]([]*ts.Triple))

	for i := uint32(0); i < amount; i++ {
		randInt64Slice, err := utils.GetRandInt64Slice(2, tripleRandMin, tripleRandMax)
		if err != nil {
			errText := "乱数取得に失敗"
			logger.Error(errText)
			return nil, errors.New(errText)
		}

		a := randInt64Slice[0]
		b := randInt64Slice[1]
		c := a * b

		aShares, err := sharize(a, cs.Conf.PartyNum)
		if err != nil {
			return nil, err
		}
		bShares, err := sharize(b, cs.Conf.PartyNum)
		if err != nil {
			return nil, err
		}
		cShares, err := sharize(c, cs.Conf.PartyNum)
		if err != nil {
			return nil, err
		}

		// partyIdは1-index
		for partyId := uint32(1); partyId <= cs.Conf.PartyNum; partyId++ {
			t := ts.Triple{
				A: aShares[partyId-1],
				B: bShares[partyId-1],
				C: cShares[partyId-1],
			}
			ret[partyId] = append(ret[partyId], &t)
		}
	}

	return ret, nil
}

func GetTriples(jobId uint32, partyId uint32, amount uint32) ([]*ts.Triple, error) {
	Db.Mux.Lock()
	defer Db.Mux.Unlock()

	if len(Db.Triples[jobId]) == 0 {
		newTriples, err := GenerateTriples(amount)
		if err != nil {
			return nil, err
		}

		Db.Triples[jobId] = newTriples
	}

	var triples []*ts.Triple
	_, ok := Db.Triples[jobId][partyId]

	// とあるパーティの複数回目のリクエストが、他パーティより先行されても対応できるように全パーティに triple をappendする
	if !ok {
		newTriples, err := GenerateTriples(amount)
		if err != nil {
			return nil, err
		}

		// partyIdは1-index
		for partyId := uint32(1); partyId <= cs.Conf.PartyNum; partyId++ {
			_, ok := Db.Triples[jobId][partyId]
			if ok {
				Db.Triples[jobId][partyId] = append(Db.Triples[jobId][partyId], newTriples[partyId]...)
			} else {
				Db.Triples[jobId][partyId] = newTriples[partyId]
			}
		}
	}

	triples = Db.Triples[jobId][partyId][:amount]
	Db.Triples[jobId][partyId] = Db.Triples[jobId][partyId][amount:]
	if len(Db.Triples[jobId][partyId]) == 0 {
		delete(Db.Triples[jobId], partyId)
	}

	if len(triples) == 0 {
		errText := "すでに取得済みのリソースがリクエストされた"
		logger.Error(errText)
		return nil, errors.New(errText)
	}

	// 全て配り終わったら削除
	if len(Db.Triples[jobId]) == 0 {
		delete(Db.Triples, jobId)
	}

	return triples, nil
}
