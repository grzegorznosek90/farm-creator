package main


type Measurements struct{
	E    float64 `json:"E"`  //energia z okresu Ts
	Ich  CmpVal  `json:"Ich"`
	Kch  CmpVal  `json:"Kch"`
	Pch  CmpVal  `json:"Pch"`
	Temp CmpVal  `json:"Temp"`
	Uch  CmpVal  `json:"Uch"`
}

type SMPV struct {
	SMPV_num int `json:"SMPV_num"`
	String_num int `json:"String_num"`
	Inv_num int `json:"Inv_num"`
	Measurements Measurements `json:"measurements"`
}


type SMPVS struct{
	IP_addr []int  `json:"IP_addr"`
	MAC_addr []int  `json:"MAC_addr"`
	Data []*SMPV `json:"data"`
}

type CmpVal struct {
	V float64 `json:"v"`
	Q int     `json:"q"`
}


type Sample struct {
	Timestamp int64   `json:"timestamp"`
	Value     string `json:"value"`
}
