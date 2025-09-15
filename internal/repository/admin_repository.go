package repository

import (
	"database/sql"
	"fmt"
	"time"

	"student-management-system/internal/domain"

	"github.com/sirupsen/logrus"
)

type AdminRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewAdminRepository(db *sql.DB, logger *logrus.Logger) *AdminRepository {
	return &AdminRepository{
		db:     db,
		logger: logger,
	}
}

// CreateAdmin 创建管理员
func (r *AdminRepository) CreateAdmin(admin *domain.Admin) error {
	query := `
		INSERT INTO admins (account, password, name, phone, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRow(query, admin.Account, admin.Password, admin.Name,
		admin.Phone, admin.Email, now, now).Scan(&admin.ID)

	if err != nil {
		r.logger.WithError(err).Error("Failed to create admin")
		return fmt.Errorf("failed to create admin: %v", err)
	}

	admin.CreatedAt = now
	admin.UpdatedAt = now

	r.logger.WithField("admin_id", admin.ID).Info("Admin created successfully")
	return nil
}

// GetAdminByID 根据ID获取管理员
func (r *AdminRepository) GetAdminByID(id int) (*domain.Admin, error) {
	query := `
		SELECT id, account, password, name, phone, email, created_at, updated_at
		FROM admins
		WHERE id = $1
	`

	admin := &domain.Admin{}
	err := r.db.QueryRow(query, id).Scan(
		&admin.ID, &admin.Account, &admin.Password, &admin.Name,
		&admin.Phone, &admin.Email, &admin.CreatedAt, &admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin not found")
		}
		r.logger.WithError(err).Error("Failed to get admin by ID")
		return nil, fmt.Errorf("failed to get admin: %v", err)
	}

	return admin, nil
}

// GetAdminByAccount 根据账号获取管理员
func (r *AdminRepository) GetAdminByAccount(account string) (*domain.Admin, error) {
	query := `
		SELECT id, account, password, name, phone, email, created_at, updated_at
		FROM admins
		WHERE account = $1
	`

	admin := &domain.Admin{}
	err := r.db.QueryRow(query, account).Scan(
		&admin.ID, &admin.Account, &admin.Password, &admin.Name,
		&admin.Phone, &admin.Email, &admin.CreatedAt, &admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin not found")
		}
		r.logger.WithError(err).Error("Failed to get admin by account")
		return nil, fmt.Errorf("failed to get admin: %v", err)
	}

	return admin, nil
}

// UpdateAdmin 更新管理员信息
func (r *AdminRepository) UpdateAdmin(admin *domain.Admin) error {
	query := `
		UPDATE admins 
		SET account = $1, password = $2, name = $3, phone = $4, email = $5, updated_at = $6
		WHERE id = $7
	`

	now := time.Now()
	result, err := r.db.Exec(query, admin.Account, admin.Password, admin.Name,
		admin.Phone, admin.Email, now, admin.ID)

	if err != nil {
		r.logger.WithError(err).Error("Failed to update admin")
		return fmt.Errorf("failed to update admin: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("admin not found")
	}

	admin.UpdatedAt = now

	r.logger.WithField("admin_id", admin.ID).Info("Admin updated successfully")
	return nil
}

// DeleteAdmin 删除管理员
func (r *AdminRepository) DeleteAdmin(id int) error {
	query := `DELETE FROM admins WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete admin")
		return fmt.Errorf("failed to delete admin: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("admin not found")
	}

	r.logger.WithField("admin_id", id).Info("Admin deleted successfully")
	return nil
}

// ListAdmins 获取管理员列表
func (r *AdminRepository) ListAdmins(page, pageSize int) ([]*domain.Admin, int, error) {
	// 获取总数
	countQuery := `SELECT COUNT(*) FROM admins`
	var total int
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		r.logger.WithError(err).Error("Failed to count admins")
		return nil, 0, fmt.Errorf("failed to count admins: %v", err)
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	query := `
		SELECT id, account, password, name, phone, email, created_at, updated_at
		FROM admins
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		r.logger.WithError(err).Error("Failed to list admins")
		return nil, 0, fmt.Errorf("failed to list admins: %v", err)
	}
	defer rows.Close()

	var admins []*domain.Admin
	for rows.Next() {
		admin := &domain.Admin{}
		err := rows.Scan(
			&admin.ID, &admin.Account, &admin.Password, &admin.Name,
			&admin.Phone, &admin.Email, &admin.CreatedAt, &admin.UpdatedAt,
		)
		if err != nil {
			r.logger.WithError(err).Error("Failed to scan admin row")
			return nil, 0, fmt.Errorf("failed to scan admin: %v", err)
		}
		admins = append(admins, admin)
	}

	if err = rows.Err(); err != nil {
		r.logger.WithError(err).Error("Error iterating admin rows")
		return nil, 0, fmt.Errorf("error iterating rows: %v", err)
	}

	return admins, total, nil
}
