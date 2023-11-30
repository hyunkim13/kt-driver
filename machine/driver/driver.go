package driver

// 23.11.01 hskim
// KT NodeDriver 개발.
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/docker/machine/libmachine/drivers"

	// "github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/state"
)

const (
	defaultSSHUser       = "ubuntu"
	defaultSSHPort       = 22
	defaultActiveTimeout = 400
)

type Driver struct {
	*drivers.BaseDriver
	ApiEndpoint    string
	ActiveTimeout  int
	FlavorId       string
	FlavorName     string
	NetworkId      string
	NetworkName    string
	ImageId        string
	ImageName      string
	KeyPairName    string
	PrivateKeyFile string
	UserId         string
	UserPassword   string
	VMId           string
	SSHUser        string
	SSHPort        int
	Token          string
	IPAddress      string
}

// NodeTemplate에 보이는 화면,Flag를 설정할 수 있다.
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	fmt.Println("GetCreateFlags funciton...")
	// log.Debug("GetCreateFlags funciton...")
	return []mcnflag.Flag{

		mcnflag.StringFlag{
			EnvVar: "KT_API_ENDPOINT_URL",
			Name:   "kt-api-endpoint-url",
			Usage:  "KT API ENDPOINT URL",
			Value:  "",
		},
		mcnflag.IntFlag{
			EnvVar: "KT_Active_Timeout",
			Name:   "kt-active-timeout",
			Usage:  "KT Active Timeout",
			Value:  defaultActiveTimeout,
		},
		mcnflag.StringFlag{
			EnvVar: "KT_FLAVOR_ID",
			Name:   "kt-flavor-id",
			Usage:  "KT Flavor id to use for the instance",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "KT_FLAVOR_NAME",
			Name:   "kt-flavor-name",
			Usage:  "KT Flavor name to use for the instance",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "KT_NETWORK_ID",
			Name:   "kt-network-id",
			Usage:  "KT IaaS network id the machine will be connected on",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "KT_NETWORK_NAME",
			Name:   "kt-network-name",
			Usage:  "KT IaaS network name the machine will be connected on",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "KT_Image_Id",
			Name:   "kt-image-id",
			Usage:  "KT image id to use for the instance",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "KT_Image_Name",
			Name:   "kt-image-name",
			Usage:  "KT image name to use for the instance",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "KT_SSH_KEYPAIR_NAME",
			Name:   "kt-ssh-keypair-name",
			Usage:  "KT keypair to use to SSH to the instance",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "KT_PRIVATE_KEY_FILE",
			Name:   "kt-private-key-file",
			Usage:  "Private keyfile to use for SSH (absolute path)",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "KT_User_Id",
			Name:   "kt-user-id",
			Usage:  "KT user-id",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "KT_USER_PASSWORD",
			Name:   "kt-user-password",
			Usage:  "KT user-password",
			Value:  "",
		},
		mcnflag.StringFlag{
			EnvVar: "KT_SSH_USER",
			Name:   "kt-ssh-user",
			Usage:  "KT SSH user",
			Value:  defaultSSHUser,
		},
		mcnflag.IntFlag{
			EnvVar: "KT_SSH_PORT",
			Name:   "kt-ssh-port",
			Usage:  "KT SSH port",
			Value:  defaultSSHPort,
		},
	}
}

// KT 드라이버를 생성하고 반환
func NewDriver() *Driver {
	// log.Debug("NewDriver function...")
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			SSHUser: defaultSSHUser,
			SSHPort: defaultSSHPort,
		},
	}
}

// KT API 클라이언트 생성(토큰발급)
// NodeTemplate에서 ID와 PW를 받아 KT API를 이용해서 토큰 발급
func (d *Driver) getClient() (string, error) {
	fmt.Println("getClient funciton...")
	// log.Debug("getClient funciton...")
	url := d.ApiEndpoint + "/d1/identity/auth/tokens"
	method := "POST"
	data := `{"auth": {"identity": {"methods":["password"],"password":{"user":{"domain": {"id": "default"},"name": "` + d.UserId + `","password":"` + d.UserPassword + `"}}},"scope": {"project": {"domain": {"id": "default"},"name": "` + d.UserId + `"}}}}`
	fmt.Println("getClient url: ", method, url)
	// fmt.Println("getClient data: ", data)
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return "", errors.New("Create Token req is nil")
	}
	// fmt.Println("getClient req: ", req)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Error Creamaking Post Request:", err)
		return "", errors.New("Create Token response is nil")
	}

	defer req.Body.Close()

	// Response의 특정 Key값만 추출하고 리턴
	token := response.Header.Get("X-Subject-Token")
	fmt.Println("X-Subject-Token: ", token)
	d.Token = token

	fmt.Println("getClient funciton End...")
	// log.Debug("getClient funciton End...")
	return token, nil
}

// ssh와 함께 사용할 VM 이름 반환
func (d *Driver) GetSSHHostname() (string, error) {
	fmt.Println("GetSSHHostname function...")
	return d.GetIP()
}

// KT 드라이버 이름을 리턴
func (d *Driver) DriverName() string {
	// log.Debug("DriverName function...")
	return "kt"
}

// CreateFlags에서 반환된 개체로 드라이버 구성
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	fmt.Println("SetConfigFromFlags funciton...")
	// log.Debug("SetConfigFromFlags funciton...")
	d.ApiEndpoint = flags.String("kt-api-endpoint-url")
	d.ActiveTimeout = flags.Int("kt-active-timeout")
	d.FlavorId = flags.String("kt-flavor-id")
	d.FlavorName = flags.String("kt-flavor-name")
	d.NetworkId = flags.String("kt-network-id")
	d.NetworkName = flags.String("kt-network-name")
	d.ImageId = flags.String("kt-image-id")
	d.ImageName = flags.String("kt-image-name")
	d.KeyPairName = flags.String("kt-ssh-keypair-name")
	d.PrivateKeyFile = flags.String("kt-private-key-file")
	d.UserId = flags.String("kt-user-id")
	d.UserPassword = flags.String("kt-user-password")
	d.SSHUser = flags.String("kt-ssh-user")
	d.SSHPort = flags.Int("kt-ssh-port")

	fmt.Println("SetConfigFromFlags funciton End...")
	// log.Debug("SetConfigFromFlags funciton End...")
	return nil
}

func (d *Driver) GetURL() (string, error) {
	fmt.Println("GetURL funciton...")
	// log.Debug("GetURL funciton...")

	ip, err := d.GetIP()

	fmt.Println("GetURL ip", ip)
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", nil
	}

	fmt.Println("GetURL funciton End...")
	// log.Debug("GetURL funciton End...")
	return fmt.Sprintf("tcp://%s:%d", ip, 2376), nil
}

// VM의 IP를 추출, VM 이름을 이용하여 상세조회를 해야하고 ip를 추출
func (d *Driver) GetIP() (string, error) {
	fmt.Println("KT GetIP funciton...")
	// log.Debug("KT GetIP funciton...")

	fmt.Println("d.IPAddress: ", d.IPAddress)
	// log.Debug("d.IPAddress: ", d.IPAddress)

	if d.IPAddress != "" {
		return d.IPAddress, nil
	}

	// token, err := d.getClient()
	// fmt.Println("GETIP d.vmId: ", d.VMId)
	// url := d.ApiEndpoint + `/d1/server/servers/` + d.VMId

	fmt.Println("GetIP2 funciton...")
	// log.Debug("GetIP3 funciton...")

	vmId, err := d.getVMId(d.GetMachineName())

	fmt.Println("GETIP d.vmId or vmId: ", vmId)
	// log.Debug("GETIP2 d.vmId or vmId: ", vmId)
	url := d.ApiEndpoint + `/d1/server/servers/` + vmId
	method := "GET"
	fmt.Println("GetIP url: ", method, url)

	req, err := http.NewRequest(method, url, strings.NewReader(""))
	if err != nil {
		fmt.Errorf("Error GetIP Request:", err)
	}
	req.Header.Set("X-Auth-Token", d.Token)
	fmt.Println("GetIP token: ", d.Token)
	client := &http.Client{}
	response, err := client.Do(req)
	_ = response
	if err != nil {
		fmt.Errorf("Error Creamaking Post Request:", err)
	}

	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	resBytes := []byte(resBody)
	// fmt.Println("GetIP response", response)
	var jsonRes map[string]interface{}
	var addr string

	_ = json.Unmarshal(resBytes, &jsonRes)
	detail_map := jsonRes["server"].(map[string]interface{})
	addresses := detail_map["addresses"].(map[string]interface{})
	dmz, ok := addresses["DMZ"].([]interface{})
	fmt.Println("GETIP dmz: ", dmz)

	if !ok {
		fmt.Println("DMZ 키가 존재하지 않거나 배열이 아닙니다.")
	}

	if len(dmz) > 0 {
		firstItem, ok := dmz[0].(map[string]interface{})
		if !ok {
			fmt.Println("DMZ 배열의 첫 번째 요소가 맵이 아닙니다.")
		}

		addr, ok = firstItem["addr"].(string)
		if !ok {
			fmt.Println("addr 키가 존재하지 않거나 문자열이 아닙니다.")
		}

	} else {
		fmt.Errorf("dmz가 없습니다")
	}

	d.IPAddress = addr

	fmt.Println("addr: ", addr)

	fmt.Println("GetIP funciton End...")
	// log.Debug("GetIP funciton End...")
	return "", fmt.Errorf("No IP found for the machine")
}

// VM을 생성 후 상태체크를 해야함
func (d *Driver) GetState() (state.State, error) {
	hostname := d.GetMachineName()

	// token, err := d.getClient()

	// if err != nil {
	// 	return state.Error, err
	// }
	fmt.Println("GetState hostname: ", hostname)

	id, err := d.getVMId(hostname)

	url := d.ApiEndpoint + "/d1/server/servers/" + id
	method := "GET"
	fmt.Println("GetState url: ", method, url)
	req, err := http.NewRequest(method, url, nil)
	// fmt.Println("GetState req: " , req)
	if err != nil {
		fmt.Errorf("Error Request server list Request:", err)
	}
	req.Header.Set("X-Auth-Token", d.Token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Error Stauts GET Request:", err)
	}
	// log.Debug("response.Body: ", response.Body)
	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)
	resBytes := []byte(resBody)
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(resBytes, &jsonRes)

	status_detail_map := jsonRes["server"].(map[string]interface{})
	status := status_detail_map["status"].(string)

	fmt.Println("GetState function End...")
	// log.Debug("GetState function End...")
	fmt.Println("status: ", status)
	// log.Debug("status: ", status)

	switch status {
	case "ACTIVE":
		return state.Running, nil
	case "PAUSED":
		return state.Paused, nil
	case "SUSPENDED":
		return state.Saved, nil
	case "SHUTOFF":
		return state.Stopped, nil
	case "BUILDING":
		return state.Starting, nil
	case "ERROR":
		return state.Error, nil
	}
	return state.None, nil
}

// VM을 생성
func (d *Driver) Create() error {
	// log.Debug("Create function...")
	fmt.Println("Create function...")

	// 특정 경로에 Keyfile을 생성
	if d.PrivateKeyFile != "" {
		d.createSSHKey()
	}

	// time.Sleep(10 * time.Second)

	// log.Debug("[KT.go] Create funciton PrivateKeyFile: ", d.KeyPairName)
	// log.Info("[KT.go] Create funciton PrivateKeyFile: ", d.KeyPairName)
	// time.Sleep(10 * time.Second)

	// log.Debug("[KT.go] Create funciton &PrivateKeyFile: ", &d.KeyPairName)
	// log.Info("[KT.go] Create funciton &PrivateKeyFile: ", &d.KeyPairName)
	// time.Sleep(10 * time.Second)

	// log.Debug("[KT.go] Create funciton d.PrivateKeyFile: ", d.PrivateKeyFile)
	// log.Info("[KT.go] Create funciton d.PrivateKeyFile: ", d.PrivateKeyFile)

	// log.Debug("[KT.go] Create funciton &d.PrivateKeyFile: ", &d.PrivateKeyFile)
	// log.Info("[KT.go] Create funciton &d.PrivateKeyFile: ", &d.PrivateKeyFile)
	// time.Sleep(10 * time.Second)

	if d.KeyPairName != "" {
		fmt.Println("[KT.go]Create funciton d.KeyPairName if ...")
		// log.Debug("[KT.go]Create funciton d.KeyPairName if ...")
		if err := d.loadSSHKey(); err != nil {
			return err
		}
	} else {
		d.KeyPairName = fmt.Sprintf("%s-%s", d.MachineName, mcnutils.GenerateRandomID())
		fmt.Println("Create funciton d.KeyPairName else ...", d.KeyPairName)
		// log.Debug("Create funciton d.KeyPairName else ...", d.KeyPairName)
		if err := d.createSSHKey(); err != nil {
			return err
		}
	}

	token, err := d.getClient()
	if err != nil {
		return errors.New("Create Token req is nil")
	}
	// fmt.Println("Create token: ", token)
	hostname := d.GetMachineName()
	// 실질적 VM 생성
	id, err := d.custom_createVM(hostname, token)
	fmt.Println("Create funcion VMid: ", id)

	fmt.Println("Create function End...")
	// log.Debug("Create function End...")
	return nil
}

// VM 시작
func (d *Driver) Start() error {
	fmt.Println("Start funciton...")
	// log.Debug("Start funciton...")

	fmt.Println("Start funciton End...")
	// log.Debug("Start funciton End...")
	return nil
}

// VM 삭제
func (d *Driver) Stop() error {
	fmt.Println("Stop funciton...")
	// log.Debug("Stop funciton...")
	hostname := d.GetMachineName()

	fmt.Println("Get status for KT instance...")
	token, err := d.getClient()
	if err == nil {
		fmt.Println("Authentication Success...")
	}

	vmId, err := d.getVMId(hostname)
	if err != nil {
		fmt.Errorf("Get VM Id is null")
	}

	// KT API VM 강제 삭제
	url := d.ApiEndpoint + "/d1/server/servers/" + vmId + "/action"
	method := "POST"
	data := `{"forceDelete": null}`

	fmt.Println("Stop url: ", method, url)
	req, err := http.NewRequest(method, url, strings.NewReader(data))

	if err != nil {
		fmt.Errorf("Error Creating Request:", err)
	}
	req.Header.Set("X-Auth-Token", token)
	if err != nil {
		fmt.Errorf("Error Creamaking Post Request:", err)
		// return
	}
	client := &http.Client{}
	response, err := client.Do(req)

	defer req.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 300 {
		fmt.Errorf("unable to Stop VM ")
	}

	fmt.Println("Stop funciton End...")
	// log.Debug("Stop funciton End...")
	return err
}

// // 호스트 재시작
func (d *Driver) Restart() error {
	fmt.Println("Restart funciton...")
	// log.Debug("Restart funciton...")

	fmt.Println("Restart funciton End...")
	// log.Debug("Restart funciton End...")
	return nil
}

// 호스트 강제 종료
func (d *Driver) Kill() error {
	fmt.Println("Kill funciton...")
	// log.Debug("Kill funciton...")

	token, err := d.getClient()
	vmId := d.GetMachineName()

	// KT API VM 삭제
	url := d.ApiEndpoint + "/d1/server/servers/" + vmId + `/action`
	method := "POST"
	data := `{"forceDelete": null}`
	req, err := http.NewRequest(method, url, strings.NewReader(data))

	fmt.Println("Remove url: ", method, url)

	if err != nil {
		fmt.Errorf("Error Creating Request:", err)
	}
	req.Header.Set("X-Auth-Token", token)
	if err != nil {
		fmt.Errorf("Error Creamaking Post Request:", err)
		// return
	}
	client := &http.Client{}
	response, err := client.Do(req)

	defer req.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 300 {
		fmt.Errorf("unable to Stop VM ")
	}

	fmt.Println("Kill funciton End...")
	// log.Debug("Kill funciton End...")
	return d.Stop()

}

// VM 삭제
func (d *Driver) Remove() error {
	fmt.Println("Remove funciton...")
	// log.Debug("Remove funciton...")
	token, err := d.getClient()

	vmId, err := d.getVMId(d.GetMachineName())
	if err != nil {
		fmt.Errorf("Remove Fail, vmId is nil: ", err)
	}

	// KT API VM 삭제
	url := d.ApiEndpoint + "/d1/server/servers/" + vmId
	method := "DELETE"
	req, err := http.NewRequest(method, url, nil)

	fmt.Println("Remove url: ", method, url)

	if err != nil {
		fmt.Errorf("Error Creating Request:", err)
	}
	req.Header.Set("X-Auth-Token", token)
	if err != nil {
		fmt.Errorf("Error Creamaking Post Request:", err)
		// return
	}
	client := &http.Client{}
	response, err := client.Do(req)

	defer req.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 300 {
		fmt.Errorf("unable to Stop VM ")
	}
	fmt.Println("Remove funciton End...")
	// log.Debug("Remove funciton End...")
	return nil
}

// ssh가 필요한 Key파일 읽기, public key는 없고 privatekey만 있다.
func (d *Driver) loadSSHKey() error {
	fmt.Println("Loading Key Pair", d.KeyPairName)
	// log.Debug("Loading Key Pair", d.KeyPairName)

	fmt.Println("loadSSHKey funciton d.PrivateKeyFile", d.PrivateKeyFile)
	// log.Debug("loadSSHKey funciton d.PrivateKeyFile", d.PrivateKeyFile)

	fmt.Println("loadSSHKey funciton privateSSHKeyPath...", d.privateSSHKeyPath())
	// log.Debug("loadSSHKey funciton privateSSHKeyPath...", d.privateSSHKeyPath())

	fmt.Println("loadSSHKey funciton publicSSHKeyPath...", d.publicSSHKeyPath())
	// log.Debug("loadSSHKey funciton publicSSHKeyPath...", d.publicSSHKeyPath())

	keyString := d.PrivateKeyFile
	// 문자열을 ASCII 코드로 변환하여 저장할 슬라이스 생성
	asciiBytes := make([]byte, len(keyString))

	// 문자열을 ASCII 코드로 변환
	for i, char := range keyString {
		asciiBytes[i] = byte(char)
	}

	fmt.Println("loadSSHKey funciton asciiBytes...", asciiBytes)
	// log.Debug("loadSSHKey funciton asciiBytes...", asciiBytes)
	fmt.Println(reflect.TypeOf(asciiBytes))
	// log.Debug(reflect.TypeOf(asciiBytes))

	privateKey, err := ioutil.ReadFile(d.PrivateKeyFile)
	if err != nil {
		return err
	}

	fmt.Println("loadSSHKey funciton privateKey...", privateKey)
	// log.Debug("loadSSHKey funciton privateKey...", privateKey)

	// if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
	// 	return err
	// }

	// publicKey, err := os.ReadFile(fmt.Sprintf("%s.pub", d.GetSSHKeyPath()))
	// // publicKey, err := ioutil.ReadFile(d.GetSSHKeyPath())
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("loadSSHKey funciton publicKey...", publicKey)
	// log.Debug("loadSSHKey funciton publicKey...", publicKey)

	if err := ioutil.WriteFile(d.privateSSHKeyPath(), privateKey, 0600); err != nil {
		return err
	}
	// if err := ioutil.WriteFile(d.publicSSHKeyPath(), publicKey, 0600); err != nil {
	// 	return err
	// }
	return nil
}

// Key파일 경로
func (d *Driver) privateSSHKeyPath() string {
	return d.GetSSHKeyPath()
}

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}

// 클러스터 생성 전 파라미터가 맞게 들어갔는지 확인
// 없는 부분이 있으면 return err
func (d *Driver) PreCreateCheck() error {
	// log.Debug("PreCreateCheck function...")
	fmt.Println("PreCreateCheck function...")
	if d.ApiEndpoint == "" {
		return fmt.Errorf("ApiEndpoint is nil")
	}
	if d.FlavorId == "" && d.FlavorName == "" {
		return fmt.Errorf("FlavorId or FlavorName is nil")
	}
	if d.NetworkId == "" && d.NetworkName == "" {
		return fmt.Errorf("NetworkId or NetworkName is nil")
	}
	if d.ImageId == "" && d.ImageName == "" {
		return fmt.Errorf("ImageId or ImageName is nil")
	}
	if d.KeyPairName == "" && d.PrivateKeyFile == "" {
		fmt.Println("PreCreateCheck KeyPairName...", d.KeyPairName)
		fmt.Println("PreCreateCheck PrivateKeyFile...", d.PrivateKeyFile)
		return fmt.Errorf("KeyPairName or PrivateKeyFile is nil")
	}
	if d.UserId == "" && d.UserPassword == "" {
		return fmt.Errorf("UserId or UserPassword is nil")
	}
	if d.SSHUser == "" && d.SSHPort == 0 {
		d.SSHUser = defaultSSHUser
		d.SSHPort = defaultSSHPort
	}

	fmt.Println("PreCreateCheck funciton End...")
	// log.Debug("PreCreateCheck funciton End...")
	return nil
}

// VM이름을 가지고 vmid를 검색, vmIP를 추출하기 위해 id가 필요하다.
func (d *Driver) getVMId(hostname string) (string, error) {
	fmt.Println("getVMId funciton...")
	// log.Debug("getVMId funciton...")
	// token, err := d.getClient()
	url := d.ApiEndpoint + "/d1/server/servers"
	method := "GET"

	fmt.Println("getVMId url: ", method, url)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Errorf("Error Request server list Request:", err)
	}
	req.Header.Set("X-Auth-Token", d.Token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Error Creamaking Post Request:", err)
		return "response nil", nil
	}
	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)
	resBytes := []byte(resBody) // Response를 decode
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(resBytes, &jsonRes)

	var data map[string]interface{}

	// fmt.Println("getVMId ResponseBody: ", string(resBytes))
	fmt.Println("getVMId hostname: ", hostname)

	decoder := json.NewDecoder(strings.NewReader(string(resBytes)))
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "getVMId functioon Json Decode", nil
	}

	servers, found := data["servers"].([]interface{})
	var vmId string

	if !found {
		fmt.Println("Key 'servers' not found")
		return "getVMId functioon 'server's key is not found", nil
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
		if name == hostname {
			vmId, found = serverMap["id"].(string)
			if !found {
				fmt.Println("Find!!")
				continue
			}
			break
			fmt.Println("VM id:", vmId)
		}
	}

	fmt.Println("VM id: ", vmId)

	fmt.Println("getVMId funciton End...")
	// log.Debug("getVMId funciton End...")

	return vmId, nil
}

// VM을 생성
func (d *Driver) custom_createVM(hostname string, token string) (string, error) {
	fmt.Println("custom_createVM funciton...")
	// log.Debug("custom_createVM funciton...")

	url := d.ApiEndpoint + `/d1/server/servers`
	method := "POST"
	data := `{"server":{"name": "` + hostname + `", "key_name": "` + d.KeyPairName + `","flavorRef": "`
	data += d.FlavorId + `","availability_zone":"DX-M1","networks":[{"uuid": "`
	data += d.NetworkId + `"}],"block_device_mapping_v2":[{"destination_type": "volume","boot_index": "0","source_type": "image","volume_size": 50,"uuid": "`
	data += d.ImageId + `"}]}}`

	fmt.Println("Create data: ", data)
	fmt.Println("Create token: ", token)

	req, err := http.NewRequest(method, url, strings.NewReader(data))

	// log.Debug("custom_createVM url: ", method, url)
	// fmt.Println("custom_createVM url: ", method, url)
	// time.Sleep(2 * time.Second)
	// log.Debug("custom_createVM hostname: ", hostname)
	// fmt.Println("custom_createVM hostname: ", hostname)
	// time.Sleep(2 * time.Second)
	// log.Debug("custom_createVM KeyPairName: ", d.KeyPairName)
	// fmt.Println("custom_createVM KeyPairName: ", d.KeyPairName)
	// time.Sleep(2 * time.Second)
	// log.Debug("custom_createVM FlavorId: ", d.FlavorId)
	// fmt.Println("custom_createVM FlavorId: ", d.FlavorId)
	// time.Sleep(2 * time.Second)
	// log.Debug("custom_createVM NetworkId: ", d.NetworkId)
	// fmt.Println("custom_createVM NetworkId: ", d.NetworkId)
	// time.Sleep(2 * time.Second)
	// log.Debug("custom_createVM ImageId: ", d.ImageId)
	// fmt.Println("custom_createVM ImageId: ", d.ImageId)
	// time.Sleep(2 * time.Second)

	// fmt.Println("Create req: ", req)
	if err != nil {
		fmt.Errorf("Error Creating Request:", err)
	}
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{}
	response, err := client.Do(req)
	_ = response
	if err != nil {
		fmt.Errorf("Error Creamaking Post Request:", err)
	}

	if response.StatusCode < 200 || response.StatusCode > 300 {
		return "", fmt.Errorf("Create VM Error, Check FlavorId or NetworkId, ImageId")
	}
	defer req.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	resBytes := []byte(resBody)
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(resBytes, &jsonRes)
	id_detail_map := jsonRes["server"].(map[string]interface{})
	id := id_detail_map["id"].(string)
	d.VMId = id
	fmt.Println("Create id: ", id)

	fmt.Println("custom_createVM funciton End...")
	// log.Debug("custom_createVM funciton End...")

	return id, err
}

// private key를 사전에 생성, 이 부분이 없으면 vm 생성이 안됨, custom으로 만들어진 부분.
func (d *Driver) createSSHKey() error {
	// log.Debug("[KT.go]createSSHKey &d.KeyPairName .. map[string]string...: ", &d.KeyPairName, map[string]string{"Name": d.KeyPairName})
	fmt.Println("[KT.go]createSSHKey &d.KeyPairName .. map[string]string...: ", &d.KeyPairName, map[string]string{"Name": d.KeyPairName})

	// log.Debug("privatesshkeypath: ", d.privateSSHKeyPath())
	fmt.Println("privatesshkeypath: ", d.privateSSHKeyPath())

	// log.Debug("publicSSHKeyPath: ", d.publicSSHKeyPath())
	fmt.Println("publicSSHKeyPath: ", d.publicSSHKeyPath())

	keyString := d.PrivateKeyFile
	// 문자열을 ASCII 코드로 변환하여 저장할 슬라이스 생성
	asciiBytes := make([]byte, len(keyString))

	// 문자열을 ASCII 코드로 변환
	for i, char := range keyString {
		asciiBytes[i] = byte(char)
	}

	// 변환된 ASCII 코드를 파일로 저장
	if err := ioutil.WriteFile(d.privateSSHKeyPath(), asciiBytes, 0600); err != nil {
		return err
	}

	// log.Debug("Save keyfile: asciiBytes", asciiBytes)
	// fmt.Println("Save keyfile: asciiBytes", asciiBytes)
	// // time.Sleep(10 * time.Second)

	// if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
	// 	return err
	// }
	// publicKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
	// fmt.Println("createPrivateKey funciton publicKey...", publicKey)
	// log.Debug("createPrivateKey funciton publicKey...", publicKey)
	// if err != nil {
	// 	return err
	// }
	// if err := ioutil.WriteFile(d.GetSSHKeyPath()+".pub", publicKey, 0600); err != nil {
	// 	return err
	// }

	// time.Sleep(10 * time.Second)

	// 	log.Debug("[KT.go] createPrivateKey publicKey: ", publicKey)
	// fmt.Println("[KT.go] createPrivateKey publicKey::", publicKey)
	// if publicKey != nil {
	// 	log.Errorf("Error generating ssh key")
	// }

	// cmd := exec.Command("ssh-keygen", "-f", "/opt/jail/"+d.GetSSHKeyPath(), "-y", ">", "/opt/jail/"+d.GetSSHKeyPath()+".pub")
	// log.Debug("publicSSHKeyPath: ", cmd)
	// fmt.Println("publicSSHKeyPath: ", cmd)
	// cmd.CombinedOutput()

	return nil
}
