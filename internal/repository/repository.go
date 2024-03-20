package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository        IUserRepository
	OtpRepository         IOtpRepository
	TokenRepository       ITokenRepository
	SessionRepository     ISessionRepository
	ProductRepository     IProductRepository
	CategoryRepository    ICategoryRepository
	MediaRepository       IMediaRepository
	TransactionRepository ITransactionRepository
	AccountRepository     IAccountRepository
}

func NewRepository(db *gorm.DB) *Repository {
	ur := NewUserRepository(db)
	or := NewOtpRepository(db)
	tr := NewTokenRepository(db)
	sr := NewSessionRepository(db)
	cr := NewCategoryRepository(db)
	pr := NewProductRepository(db, cr)
	mr := NewMediaRepository(db)
	trr := NewTransactionRepository(db)
	ar := NewAccountRepository(db)

	return &Repository{
		OtpRepository:         or,
		UserRepository:        ur,
		TokenRepository:       tr,
		SessionRepository:     sr,
		ProductRepository:     pr,
		CategoryRepository:    cr,
		MediaRepository:       mr,
		TransactionRepository: trr,
		AccountRepository:     ar,
	}
}
