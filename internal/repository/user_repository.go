package repository

import (
	"go-article/internal/domain/entity"
	"go-article/internal/domain/model"
	"log"

	"gorm.io/gorm"
)

// UserRepository adalah interface yang mendefinisikan semua method untuk operasi user
type UserRepository interface {
	// Create membuat user baru, return pointer ke UserEntity agar tidak ada copy data
	Create(user entity.UserEntity) (*entity.UserEntity, error)
	// CreateWithRoles membuat user dengan role-role yang diberikan, return pointer
	CreateWithRoles(user entity.UserEntity, roleIDs []uint64) (*entity.UserEntity, error)
	// FindByEmail mencari user berdasarkan email, return pointer untuk efisiensi memory
	FindByEmail(email string) (*entity.UserEntity, error)
	// FindByID mencari user berdasarkan ID, return pointer untuk efisiensi memory
	FindByID(id uint64) (*entity.UserEntity, error)
	// GetRolesByIDs mengambil daftar role berdasarkan ID-ID yang diberikan
	GetRolesByIDs(roleIDs []uint64) ([]model.Role, error)
}

// userRepository adalah implementasi konkret dari interface UserRepository
type userRepository struct {
	// db adalah koneksi database GORM yang digunakan untuk query
	db *gorm.DB
}

// NewUserRepository adalah constructor untuk membuat instance userRepository baru
func NewUserRepository(db *gorm.DB) UserRepository {
	// Mengembalikan pointer ke struct userRepository dengan db yang diberikan
	return &userRepository{db: db}
}

// CreateWithRoles membuat user baru dengan role yang terkait
// Parameter: user adalah data user yang akan dibuat, roleIDs adalah ID-ID role yang akan dihubungkan
// Return: pointer ke UserEntity (bukan value) untuk efisiensi memory
func (u *userRepository) CreateWithRoles(user entity.UserEntity, roleIDs []uint64) (*entity.UserEntity, error) {
	// Konversi UserEntity ke User model agar compatible dengan GORM
	userModel := model.User{
		Name:     user.Name,     // Salin nama dari entity ke model
		Email:    user.Email,    // Salin email dari entity ke model
		Password: user.Password, // Salin password dari entity ke model
		Avatar:   user.Avatar,   // Salin avatar dari entity ke model
	}

	// Buat user di database menggunakan GORM Create
	err := u.db.Create(&userModel).Error
	if err != nil {
		// Jika ada error, log error dan return nil dengan error message
		log.Println("[UserRepository] CreateWithRoles:", err)
		return nil, err
	}

	// Jika ada role IDs yang diberikan, hubungkan role ke user
	if len(roleIDs) > 0 {
		// Query database untuk mendapatkan role dengan ID yang diminta
		var roles []model.Role
		err = u.db.Where("id IN ?", roleIDs).Find(&roles).Error
		if err != nil {
			// Jika error saat fetch roles, log dan return nil
			log.Println("[UserRepository] CreateWithRoles - fetching roles:", err)
			return nil, err
		}

		// Hubungkan (associate) roles ke user menggunakan GORM many-to-many relationship
		err = u.db.Model(&userModel).Association("Roles").Append(roles)
		if err != nil {
			// Jika error saat associate roles, log dan return nil
			log.Println("[UserRepository] CreateWithRoles - associating roles:", err)
			return nil, err
		}
	}

	// Konversi User model kembali ke UserEntity dengan data lengkap termasuk Roles
	result := &entity.UserEntity{
		ID:       userModel.ID,       // ID yang di-generate oleh database
		Name:     userModel.Name,     // Nama user
		Email:    userModel.Email,    // Email user
		Password: userModel.Password, // Password user (akan di-hide di JSON response)
		Avatar:   userModel.Avatar,   // Avatar URL (bisa nil)
		Roles:    userModel.Roles,    // Daftar role yang terkait dengan user
	}

	// Return pointer ke UserEntity (bukan value copy) - ini lebih efisien secara memory
	return result, nil
}

// FindByID mencari user berdasarkan ID dan mengembalikan pointer ke UserEntity
// Parameter: id adalah ID user yang dicari
// Return: pointer ke UserEntity (bukan value copy untuk efisiensi) dan error jika ada
func (u *userRepository) FindByID(id uint64) (*entity.UserEntity, error) {
	// Deklarasi variabel user dengan tipe model.User
	var user model.User
	// Query database untuk mencari user dengan ID tertentu dan preload Roles-nya
	err := u.db.Where("id = ?", id).Preload("Roles").First(&user).Error
	if err != nil {
		// Jika error (user tidak ditemukan atau error database), log error dan return nil
		log.Println("[UserRepository] FindByID:", err)
		return nil, err
	}

	// Konversi User model ke UserEntity dan return sebagai pointer
	return &entity.UserEntity{
		ID:    user.ID,    // ID user dari database
		Name:  user.Name,  // Nama user
		Email: user.Email, // Email user
		Roles: user.Roles, // Role-role yang terkait (sudah di-preload dari database)
	}, nil
}

// Create membuat user baru tanpa role dan mengembalikan pointer ke UserEntity
// Parameter: user adalah data user yang akan dibuat
// Return: pointer ke UserEntity yang baru dibuat dan error jika ada
func (u *userRepository) Create(user entity.UserEntity) (*entity.UserEntity, error) {
	// Konversi UserEntity ke User model agar compatible dengan GORM
	userModel := model.User{
		Name:     user.Name,     // Salin nama dari entity ke model
		Email:    user.Email,    // Salin email dari entity ke model
		Password: user.Password, // Salin password (sudah di-hash sebelumnya)
		Avatar:   user.Avatar,   // Salin avatar jika ada
	}

	// Buat user di database menggunakan GORM Create
	err := u.db.Create(&userModel).Error
	if err != nil {
		// Jika gagal membuat (misal: email duplikat), log error dan return nil
		log.Println("[UserRepository] Create:", err)
		return nil, err
	}

	// Konversi User model yang sudah disimpan kembali ke UserEntity
	result := &entity.UserEntity{
		ID:       userModel.ID,       // ID yang di-generate oleh database
		Name:     userModel.Name,     // Nama user
		Email:    userModel.Email,    // Email user
		Password: userModel.Password, // Password user (akan di-hide di JSON response)
		Avatar:   userModel.Avatar,   // Avatar URL
		Roles:    userModel.Roles,    // Roles (kosong karena tidak ada yang di-associate)
	}

	// Return pointer ke UserEntity baru - efisien karena tidak ada copy data yang besar
	return result, nil
}

// FindByEmail mencari user berdasarkan email dan mengembalikan pointer ke UserEntity
// Parameter: email adalah email user yang dicari
// Return: pointer ke UserEntity (dengan password untuk internal use) dan error jika ada
func (u *userRepository) FindByEmail(email string) (*entity.UserEntity, error) {
	// Deklarasi variabel user dengan tipe model.User
	var user model.User
	// Query database untuk mencari user dengan email tertentu dan preload Roles-nya
	err := u.db.Where("email = ?", email).Preload("Roles").First(&user).Error
	if err != nil {
		// Jika error (user tidak ditemukan atau error database), log error dan return nil
		log.Println("[UserRepository] FindByEmail:", err)
		return nil, err
	}

	// Konversi User model ke UserEntity dan return sebagai pointer
	return &entity.UserEntity{
		ID:       user.ID,       // ID user dari database
		Name:     user.Name,     // Nama user
		Email:    user.Email,    // Email user
		Password: user.Password, // Password user (untuk keperluan verifikasi di service)
		Avatar:   user.Avatar,   // Avatar URL
		Roles:    user.Roles,    // Role-role yang terkait dengan user
	}, nil
}

// GetRolesByIDs mengambil daftar role berdasarkan ID-ID yang diberikan
// Parameter: roleIDs adalah slice dari ID role yang ingin diambil
// Return: slice dari model.Role dan error jika ada
func (u *userRepository) GetRolesByIDs(roleIDs []uint64) ([]model.Role, error) {
	// Deklarasi slice roles untuk menampung hasil query
	var roles []model.Role
	// Query database untuk mengambil roles dengan ID yang ada di dalam roleIDs slice
	err := u.db.Where("id IN ?", roleIDs).Find(&roles).Error
	if err != nil {
		// Jika error saat query, log error dan return slice kosong dengan error
		log.Println("[UserRepository] GetRolesByIDs:", err)
		return roles, err
	}
	// Return slice roles yang berhasil diambil dari database
	return roles, nil
}
