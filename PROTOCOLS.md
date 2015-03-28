#### VNC

Name | Description
---|---
hostname | The hostname or IP address of the VNC server RConsole should connect to.
port | The port the VNC server is listening on, usually 5900 or 5900 + display number. For example, if your VNC server is serving display number 1 (sometimes written as :1), your port number here would be 5901.
password | The password to use when attempting authentication, if any. This parameter is optional.
read-only | Whether this connection should be read-only. If set to "true", no input will be accepted on the connection at all. Users will only see the desktop and whatever other users using that same desktop are doing. This parameter is optional.
swap-red-blue | If the colors of your display appear wrong (blues appear orange or red, etc.), it may be that your VNC server is sending image data incorrectly, and the red and blue components of each color are swapped. If this is the case, set this parameter to "true" to work around the problem. This parameter is optional.
color-depth | The color depth to request, in bits-per-pixel. This parameter is optional. If specified, this must be either 8, 16, 24, or 32. Regardless of what value is chosen here, if a particular update uses less than 256 colors, RConsole will always send that update as a 256-color PNG.
cursor | If set to "remote", the mouse pointer will be rendered remotely, and the local position of the mouse pointer will be indicated by a small dot.
autoretry | The number of times to retry connecting before giving up and returning an error. In the case of a reverse connection, this is the number of times the connection process is allowed to time out.
encodings | A space-delimited list of VNC encodings to use. The format of this parameter is dictated by libvncclient and thus doesn't really follow the form of other RConsole parameters. This parameter is optional, and libguac-client-vnc will use any supported encoding by default.Beware that this parameter is intended to be replaced with individual, encoding-specific parameters in a future release.
dest-host | The destination host to request when connecting to a VNC proxy such as UltraVNC Repeater. This is only necessary if the VNC proxy in use requires the connecting user to specify which VNC server to connect to. If the VNC proxy automatically connects to a specific server, this parameter is not necessary.
dest-port | The destination port to request when connecting to a VNC proxy such as UltraVNC Repeater. This is only necessary if the VNC proxy in use requires the connecting user to specify which VNC server to connect to. If the VNC proxy automatically connects to a specific server, this parameter is not necessary.
enable-audio | If set to "true", experimental sound support will be enabled. VNC does not support sound, but RConsole's VNC support can include sound using PulseAudio.Most Linux systems provide audio through a service called PulseAudio. This service is capable of communicating over the network. If PulseAudio is configured to allow TCP connections, RConsole can connect to your PulseAudio server and combine its audio with the graphics coming over VNC.Beware that you must disable authentication within PulseAudio in order to allow RConsole to connect, as RConsole does not yet support this. The amount of latency you will see depends largely on the network and how PulseAudio is configured.
audio-servername | The name of the PulseAudio server to connect to. This will be the hostname of the computer providing audio for your connection via PulseAudio, most likely the same as the value given for the hostname parameter.If this parameter is omitted, the default PulseAudio device will be used, which will be the PulseAudio server running on the same machine as this service.
reverse-connect	| Whether reverse connection should be used. If set to "true", instead of connecting to a server at a given hostname and port, this service will listen on the given port for inbound connections from a VNC server.
listen-timeout | If reverse connection is in use, the maximum amount of time to wait for an inbound connection from a VNC server, in milliseconds. If blank, the default value is 5000 (five seconds).


### RDP

Name | Description
---|---
hostname | The hostname or IP address of the RDP server RConsole should connect to.
port | The port the RDP server is listening on, usually 3389. This parameter is optional. If this is not specified, the default of 3389 will be used.
username | The username to use to authenticate, if any. This parameter is optional.
password | The password to use when attempting authentication, if any. This parameter is optional.
domain | The domain to use when attempting authentication, if any. This parameter is optional.
color-depth | The color depth to request, in bits-per-pixel. This parameter is optional. If specified, this must be either 8, 16, or 24. Regardless of what value is chosen here, if a particular update uses less than 256 colors, RConsole will always send that update as a 256-color PNG.
width | The width of the display to request, in pixels. This parameter is optional. If this value is not specified, the width of the connecting client display will be used instead.
height | The height of the display to request, in pixels. This parameter is optional. If this value is not specified, the height of the connecting client display will be used instead.
dpi | The desired effective resolution of the client display, in DPI. This parameter is optional. If this value is not specified, the resolution and size of the client display will be used together to determine, heuristically, an appropriate resolution for the RDP session.
disable-audio | Audio is enabled by default in both the client and in libguac-client-rdp. If you are concerned about bandwidth usage, or sound is causing problems, you can explicitly disable sound by setting this parameter to "true".
enable-printing | Printing is disabled by default, but with printing enabled, RDP users can print to a virtual printer that sends a PDF containing the document printed to the RConsole client. Enable printing by setting this parameter to "true".Printing support requires GhostScript to be installed. If this service cannot find the gs executable when printing, the print attempt will fail.
enable-drive | File transfer is disabled by default, but with file transfer enabled, RDP users can transfer files to and from a virtual drive which persists on the RConsole server. Enable file transfer support by setting this parameter to "true".Files will be stored in the directory specified by the "drive-path" parameter, which is required if file transfer is enabled.
drive-path | The directory on the RConsole server in which transfered files should be stored. This directory must be accessible by this service and both readable and writable by the user that runs this service. This parameter does not refer to a directory on the RDP server.If file transfer is not enabled, this parameter is ignored.
console	| If set to "true", you will be connected to the console (admin) session of the RDP server.
console-audio | If set to "true", audio will be explicitly enabled in the console (admin) session of the RDP server. Setting this option to "true" only makes sense if the console parameter is also set to "true".
initial-program	| The full path to the program to run immediately upon connecting. This parameter is optional.
server-layout | The server-side keyboard layout. This is the layout of the RDP server and has nothing to do with the keyboard layout in use on the client. The RConsole client is independent of keyboard layout. The RDP protocol, however, is not independent of keyboard layout, and RConsole needs to know the keyboard layout of the server in order to send the proper keys when a user is typing.Possible values are: **en-us-qwerty**: English (US) keyboard   **de-de-qwertz**: German keyboard (qwertz) **fr-fr-azerty**: French keyboard (azerty)  **failsafe**: Unknown keyboard - this option sends only Unicode events and should work for any keyboard, though not necessarily all RDP servers or applications.If your server's keyboard layout is not yet supported, this option should work in the meantime.
security | The security mode to use for the RDP connection. This mode dictates how data will be encrypted and what type of authentication will be performed, if any. By default, the server is allowed to control what type of security is used.Possible values are: **rdp**: Standard RDP encryption. This mode should be supported by all RDP servers. **nla:** Network Level Authentication. This mode requires the username and password, and performs an authentication step before the remote desktop session actually starts. If the username and password are not given, the connection cannot be made. **tls:**TLS encryption. TLS (Transport Layer Security) is the successor to SSL. **any:** Allow the server to choose the type of security. This is the default. 
ignore-cert | If set to "true", the certificate returned by the server will be ignored, even if that certificate cannot be validated. This is useful if you universally trust the server and your connection to the server, and you know that the server's certificate cannot be validated (for example, if it is self-signed).
disable-auth | If set to "true", authentication will be disabled. Note that this refers to authentication that takes place while connecting. Any authentication enforced by the server over the remote desktop session (such as a login dialog) will still take place. By default, authentication is enabled and only used when requested by the server.If you are using NLA, authentication must be enabled by definition.
remote-app | Specifies the RemoteApp to start on the remote desktop. If supported by your remote desktop server, this application, and only this application, will be visible to the user.Windows requires a special notation for the names of remote applications. The names of remote applications must be prefixed with two vertical bars. For example, if you have created a remote application on your server for notepad.exe and have assigned it the name "notepad", you would set this parameter to: "IInotepad".
remote-app-dir | The working directory, if any, for the remote application. This parameter has no effect if RemoteApp is not in use.
remote-app-args	| The command-line arguments, if any, for the remote application. This parameter has no effect if RemoteApp is not in use.
static-channels	| A comma-separated list of static channel names to open and expose as pipes. If you wish to communicate between an application running on the remote desktop and JavaScript, this is the best way to do it. RConsole will open an outbound pipe with the name of the static channel. If JavaScript needs to communicate back in the other direction, it should respond by opening another pipe with the same name.RConsole allows any number of static channels to be opened, but protocol restrictions of RDP limit the size of each channel name to 7 characters.


### SSH

Name | Description
---|---
hostname | The hostname or IP address of the SSH server RConsole should connect to.
port | The port the SSH server is listening on, usually 22. This parameter is optional. If this is not specified, the default of 22 will be used.
username | The username to use to authenticate, if any. This parameter is optional. If not specified, you will be prompted for the username upon connecting.
password | The password to use when attempting authentication, if any. This parameter is optional. If not specified, you will be prompted for your password upon connecting.
font-name | The name of the font to use. This parameter is optional. If not specified, the default of "monospace" will be used instead.
font-size | The size of the font to use, in points. This parameter is optional. If not specified, the default of 12 will be used instead.
enable-sftp | Whether file transfer should be enabled. If set to "true", the user will be allowed to upload or download files from the SSH server using SFTP. RConsole includes the guacctl utility which controls file downloads and uploads when run on the SSH server by the user over the SSH connection.
private-key | The entire contents of the private key to use for public key authentication. If this parameter is not specified, public key authentication will not be used. The private key must be in OpenSSH format, as would be generated by the OpenSSH ssh-keygen utility.
passphrase | The passphrase to use to decrypt the private key for use in public key authentication. This parameter is not needed if the private key does not require a passphrase. If the private key requires a passphrase, but this parameter is not provided, the user will be prompted for the passphrase upon connecting.



### TELNET
Name | Description
---|---
hostname | The hostname or IP address of the telnet server RConsole should connect to.
port | The port the telnet server is listening on, usually 23. This parameter is optional. If this is not specified, the default of 23 will be used.
username | The username to use to authenticate, if any. This parameter is optional. If not specified, or not supported by the telnet server, the login process on the telnet server will prompt you for your credentials. For this to work, your telnet server must support the NEW-ENVIRON option, and the telnet login process must pay attention to the USER environment variable. Most telnet servers satisfy this criteria.
username-regex | The regular expression to use when waiting for the username prompt. This parameter is optional. If not specified, a reasonable default built into RConsole will be used. The regular expression must be written in the POSIX ERE dialect (the dialect typically used by egrep).
password | The password to use when attempting authentication, if any. This parameter is optional. If specified, your password will be typed on your behalf when the password prompt is detected.
password-regex | The regular expression to use when waiting for the password prompt. This parameter is optional. If not specified, a reasonable default built into RConsole will be used. The regular expression must be written in the POSIX ERE dialect (the dialect typically used by egrep).
font-name | The name of the font to use. This parameter is optional. If not specified, the default of "monospace" will be used instead.
font-size | The size of the font to use, in points. This parameter is optional. If not specified, the default of 12 will be used instead.

