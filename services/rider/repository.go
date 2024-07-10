package rider

import (
	"log"

	"github.com/Ayobami6/pickitup/models"
	"github.com/Ayobami6/pickitup/services/rider/dto"
	"gorm.io/gorm"
)

type riderRepositoryImpl struct {
	db *gorm.DB
}

func NewRiderRepositoryImpl(db *gorm.DB) *riderRepositoryImpl {
	err := db.AutoMigrate(&models.Rider{}, &models.Review{})
	if err!= nil {
        log.Fatal(err)
    }
    return &riderRepositoryImpl{db: db}
}

// overide the interface methods

func (r *riderRepositoryImpl) CreateRider(rider *models.Rider) error {
    return r.db.Create(rider).Error
}

func (r *riderRepositoryImpl) GetRiders() (rider []dto.RiderListResponse, err error) {
	riders := []models.Rider{}
	var parsedRiders []dto.RiderListResponse

	res := r.db.Find(&riders)
	if res.Error!= nil {
        return nil, res.Error
    }
	for i := range riders {
		rider := dto.RiderListResponse{
			ID: riders[i].ID,
			FirstName: riders[i].FirstName,
			LastName: riders[i].LastName,
            RiderID:  riders[i].RiderID,
			BikeNumber: riders[i].BikeNumber,
			Address: riders[i].Address,
            SuccessfulRides: riders[i].SuccessfulRides,
			Rating: riders[i].Rating,
			CurrentLocation: riders[i].CurrentLocation,
			Level: riders[i].Level,
		}
		parsedRiders = append(parsedRiders, rider)
		
	}
	
	return parsedRiders, nil
}

func (r *riderRepositoryImpl) GetRider(id int) (models.Rider, error) {
	var rider models.Rider
    res := r.db.First(&rider, id)
    if res.Error!= nil {
        return rider, res.Error
    }
    return rider, nil
}

func (r *riderRepositoryImpl)CreateRating(Id uint)(string, error){
	// constraints only a user not rider can submit rating
	// this should be sent by and authenticated user to get their from the context
	//
    // TODO: implement rating logic
    return "Rating submitted successfully", nil
	// 
}
