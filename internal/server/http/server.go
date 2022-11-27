package internalhttp

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	ms "github.com/PalPalych7/OtusProjectWork/internal/mainstructs"
)

type Server struct {
	myCtx     context.Context
	myStorage ms.Storage
	myLogger  ms.Logger
	HTTPConf  ms.HTTPConf
	myHTTP    http.Server
}

func NewServer(ctx context.Context, app ms.Storage, httpConf ms.HTTPConf, myLogger ms.Logger) *Server {
	return &Server{myCtx: ctx, myStorage: app, myLogger: myLogger, HTTPConf: httpConf}
}

func getBodyRaw(reqBody io.ReadCloser) []byte {
	raw, err := ioutil.ReadAll(reqBody)
	if err != nil {
		return nil
	}
	defer reqBody.Close()
	return raw
}

func (s *Server) Serve() error {
	s.myHTTP.Addr = ":" + s.HTTPConf.Port
	mux := http.NewServeMux()
	mux.HandleFunc("/AddBannerSlot", s.AddBannerSlot)
	mux.HandleFunc("/GetBannerForSlot", s.GetBannerForSlot)
	mux.HandleFunc("/BannerClick", s.BannerClick)
	mux.HandleFunc("/DelBannerSlot", s.DelBannerSlot)

	server := &http.Server{
		Addr:              s.myHTTP.Addr,
		ReadHeaderTimeout: time.Second * time.Duration(s.HTTPConf.ReadHeaderTimeout),
		Handler:           s.loggingMiddleware(mux),
	}

	err := server.ListenAndServe()
	if err != nil {
		s.myLogger.Error(err)
	}
	return err
}

func (s *Server) Stop() error {
	err := s.myHTTP.Shutdown(s.myCtx)
	return err
}

func (s *Server) AddBannerSlot(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("AddBannerSlot")
	myCtx, cancel := context.WithTimeout(s.myCtx, time.Second*time.Duration(s.HTTPConf.TimeOutSec))
	defer cancel()
	myRaw := getBodyRaw(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	mySB := ms.SlotBanner{}
	if err := json.Unmarshal(myRaw, &mySB); err != nil {
		s.myLogger.Error("Error json.Unmarshal - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.myStorage.AddBannerSlot(myCtx, mySB.SlotID, mySB.BannerID)
	s.myLogger.Debug("result:", myErr)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) BannerClick(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("BannerClick")
	myCtx, cancel := context.WithTimeout(s.myCtx, time.Second*time.Duration(s.HTTPConf.TimeOutSec))
	defer cancel()

	myRaw := getBodyRaw(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myFBC := ms.ForBannerClick{}
	if err := json.Unmarshal(myRaw, &myFBC); err != nil {
		s.myLogger.Debug(myRaw)
		s.myLogger.Error("Error json.Unmarshal - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.myStorage.BannerClick(myCtx, myFBC.SlotID, myFBC.BannerID, myFBC.SocGroupID)
	s.myLogger.Debug("result:", myErr)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) DelBannerSlot(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("DelBannerSlot")
	myCtx, cancel := context.WithTimeout(s.myCtx, time.Second*time.Duration(s.HTTPConf.TimeOutSec))
	defer cancel()
	myRaw1 := getBodyRaw(req.Body)
	if myRaw1 == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	mySB := ms.SlotBanner{}
	if err1 := json.Unmarshal(myRaw1, &mySB); err1 != nil {
		s.myLogger.Debug(myRaw1)
		s.myLogger.Error("Error json.Unmarshal - " + err1.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.myStorage.DelBannerSlot(myCtx, mySB.SlotID, mySB.BannerID)
	if myErr != nil {
		s.myLogger.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) GetBannerForSlot(rw http.ResponseWriter, req *http.Request) {
	s.myLogger.Info("GetBannerForSlot")
	myCtx, cancel := context.WithTimeout(s.myCtx, time.Second*time.Duration(s.HTTPConf.TimeOutSec))
	defer cancel()
	myRaw := getBodyRaw(req.Body)
	if myRaw == nil {
		s.myLogger.Error("Request body processing error")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myFGBS := ms.ForGetBanner{}
	if err := json.Unmarshal(myRaw, &myFGBS); err != nil {
		s.myLogger.Debug(myRaw)
		s.myLogger.Error("Error json.Unmarshal - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myGetBannerStruct, myEr := s.myStorage.GetBannerForSlot(myCtx, myFGBS.SlotID, myFGBS.SocGroupID)
	s.myLogger.Debug(myGetBannerStruct, myEr)
	if myEr != nil {
		s.myLogger.Error(myEr)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rawResp, err3 := json.Marshal(myGetBannerStruct)
		if err3 == nil {
			rw.Write(rawResp)
		}
	}
}
