package store

import (
	"apietherscan/internal/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"sync"
)

type Store interface {
	InsertTransInfo([]interface{}) error
	GetTransInfo(int) ([]model.TransInfo, error)
}

type store struct {
	db *mgo.Session
	m  sync.RWMutex
}

func NewStore(db *mgo.Session) *store {
	return &store{db: db}
}

func (s *store) InsertTransInfo(data []interface{}) error {
	s.m.Lock()
	if err := s.db.DB("e").C("t").Insert(data...); err != nil {
		return err
	}
	s.m.Unlock()
	return nil
}

func (s *store) GetTransInfo(block int) ([]model.TransInfo, error) {
	var ti []model.TransInfo

	q := bson.M{
		"numblock": block,
	}

	s.m.RLock()
	if err := s.db.DB("e").C("t").Find(q).All(&ti); err != nil {
		return nil, err
	}
	s.m.RUnlock()

	return ti, nil
}
