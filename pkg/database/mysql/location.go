package mysql

var location = map[int]struct {
	Address   string
	Latitude  float64
	Longitude float64
}{1: {
	Address:   "Jl. Veteran, Ketawanggede, Kec. Lowokwaru, Kota Malang, Jawa Timur 65145",
	Latitude:  -7.9522805705729835,
	Longitude: 112.6126959150791,
}, 2: {
	Address:   "2JX6+4W8, Jl. MT. Haryono, Ketawanggede, Kec. Lowokwaru, Kota Malang, Jawa Timur 65145",
	Latitude:  -7.952184884438989,
	Longitude: 112.61231999408585,
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
