package mongo

import (
	"math"
	"sync"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/opensds/go-panda/backend/pkg/model"
)

type mongoRepository struct {
	session *mgo.Session
}

var defaultDBName = "go-panda"
var defaultCollection = "backends"
var mutex sync.Mutex
var mongoRepo = &mongoRepository{}

func Init(host string) *mongoRepository {
	mutex.Lock()
	defer mutex.Unlock()

	if mongoRepo.session != nil {
		return mongoRepo
	}

	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	mongoRepo.session = session
	return mongoRepo
}

// The implementation of Repository

func (repo *mongoRepository) CreateBackend(backend *model.Backend) (*model.Backend, error) {
	session := repo.session.Copy()
	defer session.Close()

	if backend.Id == "" {
		backend.Id = bson.NewObjectId()
	}

	err := session.DB(defaultDBName).C(defaultCollection).Insert(backend)
	if err != nil {
		return nil, err
	}
	return backend, nil
}

func (repo *mongoRepository) DeleteBackend(id string) error {
	session := repo.session.Copy()
	defer session.Close()
	return session.DB(defaultDBName).C(defaultCollection).RemoveId(bson.ObjectIdHex(id))
}

func (repo *mongoRepository) UpdateBackend(backend *model.Backend) (*model.Backend, error) {
	session := repo.session.Copy()
	defer session.Close()

	err := session.DB(defaultDBName).C(defaultCollection).UpdateId(backend.Id, backend)
	if err != nil {
		return nil, err
	}
	return backend, nil
}

func (repo *mongoRepository) GetBackend(id string) (*model.Backend, error) {
	session := repo.session.Copy()
	defer session.Close()

	var backend = &model.Backend{}
	collection := session.DB(defaultDBName).C(defaultCollection)
	err := collection.FindId(bson.ObjectIdHex(id)).One(backend)
	if err != nil {
		return nil, err
	}
	return backend, nil
}

func (repo *mongoRepository) ListBackend(limit, offset int, query interface{}) ([]*model.Backend, error) {

	session := repo.session.Copy()
	defer session.Close()

	if limit == 0 {
		limit = math.MinInt32
	}
	var backends []*model.Backend
	err := session.DB(defaultDBName).C(defaultCollection).Find(query).Skip(offset).Limit(limit).All(&backends)
	if err != nil {
		return nil, err
	}
	return backends, nil
}

func (repo *mongoRepository) Close() {
	repo.session.Close()
}
