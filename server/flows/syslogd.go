package flows

// Runs a syslog daemon to receive logs from the devices
import (
	"fmt"
	"regexp"
	"server/logger"
	"server/messages"

	"github.com/cnaude/go-syslog/syslog/v3"
)

var allRegex = `\S* (\S*) IN=(\S*) OUT=(\S*) MAC source = (\S\S:\S\S:\S\S:\S\S:\S\S:\S\S) MAC dest = (\S\S:\S\S:\S\S:\S\S:\S\S:\S\S) proto = (\S*)(?: IP SRC=)?(\S*)?(?: IP DST=)?([0-9\.]*)?(?:[\,])?(?: IP tos=)?(?:\S*)?(?: IP proto=)?(\S*)?(?: SPT=)?(\S*)?(?: DPT=)?(\S*)?`

func Syslogd() {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.Automatic)
	server.SetHandler(handler)
	server.ListenUDP("0.0.0.0:8514")
	server.Boot()

	var allSet = regexp.MustCompile(allRegex)

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			logger.Log("trace", "Syslog client", fmt.Sprintf("%s", logParts["client"]))

			result := allSet.FindStringSubmatch(fmt.Sprintf("%s", logParts["content"]))

			logentry := messages.FlowEntry{Hostname: logParts["hostname"].(string), IpAddress: logParts["client"].(string)}
			if len(result) == 12 {
				logentry.Prefix = result[1]
				logentry.IfaceIn = result[2]
				logentry.IfaceOut = result[3]
				logentry.SrcMac = result[4]
				logentry.DstMac = result[5]
				logentry.HwProto = result[6]
				logentry.SrcIp = result[7]
				logentry.DstIp = result[8]
				logentry.IpProto = result[9]
				logentry.SrcPort = result[10]
				logentry.DstPort = result[11]
				logger.Log("trace", "Log entry", fmt.Sprintf("%+v", logentry))
				addEntry(logentry)
			}
		}
	}(channel)

	server.Wait()
}

// \[[0-9]*[.][0-9]*\] (\S*) IN=(\S*) OUT=(\S*) MAC source = (\S\S:\S\S:\S\S:\S\S:\S\S:\S\S) MAC dest = (\S\S:\S\S:\S\S:\S\S:\S\S:\S\S) proto = (\S*) IP SRC=(\S*) IP DST=(\S*)\, IP tos=(\S*) IP proto=(\S*)(?: SPT=)?(\S*)?(?: DPT=)?(\S*)?
