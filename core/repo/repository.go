package repo

import (
	"context"

	"gorm.io/gorm"
)

// Repository adalah generic repository untuk entity T
type Repository[T any] struct {
	DB *gorm.DB
}

// ==========================
// Query Options
// ==========================

type QueryOption func(*gorm.DB) *gorm.DB

func WithWhere(query string, args ...any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

func WithOrder(order string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	}
}

func WithLimit(limit int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func WithOffset(offset int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func WithPreload(field string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(field)
	}
}

func WithSearch(field string, keyword string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" {
			return db
		}
		return db.Where(field+" LIKE ?", "%"+keyword+"%")
	}
}

func WithEqual(field string, value any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if value == nil {
			return db
		}
		return db.Where(field+" = ?", value)
	}
}

func WithBetween(field string, from, to any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if from == nil || to == nil {
			return db
		}
		return db.Where(field+" BETWEEN ? AND ?", from, to)
	}
}

// ==========================
// CRUD & Query
// ==========================

// FindAll: ambil semua record tanpa kondisi
func (r *Repository[T]) FindAll(db *gorm.DB, out *[]T) error {
	return db.Model(new(T)).Find(out).Error
}

// FindMany: ambil banyak record dengan query options
func (r *Repository[T]) FindMany(db *gorm.DB, out *[]T, opts ...QueryOption) error {
	query := db.Model(new(T))
	for _, opt := range opts {
		query = opt(query)
	}
	return query.Find(out).Error
}

// FindWithPagination: ambil data + total count
func (r *Repository[T]) FindWithPagination(
	db *gorm.DB,
	out *[]T,
	page, pageSize int,
	opts ...QueryOption,
) (total int64, err error) {
	query := db.Model(new(T))
	for _, opt := range opts {
		query = opt(query)
	}

	// hitung total
	if err = query.Count(&total).Error; err != nil {
		return
	}

	// pagination
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	err = query.Find(out).Error
	return
}

// FindOne: ambil single record
func (r *Repository[T]) FindOne(db *gorm.DB, out *T, opts ...QueryOption) error {
	query := db.Model(new(T))
	for _, opt := range opts {
		query = opt(query)
	}
	return query.Take(out).Error
}

// Create: insert 1 record
func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

// CreateBulk: insert banyak record
func (r *Repository[T]) CreateBulk(db *gorm.DB, entities []T) error {
	return db.Create(&entities).Error
}

// Update: update full entity (by primary key)
func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

// UpdateOne: update sebagian field berdasarkan kondisi
func (r *Repository[T]) UpdateOne(db *gorm.DB, updates map[string]interface{}, opts ...QueryOption) error {
	query := db.Model(new(T))
	for _, opt := range opts {
		query = opt(query)
	}
	return query.Updates(updates).Error
}

// UpdateBulk: update banyak record dengan kondisi
func (r *Repository[T]) UpdateBulk(db *gorm.DB, updates map[string]interface{}, opts ...QueryOption) error {
	query := db.Model(new(T))
	for _, opt := range opts {
		query = opt(query)
	}
	return query.Updates(updates).Error
}

// Delete: hapus 1 entity
func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

// Count: hitung jumlah record dengan kondisi
func (r *Repository[T]) Count(db *gorm.DB, opts ...QueryOption) (int64, error) {
	var total int64
	query := db.Model(new(T))
	for _, opt := range opts {
		query = opt(query)
	}
	err := query.Count(&total).Error
	return total, err
}

// ==========================
// Transaction
// ==========================

func ExecuteInTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// jalankan function yang dikasih
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
