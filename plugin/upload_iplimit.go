package plugin

import (
	"fqimg/lib"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

var (
	uploadAllowedIps []interface{}
	uploadDenyIps    []interface{}
	ipLookups        []string
)

func Plugin_upload_iplimit(h http.HandlerFunc) http.HandlerFunc {
	parseIp(lib.Conf.UploadIpAllowed)
	parseIp(lib.Conf.UploadIpDeny)

	return func(writer http.ResponseWriter, request *http.Request) {
		if len(uploadAllowedIps) > 0 || len(uploadDenyIps) > 0 {
			if ipPass(request) == false {
				writer.WriteHeader(http.StatusInternalServerError)
				io.WriteString(writer, "your ip address is forbiden to upload")
				return
			}
		} else {
			h.ServeHTTP(writer, request)
		}
	}
}

func findIp(req *http.Request) net.IP {
	if len(lib.Conf.UploadIpLookups) > 0 {
		var ip string
		for _, value := range lib.Conf.UploadIpLookups {
			if value == "RemoteAddr" {
				adds := strings.Split(req.RemoteAddr, ":")
				ip = adds[0]
			} else if value == "X-Forwarded-For" {
				adds := strings.Split(req.Header.Get("X-Forwarded-For"), ",")
				ip = adds[0]
			} else {
				ip = req.Header.Get(value)
			}
			if ip != "" {
				break
			}
		}
		ipout := net.ParseIP(ip)
		if ipout == nil {
			log.Fatal("uploadIpLookups  configure error that no client ip can be found, and the default ip returned")
			return __getip(req)
		}

		return ipout
	} else {
		return __getip(req)
	}
}

func __getip(req *http.Request) net.IP {
	adds := strings.Split(req.RemoteAddr, ":")
	ipout := net.ParseIP(adds[0])
	if ipout == nil {
		log.Panic("ip parse error")
	}
	return ipout
}

func ipPass(req *http.Request) bool {
	//adds := strings.Split(req.RemoteAddr, ":")
	//ip := net.ParseIP(adds[0]);
	//log.Print("req remote ip :", ip)
	ip := findIp(req)
	if len(uploadAllowedIps) > 0 {

		for i := 0; i < len(uploadAllowedIps); i++ {
			switch uploadAllowedIps[i].(type) {
			case net.IPNet:
				ipnet := uploadAllowedIps[i].(net.IPNet)
				if ipnet.Contains(ip) {
					//log.Println(ip.String()," match allow: ",ipnet.String())
					return true
				}
			case net.IP:
				thisip := uploadAllowedIps[i].(net.IP)
				if thisip.Equal(ip) {
					//log.Println(ip.String()," match allow: ",thisip.String())
					return true
				}
			}
		}
		//log.Println(ip.String()," miss all allow: ",conf.UploadAllowedInterface)
		return false
	}

	if len(uploadDenyIps) > 0 {
		for i := 0; i < len(uploadDenyIps); i++ {
			switch uploadDenyIps[i].(type) {
			case net.IPNet:
				ipnet := uploadDenyIps[i].(net.IPNet)
				if ipnet.Contains(ip) {
					//log.Println(ip.String()," match deny: ",ipnet.String())
					return false
				}
			case net.IP:
				thisip := uploadDenyIps[i].(net.IP)
				if thisip.Equal(ip) {
					//log.Println(ip.String()," match deny: ",thisip.String())
					return false
				}
			}
		}
		//log.Println(ip.String()," miss all deny: ",conf.UploadDenyInterface)
		return true
	}

	//log.Println(ip.String()," miss all ip filter")
	return true
}

func parseIp(ips []string) []interface{} {
	re := make([]interface{}, len(ips))
	for i := 0; i < len(ips); i++ {
		if strings.Contains(ips[i], "/") {
			parts := strings.Split(ips[i], "/")
			netipTemp := parts[0]
			masklenStrTemp := parts[1]
			netip := net.ParseIP(netipTemp)
			maskLenTemp, err1 := strconv.Atoi(masklenStrTemp)
			if err1 != nil {
				log.Fatalln(err1)
				return nil
			}
			mask := net.CIDRMask(maskLenTemp, 32)
			re[i] = net.IPNet{netip, mask}
		} else {
			re[i] = net.ParseIP(ips[i])
		}
	}
	return re
}
