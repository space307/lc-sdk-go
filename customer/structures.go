package customer

import (
	"encoding/json"
	"time"

	"github.com/livechat/lc-sdk-go/v4/objects"
)

func unmarshalOptionalRawField(source json.RawMessage, target interface{}) error {
	if source != nil {
		return json.Unmarshal(source, target)
	}
	return nil
}

// Form struct describes schema of custom form (e-mail, prechat or postchat survey).
type Form struct {
	ID     string `json:"id"`
	Fields []struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Label    string `json:"label"`
		Required bool   `json:"required"`
		Options  []struct {
			ID    string `json:"id"`
			Type  int    `json:"group_id"`
			Label string `json:"label"`
		} `json:"options"`
	} `json:"fields"`
}

// PredictedAgent is an agent returned by GetPredictedAgent method.
type PredictedAgent struct {
	Agent struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar"`
		IsBot     bool   `json:"is_bot"`
		JobTitle  string `json:"job_title"`
		Type      string `json:"type"`
	} `json:"agent"`
	Queue bool `json:"queue"`
}

// URLInfo contains some OpenGraph info of the URL.
type URLInfo struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	URL              string `json:"url"`
	ImageURL         string `json:"image_url"`
	ImageOriginalURL string `json:"image_original_url"`
	ImageWidth       int    `json:"image_width"`
	ImageHeight      int    `json:"image_height"`
}

type DynamicConfiguration struct {
	GroupID             int    `json:"group_id"`
	ClientLimitExceeded bool   `json:"client_limit_exceeded"`
	DomainAllowed       bool   `json:"domain_allowed"`
	ConfigVersion       string `json:"config_version"`
	LocalizationVersion string `json:"localization_version"`
	Language            string `json:"language"`
}

type ConfigButton struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	OnlineValue  string `json:"online_value"`
	OfflineValue string `json:"offline_value"`
}

type Configuration struct {
	Buttons        []ConfigButton               `json:"buttons"`
	TicketForm     *Form                        `json:"ticket_form,omitempty"`
	PrechatForm    *Form                        `json:"prechat_form,omitempty"`
	AllowedDomains []string                     `json:"allowed_domains,omitempty"`
	Integrations   map[string]map[string]string `json:"integrations"`
	Properties     struct {
		Group   objects.Properties `json:"group"`
		License objects.Properties `json:"license"`
	} `json:"properties"`
}

type eventSpecific struct {
	Text            json.RawMessage `json:"text"`
	Fields          json.RawMessage `json:"fields"`
	ContentType     json.RawMessage `json:"content_type"`
	URL             json.RawMessage `json:"url"`
	Width           json.RawMessage `json:"width"`
	Height          json.RawMessage `json:"height"`
	Name            json.RawMessage `json:"name"`
	TemplateID      json.RawMessage `json:"template_id"`
	Elements        json.RawMessage `json:"elements"`
	Postback        json.RawMessage `json:"postback"`
	AlternativeText json.RawMessage `json:"alternative_text"`
}

// Event represents base of all LiveChat chat events.
//
// To get speficic event type's structure, call appropriate function based on Event's Type.
type Event struct {
	ID         string             `json:"id,omitempty"`
	CustomID   string             `json:"custom_id,omitempty"`
	CreatedAt  time.Time          `json:"created_at,omitempty"`
	AuthorID   string             `json:"author_id"`
	Properties objects.Properties `json:"properties,omitempty"`
	Recipients string             `json:"recipients,omitempty"`
	Type       string             `json:"type,omitempty"`
	eventSpecific
}

// FilledForm represents LiveChat filled form event.
type FilledForm struct {
	Fields []struct {
		Label string `json:"label"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"fields"`
	Event
}

// FilledForm function converts Event object to FilledForm object if Event's Type is "filled_form".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) FilledForm() *FilledForm {
	if e.Type != "filled_form" {
		return nil
	}
	var f FilledForm

	f.Event = *e
	if err := json.Unmarshal(e.Fields, &f.Fields); err != nil {
		return nil
	}
	return &f
}

// Postback represents postback data in LiveChat message event.
type Postback struct {
	ID       string `json:"id"`
	ThreadID string `json:"thread_id"`
	EventID  string `json:"event_id"`
	Type     string `json:"type,omitempty"`
	Value    string `json:"value,omitempty"`
}

// Message represents LiveChat message event.
type Message struct {
	Event
	Text     string    `json:"text,omitempty"`
	Postback *Postback `json:"postback,omitempty"`
}

// Message function converts Event object to Message object if Event's Type is "message".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) Message() *Message {
	if e.Type != "message" {
		return nil
	}
	var m Message

	m.Event = *e
	if err := json.Unmarshal(e.Text, &m.Text); err != nil {
		return nil
	}
	if err := unmarshalOptionalRawField(e.Postback, &m.Postback); err != nil {
		return nil
	}
	return &m
}

// SystemMessage represents LiveChat system message event.
type SystemMessage struct {
	Event
	Type     string            `json:"system_message_type,omitempty"`
	Text     string            `json:"text,omitempty"`
	TextVars map[string]string `json:"text_vars,omitempty"`
}

// File represents LiveChat file event
type File struct {
	Event
	ContentType     string `json:"content_type"`
	URL             string `json:"url"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	Name            string `json:"name"`
	AlternativeText string `json:"alternative_text"`
}

// File function converts Event object to File object if Event's Type is "file".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) File() *File {
	if e.Type != "file" {
		return nil
	}
	var f File
	f.Event = *e
	if err := json.Unmarshal(e.ContentType, &f.ContentType); err != nil {
		return nil
	}
	if err := json.Unmarshal(e.URL, &f.URL); err != nil {
		return nil
	}
	if err := unmarshalOptionalRawField(e.Width, &f.Width); err != nil {
		return nil
	}
	if err := unmarshalOptionalRawField(e.Height, &f.Height); err != nil {
		return nil
	}
	if err := json.Unmarshal(e.Name, &f.Name); err != nil {
		return nil
	}
	if err := unmarshalOptionalRawField(e.AlternativeText, &f.AlternativeText); err != nil {
		return nil
	}

	return &f
}

// RichMessage represents LiveChat rich message event
type RichMessage struct {
	Event
	TemplateID string               `json:"template_id"`
	Elements   []RichMessageElement `json:"elements"`
}

// RichMessageElement represents element of LiveChat rich message
type RichMessageElement struct {
	Buttons  []RichMessageButton `json:"buttons"`
	Title    string              `json:"title"`
	Subtitle string              `json:"subtitle"`
	Image    *RichMessageImage   `json:"image,omitempty"`
}

// RichMessageButton represents button in LiveChat rich message
type RichMessageButton struct {
	Text       string   `json:"text"`
	Type       string   `json:"type"`
	UserIds    []string `json:"user_ids"`
	Value      string   `json:"value"`
	PostbackID string   `json:"postback_id"`
	// Allowed values: compact, full, tall
	WebviewHeight string `json:"webview_height"`
	// Allowed values: new, current
	Target string `json:"target,omitempty"`
}

// RichMessageImage represents image in LiveChat rich message
type RichMessageImage struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	ContentType     string `json:"content_type"`
	Size            int    `json:"size"`
	Width           int    `json:"width,omitempty"`
	Height          int    `json:"height,omitempty"`
	AlternativeText string `json:"alternative_text,omitempty"`
}

// RichMessage function converts Event object to RichMessage object if Event's Type is "rich_message".
// If Type is different or Event is malformed, then it returns nil.
func (e *Event) RichMessage() *RichMessage {
	if e.Type != "rich_message" {
		return nil
	}
	var rm RichMessage

	rm.Event = *e
	if err := json.Unmarshal(e.TemplateID, &rm.TemplateID); err != nil {
		return nil
	}

	if err := json.Unmarshal(e.Elements, &rm.Elements); err != nil {
		return nil
	}

	return &rm
}
