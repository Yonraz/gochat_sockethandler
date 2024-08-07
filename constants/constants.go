package constants

type Queues string
type RoutingKey string
type Exchange string
type MessageType string

const (
	UserEventsExchange    Exchange = "UserEventsExchange"
	MessageEventsExchange Exchange = "MessageEventsExchange"
)

const (
	UserLoggedInKey     RoutingKey = "user.logged.in"
	UserSignedoutKey    RoutingKey = "user.signed.out"
	MessageSentKey      RoutingKey = "message.sent"
	MessageDeliveredKey RoutingKey = "message.delivered"
	MessageReadKey      RoutingKey = "message.read"
)

const (
	Notifications_UserLoginQueue        Queues = "NOTIFICATIONS_SRV_UserLoginQueue"
	Notifications_UserSignoutQueue      Queues = "NOTIFICATIONS_SRV_UserSignoutQueue"
	Notifications_MessageSentQueue      Queues = "NOTIFICATIONS_SRV_MessageSentQueue"
	Notifications_MessageDeliveredQueue Queues = "NOTIFICATIONS_SRV_MessageDeliveredQueue"
	Notifications_MessageReadQueue      Queues = "NOTIFICATIONS_SRV_MessageReadQueue"
)

const (
	Messages_UserLoginQueue        Queues = "MESSAGES_SRV_UserLoginQueue"
	Messages_UserSignoutQueue      Queues = "MESSAGES_SRV_UserSignoutQueue"
	Messages_MessageSentQueue      Queues = "MESSAGES_SRV_MessageSentQueue"
	Messages_MessageDeliveredQueue Queues = "MESSAGES_SRV_MessageDeliveredQueue"
	Messages_MessageReadQueue      Queues = "MESSAGES_SRV_MessageReadQueue"
)

const (
	MessageUpdate MessageType = "message.update"
	MessageCreate MessageType = "message.create"
)