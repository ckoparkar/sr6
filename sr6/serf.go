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
		// if we are leader, add host to ambari
		if s.isLeader() {
			// 1) bootstrap new node (POST http://localhost:8080/api/v1/bootstrap)
			// 2) poll till its done
			// 3) add host (POST /clusters/:clusterName/hosts/:hostName)
		}
	}
}

func (s *Server) nodeLeave(me serf.MemberEvent) {
	for _, m := range me.Members {
		if err := s.hosts.remove(m.Addr.String(), m.Name); err != nil {
			log.Printf("[ERR] Couldn't remove host , %#v", err)
		}
	}
}
