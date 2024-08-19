package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"zt-event-logger/pkg/db"

	"github.com/zerotier/ztchooks"
)

// Processor is an interface to mimic how an event processor should look like
type Processor interface {
	// Process processes the given webhook event payload
	Process(payload []byte, opts ...SignatureOpt) (*ztchooks.HookBase, error)
}

type processor struct {
	// dbClient is a DB client to process the given payload
	dbClient db.DB
	// signature is to check whether or not a payload is valid
	signature string
	// pre shared key (a.k.a psk) is a unique key to validate the payload
	preSharedKey string
}

func (p *processor) Process(payload []byte, opts ...SignatureOpt) (*ztchooks.HookBase, error) {
	// if signature info is provided via SignatureOpt(s) then it will be used to validate the
	// payload later
	p.processSignatureOpts(opts...)
	defer p.cleanSignatureInfo()

	err := p.verify(payload)
	if err != nil {
		return nil, fmt.Errorf("signature verification failed: %v", err)
	}

	// hook type is fetched
	hType, err := ztchooks.GetHookType(payload)
	if err != nil {
		return nil, fmt.Errorf("error when fetching the hook type: %v", err)
	}

	var hb ztchooks.HookBase
	var event db.Event

	switch hType {
	case ztchooks.NETWORK_JOIN:
		var nmj ztchooks.NewMemberJoined
		if err := json.Unmarshal(payload, &nmj); err != nil {
			return nil, fmt.Errorf("error marshalling NETWORK_JOIN event: %s", err)
		}

		hb = nmj.HookBase
		event = p.convertNewMemberJoinedToDBEvent(&nmj)

	case ztchooks.NETWORK_CREATED:
		var nc ztchooks.NetworkCreated
		if err := json.Unmarshal(payload, &nc); err != nil {
			return nil, fmt.Errorf("error marshalling NETWORK_CREATED event: %s", err)
		}

		hb = nc.HookBase
		event = p.convertNetworkCreatedToDBEvent(&nc)
	case ztchooks.NETWORK_CONFIG_CHANGED:
		var ncc ztchooks.NetworkConfigChanged
		if err := json.Unmarshal(payload, &ncc); err != nil {
			return nil, fmt.Errorf("error marshalling NETWORK_CONFIG_CHANGED event: %s", err)
		}

		hb = ncc.HookBase
		event = p.convertNetworkConfigChangedToDBEvent(&ncc)
	default:
		return nil, errors.New("unhandled event type")
	}

	err = p.dbClient.Insert(&event)
	if err != nil {
		return nil, fmt.Errorf("error inserting event into database: %s", err)
	}

	return &hb, nil
}

// processSignatureOpts processes all SignatureOpts if given
func (p *processor) processSignatureOpts(opts ...SignatureOpt) {
	for _, opt := range opts {
		opt(p)
	}
}

// cleanSignatureInfo cleans the signature and preSharedKey for the upcoming process requests
func (p *processor) cleanSignatureInfo() {
	p.signature = ""
	p.preSharedKey = ""
}

// verify verifies the payload
func (p *processor) verify(payload []byte) error {
	if p.signature != "" && p.preSharedKey != "" {
		err := ztchooks.VerifyHookSignature(p.preSharedKey, p.signature, payload, ztchooks.DefaultTolerance)
		if err != nil {
			return err
		}
	}
	return nil
}

// convertNewMemberJoinedToDBEvent converts NewMemberJoined webhook event to a db.Event struct
func (p *processor) convertNewMemberJoinedToDBEvent(nmj *ztchooks.NewMemberJoined) db.Event {
	return db.Event{
		HookID:    nmj.HookID,
		OrgID:     nmj.OrgID,
		HookType:  string(nmj.HookType),
		NetworkID: nmj.NetworkID,
		MemberID:  nmj.MemberID,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}

// convertNetworkCreatedToDBEvent converts NetworkCreated webhook event to a db.Event struct
func (p *processor) convertNetworkCreatedToDBEvent(nc *ztchooks.NetworkCreated) db.Event {
	return db.Event{
		HookID:        nc.HookID,
		OrgID:         nc.OrgID,
		HookType:      string(nc.HookType),
		NetworkID:     nc.NetworkID,
		UserID:        nc.UserID,
		UserEmail:     nc.UserEmail,
		NetworkConfig: nc.NetworkConfig,
		Metadata:      nc.NetworkMetadata,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
}

// convertNetworkConfigChangedToDBEvent converts NetworkConfigChanged webhook event to a db.Event struct
func (p *processor) convertNetworkConfigChangedToDBEvent(ncc *ztchooks.NetworkConfigChanged) db.Event {
	return db.Event{
		HookID:    ncc.HookID,
		OrgID:     ncc.OrgID,
		HookType:  string(ncc.HookType),
		NetworkID: ncc.NetworkID,
		UserID:    ncc.UserID,
		UserEmail: ncc.UserEmail,
		OldConfig: ncc.OldConfig,
		NewConfig: ncc.NewConfig,
		Metadata:  ncc.NetworkMetadata,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}

// SignatureOpt is a way to extend the functionality of the NewProcessor func
type SignatureOpt func(p *processor)

// WithSignatureInfo is passed to NewProcessor when a signature for validating the webhook payload
// is needed
func WithSignatureInfo(signature, preSharedKey string) SignatureOpt {
	return func(p *processor) {
		p.signature = signature
		p.preSharedKey = preSharedKey
	}
}

// NewProcessor makes a Processor client
func NewProcessor(dbClient db.DB) (Processor, error) {
	p := &processor{dbClient: dbClient}
	return p, nil
}
