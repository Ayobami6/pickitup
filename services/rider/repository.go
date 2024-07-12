package rider

import (
	"fmt"
	"log"
	"math"
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

    var reviewResponse []dto.ReviewResponse
    reviews, err := r.GetRiderReviews(uint(id))
    if err!= nil {
        log.Println(err)
    }

    if reviews != nil {
        for _, review := range reviews {
            reviewResponse = append(reviewResponse, dto.ReviewResponse{
                Rating: review.Rating,
                Comment: review.Comment,
            })
        }

    } else {
        reviewResponse = []dto.ReviewResponse{}

    }
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
        Reviews: reviewResponse,
        MinimumCharge: rider.MinimumCharge,
        MaximumCharge: rider.MaximumCharge,
	}
    return response, nil
}


func (r *riderRepositoryImpl) GetRiderByUserID(userID uint) (*models.Rider, error){
	var rider models.Rider
    res := r.db.Where(&models.Rider{UserID: userID}).First(&rider)
    if res.Error!= nil {
        return nil, res.Error
    }
    return &rider, nil
}

func (r *riderRepositoryImpl) GetRiderByID(id uint) (*models.Rider, error) {
	var rider models.Rider
	res := r.db.First(&rider, id)
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

func (r *riderRepositoryImpl)UpdateRating(riderID uint)(error){
	var rider models.Rider
	res := r.db.Where(&models.Rider{ID: riderID}).First(&rider)
	if res.Error!= nil {
        return res.Error
    }
	// get all ratings for the rider 
	var reviews []models.Review
    res = r.db.Where(&models.Review{RiderID: riderID}).Find(&reviews)
    if res.Error!= nil {
        return res.Error
    }

    // calculate new average rating
    var totalRating float64 = 0
    for _, review := range reviews {
        totalRating += review.Rating
    }
    newRating := totalRating / float64(len(reviews))
    // round to 1 decimal place
    newRating = math.Round(newRating*10) / 10

    // update rider rating
    rider.Rating = newRating
    res = r.db.Save(&rider)
    if res.Error!= nil {
        return res.Error
    }

    // send notification to all users who rated the rider
    // sendNotificationToRaterUsers(riderID, newRating)
    // send notification to all users who requested rides from the rider
    // sendNotification
	

    return  nil
	// 
}

func (r *riderRepositoryImpl)UpdateMinAndMaxCharge(minCharge float64, maxCharge float64, userID uint) (error) {
    res := r.db.Where(&models.Rider{UserID: userID}).Updates(&models.Rider{MinimumCharge: minCharge, MaximumCharge: maxCharge})
    if res.Error != nil {
        return res.Error
    }
    return nil
}

func (r *riderRepositoryImpl) GetRiderReviews(riderID uint) ([]models.Review, error) {
    var reviews []models.Review
    res := r.db.Where(&models.Review{RiderID: riderID}).Find(&reviews)
    if res.Error!= nil {
        return nil, res.Error
    }
    return reviews, nil
}