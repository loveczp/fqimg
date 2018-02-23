package plugin

import (
	"net/http"
	"strings"
	"net"
	"github.com/loveczp/fqimg/lib"
	"io"
	"log"
	"strconv"
)

var (
	uploadAllowedInterface []interface{}
	uploadDenyInterface    []interface{}
)

func Plugin_upload_ipauth(h http.HandlerFunc, conf lib.Config) http.HandlerFunc {
	parseIp(conf.UploadAllowed)
	parseIp(conf.UploadDeny)
	return func(writer http.ResponseWriter, request *http.Request) {
		if ipPass(request, conf) == false {
			writer.WriteHeader(http.StatusInternalServerError)
			io.WriteString(writer, "your ip address is forbiden to upload");
			return
		}
	}
}

func ipPass(req *http.Request, config lib.Config) bool {
	adds := strings.Split(req.RemoteAddr, ":")
	ip := net.ParseIP(adds[0]);
	//log.Print("req remote ip :", ip)
	if len(uploadAllowedInterface) > 0 {

		for i := 0; i < len(uploadAllowedInterface); i++ {
			switch uploadAllowedInterface[i].(type) {
			case net.IPNet:
				ipnet := uploadAllowedInterface[i].(net.IPNet)
				if ipnet.Contains(ip) {
					//log.Println(ip.String()," match allow: ",ipnet.String())
					return true;
				}
			case net.IP:
				thisip := uploadAllowedInterface[i].(net.IP)
				if thisip.Equal(ip) {
					//log.Println(ip.String()," match allow: ",thisip.String())
					return true;
				}
			}
		}
		//log.Println(ip.String()," miss all allow: ",conf.UploadAllowedInterface)
		return false;
	}

	if len(uploadDenyInterface) > 0 {
		for i := 0; i < len(uploadDenyInterface); i++ {
			switch uploadDenyInterface[i].(type) {
			case net.IPNet:
				ipnet := uploadDenyInterface[i].(net.IPNet)
				if ipnet.Contains(ip) {
					//log.Println(ip.String()," match deny: ",ipnet.String())
					return false;
				}
			case net.IP:
				thisip := uploadDenyInterface[i].(net.IP)
				if thisip.Equal(ip) {
					//log.Println(ip.String()," match deny: ",thisip.String())
					return false;
				}
			}
		}
		//log.Println(ip.String()," miss all deny: ",conf.UploadDenyInterface)
		return true;
	}

	//log.Println(ip.String()," miss all ip filter")
	return true;
}

func parseIp(ips []string) []interface{} {
	re := make([]interface{}, len(ips));
	for i := 0; i < len(ips); i++ {
		if strings.Contains(ips[i], "/") {
			parts := strings.Split(ips[i], "/");
			netipTemp := parts[0];
			masklenStrTemp := parts[1];
			netip := net.ParseIP(netipTemp);
			maskLenTemp, err1 := strconv.Atoi(masklenStrTemp);
			if err1 != nil {
				log.Fatalln(err1);
				return nil;
			}
			mask := net.CIDRMask(maskLenTemp, 32)
			re[i] = net.IPNet{netip, mask}
		} else {
			re[i] = net.ParseIP(ips[i]);
		}
	}
	return re
}
