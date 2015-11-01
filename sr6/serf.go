package sr6

import (
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
		if err := s.hosts.update(m.Addr.String(), m.Name); err != nil {
			log.Printf("[ERR] Couldn't update hosts file, %#v", err)
		}
	}
}
