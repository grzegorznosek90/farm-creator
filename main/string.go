package main

import (
	"log"
)

type Sstring struct {
	SMPV_Uch CmpVal
	SMPV_Ich CmpVal
	SMPV_Pch CmpVal
	SMPV_E   CmpVal

	SMPV_num         int
	SMPV_ok_num      int
	SMPV_error_num   int
	SMPV_offline_num int

	IgtIdSys        string    `json:"iGT_id_sys"`
	IgtIdFarm       string    `json:"iGT_id_farm"`
	IgtIdGroup      string    `json:"iGT_id_group"`
	IgtIdString     string    `json:"iGT_id_string"`

}


func (s *Subscriber) SendString(sample *Sstring) error {
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

func updateSubscribersStr(sample *Sstring) {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("Subscribers: %d", subs.Len())
	e := subs.Front()
	for e != nil {
		s := e.Value.(*Subscriber)
		if err := s.SendString(sample); err != nil {
			log.Printf("Subscriber failed, err: %v", err)
			subs.Remove(e)
		}
		e = e.Next()
	}
}
