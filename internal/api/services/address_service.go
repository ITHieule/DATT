package services

import (
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"
)

type AddressService struct {
	*BaseService
}

var Address = &AddressService{}

func (s *AddressService) GetAddressByUserID(req *types.ShippingAddress) ([]types.ShippingAddress, error) {
	var addresses []types.ShippingAddress

	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	query := `SELECT * FROM shipping_addresses WHERE user_id = ?`

	err = db.Raw(query,
		req.UserID,
	).Scan(&addresses).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return addresses, nil
}

func (s *AddressService) CreateAddressByUserID(requestParams *types.ShippingAddress) (*types.ShippingAddress, error) {
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	if err := db.Create(requestParams).Error; err != nil {
		fmt.Println("Insert execution error:", err)
		return nil, err
	}

	return requestParams, nil
}

func (s *AddressService) UpdateAddressByUserID(requestParams *types.ShippingAddress) error {
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	query := `
        UPDATE shipping_addresses 
        SET province = ?, district = ?, ward = ?, postal_code = ?, latitude = ?, longitude = ?, is_default = ? 
        WHERE user_id = ? AND id = ?
    `
	err = db.Exec(query,
		requestParams.Province,
		requestParams.District,
		requestParams.Ward,
		requestParams.PostalCode,
		requestParams.Latitude,
		requestParams.Longitude,
		requestParams.IsDefault,
		requestParams.UserID,
		requestParams.ID,
	).Error
	if err != nil {
		fmt.Println("Update execution error:", err)
		return err
	}

	return nil
}

func (s *AddressService) DeleteAddressByUserID(requestParams *types.ShippingAddress) error {
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	var exists bool
	err = db.Raw("SELECT EXISTS(SELECT 1 FROM shipping_addresses WHERE id = ? AND user_id = ?)", requestParams.ID, requestParams.UserID).Scan(&exists).Error
	if err != nil {
		fmt.Println("Error checking address existence:", err)
		return err
	}

	if !exists {
		return fmt.Errorf("address with id %d and user_id %d not found", requestParams.ID, requestParams.UserID)
	}
	query := `
        DELETE FROM shipping_addresses WHERE id = ? AND user_id = ?
    `
	result := db.Exec(query, requestParams.ID, requestParams.UserID)
	if result.Error != nil {
		fmt.Println("Delete execution error:", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no address deleted for id %d and user_id %d", requestParams.ID, requestParams.UserID)
	}

	return nil
}
