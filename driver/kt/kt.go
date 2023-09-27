package kt


import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/drivers/plugin"
)

const (
	defaultSSHUser       = "ubuntu"
	defaultSSHPort       = 22
	defaultActiveTimeout = 400
)

type Driver struct {
	*drivers.BaseDriver
	ApiEndpoint          string
    ActiveTimeout        int
	FlavorId             string
	FlavorName           string
	NetworkId            string
	NetworkName          string
	ImageId              string
	ImageName            string
	KeyPairName          string
	UserId               string
	UserPassword         string
	client               Client
}

// NodeTemplate에 보이는 화면, VM 설정할 수 잇는 Flag
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
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

// KT 드라이버의 새 인스턴스를 생성하고 반환
func NewDriver(hostname) *Driver {
	return &Driver{
		BaseDriver:   &drivers.BaseDriver{
			SSHUser:      defaultSSHUser,
			SSHPort:      defaultSSHPort,
			MachineName:  hostname,
		},
	}
}

// KT API 클라이언트를 생성(토큰발급)
func (d *Driver) getClient() (token, error) {
	url := ApiEndpoint + `/d1/identity/auth/tokens`
	method := "POST"
	data := `{"auth": {"identity": {"methods":["password"],"password":{"user":{"domain": {"id": "default"},"name": "` + UserId + `","password":"`+UserPassword+`"}}},"scope": {"project": {"domain": {"id": "default"},"name": "`+userId+`"}}}}`
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil{
		return fmt.Errorf("Error Creating Request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil{
		fmt.Errorf("Error Creamaking Post Request:", err)
		return
	}

	defer req.Body.Close()

	// Response의 특정 Key값만 추출하고 리턴
	token := response.Header.Get("X-Subject-Token")
	fmt.Println("X-Subject-Token: ", token)

	return 
}

// KT 드라이버 이름을 리턴
func (d *Driver) DriverName() string {
	return "kt"
}

// CreateFlags에서 반환된 개체로 드라이버 구성
func (d *Driver) setConfigFromFlags(flags drivers.DriverOptions) error {
	d.ApiEndpoint = flags.String("kt-api-endpoint-url")
	d.ActiveTimeout = flags.String("kt-active-timeout")
	d.FlavorId = flags.String("kt-flavor-id")
	d.FlavorName = flags.String("kt-flavor-name")
	d.NetworkId = flags.String("kt-network-id")
	d.NetworkName = flags.String("kt-network-name")
	d.ImageId = flags.String("kt-image-id")
	d.ImageName = flags.String("kt-image-name")
	d.KeyPairName = flags.String("kt-ssh-keypair-name")
	d.UserId = flags.String("kt-user-id")
	d.UserPassword = flags.String("kt-user-password")
	d.SSHUser = flags.String("kt-ssh-user")
	d.SSHPort = flags.String("kt-ssh-port")
}


// 드라이버를 구성하여 호스트를 생성
func(d *Driver) create() (state.State, error) {
	log.Info("Create Machine...")

	client, err := d.getClient()
	if err != nil {
		return err
	}

	token, err := d.getClient()
	url := ApiEndpoint + `d1/server/servers`
	method := "POST"
	data := `{"server":{"name": "` + d.hostname +`", "key_name": "` + d.KeyPairName +`","flavorRef": "` \
	        + FlavorId + `","availability_zone":"DX-M1","networks":[{"uuid": "` \
	        + d.NetworkId + `"}],"block_device_mapping_v2":[{"destination_type": "volume","boot_index": "0","source_type": "image","volume_size": 50,"uuid": "` \
			+ ImageId +`"}]}}`

	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil{
		fmt.Errorf("Error Creating Request:", err)
	}
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil{
		fmt.Errorf("Error Creamaking Post Request:", err)
		return
	}

	defer req.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 300 {
		return state.Error, err
	} else{
		d.Status = response.Status

		switch d.Status {
			case ""
		}
	} else {

	}


}

// ssh와 함께 사용할 VM 이름 반환
func GetSSHHostname() {
}


// 호스트의 상태를 반환
func (d *Driver) GetState() (state.State, error){
	log.Debug("Get status for OpenStack instance...")
	client, err := d.getClient()

	if err != nil {
		return state.Error, err
	}


	switch s {
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

// 호스트 강제 종료
func Kill() {

}

// 호스트 삭제
func Remove() {

}

// 호스트 재시작
func Restart() {

}

// 호스트 성공
func Start() error{

}

// 호스트 중지
func Stop() {

}

// 드라이버를 생성할 준비가 되었는지 확인
func PreCreateCheck() error {
	if d.ApiEndpoint == "" {
		return fmt.Errorf("ApiEndpoint is nil")
	}
	if d.FlavorId == "" && d.FlavorName == ""{
		return fmt.Errorf("FlavorId or FlavorName is nil")
	}
	if d.NetworkId == "" && d.NetworkName == ""{
		return fmt.Errorf("NetworkId or NetworkName is nil")
	}
	if d.ImageId == "" && d.ImageName == ""{
		return fmt.Errorf("ImageId or ImageName is nil")
	}
	if d.KeyPairName == "" {
		return fmt.Errorf("KeyPairName is nil")
	}
	if d.UserId == "" && d.UserPassword == ""{
		return fmt.Errorf("UserId or UserPassword is nil")
	}
	if d.SSHUser == "" && d.SSHPort == ""{
		d.SSHUser = defaultSSHUser
		d.SSHPort = defaultSSHPort
	}
}


// SSH에 사용할 사용자 이름 반환
func (d *Driver) getSSHUsername() string {
	return defaultSSHUser
}



