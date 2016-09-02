package main

import (
	"log"
)

type Farm struct {
	SMPV_Uch             CmpVal
	SMPV_Ich             CmpVal
	SMPV_Pch             CmpVal
	SMPV_E               CmpVal
	SMPV_num             int
	SMPV_ok_num          int
	SMPV_error_num       int
	SMPV_offline_num     int
	STRING_L_num         int
	STRING_L_ok_num      int
	STRING_L_error_num   int
	STRING_L_offline_num int
	INV_L_num            int
	INV_L_ok_num         int
	INV_L_error_num      int
	INV_L_offline_num    int
	IgtIdFarm        string    `json:"iGT_id_sys"`
}



	func (s *Subscriber) SendFarm(sample *Farm) error {
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

	func updateSubscribersFarm(sample *Farm) {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("Subscribers: %d", subs.Len())
	e := subs.Front()
	for e != nil {
		s := e.Value.(*Subscriber)
		if err := s.SendFarm(sample); err != nil {
			log.Printf("Subscriber failed, err: %v", err)
			subs.Remove(e)
		}
		e = e.Next()
	}
}

func updateFarms(){
	farm := farmsToSend.Front()
	for farm != nil {
		frm := farm.Value.(*Farm)

			frm.SMPV_Ich.V = frm.SMPV_Ich.V/float64(frm.SMPV_num)

		farm = farm.Next()
	}
}
