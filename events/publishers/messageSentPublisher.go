package publishers

import "github.com/yonraz/gochat_sockethandler/constants"


func (p *Publisher) MessageEvent(routingKey constants.RoutingKey, body interface{}) error {
	return p.Publish(constants.MessageEventsExchange, routingKey, body)
}