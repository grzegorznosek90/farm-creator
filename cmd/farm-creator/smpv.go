package main

import (
	"fmt"
	"time"
	"log"
)

type SMPVToSend struct {
	Timestamp time.Time `json:"timestamp"`
	E    float64 `json:"E"`
	Ich  CmpVal  `json:"Ich"`
	Kch  CmpVal  `json:"Kch"`
	Pch  CmpVal  `json:"Pch"`
	Temp CmpVal  `json:"Temp"`
	Uch  CmpVal  `json:"Uch"`
	IgtIdSmpv      string    `json:"iGT_id_SMPV"`
	IgtIdString    string    `json:"iGT_id_string"`
	IgtIdGroup     string    `json:"iGT_id_inv"`
	IgtIdFarm      string    `json:"iGT_id_farm"`
}

func (s *Subscriber) SendPv(sample *SMPVToSend) error {
	if err := s.enc.Encode(sample); err != nil {
		return err
	}
	if err := s.w.WriteByte(0); err != nil {
		return err
	}
	if err := s.w.Flush(); err != nil {
		return err
	}
	return nil
}

func updateSubscribersPv(sample *SMPVToSend) {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("Subscribers: %d", subsOut.Len())
	e := subsOut.Front()
	for e != nil {
		s := e.Value.(*Subscriber)
		if err := s.SendPv(sample); err != nil {
			log.Printf("Subscriber failed, err: %v", err)
			subsOut.Remove(e)
		}
		e = e.Next()
	}
}

func buildStructure(smpv *SMPV) {


		smpvName := "SMPV-" +
			fmt.Sprintf("%06d", smpv.Inv_num) +
			"-" +
			fmt.Sprintf("%03d", smpv.String_num) +
			"-" +
			fmt.Sprintf("%03d", smpv.SMPV_num)
		stringName := "STRING-" +
			fmt.Sprintf("%06d", smpv.Inv_num) +
			"-" +
			fmt.Sprintf("%03d", smpv.String_num)

		invName := "INV-" + fmt.Sprintf("%06d", smpv.Inv_num)
		farmName := "Farm-1";

		toSend := new(SMPVToSend)
		toSend.IgtIdSmpv = smpvName;
		toSend.IgtIdString = stringName;
		toSend.IgtIdGroup = invName;
		toSend.IgtIdFarm = farmName;

		toSend.E = smpv.Measurements.E
		toSend.Ich = smpv.Measurements.Ich
		toSend.Uch = smpv.Measurements.Uch
		toSend.Kch = smpv.Measurements.Kch
		toSend.Pch = smpv.Measurements.Pch
		toSend.Temp = smpv.Measurements.Temp


		smpvToSend.PushBack(toSend)

		toSendString := new(Sstring)
		string := stringsToSend.Front()
		for string != nil {
			str := string.Value.(*Sstring)
			if invName == str.IgtIdGroup && stringName == str.IgtIdString {
				toSendString = str
				stringsToSend.Remove(string)
				break
			}
			string = string.Next()
		}

		toSendString.IgtIdString = stringName;
		toSendString.IgtIdGroup = invName;
		toSendString.IgtIdFarm = farmName;

		toSendString.SMPV_Ich.V += smpv.Measurements.Ich.V
		toSendString.SMPV_Pch.V += smpv.Measurements.Pch.V
		toSendString.SMPV_Uch.V += smpv.Measurements.Uch.V
		toSendString.SMPV_E.V += smpv.Measurements.E
		if smpv.Measurements.Uch.Q ==192 && smpv.Measurements.Ich.Q==192 {
					toSendString.SMPV_ok_num++
				} else {
					toSendString.SMPV_error_num++
				}
		toSendString.SMPV_num = toSendString.SMPV_ok_num + toSendString.SMPV_error_num
		stringsToSend.PushBack(toSendString)



		toSendInverter := new(Inverter)
		inverter := invertersToSend.Front()
		for inverter != nil {
			inv := inverter.Value.(*Inverter)
			if invName == inv.IgtIdGroup {
				toSendInverter = inv
				invertersToSend.Remove(inverter)
				break
			}
			inverter = inverter.Next()
		}

		toSendInverter.IgtIdGroup = invName;
		toSendInverter.IgtIdFarm = farmName;

		toSendInverter.SMPV_Ich.V += smpv.Measurements.Ich.V
		toSendInverter.SMPV_Pch.V += smpv.Measurements.Pch.V
		toSendInverter.SMPV_Uch.V += smpv.Measurements.Uch.V
		toSendInverter.SMPV_E.V += smpv.Measurements.E
		if smpv.Measurements.Uch.Q ==192 && smpv.Measurements.Ich.Q==192 {
					toSendInverter.SMPV_ok_num++
				} else {
					toSendInverter.SMPV_error_num++
				}

		toSendInverter.SMPV_num = toSendInverter.SMPV_ok_num + toSendInverter.SMPV_error_num
	  invertersToSend.PushBack(toSendInverter)

		toSendFarm := new(Farm)
		farm := farmsToSend.Front()
		for farm != nil {
			farmObj := farm.Value.(*Farm)
			if farmName == farmObj.IgtIdFarm {
				toSendFarm = farmObj
				farmsToSend.Remove(farm)
				break
			}
			farm = farm.Next()
		}

		toSendFarm.IgtIdFarm = farmName;

		toSendFarm.SMPV_Ich.V += smpv.Measurements.Ich.V
		toSendFarm.SMPV_Pch.V += smpv.Measurements.Pch.V
		toSendFarm.SMPV_Uch.V += smpv.Measurements.Uch.V
		toSendFarm.SMPV_E.V += smpv.Measurements.E
			if smpv.Measurements.Uch.Q ==192 && smpv.Measurements.Ich.Q==192 {
					toSendFarm.SMPV_ok_num++
				} else {
					toSendFarm.SMPV_error_num++
				}

		toSendFarm.SMPV_num = toSendFarm.SMPV_ok_num + toSendFarm.SMPV_error_num
		farmsToSend.PushBack(toSendFarm)


}
