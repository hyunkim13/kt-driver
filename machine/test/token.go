package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)


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
	client := &http.Client{}
	response, err := client.Do(req)

	defer req.Body.Close()

    resBody, _ := ioutil.ReadAll(response.Body)
	resBytes := []byte(resBody)
	fmt.Println("resBody: ", resBody)
	fmt.Println("resBytes: ", resBytes)
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(resBytes, &jsonRes)

	fmt.Println("jsonRes: ", jsonRes)
	detail_map := jsonRes["server"].(map[string]interface{})
	id := detail_map["id"].(string)
	
	
	fmt.Println("id: ", id)

	// uuid := response.Body.get("")
}

func getState() {
	token, err := token()
	url := "https://api.ucloudbiz.olleh.com" + `/d1/server/servers/d4e2deae-c3f8-446c-b949-e90a0f75b0dd`
	method := "GET"
	
	req, err := http.NewRequest(method, url, nil)
	if err != nil{
		fmt.Errorf("Error Request server list Request:", err)
	}
	req.Header.Set("X-Auth-Token", token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil{
		fmt.Errorf("Error Creamaking Post Request:", err)
		return 
	}
	fmt.Println("response.Body: ", response.Body)
	defer response.Body.Close()
	
	resBody, _ := ioutil.ReadAll(response.Body)
	resBytes := []byte(resBody)
	fmt.Println("resBody: ", resBody)
	fmt.Println("resBytes: ", resBytes)
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(resBytes, &jsonRes)

	fmt.Println("jsonRes: ", jsonRes)
	detail_map := jsonRes["server"].(map[string]interface{})
	status := detail_map["status"].(string)

	fmt.Println("status: ", status)
}

func getVMId(){

	token, err := token()
	url := "https://api.ucloudbiz.olleh.com" + `/d1/server/servers`
	method := "GET"
	
	req, err := http.NewRequest(method, url, nil)
	if err != nil{
		fmt.Errorf("Error Request server list Request:", err)
	}
	req.Header.Set("X-Auth-Token", token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil{
		fmt.Errorf("Error Creamaking Post Request:", err)
		return 
	}
	fmt.Println("response.Body: ", response.Body)
	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)
	resBytes := []byte(resBody)
	fmt.Println("resBody: ", resBody)
	fmt.Println("resBytes: ", resBytes)
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(resBytes, &jsonRes)

	var data map[string]interface{}
	var id string
	
	decoder := json.NewDecoder(strings.NewReader(string(resBytes)))
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	servers, found := data["servers"].([]interface{})
	if !found {
		fmt.Println("Key 'servers' not found")
		return
	}

	

	for _, server := range servers {
		serverMap, isMap := server.(map[string]interface{})
		if !isMap {
			fmt.Println("Invalid server data")
			continue
		}

		name, found := serverMap["name"].(string)
		if !found {
			continue
		}

		if name == "SaaSify01" {
			id, found := serverMap["id"].(string)
			if !found {
				fmt.Println("ID not found for 'SaaSify01'")
				continue
			}
			fmt.Println("ID for 'SaaSify01':", id)
			break
			fmt.Println("id:", id)
		}
		
	}
	fmt.Println("id:", id)

	
}

func main(){
	// createVM()
	// getState()
	// getVMId()
	token()
}
