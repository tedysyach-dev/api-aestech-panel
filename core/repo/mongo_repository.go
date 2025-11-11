package repo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ==========================
// Generic MongoRepository
// ==========================

// Contoh Entity MongoDb
//
//	type User struct {
//		ID        string `bson:"_id,omitempty"`
//		Name      string `bson:"name"`
//		Email     string `bson:"email"`
//		Role      string `bson:"role"`
//		CreatedAt int64  `bson:"created_at"`
//	}
type MongoRepository[T any] struct {
	Collection *mongo.Collection
}

func NewMongoRepository[T any](db *mongo.Database, collection string) *MongoRepository[T] {
	return &MongoRepository[T]{
		Collection: db.Collection(collection),
	}
}

// ==========================
// Query Options
// ==========================

type MongoQueryOption func(*queryConfig)

type queryConfig struct {
	filter    bson.M
	sort      bson.D
	limit     int64
	skip      int64
	lookups   []bson.D
	searchKey string
	searchVal string
}

// Filter
func MWithFilter(filter bson.M) MongoQueryOption {
	return func(q *queryConfig) {
		q.filter = filter
	}
}

// Sort
func MWithSort(sort bson.D) MongoQueryOption {
	return func(q *queryConfig) {
		q.sort = sort
	}
}

// Limit
func MWithLimit(limit int64) MongoQueryOption {
	return func(q *queryConfig) {
		q.limit = limit
	}
}

// Skip / Offset
func MWithSkip(skip int64) MongoQueryOption {
	return func(q *queryConfig) {
		q.skip = skip
	}
}

// Search
func MWithSearch(field, keyword string) MongoQueryOption {
	return func(q *queryConfig) {
		q.searchKey = field
		q.searchVal = keyword
	}
}

// Lookup
func MWithLookup(stage bson.D) MongoQueryOption {
	return func(q *queryConfig) {
		q.lookups = append(q.lookups, stage)
	}
}

// ==========================
// CRUD Operations
// ==========================

// Create
func (r *MongoRepository[T]) Create(ctx context.Context, entity *T) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, entity)
}

// CreateBulk
func (r *MongoRepository[T]) CreateBulk(ctx context.Context, entities []T) (*mongo.InsertManyResult, error) {
	docs := make([]interface{}, len(entities))
	for i, e := range entities {
		docs[i] = e
	}
	return r.Collection.InsertMany(ctx, docs)
}

// FindMany
func (r *MongoRepository[T]) FindMany(ctx context.Context, out *[]T, opts ...MongoQueryOption) (int64, error) {
	cfg := &queryConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// --- Build filter ---
	filter := cfg.filter
	if filter == nil {
		filter = bson.M{}
	}

	// --- Build pipeline if lookup or search used ---
	if len(cfg.lookups) > 0 || cfg.searchVal != "" {
		var pipeline mongo.Pipeline

		// Search
		if cfg.searchVal != "" {
			pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{
				cfg.searchKey: bson.M{"$regex": cfg.searchVal, "$options": "i"},
			}}})
		}

		// Lookup
		for _, stage := range cfg.lookups {
			pipeline = append(pipeline, stage)
		}

		// Filter
		if len(cfg.filter) > 0 {
			pipeline = append(pipeline, bson.D{{Key: "$match", Value: cfg.filter}})
		}

		// Sort
		if len(cfg.sort) > 0 {
			pipeline = append(pipeline, bson.D{{Key: "$sort", Value: cfg.sort}})
		}

		// Skip & Limit
		if cfg.skip > 0 {
			pipeline = append(pipeline, bson.D{{Key: "$skip", Value: cfg.skip}})
		}
		if cfg.limit > 0 {
			pipeline = append(pipeline, bson.D{{Key: "$limit", Value: cfg.limit}})
		}

		// Count total (tanpa skip & limit)
		countPipeline := append(pipeline[:len(pipeline)-2], bson.D{{Key: "$count", Value: "total"}})
		countCursor, err := r.Collection.Aggregate(ctx, countPipeline)
		if err != nil {
			return 0, err
		}
		var countResult []bson.M
		if err := countCursor.All(ctx, &countResult); err != nil {
			return 0, err
		}

		var total int64
		if len(countResult) > 0 {
			total = countResult[0]["total"].(int64)
		}

		cursor, err := r.Collection.Aggregate(ctx, pipeline)
		if err != nil {
			return 0, err
		}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, out); err != nil {
			return 0, err
		}

		return total, nil
	}

	// --- Simple find ---
	findOpts := options.Find()
	if cfg.limit > 0 {
		findOpts.SetLimit(cfg.limit)
	}
	if cfg.skip > 0 {
		findOpts.SetSkip(cfg.skip)
	}
	if len(cfg.sort) > 0 {
		findOpts.SetSort(cfg.sort)
	}

	cursor, err := r.Collection.Find(ctx, filter, findOpts)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, out); err != nil {
		return 0, err
	}

	total, err := r.Collection.CountDocuments(ctx, filter)
	return total, err
}

// FindOne
func (r *MongoRepository[T]) FindOne(ctx context.Context, out *T, filter bson.M) error {
	return r.Collection.FindOne(ctx, filter).Decode(out)
}

// UpdateOne
func (r *MongoRepository[T]) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	if len(filter) == 0 {
		return errors.New("filter kosong")
	}
	_, err := r.Collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	return err
}

// UpdateMany
func (r *MongoRepository[T]) UpdateMany(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := r.Collection.UpdateMany(ctx, filter, bson.M{"$set": update})
	return err
}

// DeleteOne
func (r *MongoRepository[T]) DeleteOne(ctx context.Context, filter bson.M) error {
	_, err := r.Collection.DeleteOne(ctx, filter)
	return err
}

// DeleteMany
func (r *MongoRepository[T]) DeleteMany(ctx context.Context, filter bson.M) error {
	_, err := r.Collection.DeleteMany(ctx, filter)
	return err
}

// ==========================
// Pagination Helper
// ==========================

type PaginatedResult[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

func (r *MongoRepository[T]) FindWithPagination(
	ctx context.Context,
	page int,
	limit int,
	opts ...MongoQueryOption,
) (*PaginatedResult[T], error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// apply skip & limit
	skip := int64((page - 1) * limit)

	cfg := &queryConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// tambahkan skip & limit ke query
	opts = append(opts, MWithSkip(skip), MWithLimit(int64(limit)))

	// jalankan FindMany dengan total count
	var results []T
	total, err := r.FindMany(ctx, &results, opts...)
	if err != nil {
		return nil, err
	}

	return &PaginatedResult[T]{
		Data:  results,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
