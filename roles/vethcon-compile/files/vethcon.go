/*

Beautifully crafted by Tomofumi Hayashi (@s1061123 on github)

*/
package main

import ("os"
  "fmt"
  "net"
  "github.com/vishvananda/netlink"
  "github.com/vishvananda/netns"

  "github.com/docker/docker/client"
//  "github.com/docker/docker/api/types"
//  "github.com/docker/docker/api/types/container"
  "golang.org/x/net/context"
)

func makeVethPair(name, peer string, mtu int) (netlink.Link, error) {
        veth := &netlink.Veth{
                LinkAttrs: netlink.LinkAttrs{
                        Name:  name,
                        Flags: net.FlagUp,
                        MTU:   mtu,
                },
                PeerName: peer,
        }

        if err := netlink.LinkAdd(veth); err != nil {
                return nil, err
        }
        return veth, nil
}

func getVethPair(name1 string, name2 string) (veth1 netlink.Link, veth2 netlink.Link, err error) {

  veth1, err = makeVethPair(name1, name2, 1500)
  if err != nil {
    switch {
    case os.IsExist(err):
      err = fmt.Errorf("container veth name provided (%v) already exists\n", name1)
      return
    default:
      err = fmt.Errorf("failed to make veth pair: %v\n", err)
      return
    }
  }

  veth2, err = netlink.LinkByName(name2)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to lookup %q: %v\n", name2, err)
  }

  return
}

func getContainerNS(containerId string) (namespace string, err error) {
  ctx := context.Background()
  cli, err := client.NewEnvClient()
  if err != nil {
    panic(err)
  }

  cli.UpdateClientVersion("1.24")

  json, err := cli.ContainerInspect(ctx, containerId)
  if err != nil {
    err = fmt.Errorf("failed to get container info: %v\n", err)
  }
  namespace = json.NetworkSettings.NetworkSettingsBase.SandboxKey
  return
}


func main() {

  if len(os.Args) < 5 {
    fmt.Printf("Usage: %s <Container1> <Container2> <Link1Name> <Link2Name>\n", os.Args[0])
    os.Exit(1)
  }

  containerId1 := os.Args[1]
  containerId2 := os.Args[2]
  linkName1 := os.Args[3]
  linkName2 := os.Args[4]

  nspath1, err := getContainerNS(containerId1)
  if err != nil {
    fmt.Fprintf(os.Stderr, "%v", err)
  }

  nspath2, err := getContainerNS(containerId2)
  if err != nil {
    fmt.Fprintf(os.Stderr, "%v", err)
  }

  nsh1, err := netns.GetFromPath(nspath1)
  if err != nil {
    fmt.Fprintf(os.Stderr, "%v", err)
  }
  defer nsh1.Close()
  
  nsh2, err := netns.GetFromPath(nspath2)
  if err != nil {
    fmt.Fprintf(os.Stderr, "%v", err)
  }
  defer nsh2.Close()

  veth1, veth2, err := getVethPair(linkName1, linkName2)
  if err != nil {
    fmt.Fprintf(os.Stderr, "%v", err)
  }

  if err = netlink.LinkSetUp(veth1); err != nil {
    fmt.Fprintf(os.Stderr,
      "Failed to set %q up: %v", linkName1, err)
  }
  if err = netlink.LinkSetUp(veth2); err != nil {
    fmt.Fprintf(os.Stderr,
      "Failed to set %q up: %v", linkName2, err)
  }

  if err = netlink.LinkSetNsFd(veth1, int(nsh1)); err != nil {
    fmt.Fprintf(os.Stderr, "%v", err)
  }

  if err = netlink.LinkSetNsFd(veth2, int(nsh2)); err != nil {
    fmt.Fprintf(os.Stderr, "%v", err)
  }

}