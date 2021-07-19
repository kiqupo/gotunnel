# Socks5协议
## Socks5协议解析之授权认证
如果要与socks5服务器建立TCP连接，客户端需要先发起请求来对协议的版本及其认证方式。这里就是客户端请求服务器的请求格式：

           +----+----------+----------+
           |VER | NMETHODS | METHODS  |
           +----+----------+----------+
           | 5  |    1     | 1 to 255 |
           +----+----------+----------+
VER这里指定的就是socks的版本，如果使用的是socks5的话这里的值需要是0x05，其固定长度为1个字节；
第二个字段NMETHODS表示第三个字段METHODS的长度，它的长度也是1个字节；
第三个METHODS表示客户端支持的验证方式，可以有多种，他的尝试是1-255个字节；
目前支持的验证方式一共有：

**X’00’ **NO AUTHENTICATION REQUIRED（不需要验证）
**X’01’ **GSSAPI
**X’02’** USERNAME/PASSWORD（用户名、密码认证）
**X’03’ ** to X’7F’ IANA ASSIGNED （ 由IANA分配（保留）
**X’80’** to X’FE’ RESERVED FOR PRIVATE METHODS 私人方法保留）
**X’FF’** NO ACCEPTABLE METHODS（都不支持，没法连接了）
以上的都是十六进制常量，比如X’00’表示十六进制0x00。
当服务器端收到了客户端的请求之后，就会在响应客户端的请求，服务端需要客户端提供哪种验证方式的信息。

                 +----+--------+
                 |VER | METHOD |
                 +----+--------+
                 | 5  |   0    |
                 +----+--------+
第一个字段VER代表Socket的版本，Soket5默认为0x05，其值长度为1个字节
第二个字段METHOD代表需要服务端需要客户端按照此验证方式提供验证信息，其值长度为1个字节，选择为上面的六种验证方式。

## 服务端判断并返回认证方法
``` golang
func Socks5Auth(client net.Conn) (err error) {
	buf := make([]byte, 256)

	// 读取 VER 和 NMETHODS
	n, err := io.ReadFull(client, buf[:2])
	if n != 2 {
		return errors.New("reading header: " + err.Error())
	}

	ver, nMethods := int(buf[0]), int(buf[1])
	if ver != 5 {
		return errors.New("invalid version")
	}

	// 读取 METHODS 列表
	n, err = io.ReadFull(client, buf[:nMethods])
	if n != nMethods {
		return errors.New("reading methods: " + err.Error())
	}

	//无需认证
	n, err = client.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		return errors.New("write rsp err: " + err.Error())
	}
	return nil
}
``` 

## Socks5协议解析之建立连接
Socket5的客户端和服务端进行双方授权验证通过之后，就开始建立连接了。连接由客户端发起，告诉Sokcet服务端客户端需要访问哪个远程服务器，其中包含，远程服务器的地址和端口，地址可以是IP4，IP6，也可以是域名。

	+----+-----+-------+------+----------+----------+
	|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
**VER**代表Socket协议的版本，Soket5默认为0x05，其值长度为1个字节
**CMD**代表客户端请求的类型，值长度也是1个字节，有三种类型
1. CONNECT X’01’
2. BIND X’02’
3. UDP ASSOCIATE X’03’

**RSV**保留字，值长度为1个字节
**ATYP**代表请求的远程服务器地址类型，值长度1个字节，有三种类型
1. IP V4 address: X’01’
2. DOMAINNAME: X’03’
3. IP V6 address: X’04’

**DST.ADDR**代表远程服务器的地址，根据ATYP进行解析，值长度不定。
**DST.PORT**代表远程服务器的端口，要访问哪个端口的意思，值长度2个字节
接着客户端把要请求的远程服务器的信息都告诉Socket5代理服务器了，那么Socket5代理服务器就可以和远程服务器建立连接了，不管连接是否成功等，都要给客户端回应，其回应格式为：

	+----+-----+-------+------+----------+----------+
	|VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
**VER**代表Socket协议的版本，Soket5默认为0x05，其值长度为1个字节
**REP**代表响应状态码，值长度也是1个字节，有以下几种类型
**X’00’** succeeded
**X’01’** general SOCKS server failure
**X’02’** connection not allowed by ruleset
**X’03’** Network unreachable
**X’04’** Host unreachable
**X’05’** Connection refused
**X’06’** TTL expired
**X’07’** Command not supported
**X’08’** Address type not supported
**X’09’** to X’FF’ unassigned
**RSV**保留字，值长度为1个字节
**ATYP**代表请求的远程服务器地址类型，值长度1个字节，有三种类型
1. IP V4 address: X’01’
2. DOMAINNAME: X’03’
3. IP V6 address: X’04’

**BND.ADDR**表示绑定地址，值长度不定。
**BND.PORT**表示绑定端口，值长度2个字节
服务端响应客户端连接成功就会返回如下的数据给客户端。
``` golang
n, err = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
``` 