package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/netip"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/voldmir/adnettest/internal/dns"
	"github.com/voldmir/adnettest/internal/krb"
	"github.com/voldmir/adnettest/internal/ldap"
	"github.com/voldmir/adnettest/internal/netbios"
	"github.com/voldmir/adnettest/internal/ntp"
	"github.com/voldmir/adnettest/internal/proto"
	"github.com/voldmir/adnettest/internal/rsync"
	"github.com/voldmir/adnettest/internal/smb"
	"github.com/voldmir/adnettest/internal/www"
)

/*
   Порт сервера      Служба
   -----------------------------
   53/TCP/UDP        DNS !
   88/TCP/UDP        Kerberos !
   123/UDP           W32Time !
   135/TCP           Сопоставитель конечных точек RPC
   137/UDP           netbios-ns !
   138/UDP           netbios-dgm !
   139/TCP           netbios-ssn !
   389/TCP/UDP       LDAP !
   445/TCP           SMB !
   464/TCP/UDP       Изменение пароля в Kerberos !
   636/TCP           LDAP SSL !
   3268/TCP          LDAP GC !
   3269/TCP          LDAP GC SSL !
   49152-65535/TCP   FRS RPC
*/
/** ******************************************** **/

func main() {
	var test_user, test_computer_name, domain, dnsServer string

	flag.StringVar(&domain, "domain", "", "-domain suffix domain")
	flag.StringVar(&dnsServer, "dns", "", "-dns <ip addr server[:port]>")

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "%s [-domain <example.com>] [host.example.com] [192.168.0.1] [192.168.0.2] [192.168.0.3] ...\n", filepath.Base(os.Args[0]))

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "    %v\n", f.Usage)
		})
		fmt.Fprint(os.Stderr, `The program tries to find the domain suffix in the PC name if the -domain parameter is not specified.
If the -domain parameter is specified and the server address is not passed, then an attempt will be made to determine the domain controllers from the SRV record.`)
	}

	flag.Parse()
	servers := flag.Args()

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	userInfo, err := user.Current()
	if err != nil {
		panic(err)
	}

	if strings.Contains(userInfo.Username, "\\") {
		test_user = strings.SplitN(userInfo.Username, "\\", 2)[1]
	} else {
		test_user = strings.SplitN(userInfo.Username, "@", 2)[0]
	}

	if strings.Contains(hostname, ".") {
		s := strings.SplitN(hostname, ".", 2)
		test_computer_name, domain = s[0], s[1]
	} else {
		test_computer_name = hostname
		if domain == "" {
			panic("used flag -domain <suffix domain>")
		}
	}

	var resolver *net.Resolver

	if dnsServer != "" {
		if !strings.Contains(dnsServer, ":") {
			dnsServer = fmt.Sprintf("%s:53", dnsServer)
		}

		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, "udp", dnsServer)
			},
		}
	} else {
		resolver = net.DefaultResolver
	}

	ctx := context.Background()

	if len(servers) < 1 {
		servers, err = getControllersDomainBySyffixDomain(context.WithoutCancel(ctx), domain, resolver)
		if err != nil {
			panic(err.Error())
		}
	}

	servers_ips, err := resolveAddress(context.WithoutCancel(ctx), servers, resolver)
	if err != nil {
		panic(err.Error())
	}

	if len(servers_ips) < 1 {
		panic("no servers presented")
	}

	services := []proto.Proto{
		&dns.DNSPacket{},  // 53
		&www.HTTPPacket{}, // 80
		&krb.KerberosPacket{ // 88
			Realm: domain,
			SPN:   test_user,
		},
		&ntp.NTPPacket{},      //123
		&netbios.NBNSPacket{}, //137
		&netbios.NBDgramPacket{ // 138
			ComputerName: test_computer_name,
			GroupName:    strings.Split(domain, ".")[0],
		},
		&netbios.NBSSNPacket{}, // 139
		&ldap.LDAPPacket{},     // 389
		&www.HTTPSPacket{},     // 443
		&smb.SMBPacket{},       // 445
		&krb.KpaswdPacket{ // 464
			Realm: domain,
			SPN:   test_user,
		},
		&ldap.LDAPSPacket{},   // 636
		&rsync.RSYNCPacket{},  // 873
		&ldap.LDAPGCPacket{},  // 3268
		&ldap.LDAPGCSPacket{}, // 3269
	}

	for _, srv := range servers_ips {
		fmt.Println("")
		fmt.Println("Host: ", srv)
		for _, s := range services {
			proto.TestService(srv, s, 5)
		}
	}

	fmt.Println("")
	os.Exit(0)
}

func checkIpAddress(address string) (string, error) {
	addr_ip, err := netip.ParseAddr(address)
	if err == nil {
		return addr_ip.String(), nil
	}

	return "", fmt.Errorf("invalid address %s", address)
}

func resolveAddress(ctx context.Context, addresses []string, r *net.Resolver) ([]string, error) {
	var regex = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	var addressesPrepare []string

	for _, address := range addresses {

		if regex.MatchString(address) {
			ip, err := checkIpAddress(address)
			if err != nil {
				return nil, err
			}
			addressesPrepare = append(addressesPrepare, ip)
		} else {
			addrResolv, err := r.LookupHost(ctx, address)
			if err != nil {
				return nil, err
			}
			addressesPrepare = append(addressesPrepare, addrResolv...)
		}
	}

	return addressesPrepare, nil
}

func getControllersDomainBySyffixDomain(ctx context.Context, syffix string, r *net.Resolver) ([]string, error) {
	var targets []string

	_, addrs, err := r.LookupSRV(ctx, "ldap", "tcp", syffix)
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		host, _ := strings.CutSuffix(addr.Target, ".")
		targets = append(targets, host)
	}

	return targets, nil
}
