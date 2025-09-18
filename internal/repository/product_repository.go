package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ojihalawa/daily-coffee-api.git/internal/entity"
	"github.com/ojihalawa/daily-coffee-api.git/internal/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductRepository struct {
	*Repository[entity.Product]
	Log         *logrus.Logger
	RedisClient *redis.Client
}

func NewProductRepository(log *logrus.Logger, redisClient *redis.Client) *ProductRepository {
	// return &ProductRepository{
	// 	Log:         log,
	// 	RedisClient: redisClient,
	// }
	return &ProductRepository{
		Repository: &Repository[entity.Product]{ // inisialisasi embedded
			RedisClient: redisClient,
		},
		Log: log,
	}
}

func (r *ProductRepository) ExistsByName(db *gorm.DB, name string) (bool, error) {
	var count int64
	err := db.Model(&entity.Product{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (r *ProductRepository) FindSpecialProduct(db *gorm.DB) (*entity.Product, error) {
	var p entity.Product
	if err := db.Where("is_special = ?", true).Take(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository[T]) FindAllWithRedis(db *gorm.DB, entities *[]T, pagination *utils.PaginationRequest) (int64, error) {
	var total int64
	entityName := fmt.Sprintf("%T", new(T))
	cacheKey := fmt.Sprintf("%s:all:page:%d:limit:%d:search:%s:order:%s:%s",
		entityName,
		pagination.Page,
		pagination.Limit,
		pagination.Search,
		pagination.OrderBy,
		pagination.SortBy,
	)

	// 2. Cek Redis
	val, err := r.RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		// ada di cache → unmarshal dan return
		if err := json.Unmarshal([]byte(val), entities); err == nil {
			fmt.Println("✅ Data dari Redis cache")
			// total juga bisa cache jika perlu, untuk simplicity kita skip cache total
			return int64(len(*entities)), nil
		}
	} else if err != redis.Nil {
		// Redis error
		return 0, err
	}

	query := db.Model(new(T))

	// cek apakah entity implement Searchable
	if s, ok := any(new(T)).(utils.Searchable); ok && pagination.Search != "" {
		fields := s.SearchFields()
		conditions := make([]string, len(fields))
		args := make([]interface{}, len(fields))
		for i, f := range fields {
			conditions[i] = f + " ILIKE ?"
			args[i] = "%" + pagination.Search + "%"
		}
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	// order
	if pagination.OrderBy != "" {
		order := pagination.OrderBy
		if pagination.SortBy != "" {
			order += " " + pagination.SortBy
		}
		query = query.Order(order)
	}

	// count total data
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}

	// paging
	offset := (pagination.Page - 1) * pagination.Limit
	if err := query.Offset(offset).Limit(pagination.Limit).Find(entities).Error; err != nil {
		return 0, err
	}

	data, _ := json.Marshal(entities)
	// TTL bisa diubah sesuai kebutuhan
	if err := r.RedisClient.Set(context.Background(), cacheKey, data, 10*time.Minute).Err(); err != nil {
		fmt.Println("⚠️ Gagal simpan cache:", err)
	}

	fmt.Println("💾 Data dari DB, disimpan ke Redis")

	return total, nil
}
