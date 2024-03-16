package mysql

import (
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"intern-bcc-2024/entity"
	"log"
	"math/rand"
	"strings"
)

var location = map[int]struct {
	Address   string
	Latitude  float64
	Longitude float64
}{1: {
	Address:   "Jl. MT. Haryono XI No.363, Dinoyo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65144",
	Latitude:  -7.942065681365051,
	Longitude: 112.61062437613042,
}, 2: {
	Address:   "Jl. MT. Haryono Gg. 4 No.816, Dinoyo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65144",
	Latitude:  -7.942894502907431,
	Longitude: 112.60956222141232,
}, 3: {
	Address:   "Jl. Gajayana Gg. 1 No.743A, Dinoyo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65144",
	Latitude:  -7.943456347355645,
	Longitude: 112.60821171081696,
}, 4: {
	Address:   "Gg. Tower, Dinoyo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65144",
	Latitude:  -7.945411506819059,
	Longitude: 112.6072032002765,
}, 5: {
	Address:   "Jl. MT. Haryono Gg. 4 No.816, Dinoyo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65144",
	Latitude:  -7.942903800602881,
	Longitude: 112.6096064595539,
}, 6: {
	Address:   "Jl. MT. Haryono XI No.363, Dinoyo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65144",
	Latitude:  -7.9419793457188925,
	Longitude: 112.61060424125883,
}, 7: {
	Address:   "Jl. MT. Haryono Jl. Mt. Haryono Gg. 5 No.277, Dinoyo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65145",
	Latitude:  -7.944072648681541,
	Longitude: 112.61165566722771,
}, 8: {
	Address:   "Jl. Kendalsari Bar. I B No.10, Tulusrejo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65141",
	Latitude:  -7.944590421047384,
	Longitude: 112.62242513963677,
}, 9: {
	Address:   "KOST PUTRI .Bu tatik, Jl. Veteran Dalam No.7b, Sumbersari, Kec. Lowokwaru, Kota Malang, Jawa Timur 65145",
	Latitude:  -7.957048630211308,
	Longitude: 112.61595652171471,
}, 10: {
	Address:   "Jl. Bareng Kulon VI No.963, Bareng, Kec. Klojen, Kota Malang, Jawa Timur 65116",
	Latitude:  -7.975334777038475,
	Longitude: 112.61747471858945,
}, 11: {
	Address:   "Jl. Bareng Tengah V No.841-842, Bareng, Kec. Klojen, Kota Malang, Jawa Timur 65116",
	Latitude:  -7.976992283876288,
	Longitude: 112.62052170788189,
}, 12: {
	Address:   "Jl. Terusan Ijen No.859a, Bareng, Kec. Klojen, Kota Malang, Jawa Timur 65116",
	Latitude:  -7.9784672936069585,
	Longitude: 112.61960302850574,
}, 13: {
	Address:   "Jl. Sumbersari Gg. 4 No.251d, Sumbersari, Kec. Lowokwaru, Kota Malang, Jawa Timur 65145",
	Latitude:  -7.954803979188259,
	Longitude: 112.61036940137853,
}, 14: {
	Address:   "Jl. Garbis No.2, Bareng, Kec. Klojen, Kota Malang, Jawa Timur 65116",
	Latitude:  -7.977829973215653,
	Longitude: 112.61333689089392,
}, 15: {
	Address:   "Jl. Pisang Candi Barat No.18k, RT.03/RW.04, Pisang Candi, Sukun, Malang City, East Java 65146",
	Latitude:  -7.975850919261359,
	Longitude: 112.61066108073118,
}, 16: {
	Address:   "Jl. Pisang Candi Barat No.77, Pisang Candi, Kec. Sukun, Kota Malang, Jawa Timur 65146",
	Latitude:  -7.97221144764025,
	Longitude: 112.60674897863714,
}, 17: {
	Address:   "Jl. Pisang Agung No.48, Pisang Candi, Kec. Sukun, Kota Malang, Jawa Timur 65146",
	Latitude:  -7.973402246038809,
	Longitude: 112.60559736420383,
}, 18: {
	Address:   "Merjosari, Kec. Lowokwaru, Kota Malang, Jawa Timur 65144",
	Latitude:  -7.944111300516476,
	Longitude: 112.59440141448498,
}, 19: {
	Address:   "Jl. Kanjuruhan Asri No.5, Tlogomas, Kec. Lowokwaru, Kota Malang, Jawa Timur 65144",
	Latitude:  -7.938384357424094,
	Longitude: 112.59491813910834,
}, 20: {
	Address:   "Jl. Joyo Agung No.184a, Merjosari, Kec. Lowokwaru, Kota Malang, Jawa Timur 65144",
	Latitude:  -7.937848213960435,
	Longitude: 112.58574012518481,
}, 21: {
	Address:   "tirtorahayu gg XI A nomor D, 1, Dusun Bend., Landungsari, Kec. Dau, Kabupaten Malang, Jawa Timur 65151",
	Latitude:  -7.930195908414379,
	Longitude: 112.59427838493384,
}, 22: {
	Address:   "Gg. X No.15, Dusun Rambaan, Landungsari, Kec. Dau, Kabupaten Malang, Jawa Timur 65151",
	Latitude:  -7.926613410807994,
	Longitude: 112.59531183418055,
}, 23: {
	Address:   "Jl. Tirto Taruno Gg. I No.2, Dusun Bend., Landungsari, Kec. Dau, Kabupaten Malang, Jawa Timur 65151",
	Latitude:  -7.928002546227752,
	Longitude: 112.59280202862845,
}, 24: {
	Address:   "Jl. Tirto Taruno.6C/3B, Jl. Tirto Taruno No.6C, Dusun Bend., Landungsari, Kec. Dau, Kabupaten Malang, Jawa Timur 65151",
	Latitude:  -7.92866055610922,
	Longitude: 112.59255596928398,
}, 25: {
	Address:   "Gg. II No.16, RT.01/RW.07. Klandungan, Dusun Klandungan, Landungsari, Kec. Dau, Kabupaten Malang, Jawa Timur 65151",
	Latitude:  -7.9283437366680305,
	Longitude: 112.59154712597172,
}, 26: {
	Address:   "Jl. Semangka No.12, Dermo, Mulyoagung, Kec. Dau, Kabupaten Malang, Jawa Timur 65151",
	Latitude:  -7.926978973209289,
	Longitude: 112.59233451587396,
}, 27: {
	Address:   "Jl. Bukit Cemara Tujuh Blk. GG No.39, Dermo, Mulyoagung, Kec. Dau, Kabupaten Malang, Jawa Timur 65151",
	Latitude:  -7.925151157915064,
	Longitude: 112.59117803695503,
}, 28: {
	Address:   "Dekat UMM kampus 3, pascarasjana UIN,dan unisma, Perumahan, Jl. Perumahan Taman Embrong Anyar II blok e.1, Jetis, Mulyoagung, Kec. Dau, Kabupaten Malang, Jawa Timur 65151",
	Latitude:  -7.921422389483643,
	Longitude: 112.59026761736729,
}, 29: {
	Address:   "Jl. Bukit Cemara Tujuh No.3 Blok D, Jetis, Mulyoagung, Kec. Dau, Kota Malang, Jawa Timur 65151",
	Latitude:  -7.921812327626196,
	Longitude: 112.59230990992629,
}, 30: {
	Address:   "Blok 03.77, Tlogomas, Lowokwaru, Malang City, East Java 65151",
	Latitude:  -7.921950006367093,
	Longitude: 112.59467588231911,
}, 31: {
	Address:   "Jl. Bukit Cemara Tujuh No.mor 9 Blok 1, Tlogomas, Kec. Lowokwaru, Kota Malang, Jawa Timur 65144",
	Latitude:  -7.922030652403649,
	Longitude: 112.59530691113017,
}, 32: {
	Address:   "Jl. Dieng Atas Dsn. Kunci, RT.01/RW.03, Kunci, Kalisongo, Kec. Dau, Kabupaten Malang, Jawa Timur 65151",
	Latitude:  -7.967100268345762,
	Longitude: 112.59187906887661,
}, 33: {
	Address:   "Jl. Raya Candi VC, Karangbesuki, Kec. Sukun, Kota Malang, Jawa Timur 65149",
	Latitude:  -7.955709810511595,
	Longitude: 112.60020464600802,
}, 34: {
	Address:   "Jl. Sunan Muria 6, Dinoyo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65149",
	Latitude:  -7.95363385624499,
	Longitude: 112.60459949463291,
}, 35: {
	Address:   "Jalan Simpang Sunan Kalijaga No. 16B, Lowokwaru, Dinoyo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65145",
	Latitude:  -7.95229940091418,
	Longitude: 112.60385445353413,
}, 36: {
	Address:   "1A No, Jl. Simpang Sunan Kalijaga No.1, Dinoyo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65149",
	Latitude:  -7.9523307999130175,
	Longitude: 112.60567742643543,
}, 37: {
	Address:   "2JW4+MHX, Jl. Sunan Muria I, Dinoyo, Kec. Lowokwaru, Kota Malang, Jawa Timur 65149",
	Latitude:  -7.953052976222885,
	Longitude: 112.60643831947247,
}, 38: {
	Address:   "Jl. Sumbersari III No.206, Sumbersari, Kec. Lowokwaru, Kota Malang, Jawa Timur 65145",
	Latitude:  -7.954151937746107,
	Longitude: 112.6094977437259,
}, 39: {
	Address:   "2JV5+9VX, Sumbersari, Kec. Lowokwaru, Kota Malang, Jawa Timur 65145",
	Latitude:  -7.956381250618374,
	Longitude: 112.60979893055308,
}, 40: {
	Address:   "Jl. Bend. Bening No.11, Sumbersari, Kec. Lowokwaru, Kota Malang, Jawa Timur 65145",
	Latitude:  -7.9585948521088685,
	Longitude: 112.60967211504689,
},
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

		username := strings.ToLower(strings.Replace(fullName, " ", "", 1))
		hashPassword, err := bcrypt.GenerateFromPassword([]byte("12345678"), 10)
		if err != nil {
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
			ID:          uuid.New(),
			Name:        "Elektronik",
			UrlCategory: "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Elektronik-cgenk5n1n80qwtzi28uh4doufb64kq.png",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Pakaian dan Aksesoris",
			UrlCategory: "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Pakaian%20dan%20Aksesoris-bujd1dny936mccqqtsg5bgdbkjrd5p.png",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Perabotan dan Dekorasi Kamar",
			UrlCategory: "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Perabotan%20dan%20Dekorasi%20Kamar-1mxhxw7rr7zj5zwaxl7uv3srfsa5lb.png",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Kendaraan",
			UrlCategory: "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Kendaraan-0m31a5v9glcmxp5mmcyrzoj2kmg7dk.png",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Alat Tulis dan Perlengkapan kantor",
			UrlCategory: "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Alat%20Tulis%20dan%20Perlengkapan%20Kantor-6mqpkyophszmho2vgs3fd3i8y6ydw1.png",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Perlengkapan Kuliah",
			UrlCategory: "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Perlengkapan%20Kuliah-t3d7leouhijnr5ya6076iv5d848uaz.png",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Kesehatan dan Kecantikan",
			UrlCategory: "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/kesehatan%20dan%20Kecantikan-izqrjov5jfw7aoa2hug8stgnpbuon5.png",
		},
		&entity.Category{
			ID:          uuid.New(),
			Name:        "Buku dan Bahan Ajar",
			UrlCategory: "https://hpmvwwpbkmlzqfmmbjbd.supabase.co/storage/v1/object/public/intern-bcc-2024/Categories/Buku%20dan%20Bahan%20Ajar-ich0womggrlxuskvowtvd7kpsf0c0p.png",
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
