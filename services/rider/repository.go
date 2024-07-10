package rider

import (
	"fmt"
	"log"
	"net/http"

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

func (r *riderRepositoryImpl) GetRiders(req *http.Request) (rider []dto.RiderListResponse, err error) {
	riders := []models.Rider{}
	var parsedRiders []dto.RiderListResponse

	res := r.db.Find(&riders)
	if res.Error!= nil {
        return nil, res.Error
    }
	domain := getDomainURL(req)
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
			SelfUrl: fmt.Sprintf("%s/riders/%d", domain, riders[i].ID),
		}
		parsedRiders = append(parsedRiders, rider)
		
	}
	
	return parsedRiders, nil
}

func (r *riderRepositoryImpl) GetRider(id int, req *http.Request) (dto.RiderResponse, error) {
	var rider models.Rider
    res := r.db.First(&rider, id)
	if res.Error!= nil {
        return dto.RiderResponse{}, res.Error
    }
	domain := getDomainURL(req)
	var selfUrl = fmt.Sprintf("%s/riders/%d", domain, rider.ID)
	response := dto.RiderResponse{
		ID: rider.ID,
        FirstName: rider.FirstName,
        LastName: rider.LastName,
        RiderID:  rider.RiderID,
        BikeNumber: rider.BikeNumber,
        Address: rider.Address,
        SuccessfulRides: rider.SuccessfulRides,
        Rating: rider.Rating,
        CurrentLocation: rider.CurrentLocation,
        Level: rider.Level,
        SelfUrl: selfUrl,
	}
    return response, nil
}

func (r *riderRepositoryImpl)CreateRating(Id uint)(string, error){
	// constraints only a user not rider can submit rating
	// this should be sent by and authenticated user to get their from the context
	//
    // TODO: implement rating logic
    return "Rating submitted successfully", nil
	// 
}


func (r *riderRepositoryImpl) GetRiderByID(Id int) (*models.Rider, error){
	var rider models.Rider
    res := r.db.First(&rider, Id)
    if res.Error!= nil {
        return nil, res.Error
    }
    return &rider, nil
}

// get domain function
func getDomainURL(r *http.Request) string {
    scheme := "http"
    if r.TLS != nil {
        scheme = "https"
    }
    return scheme + "://" + r.Host
}