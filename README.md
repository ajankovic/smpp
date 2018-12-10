# SMPP 3.4 Library

[![GoDoc Badge]][GoDoc] [![GoReportCard Badge]][GoReportCard] [![Build Status](https://travis-ci.com/ajankovic/smpp.svg?branch=master)](https://travis-ci.com/ajankovic/smpp)

**smpp** is library contains implementation of [SMPP 3.4 protocol](http://opensmpp.org/specs/smppv34_gsmumts_ig_v10.pdf).

It allows easier creation of SMPP clients and servers by providing utilities for PDU and session handling.

## Project Status

Although usable, project is still to be considered as _WORK IN PROGRESS_ until it reaches it's first official version. Util then breaking API to fix possible design flaws is possible.

## Project goals

- Provide API for manipulating and observing PDU data.
- Provide API for SMPP session handling.
- Provide reference SMPP SMSC Server implementation.
- Enforce SMPP specification integrity (to the most useful degree).

## Feature description

- [X] Smpp protocol integrity is enforced with the use of session states.
- [X] Allow communication over tcp protocol.
- [X] Pdu data structure should have access to all of it's elements.
- [X] When client and server are connected session is created.
- [X] When connection is terminated or unbinding is finished session is closed.
- [X] Closing can be completed only after graceful handling of all remaining operations.
- [X] Request window is the number of sent requests during a session that are waiting for the matching response from the other side.
- [X] Size of the window is configurable if window is closed send should fail.
- [X] Sending should timeout if response is not received in configured amount of time.
- [X] Sender should wait for responses only for requests that have matching responses defined by the spec.
- [X] Client should have an option to listen for outbind requests from server.
- [X] Server should be able to rate limit client's requests.
- [X] Session should handle sequence numbers.
- [X] Provide logging for critical paths.
- [X] Sessions should be uniquely identifiable.
- [ ] Helpers for sending enquire_link in regular intervals.
- [ ] If an SMPP entity receives an unrecognized PDU/command, it must return a generic_nack PDU indicating an invalid command_id in the command_status field of the header.
- [ ] Provide stats about running session(s):

  - Open sessions
  - Type of sessions
  - Session send window size
  - Session receive window size
  - Session running time
  - Average send/response times
  - Rate of failures
- [ ] Support all PDU commands defined by the specification.
- [ ] Helper functions for common tasks.

## Installation

You can use _go get_:

    go get -u github.com/ajankovic/smpp

## Usage

In order to do any kind of interaction you first need to create an SMPP [Session](https://godoc.org/github.com/ajankovic/smpp#Session). Session is the main carrier of the protocol and enforcer of the specification rules.

Naked session can be created with:

    // You must provide already established connection and configuration struct.
    sess := smpp.NewSession(conn, conf)

But it's much more convenient to use helpers that would do the binding with the remote SMSC and return you session prepared for sending:

    // Bind with remote server by providing config structs.
    sess, err := smpp.BindTRx(sessConf, bindConf)

And once you have the session it can be used for sending PDUs to the binded peer.

    sm := smpp.SubmitSm{
        SourceAddr:      "11111111",
        DestinationAddr: "22222222",
        ShortMessage:    "Hello from SMPP!",
    }
    // Session can then be used for sending PDUs.
    resp, err := sess.Send(p)

Session that is no longer used must be closed:

    sess.Close()

If you want to handle incoming requests to the session specify SMPPHandler in session configuration when creating new session similarly to HTTPHandler from _net/http_ package:

    conf := smpp.SessionConf{
        Handler: smpp.HandlerFunc(func(ctx *smpp.Context) {
            switch ctx.CommandID() {
            case pdu.UnbindID:
                ubd, err := ctx.Unbind()
                if err != nil {
                    t.Errorf(err.Error())
                }
                resp := ubd.Response()
                if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
                    t.Errorf(err.Error())
                }
            }
        }),
    }

Detailed examples for SMPP client and server can be found in the [examples dir](https://github.com/ajankovic/smpp/tree/master/parser).

[GoDoc]: https://godoc.org/github.com/ajankovic/smpp
[GoDoc Badge]: https://godoc.org/github.com/ajankovic/smpp?status.svg
[GoReportCard]: https://goreportcard.com/report/github.com/ajankovic/smpp
[GoReportCard Badge]: https://goreportcard.com/badge/github.com/ajankovic/smpp