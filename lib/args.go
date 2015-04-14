package lib

type VNCArgs struct {
	UUID       string `form:"-" bson:"uuid"`
	Type       string `form:"-" bson:"type"`
	Hostname   string `valid:"Required" form:"hostname" bson:"hostname"`
	Port       string `valid:"Required; Numeric" form:"port" bson:"port"`
	Width      string `valid:"Required; Numeric" form:"width" bson:"width"`
	Height     string `valid:"Required; Numeric" form:"height" bson:"height"`
	DPI        string `valid:"Required; Numeric" form:"dpi" bson:"dpi"`
	ColorDepth string `valid:"Required; Numeric" form:"color-depth" bson:"color-depth"`
}

type RDPArgs struct {
	UUID          string `form:"-" bson:"uuid"`
	Type          string `form:"-" bson:"type"`
	Hostname      string `valid:"Required" form:"hostname" bson:"hostname"`
	Port          string `valid:"Required; Numeric" form:"port" bson:"port"`
	Width         string `valid:"Required; Numeric" form:"width" bson:"width"`
	Height        string `valid:"Required; Numeric" form:"height" bson:"height"`
	DPI           string `valid:"Required; Numeric" form:"dpi" bson:"dpi"`
	ColorDepth    string `valid:"Required; Numeric" form:"color-depth" bson:"color-depth"`
	Console       string `form:"console" bson:"console"`
	IntialProgram string `form:"initial-program" bson:"initial-program"`
	RemoteApp     string `form:"remote-app" bson:"remote-app"`
	RemoteAppDirs string `form:"remote-app-dirs" bson:"remote-app-dirs"`
	RemoteAppArgs string `form:"remote-app-args" bson:"remote-app-args"`
}

type SSHArgs struct {
	UUID       string `form:"-" bson:"uuid"`
	Type       string `form:"-" bson:"type"`
	Hostname   string `valid:"Required" form:"hostname" bson:"hostname"`
	Port       string `valid:"Required; Numeric" form:"port" bson:"port"`
	Width      string `valid:"Required; Numeric" form:"width" bson:"width"`
	Height     string `valid:"Required; Numeric" form:"height" bson:"height"`
	DPI        string `valid:"Required; Numeric" form:"dpi" bson:"dpi"`
	PrivateKey string `form:"private-key" bson:"private-key"`
	Passphrase string `form:"passphrase" bson:"passphrase"`
}

type TELNETArgs struct {
	UUID          string `form:"-" bson:"uuid"`
	Type          string `form:"-" bson:"type"`
	Hostname      string `valid:"Required" form:"hostname" bson:"hostname"`
	Port          string `valid:"Required; Numeric" form:"port" bson:"port"`
	Width         string `valid:"Required; Numeric" form:"width" bson:"width"`
	Height        string `valid:"Required; Numeric" form:"height" bson:"height"`
	DPI           string `valid:"Required; Numeric" form:"dpi" bson:"dpi"`
	UsernameRegex string `form:"username-regex" bson:"username-regex"`
	PasswordRegex string `form:"password-regex" bson:"password-regex"`
}

type SPICEArgs struct {
	UUID     string `form:"-" bson:"uuid"`
	Type     string `form:"-" bson:"type"`
	Hostname string `valid:"Required" form:"hostname" bson:"hostname"`
	Port     string `valid:"Required; Numeric" form:"port" bson:"port"`
	Password string `form:"password" bson:"password"`
}

type LIBVIRTArgs struct {
	UUID     string `form:"-" bson:"uuid"`
	Type     string `form:"-" bson:"type"`
	Hostname string `valid:"Required" form:"hostname" bson:"hostname"`
	Port     string `valid:"Required; Numeric" form:"port" bson:"port"`
	VM       string `valid:"Required" form:"vm" bson:"vm"`
	Shared   string `valid:"Required" form:"shared" bson:"shared"`
}
