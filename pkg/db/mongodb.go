package db

import (
	"crypto/tls"
	"github.com/globalsign/mgo"
	"net"
)

func NewMongo(pass string) (*mgo.Session, error) {
	conn, err := createConnection(pass)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}

func createConnection(pass string) (*mgo.Session, error) {
	dialInfo := mgo.DialInfo{Addrs: []string{
		"etherscan-shard-00-00.iibcf.mongodb.net:27017",
		"etherscan-shard-00-01.iibcf.mongodb.net:27017",
		"etherscan-shard-00-02.iibcf.mongodb.net:27017",
	},
		Username: "ether",
		Password: pass}

	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	return mgo.DialWithInfo(&dialInfo)
}
