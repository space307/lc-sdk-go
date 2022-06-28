package agent

import (
	"encoding/json"
	"time"

	"github.com/space307/lc-sdk-go/v4/objects"
)

type postback struct {
	ID      string `json:"id"`
	Toggled bool   `json:"toggled"`
}

type ban struct {
	Days uint `json:"days"`
}

type InitialChat struct {
	objects.InitialChat
	Users []*User `json:"users,omitempty"`
}

// MulticastRecipients aggregates Agent and Customer recipients that multicast should be sent to
type MulticastRecipients struct {
	Agents    *MulticastRecipientsAgents    `json:"agents,omitempty"`
	Customers *MulticastRecipientsCustomers `json:"customers,omitempty"`
}

// MulticastRecipientsAgents represents recipients for multicast to agents
type MulticastRecipientsAgents struct {
	Groups []uint   `json:"groups,omitempty"`
	IDs    []string `json:"ids,omitempty"`
	All    *bool    `json:"all,omitempty"`
}

// MulticastRecipientsCustomers represents recipients for multicast to customers
type MulticastRecipientsCustomers struct {
	IDs []string `json:"ids,omitempty"`
}

type transferTarget struct {
	Type string        `json:"type"`
	IDs  []interface{} `json:"ids"`
}

type routingStatusesFilter struct {
	GroupIDs []int `json:"group_ids,omitempty"`
}

type AgentsForTransfer []struct {
	AgentID          string `json:"agent_id"`
	TotalActiveChats uint   `json:"total_active_chats"`
}

// User represents base of both Customer and Agent
//
// To get speficic user type's structure, call Agent() or Customer() (based on Type value).
type User struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	Avatar         string    `json:"avatar"`
	Email          string    `json:"email"`
	Present        bool      `json:"present"`
	EventsSeenUpTo time.Time `json:"events_seen_up_to"`
	userSpecific
}

type userSpecific struct {
	RoutingStatus              json.RawMessage `json:"routing_status"`
	LastVisit                  json.RawMessage `json:"last_visit"`
	Statistics                 json.RawMessage `json:"statistics"`
	AgentLastEventCreatedAt    json.RawMessage `json:"agent_last_event_created_at"`
	CustomerLastEventCreatedAt json.RawMessage `json:"customer_last_event_created_at"`
	SessionFields              json.RawMessage `json:"session_fields"`
	Followed                   json.RawMessage `json:"followed"`
	Online                     json.RawMessage `json:"online"`
	State                      json.RawMessage `json:"state"`
	GroupIDs                   json.RawMessage `json:"group_ids"`
	EmailVerified              json.RawMessage `json:"email_verified"`
	CreatedAt                  json.RawMessage `json:"created_at"`
	Visibility                 json.RawMessage `json:"visibility"`
}

// Agent function converts User object to Agent object if User's Type is "agent".
// If Type is different or User is malformed, then it returns nil.
func (u *User) Agent() *Agent {
	if u.Type != "agent" {
		return nil
	}
	var a Agent

	a.User = u
	if err := json.Unmarshal(u.RoutingStatus, &a.RoutingStatus); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.Visibility, &a.Visibility); err != nil {
		return nil
	}
	return &a
}

// Customer function converts User object to Customer object if User's Type is "customer".
// If Type is different or User is malformed, then it returns nil.
func (u *User) Customer() *Customer {
	if u.Type != "customer" {
		return nil
	}
	var c Customer

	c.User = u
	if err := json.Unmarshal(u.LastVisit, &c.LastVisit); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.Statistics, &c.Statistics); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.AgentLastEventCreatedAt, &c.AgentLastEventCreatedAt); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.CustomerLastEventCreatedAt, &c.CustomerLastEventCreatedAt); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.EmailVerified, &c.EmailVerified); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.CreatedAt, &c.CreatedAt); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.Followed, &c.Followed); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.Online, &c.Online); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.State, &c.State); err != nil {
		return nil
	}
	if err := json.Unmarshal(u.SessionFields, &c.SessionFields); err != nil {
		return nil
	}
	return &c
}

// Agent represents LiveChat agent.
type Agent struct {
	*User
	RoutingStatus string `json:"routing_status"`
	Visibility    string `json:"visibility,omitempty"`
}

// Customer represents LiveChat customer.
type Customer struct {
	*User
	EmailVerified bool          `json:"email_verified"`
	LastVisit     objects.Visit `json:"last_visit"`
	Statistics    struct {
		VisitsCount            int `json:"visits_count"`
		ThreadsCount           int `json:"threads_count"`
		ChatsCount             int `json:"chats_count"`
		PageViewsCount         int `json:"page_views_count"`
		GreetingsShownCount    int `json:"greetings_shown_count"`
		GreetingsAcceptedCount int `json:"greetings_accepted_count"`
	} `json:"statistics"`
	AgentLastEventCreatedAt    time.Time           `json:"agent_last_event_created_at"`
	CustomerLastEventCreatedAt time.Time           `json:"customer_last_event_created_at"`
	CreatedAt                  time.Time           `json:"created_at"`
	SessionFields              []map[string]string `json:"session_fields"`
	Followed                   bool                `json:"followed"`
	Online                     bool                `json:"online"`
	State                      string              `json:"state"`
	GroupIDs                   []int               `json:"group_ids"`
}

// Chat represents LiveChat chat.
type Chat struct {
	ID         string             `json:"id,omitempty"`
	Properties objects.Properties `json:"properties,omitempty"`
	Access     objects.Access     `json:"access,omitempty"`
	Thread     objects.Thread     `json:"thread,omitempty"`
	Threads    []objects.Thread   `json:"threads,omitempty"`
	IsFollowed bool               `json:"is_followed,omitempty"`
	Agents     map[string]*Agent
	Customers  map[string]*Customer
}

// Users function returns combined list of Chat's Agents and Customers.
func (c *Chat) Users() []*User {
	u := make([]*User, 0, len(c.Agents)+len(c.Customers))
	for _, a := range c.Agents {
		u = append(u, a.User)
	}
	for _, cu := range c.Customers {
		u = append(u, cu.User)
	}

	return u
}

// UnmarshalJSON implements json.Unmarshaler interface for Chat.
func (c *Chat) UnmarshalJSON(data []byte) error {
	type ChatAlias Chat
	var cs struct {
		*ChatAlias
		Users []json.RawMessage `json:"users"`
	}

	if err := json.Unmarshal(data, &cs); err != nil {
		return err
	}

	var t struct {
		Type string `json:"type"`
	}

	*c = (Chat)(*cs.ChatAlias)
	c.Agents = make(map[string]*Agent)
	c.Customers = make(map[string]*Customer)
	for _, u := range cs.Users {
		if err := json.Unmarshal(u, &t); err != nil {
			return err
		}
		switch t.Type {
		case "agent":
			var a Agent
			if err := json.Unmarshal(u, &a); err != nil {
				return err
			}
			c.Agents[a.ID] = &a
		case "customer":
			var cu Customer
			if err := json.Unmarshal(u, &cu); err != nil {
				return err
			}
			c.Customers[cu.ID] = &cu
		}
	}

	return nil
}

// TransferChatOptions defines options for TransferChat method.
type TransferChatOptions struct {
	IgnoreRequesterPresence  bool
	IgnoreAgentsAvailability bool
}
