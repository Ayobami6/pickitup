package rider

type RiderRepository interface {
	CreateRider(rider *Rider) error
}