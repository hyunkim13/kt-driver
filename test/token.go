package main

import (
	"fmt"
	"strings"
	"net/http"
)

type Server struct {
	ID string `json:"id"`
}

type ResponseBody struct {
	Server Server `json:"server"`
}


func token() (string, error) {
	url := "https://api.ucloudbiz.olleh.com" + `/d1/identity/auth/tokens`
	method := "POST"
	data := `{"auth": {"identity": {"methods":["password"],"password":{"user":{"domain": {"id": "default"},"name": "` + "saasify@systeer.com" + `","password":"`+"Sysmaster77!!"+`"}}},"scope": {"project": {"domain": {"id": "default"},"name": "`+"saasify@systeer.com"+`"}}}}`
	fmt.Println(data)
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil{
		fmt.Errorf("Error Creating Request:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil{
		fmt.Errorf("Error Creamaking Post Request:", err)
	}

	defer req.Body.Close()

	// fmt.Println("Status Code:", response.Status)
	// headers := response.Header
	// for key, values := range headers {
	// 	fmt.Printf("%s: %s\n", key, values)
	// }

	token := response.Header.Get("X-Subject-Token")
	fmt.Println("X-Subject-Token: ", token)
	return token, nil
}

func createVM() {
	token, err := token()
	url := "https://api.ucloudbiz.olleh.com" + `/d1/server/servers`
	method := "POST"
	data := `{"server":{"name": "infranics3", "key_name": "SaaSifyKey","flavorRef": "7f56ce4a-5b56-4b53-be63-f4dda5216b63","availability_zone":"DX-M1","networks":[{"uuid": "a3f25a44-efaa-47d7-bdd4-b78032662d68"}],"block_device_mapping_v2":[{"destination_type": "volume","boot_index": "0","source_type": "image","volume_size": 50,"uuid": "84a10047-cbd8-4fb3-a743-85600a7b6961"}]}}`

	fmt.Println(data)
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil{
		fmt.Errorf("Error Creating Request:", err)
	}
	req.Header.Set("X-Auth-Token", token)
	if err != nil{
		fmt.Errorf("Error Creamaking Post Request:", err)
		return
	}

	defer req.Body.Close()
}

func main(){
	createVM()
}
