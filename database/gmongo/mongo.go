package gmongo

// note: time.Time precision in MongoDB is Milliseconds, but in Go it is Nanoseconds,
// if time.Time used to be _id of document, please note the precision lost!

import (
	"context"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/ginterface"
	"github.com/cryptowilliam/goutil/container/grange"
	"github.com/cryptowilliam/goutil/container/gstring"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type (
	// MongoDB connection.
	Conn struct {
		inCli *mongo.Client
	}

	// MongoDB database instance.
	Database struct {
		c    *Conn
		inDb *mongo.Database
	}

	// MongoDB collection instance.
	Coll struct {
		db     *Database
		inColl *mongo.Collection
	}

	// Collection cursor.
	Cursor struct {
		coll  *Coll
		inCur *mongo.Cursor
	}
)

func Id2String(id interface{}) string {
	switch id.(type) {
	case primitive.ObjectID:
		return id.(primitive.ObjectID).String()
	case string:
		return id.(string)
	default:
		return fmt.Sprintf("%v", id)
	}
}

// uri example: mongodb:192.168.7.11:34567
func Dial(dsn string) (*Conn, error) {
	c := Conn{}
	opts := options.Client().ApplyURI(dsn)

	ofcSess, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, gerrors.Wrap(err, dsn)
	}
	if err := ofcSess.Ping(context.Background(), nil); err != nil {
		return nil, gerrors.Wrap(err, dsn)
	}
	c.inCli = ofcSess
	return &c, nil
}

// Get all database names.
func (c *Conn) ListDatabases() ([]string, error) {
	return c.inCli.ListDatabaseNames(context.Background(), bson.D{} /*can't input nil, otherwise will return error "document is nil"*/)
}

// Whether database exists or not.
func (c *Conn) IsDatabaseExist(dbName string) (bool, error) {
	dbNames, err := c.ListDatabases()
	if err != nil {
		return false, err
	}
	return gstring.CountByValue(dbNames, dbName) > 0, nil
}

// IMPORTANT
// Various indications in practical applications,
// this will create new socket connection between app and server.
// Need official document to prove this.
func (c *Conn) Database(DBName string) *Database {
	db := Database{}
	db.c = c
	db.inDb = c.inCli.Database(DBName)
	return &db
}

// Close mongodb connection.
func (c *Conn) Close() error {
	if c.inCli == nil {
		return nil
	}
	return c.inCli.Disconnect(context.Background())
}

func (d *Database) Exist() (bool, error) {
	return d.c.IsDatabaseExist(d.inDb.Name())
}

func (d *Database) ListCollections() ([]string, error) {
	var names []string
	cur, err := d.inDb.ListCollections(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.Background()) {
		next := bsonx.Doc{}
		err := cur.Decode(&next)
		if err != nil {
			return nil, err
		}

		elem, err := next.LookupErr("name")
		if err != nil {
			return nil, err
		}

		if elem.Type() != bson.TypeString {
			return nil, fmt.Errorf("incorrect type for 'name'. got %v. want %v", elem.Type(), bson.TypeString)
		}

		elemName := elem.StringValue()
		names = append(names, elemName)
	}
	return names, nil
}

func (d *Database) Collection(collName string) *Coll {
	coll := Coll{}
	coll.db = d
	coll.inColl = d.inDb.Collection(collName)
	return &coll
}

func (d *Database) DeleteCollection(collName string) error {
	return d.inDb.Collection(collName).Drop(context.Background())
}

func (d *Database) IsCollectionExists(collName string) (bool, error) {
	if colls, err := d.ListCollections(); err != nil {
		return false, err
	} else {
		return gstring.CountByValue(colls, collName) == 1, nil
	}
}

func (d *Database) Watch() error {
	return nil
	//pipeline        := []interface{}{}
	//opts := options.ChangeStream()
	//cur, err := d.inDb.Watch(context.Background(), pipeline, opts)
}

func (d *Database) RenameCollection(from, to string) error {
	renameCmd := bson.D{
		{"renameCollection", d.inDb.Name() + "." + from},
		{"to", d.inDb.Name() + "." + to},
	}
	sr := d.inDb.Client().Database("admin").RunCommand(context.Background(), renameCmd)
	return sr.Err()
}

func (c *Coll) Exist() (bool, error) {
	return c.db.IsCollectionExists(c.inColl.Name())
}

func (c *Coll) MinFloat(floatPath string) (float64, error) {
	val, err := c.minMax(floatPath, float64(0), false)
	if err != nil {
		return float64(0), err
	}
	return val.(float64), nil
}

func (c *Coll) MaxFloat(floatPath string) (float64, error) {
	val, err := c.minMax(floatPath, float64(0), true)
	if err != nil {
		return float64(0), err
	}
	return val.(float64), nil
}

// NOTE: if collection not exist, will returne [0, false, nil]
func (c *Coll) MinInt64(intPath string) (int64, error) {
	val, err := c.minMax(intPath, int64(0), false)
	if err != nil {
		return int64(0), err
	}
	return val.(int64), nil
}

// NOTE: if collection not exist, will returne [0, false, nil]
func (c *Coll) MaxInt64(intPath string) (int64, error) {
	val, err := c.minMax(intPath, int64(0), true)
	if err != nil {
		return int64(0), err
	}
	return val.(int64), nil
}

// NOTE: if collection not exist, will returne [0, false, nil]
func (c *Coll) MinTime(intPath string) (time.Time, error) {
	val, err := c.minMax(intPath, time.Time{}, false)
	if err != nil {
		return time.Time{}, err
	}
	return val.(time.Time), nil
}

// NOTE: if collection not exist, will returne [0, false, nil]
func (c *Coll) MaxTime(intPath string) (time.Time, error) {
	val, err := c.minMax(intPath, time.Time{}, true)
	if err != nil {
		return time.Time{}, err
	}
	return val.(time.Time), nil
}

func (c *Coll) Min(itemPath string, typeSample interface{}) (interface{}, error) {
	return c.minMax(itemPath, typeSample, false)
}

func (c *Coll) Max(itemPath string, typeSample interface{}) (interface{}, error) {
	return c.minMax(itemPath, typeSample, true)
}

// NOTE: if timePath/Collection/Database doesn't exist or empty collection, returns [nil, ErrNotExist]
// NOTE: 哪怕itemPath是文档中的某个成员，比如int，Decode的时候也要用整个结构体去Decode出整个文档
// 比如 minMax("_id", ...) 的时候，尽管查询路径只是时间，可你要把*time.Time传进去cur.Decode，会报错 cannot decode invalid into a time.Time
// example: Max("_id", time.Time{})
func (c *Coll) minMax(itemPath string, typeSample interface{}, max bool) (interface{}, error) {
	opts := options.FindOptions{}
	val := int64(1)
	if max {
		val = -1
	}
	opts.SetSort(bsonx.Doc{bsonx.Elem{Key: itemPath, Value: bsonx.Int64(val)}})
	opts.SetLimit(1)
	if cur, err := c.inColl.Find(context.Background(), bson.D{} /*can't input nil, otherwise will return error "document is nil"*/, &opts); err != nil {
		return nil, err
	} else {
		doc := bsonx.Doc{}
		if cur.Next(context.Background()) {
			if err := cur.Decode(&doc); err != nil {
				return nil, err
			}

			errTypeNotMatch := gerrors.Errorf("type(%s) wanted but in database it is not", ginterface.Type(typeSample))
			errTypeNotSupported := gerrors.Errorf("type(%s) not supported in min/max query", ginterface.Type(typeSample))

			switch ginterface.Type(typeSample) {
			case "time.Time":
				tm, ok := doc.Lookup(itemPath).TimeOK()
				if !ok {
					return time.Time{}, errTypeNotMatch
				}
				return tm, nil
			case "int64":
				i64, ok := doc.Lookup(itemPath).Int64OK()
				if !ok {
					return 0, errTypeNotMatch
				}
				return i64, nil
			case "int32":
				i32, ok := doc.Lookup(itemPath).Int32OK()
				if !ok {
					return 0, errTypeNotMatch
				}
				return i32, nil
			case "float64":
				f64, ok := doc.Lookup(itemPath).DoubleOK()
				if !ok {
					return 0.0, errTypeNotMatch
				}
				return f64, nil
			default:
				return nil, errTypeNotSupported
			}
		} else { // if database/collection not exist, will enter here
			return nil, gerrors.ErrNotExist // this is not a really error
		}
	}
}

func (c *Coll) Insert(doc interface{}) error {
	_, err := c.inColl.InsertOne(context.Background(), doc)
	return err
}

func (c *Coll) InsertMany(docs ...interface{}) error {
	var ins []interface{}
	ins = append(ins, docs)
	_, err := c.inColl.InsertMany(context.Background(), ins)
	return err
}

// if not exists: insert, if exists: update
func (c *Coll) Upsert(selector interface{}, update interface{}) error {
	opts := options.FindOneAndUpdateOptions{}
	opts.SetUpsert(true)
	dr := c.inColl.FindOneAndUpdate(context.Background(), selector, update, &opts)
	err := dr.Decode(dr)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

// Upsert one entire document.
// note: time.Time precision in MongoDB is Milliseconds, but in Go it is Nanoseconds,
// if time.Time used to be _id of document, please note the precision!
func (c *Coll) UpsertEntireDoc(id interface{}, doc interface{}) error {
	b, err := bson.Marshal(doc)
	if err != nil {
		return err
	}
	mmap := bson.M{}
	err = bson.Unmarshal(b, &mmap)
	if err != nil {
		return err
	}

	exist, err := c.IdExists(id)
	if err != nil {
		return err
	}
	if exist {
		_, err = c.inColl.ReplaceOne(context.Background(), bson.M{"_id": id}, &mmap)
		return err
	} else {
		return c.UpsertFields(id, doc)
	}
}

// Upsert multiple entire documents.
// note: time.Time precision in MongoDB is Milliseconds, but in Go it is Nanoseconds,
// if time.Time used to be _id of document, please note the precision!
func (c *Coll) UpsertDocs(ids []interface{}, doc []interface{}) error {
	return gerrors.ErrNotImplemented
}

// Upsert part of one document.
// if not exists: insert, if exists: update
func (c *Coll) UpsertFields(id interface{}, doc interface{}) error {
	return c.Upsert(bson.M{"_id": id}, bson.M{"$set": doc})
}

// 未测试！
// if not exists: insert
func (c *Coll) InsertIfNotExists(selector interface{}, update interface{}) error {
	return c.Upsert(selector, bson.M{"$setOnInsert": update})
}

func (c *Coll) Remove(selector interface{}) error {
	err := error(nil)
	_, err = c.inColl.DeleteOne(context.Background(), selector)
	return err
}

func (c *Coll) RemoveCmp(path string, cmpopt Cmp, value bsonx.Val) (int64, error) {
	dr, err := c.inColl.DeleteMany(context.Background(), bsonx.Doc{{path, bsonx.Document(bsonx.Doc{{string(cmpopt), value}})}})
	if err != nil {
		return 0, err
	}
	return dr.DeletedCount, err
}

func (c *Coll) RemoveAll(selector interface{}) error {
	if selector == nil {
		selector = bson.D{} // nil filter not allowed, otherwise will return error: document is nil
	}

	err := error(nil)
	_, err = c.inColl.DeleteMany(context.Background(), selector)
	return err
}

func (c *Coll) RemoveId(id interface{}) error {
	err := error(nil)
	_, err = c.inColl.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	return err
}

func (c *Coll) RemoveStringId(stringId string) error {
	err := error(nil)
	_, err = c.inColl.DeleteOne(context.Background(), bson.M{"_id": stringId})
	return err
}

// 未测试！
func (c *Coll) RemovePartDocs(selector interface{}, toRemove interface{}) error {
	err := error(nil)
	_, err = c.inColl.UpdateMany(context.Background(), selector, bson.M{"$unset": toRemove})
	return err
}

func (c *Coll) RemoveObjectId(id primitive.ObjectID) error {
	err := error(nil)
	_, err = c.inColl.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (c *Coll) RemoveCmpIntId(id int64, cmpopt Cmp) error {
	_, err := c.inColl.DeleteMany(context.Background(), bsonx.Doc{{"_id", bsonx.Document(bsonx.Doc{{string(cmpopt), bsonx.Int64(id)}})}})
	return err
}

// count by meta data
// if Database / Coll doesn't exist, returns valid 0 and nil error
//
// WARNING:
// this API is faster but maybe incorrect
// https://stackoverflow.com/questions/30715466/incorrect-count-returned-by-mongodb-wiredtiger
// may returns error count if metadata incorrect after mongodb exception shutdown
//func (c *Coll) CountByMeta_WARNING(filter interface{}) (int64, error) {
//	return c.inColl.EstimatedDocumentCount(context.Background(), filter)
//}

// count by documents
// WARNING:
// this API is slow but correct if metadata is incorrect
func (c *Coll) Count(filter interface{}) (int64, error) {
	if filter == nil {
		filter = bson.D{} // nil filter not allowed, otherwise will return error: document is nil
	}
	return c.inColl.CountDocuments(context.Background(), filter)
}

func (c *Coll) IdExists(id interface{}) (bool, error) {
	return c.FindId(id, nil)
}

// TODO: test (again) required
// copy all data to new collection
func (c *Coll) CloneTo(to *Coll) (cloneCount int, err error) {
	fromCount, err := c.Count(nil)
	if err != nil {
		return 0, err
	}
	toCount, err := to.Count(nil)
	if err != nil {
		return 0, err
	}
	if fromCount == toCount {
		return 0, nil
	}
	if fromCount != toCount && toCount > 0 {
		if err := to.RemoveAll(nil); err != nil {
			return 0, err
		}
	}

	cur, err := c.Find(nil)
	if err != nil {
		return 0, err
	}
	for cur.Next() {
		if r, err := cur.DecodeBson(); err != nil {
			return cloneCount, err
		} else {
			if err := to.Insert(r); err != nil {
				return cloneCount, err
			} else {
				cloneCount++
			}
		}
	}
	return cloneCount, nil
}

// if Database / Coll doesn't exist, returns valid *Cursor and nil error, but cur.Next() == false
func (c *Coll) Find(query interface{}) (*Cursor, error) {
	if query == nil {
		query = bson.D{} // nil filter not allowed, otherwise will return error: document is nil
	}
	rst := Cursor{coll: c}
	err := error(nil)
	rst.inCur, err = c.inColl.Find(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return &rst, nil
}

func (c *Coll) FindId(id interface{}, decodeResult interface{}) (found bool, err error) {
	rst := Cursor{}
	rst.coll = c
	rst.inCur, err = c.inColl.Find(context.Background(), bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	if rst.Next() {
		if decodeResult != nil {
			err := rst.Decode(decodeResult)
			return true, err
		} else {
			return true, nil
		}
	} else {
		return false, nil
	}
}

// asc, true ASC, false DESC
func (c *Coll) FindEdge(limit int64, asc bool) (*Cursor, error) {
	opts := options.FindOptions{}
	sort := int64(1)
	if asc == false {
		sort = -1
	}
	opts.SetSort(bsonx.Doc{bsonx.Elem{Key: "_id", Value: bsonx.Int64(sort)}})
	opts.SetLimit(limit)
	cur, err := c.inColl.Find(context.Background(), nil, &opts)
	if err != nil {
		return nil, err
	}
	return &Cursor{coll: c, inCur: cur}, nil
}

type Cmp string

const (
	CmpLT  Cmp = "$lt"
	CmpLTE Cmp = "$lte"
	CmpGT  Cmp = "$gt"
	CmpGTE Cmp = "$gte"
)

// example:
// FindCmp("_id", mongo.CmpGTE, bsonx.T(fromTime))
func (c *Coll) FindCmp(path string, cmpopt Cmp, value bsonx.Val) (*Cursor, error) {
	return c.Find(bsonx.Doc{{path, bsonx.Document(bsonx.Doc{{string(cmpopt), value}})}})
}

func (c *Cursor) Next() bool {
	return c.inCur.Next(context.Background())
}

// TODO: test required
// get Id of current document, use to recognize some special documents which has different structure compare with other documents
func (c *Cursor) DecodeId(typeSample interface{}) (interface{}, error) {
	doc := bsonx.Doc{}

	if err := c.inCur.Decode(&doc); err != nil {
		return nil, err
	}

	errTypeNotMatch := gerrors.Errorf("type(%s) wanted but in database it is not", ginterface.Type(typeSample))
	errTypeNotSupported := gerrors.Errorf("type(%s) not supported in min/max query", ginterface.Type(typeSample))

	switch ginterface.Type(typeSample) {
	case "string":
		str, ok := doc.Lookup("_id").StringValueOK()
		if !ok {
			return time.Time{}, errTypeNotMatch
		}
		return str, nil
	case "time.Time":
		tm, ok := doc.Lookup("_id").TimeOK()
		if !ok {
			return time.Time{}, errTypeNotMatch
		}
		return tm, nil
	case "int64":
		i64, ok := doc.Lookup("_id").Int64OK()
		if !ok {
			return 0, errTypeNotMatch
		}
		return i64, nil
	case "int32":
		i32, ok := doc.Lookup("_id").Int32OK()
		if !ok {
			return 0, errTypeNotMatch
		}
		return i32, nil
	case "float64":
		f64, ok := doc.Lookup("_id").DoubleOK()
		if !ok {
			return 0.0, errTypeNotMatch
		}
		return f64, nil
	default:
		return nil, errTypeNotSupported
	}
}

// Note: 存储的时间如果是UTC时区，Decode出来的会变成Local时区，需要仔细转换时区
func (c *Cursor) Decode(result interface{}) error {
	return c.inCur.Decode(result)
}

// TODO: test required, test at CloneTo()
func (c *Cursor) DecodeBson() (bson.Raw, error) {
	doc := bson.Raw{}
	err := c.inCur.Decode(&doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// copy all data to new collection
func (mg *Conn) CopyColl(fromDb, fromColl, toDb, toColl string) (cloneCount int, err error) {
	from := mg.Database(fromDb).Collection(fromColl)
	to := mg.Database(toDb).Collection(toColl)

	fromCount, err := from.Count(nil)
	if err != nil {
		return 0, err
	}
	toCount, err := to.Count(nil)
	if err != nil {
		return 0, err
	}
	if fromCount == toCount {
		return 0, nil
	}
	if fromCount != toCount && toCount > 0 {
		if err := to.RemoveAll(nil); err != nil {
			return 0, err
		}
	}

	cur, err := from.Find(nil)
	if err != nil {
		return 0, err
	}
	for cur.Next() {
		if r, err := cur.DecodeBson(); err != nil {
			return cloneCount, err
		} else {
			if err := to.Insert(r); err != nil {
				return cloneCount, err
			} else {
				cloneCount++
			}
		}
	}
	return cloneCount, nil
}

func (c *Coll) VerifyContinuousInt64Id() (*grange.RangeFilter, error) {
	minId, err := c.MinInt64("_id")
	if err != nil {
		return nil, err
	}
	maxId, err := c.MaxInt64("_id")
	if err != nil {
		return nil, err
	}
	correctRf := grange.NewRangeFilter()
	correctRf.AddRange(grange.NewRange(minId, maxId))

	infactRf := grange.NewRangeFilter()
	cur, err := c.Find(nil)
	if err != nil {
		return nil, err
	}
	type ID struct {
		Id int64 `bson:"_id""`
	}
	id := ID{}
	for cur.Next() {
		id.Id = -1
		if err := cur.Decode(&id); err != nil {
			return nil, err
		}
		if id.Id == -1 {
			return nil, gerrors.Errorf("decoded -1 id")
		}
		infactRf.AddInt64(id.Id)
	}

	correctRf.Sub(*infactRf)
	return correctRf, nil
}
