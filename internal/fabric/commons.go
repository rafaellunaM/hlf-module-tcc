package fabric

type CA struct {
	Capacity   string `json:"capacity"`
	Name       string `json:"name"`
	EnrollId   string `json:"enrollId"`
	EnrollPw   string `json:"enrollpw"`
	Hosts      string `json:"hosts"`
	IstioPort  string `json:"istioPort"`
	User       string `json:"user"`
	Secret     string `json:"secret"`
	UserType   string `json:"userType"`
	MspID      string `json:"mspID"`
}

type Peer struct {
	StateDB        string `json:"stateDB"`
	EnrollIDpeer   string `json:"enrollIDpeer"`
	Mspid          string `json:"mspid"`
	EnrollIPWpeer  string `json:"enrollIPWpeer"`
	Capacity       string `json:"capacity"`
	Name           string `json:"name"`
	CAName         string `json:"CAName"`
	Hosts          string `json:"hosts"`
	IstioPort      string `json:"istioPort"`
	User           string `json:"user"`
	EnrollId       string `json:"enrollId"`
	EnrollPw       string `json:"enrollPw"`
	Secret         string `json:"secret"`
	UserType       string `json:"userType"`
}

type Orderer struct {
	CAName             string `json:"CAName"`
	User               string `json:"user"`
	Secret             string `json:"secret"`
	UserType           string `json:"userType"`
	EnrollID           string `json:"enrollID"`
	EnrollPW           string `json:"enrollPW"`
	Mspid              string `json:"mspid"`
	CaURL              string `json:"caURL"`
	Capacity           string `json:"capacity"`
	Name               string `json:"name"`
	IstioPort          string `json:"istioPort"`
	EnrollIDorderer    string `json:"enrollIDorderer"`
	EnrollPWorderer    string `json:"enrollPWorderer"`
	CaNameService      string `json:"ca-name-service"`
	Hosts              string `json:"hosts"`
	AdminHosts         string `json:"admin-hosts"`
}

type Channel struct {
	Name           					string `json:"name"`
	UserAdmin      					string `json:"userAdmin"`
	Secretadmin    					string `json:"secretadmin"`
	UserType       					string `json:"userType"`
	EnrollID       					string `json:"enrollID"`
	EnrollPW       					string `json:"enrollPW"`
	MspID          					string `json:"mspID"`
	Namespace      					string `json:"namespace"`
	CaNameTls      					string `json:"caNameTls"`
	CaName         					string `json:"caName"`
	FileOutput     					string `json:"fileOutput"`
	FileOutputTls  					string `json:"fileOutputTls"`
	OrderNodeHost						string `json:"orderNodeHost"`
	OrdererNodeEndpoint 	[]string `json:"ordererNodeEndpoint"`
	OrdererNodesList			[]string `json:"ordererNodesList"`
}

type JoinChannel struct {
	Namespace      					string `json:"namespace"`
	MspID          				[]string `json:"mspID"`
	FileOutputTls  				[]string `json:"fileOutputTls"`
	FabricChannelFollower []string `json:"fabricChannelFollower"`
	PeersToJoin						[][]string `json:"peersToJoin"`
	AnchorPeers						[][]string 	`json:"anchorPeers"`
	OrderNodeHost 				[][]string `json:"orderNodeHost"`
	OrdererNodesList			[][]string `json:"ordererNodesList"`
}

type Config struct {
	Orgs     []string  `json:"OrgsPeer"`
	CAs      []CA      `json:"CA"`
	Peers    []Peer    `json:"Peers"`
	Orderers []Orderer `json:"Orderer"`
	Channels []Channel `json:"Channel"`
}
