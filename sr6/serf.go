package sr6

import (
	"fmt"
	"log"

	"github.com/hashicorp/serf/serf"
)

// listenSerfEvents handles incoming serf events
func (s *Server) serfEventHandler() {
	for {
		select {
		case e := <-s.eventChLAN:
			switch e.EventType() {
			case serf.EventMemberJoin:
				s.nodeJoin(e.(serf.MemberEvent))
			case serf.EventMemberLeave:
				s.nodeLeave(e.(serf.MemberEvent))
			default:
				log.Printf("[WARN] sr6: unhandled LAN Serf Event: %#v", e)
			}
		case <-s.shutdownCh:
			return
		}
	}
}

func (s *Server) nodeJoin(me serf.MemberEvent) {
	for _, m := range me.Members {
		if err := s.hosts.add(m.Addr.String(), m.Name); err != nil {
			log.Printf("[ERR] Couldn't add host , %#v", err)
		}
	}
}

func (s *Server) nodeLeave(me serf.MemberEvent) {
	for _, m := range me.Members {
		// TODO(cskksc): remove entry from host
		fmt.Printf("%s is leaving the cluster.", m.Addr.String())
	}
}
