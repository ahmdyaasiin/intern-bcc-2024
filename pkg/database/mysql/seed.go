package mysql

import (
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"log"
	"math/rand"
)

func generateUser(db *gorm.DB) error {
	var users []*entity.User

	for i := 1; i <= 25; i++ {
		user := &entity.User{
			ID:              uuid.New(),
			Name:            faker.FirstNameMale() + " " + faker.LastName(),
			Email:           faker.Username() + "@student.ub.ac.id",
			Password:        "$2a$10$kmv/rfq1wIT9BgEwx9AMQu5wwq72kD1M8GEcumpd23oprBLnjuvPu", // 12345678
			Address:         "Jl. Veteran, Ketawanggede, Kec. Lowokwaru, Kota Malang, Jawa Timur 65145",
			Latitude:        faker.Latitude(),
			Longitude:       faker.Longitude(),
			StatusAccount:   "active",
			UrlPhotoProfile: fmt.Sprintf("https://randomuser.me/api/portraits/men/%d.jpg", randRange(1, 30)),
		}

		users = append(users, user)
	}

	for i := 1; i <= 25; i++ {
		user := &entity.User{
			ID:              uuid.New(),
			Name:            faker.FirstNameFemale() + " " + faker.LastName(),
			Email:           faker.Username() + "@student.ub.ac.id",
			Password:        "$2a$10$kmv/rfq1wIT9BgEwx9AMQu5wwq72kD1M8GEcumpd23oprBLnjuvPu", // 12345678
			Address:         "Jl. Veteran, Ketawanggede, Kec. Lowokwaru, Kota Malang, Jawa Timur 65145",
			Latitude:        faker.Latitude(),
			Longitude:       faker.Longitude(),
			StatusAccount:   "active",
			UrlPhotoProfile: fmt.Sprintf("https://randomuser.me/api/portraits/women/%d.jpg", randRange(1, 30)),
		}

		users = append(users, user)
	}

	if err := db.CreateInBatches(users, 50).Error; err != nil {
		return err
	}
	return nil
}

func generateCategory(db *gorm.DB) error {
	var categories []*entity.Category

	categories = append(categories,
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Elektronik",
			UrlCategory: "https://down-id.img.susercontent.com/file/dcd61dcb7c1448a132f49f938b0cb553",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Komputer & Aksesoris",
			UrlCategory: "https://down-id.img.susercontent.com/file/id-50009109-0bd6a9ebd0f2ae9b7e8b9ce7d89897d6",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Perawatan & Kecantikan",
			UrlCategory: "https://down-id.img.susercontent.com/file/2715b985ae706a4c39a486f83da93c4b",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Pakaian Pria",
			UrlCategory: "https://down-id.img.susercontent.com/file/04dba508f1ad19629518defb94999ef9",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Pakaian Wanita",
			UrlCategory: "https://down-id.img.susercontent.com/file/6d63cca7351ba54a2e21c6be1721fa3a",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Sepatu Pria",
			UrlCategory: "https://down-id.img.susercontent.com/file/3c8ff51aab1692a80c5883972a679168",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Jam Tangan",
			UrlCategory: "https://down-id.img.susercontent.com/file/2bdf8cf99543342d4ebd8e1bdb576f80",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Buku & Alat Tulis",
			UrlCategory: "https://down-id.img.susercontent.com/file/998c7682fd5e7a3563b2ad00aaa4e6f3",
		},
	)

	if err := db.CreateInBatches(categories, 8).Error; err != nil {
		return err
	}

	return nil
}

func generateProduct(db *gorm.DB) error {
	var products []*entity.Product

	var categories []*entity.Category
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	for _, category := range categories {
		for i := 1; i <= randRange(1, 50); i++ {
			var users *entity.User
			if err := db.Order("RAND()").Limit(1).First(&users).Error; err != nil {
				return err
			}

			product := &entity.Product{
				ID:          uuid.New(),
				UserID:      users.ID,
				CategoryID:  category.ID,
				Name:        fmt.Sprintf("%s %d", category.Name, i),
				Description: faker.Paragraph(),
				Price:       uint64(randRange(10000, 150000)),
			}

			products = append(products, product)
		}
	}

	if err := db.CreateInBatches(products, len(products)).Error; err != nil {
		return err
	}

	return nil
}

func generateMedia(db *gorm.DB) error {
	var medias []*entity.Media

	var products []*entity.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	for _, product := range products {
		for i := 0; i <= randRange(1, 5); i++ {
			media := &entity.Media{
				ID:        uuid.New(),
				ProductID: product.ID,
				UrlMedia:  fmt.Sprintf("https://randomuser.me/api/portraits/lego/%d.jpg", randRange(1, 9)),
			}

			medias = append(medias, media)
		}
	}

	if err := db.CreateInBatches(medias, len(medias)).Error; err != nil {
		return err
	}

	return nil
}

func SeedData(db *gorm.DB) {
	var totalUser int64
	if err := db.Model(&entity.User{}).Count(&totalUser).Error; err != nil {
		log.Fatalf("Error while counting user: %v", err)
	}

	if totalUser == 0 {
		if err := generateUser(db); err != nil {
			log.Fatalf("Error while generating user: %v", err)

		}
	}

	var totalCategory int64
	if err := db.Model(&entity.Category{}).Count(&totalCategory).Error; err != nil {
		log.Fatalf("Error while counting category: %v", err)
	}

	if totalUser == 0 {
		if err := generateCategory(db); err != nil {
			log.Fatalf("Error while generating category: %v", err)

		}
	}

	var totalProduct int64
	if err := db.Model(&entity.Product{}).Count(&totalProduct).Error; err != nil {
		log.Fatalf("Error while counting product: %v", err)
	}

	if totalProduct == 0 {
		if err := generateProduct(db); err != nil {
			log.Fatalf("Error while generating product: %v", err)

		}
	}

	var totalMedia int64
	if err := db.Model(&entity.Media{}).Count(&totalMedia).Error; err != nil {
		log.Fatalf("Error while counting media: %v", err)
	}

	if totalMedia == 0 {
		if err := generateMedia(db); err != nil {
			log.Fatalf("Error while generating media: %v", err)

		}
	}

}

func randRange(min, max int) int {
	return rand.Intn(max+1-min) + min
}
