package smpp

import (
	"context"
	"errors"
	"fmt"

	"github.com/ajankovic/smpp/pdu"
)

// Context represents container for SMPP request related information.
type Context struct {
	sess   *Session
	status pdu.Status
	ctx    context.Context
	req    pdu.PDU
	resp   pdu.PDU
	close  bool
}

// SystemID returns SystemID of the bounded peer that request came from.
func (ctx *Context) SystemID() string {
	return ctx.sess.conf.SystemID
}

// SessionID returns ID of the session that this context is responsible for handling this request.
func (ctx *Context) SessionID() string {
	return ctx.sess.ID()
}

// CommandID returns ID of the PDU request.
func (ctx *Context) CommandID() pdu.CommandID {
	return ctx.req.CommandID()
}

// RemoteAddr returns IP address of the bounded peer.
func (ctx *Context) RemoteAddr() string {
	return ctx.sess.remoteAddr()
}

// Context returns regular Go context.
func (ctx *Context) Context() context.Context {
	return ctx.ctx
}

// Status returns status of the current request.
func (ctx *Context) Status() pdu.Status {
	return ctx.status
}

// Respond sends pdu to the bounded peer.
func (ctx *Context) Respond(resp pdu.PDU, status pdu.Status) error {
	ctx.status = status
	ctx.resp = resp
	if resp == nil {
		return errors.New("smpp: responding with nil PDU")
	}

	ctx.sess.mu.Lock()
	if err := ctx.sess.makeTransition(resp.CommandID(), false); err != nil {
		ctx.sess.conf.Logger.ErrorF("transitioning resp pdu: %s %+v", ctx.sess, err)
		ctx.sess.mu.Unlock()
		return err
	}
	if _, err := ctx.sess.enc.Encode(resp, status); err != nil {
		ctx.sess.conf.Logger.ErrorF("error encoding pdu: %s %+v", ctx.sess, err)
		ctx.sess.mu.Unlock()
		return err
	}
	ctx.sess.conf.Logger.InfoF("sent response: %s %s %+v", ctx.sess, resp.CommandID(), resp)
	ctx.sess.mu.Unlock()

	return nil
}

// CloseSession will initiate session shutdown after handler returns.
func (ctx *Context) CloseSession() {
	ctx.close = true
}

// GenericNack returns generic request PDU as pdu.GenericNack.
func (ctx *Context) GenericNack() (*pdu.GenericNack, error) {
	if p, ok := ctx.req.(*pdu.GenericNack); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// BindRx returns generic request PDU as pdu.BindRx.
func (ctx *Context) BindRx() (*pdu.BindRx, error) {
	if p, ok := ctx.req.(*pdu.BindRx); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// BindRxResp returns generic request PDU as pdu.BindRxResp.
func (ctx *Context) BindRxResp() (*pdu.BindRxResp, error) {
	if p, ok := ctx.req.(*pdu.BindRxResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// BindTx returns generic request PDU as pdu.BindTx.
func (ctx *Context) BindTx() (*pdu.BindTx, error) {
	if p, ok := ctx.req.(*pdu.BindTx); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// BindTxResp returns generic request PDU as pdu.BindTxResp.
func (ctx *Context) BindTxResp() (*pdu.BindTxResp, error) {
	if p, ok := ctx.req.(*pdu.BindTxResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// QuerySm returns generic request PDU as pdu.QuerySm.
func (ctx *Context) QuerySm() (*pdu.QuerySm, error) {
	if p, ok := ctx.req.(*pdu.QuerySm); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// QuerySmResp returns generic request PDU as pdu.QuerySmResp.
func (ctx *Context) QuerySmResp() (*pdu.QuerySmResp, error) {
	if p, ok := ctx.req.(*pdu.QuerySmResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// SubmitSm returns generic request PDU as pdu.SubmitSm.
func (ctx *Context) SubmitSm() (*pdu.SubmitSm, error) {
	if p, ok := ctx.req.(*pdu.SubmitSm); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// SubmitSmResp returns generic request PDU as pdu.SubmitSmResp.
func (ctx *Context) SubmitSmResp() (*pdu.SubmitSmResp, error) {
	if p, ok := ctx.req.(*pdu.SubmitSmResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// DeliverSm returns generic request PDU as pdu.DeliverSm.
func (ctx *Context) DeliverSm() (*pdu.DeliverSm, error) {
	if p, ok := ctx.req.(*pdu.DeliverSm); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// DeliverSmResp returns generic request PDU as pdu.DeliverSmResp.
func (ctx *Context) DeliverSmResp() (*pdu.DeliverSmResp, error) {
	if p, ok := ctx.req.(*pdu.DeliverSmResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// Unbind returns generic request PDU as pdu.Unbind.
func (ctx *Context) Unbind() (*pdu.Unbind, error) {
	if p, ok := ctx.req.(*pdu.Unbind); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// UnbindResp returns generic request PDU as pdu.UnbindResp.
func (ctx *Context) UnbindResp() (*pdu.UnbindResp, error) {
	if p, ok := ctx.req.(*pdu.UnbindResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// ReplaceSm returns generic request PDU as pdu.ReplaceSm.
func (ctx *Context) ReplaceSm() (*pdu.ReplaceSm, error) {
	if p, ok := ctx.req.(*pdu.ReplaceSm); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// ReplaceSmResp returns generic request PDU as pdu.ReplaceSmResp.
func (ctx *Context) ReplaceSmResp() (*pdu.ReplaceSmResp, error) {
	if p, ok := ctx.req.(*pdu.ReplaceSmResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// CancelSm returns generic request PDU as pdu.CancelSm.
func (ctx *Context) CancelSm() (*pdu.CancelSm, error) {
	if p, ok := ctx.req.(*pdu.CancelSm); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// CancelSmResp returns generic request PDU as pdu.CancelSmResp.
func (ctx *Context) CancelSmResp() (*pdu.CancelSmResp, error) {
	if p, ok := ctx.req.(*pdu.CancelSmResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// BindTRx returns generic request PDU as pdu.BindTRx.
func (ctx *Context) BindTRx() (*pdu.BindTRx, error) {
	if p, ok := ctx.req.(*pdu.BindTRx); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// BindTRxResp returns generic request PDU as pdu.BindTRxResp.
func (ctx *Context) BindTRxResp() (*pdu.BindTRxResp, error) {
	if p, ok := ctx.req.(*pdu.BindTRxResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// Outbind returns generic request PDU as pdu.Outbind.
func (ctx *Context) Outbind() (*pdu.Outbind, error) {
	if p, ok := ctx.req.(*pdu.Outbind); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// EnquireLink returns generic request PDU as pdu.EnquireLink.
func (ctx *Context) EnquireLink() (*pdu.EnquireLink, error) {
	if p, ok := ctx.req.(*pdu.EnquireLink); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// EnquireLinkResp returns generic request PDU as pdu.EnquireLinkResp.
func (ctx *Context) EnquireLinkResp() (*pdu.EnquireLinkResp, error) {
	if p, ok := ctx.req.(*pdu.EnquireLinkResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// SubmitMulti returns generic request PDU as pdu.SubmitMulti.
func (ctx *Context) SubmitMulti() (*pdu.SubmitMulti, error) {
	if p, ok := ctx.req.(*pdu.SubmitMulti); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// SubmitMultiResp returns generic request PDU as pdu.SubmitMultiResp.
func (ctx *Context) SubmitMultiResp() (*pdu.SubmitMultiResp, error) {
	if p, ok := ctx.req.(*pdu.SubmitMultiResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// AlertNotification returns generic request PDU as pdu.AlertNotification.
func (ctx *Context) AlertNotification() (*pdu.AlertNotification, error) {
	if p, ok := ctx.req.(*pdu.AlertNotification); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// DataSm returns generic request PDU as pdu.DataSm.
func (ctx *Context) DataSm() (*pdu.DataSm, error) {
	if p, ok := ctx.req.(*pdu.DataSm); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}

// DataSmResp returns generic request PDU as pdu.DataSmResp.
func (ctx *Context) DataSmResp() (*pdu.DataSmResp, error) {
	if p, ok := ctx.req.(*pdu.DataSmResp); ok {
		return p, nil
	}
	return nil, fmt.Errorf("smpp: invalid cast PDU is of type %s", ctx.req.CommandID())
}
