package main

import (
	"encoding/json"
	"time"
)

func Download(rawMsg []byte) error {
	downloadReq := new(DownloadRequest)
	err := json.Unmarshal(rawMsg, downloadReq)
	if err != nil {
		Error("JSON解析失败", "err", err)
		return err
	}
	
	Info("下载文件", "file_id", downloadReq.FileID, "media_len", len(downloadReq.Media), "cdn_url", downloadReq.CDNURL)
	if downloadReqInter, ok := userID2FileMsgMap.Load(downloadReq.CDNURL); ok {
		beforeDownloadReq := downloadReqInter.(*DownloadRequest)
		if beforeDownloadReq.FilePath != "" {
			return nil
		}
		if time.Now().UnixMilli()-beforeDownloadReq.LastAppendTime > 10000000 {
			beforeDownloadReq.Media = downloadReq.Media
		} else {
			beforeDownloadReq.Media = append(beforeDownloadReq.Media, downloadReq.Media...)
		}
		
		beforeDownloadReq.LastAppendTime = time.Now().UnixMilli()
	} else {
		downloadReq.LastAppendTime = time.Now().UnixMilli()
		userID2FileMsgMap.Store(downloadReq.CDNURL, downloadReq)
	}
	
	return nil
}
