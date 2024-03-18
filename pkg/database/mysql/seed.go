package mysql

import (
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"intern-bcc-2024/pkg/mail"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func generateAccountNumberType(db *gorm.DB) error {
	var accountTypes []*entity.AccountNumberType

	accountTypes = append(accountTypes,
		&entity.AccountNumberType{
			ID:   uuid.New(),
			Name: "Gopay",
		},
		&entity.AccountNumberType{
			ID:   uuid.New(),
			Name: "Shopeepay",
		},
	)

	if err := db.CreateInBatches(accountTypes, 2).Error; err != nil {
		return err
	}

	return nil
}

func generateUser(db *gorm.DB) error {
	var users []*entity.User
	var fullName string

	for i := 1; i <= 40; i++ {
		if i%2 == 0 {
			fullName = faker.FirstNameMale() + " " + faker.LastName()
		} else {
			fullName = faker.FirstNameFemale() + " " + faker.LastName()
		}

		username := strings.ToLower(strings.Replace(strings.Replace(fullName, "'", "", 1), " ", "", 1))
		hashPassword, err := bcrypt.GenerateFromPassword([]byte("12345678"), 10)
		if err != nil {
			return err
		}

		providerCode := []string{"812", "857", "859"}

		var accountType *entity.AccountNumberType
		if err = db.Order("RAND()").Limit(1).First(&accountType).Error; err != nil {
			return err
		}

		user := &entity.User{
			ID:              uuid.New(),
			Name:            fullName,
			Email:           fmt.Sprintf("%s@student.ub.ac.id", username),
			Password:        string(hashPassword), // 12345678
			Address:         location[i].Address,
			Latitude:        location[i].Latitude,
			Longitude:       location[i].Longitude,
			StatusAccount:   "active",
			AccountNumber:   fmt.Sprintf("62%s%s", providerCode[randRange(0, len(providerCode)-1)], strconv.Itoa(randRange(11111111, 99999999))),
			AccountNumberID: accountType.ID,
			UrlPhotoProfile: fmt.Sprintf("https://randomuser.me/api/portraits/men/%d.jpg", randRange(0, 99)),
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
			ID:   uuid.New(),
			Name: "Elektronik",
			Url:  "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Elektronik-cgenk5n1n80qwtzi28uh4doufb64kq.png",
		},
		&entity.Category{
			ID:   uuid.New(),
			Name: "Pakaian dan Aksesoris",
			Url:  "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Pakaian%20dan%20Aksesoris-bujd1dny936mccqqtsg5bgdbkjrd5p.png",
		},
		&entity.Category{
			ID:   uuid.New(),
			Name: "Perabotan dan Dekorasi Kamar",
			Url:  "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Perabotan%20dan%20Dekorasi%20Kamar-1mxhxw7rr7zj5zwaxl7uv3srfsa5lb.png",
		},
		&entity.Category{
			ID:   uuid.New(),
			Name: "Kendaraan",
			Url:  "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Kendaraan-0m31a5v9glcmxp5mmcyrzoj2kmg7dk.png",
		},
		&entity.Category{
			ID:   uuid.New(),
			Name: "Alat Tulis dan Perlengkapan kantor",
			Url:  "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Alat%20Tulis%20dan%20Perlengkapan%20Kantor-6mqpkyophszmho2vgs3fd3i8y6ydw1.png",
		},
		&entity.Category{
			ID:   uuid.New(),
			Name: "Perlengkapan Kuliah",
			Url:  "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Perlengkapan%20Kuliah-t3d7leouhijnr5ya6076iv5d848uaz.png",
		},
		&entity.Category{
			ID:   uuid.New(),
			Name: "Kesehatan dan Kecantikan",
			Url:  "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/kesehatan%20dan%20Kecantikan-izqrjov5jfw7aoa2hug8stgnpbuon5.png",
		},
		&entity.Category{
			ID:   uuid.New(),
			Name: "Buku dan Bahan Ajar",
			Url:  "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Buku%20dan%20Bahan%20Ajar-ich0womggrlxuskvowtvd7kpsf0c0p.png",
		},
	)

	if err := db.CreateInBatches(categories, 8).Error; err != nil {
		return err
	}

	return nil
}

func generateProduct(db *gorm.DB) error {
	var products []*entity.Product
	var medias []*entity.Media

	var categories []*entity.Category
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	for _, category := range categories {
		for _, product := range Products[category.Name] {
			var users *entity.User
			if err := db.Order("RAND()").Limit(1).First(&users).Error; err != nil {
				return err
			}

			id := uuid.New()
			p := strings.Split(product, "|")
			price, err := strconv.Atoi(p[1])
			if err != nil {
				return err
			}

			entityProduct := &entity.Product{
				ID:          id,
				UserID:      users.ID,
				CategoryID:  category.ID,
				Name:        p[0],
				Description: p[2],
				Price:       uint64(price),
				CancelCode:  mail.GenerateSixCode(),
			}

			products = append(products, entityProduct)

			for i := 3; i < len(p); i++ {
				media := &entity.Media{
					ID:        uuid.New(),
					ProductID: id,
					Url:       p[i],
				}

				medias = append(medias, media)
			}
		}
	}

	if err := db.CreateInBatches(products, len(products)).Error; err != nil {
		return err
	}

	if err := db.CreateInBatches(medias, len(medias)).Error; err != nil {
		return err
	}

	return nil
}

func SeedData(db *gorm.DB) {
	var totalAccountType int64

	if err := db.Model(&entity.AccountNumberType{}).Count(&totalAccountType).Error; err != nil {
		log.Fatalf("Error while counting account number type: %v", err)
	}

	if totalAccountType == 0 {
		if err := generateAccountNumberType(db); err != nil {
			log.Fatalf("Error while generating account number type: %v", err)

		}
	}

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

}

func randRange(min, max int) int {
	return rand.Intn(max+1-min) + min
}
