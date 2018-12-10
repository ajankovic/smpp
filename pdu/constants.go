package pdu

//go:generate stringer -type=Status,CommandID,TagID

const (
	// MaxPDUSize is maximal size of the PDU in bytes.
	MaxPDUSize = 4096 // 4KB
)

// Status represents four byte command status.
type Status uint32

// PDU Command Status set.
const (
	StatusOK              Status = 0x00000000
	StatusInvMsgLen       Status = 0x00000001
	StatusInvCmdLen       Status = 0x00000002
	StatusInvCmdID        Status = 0x00000003
	StatusInvBnd          Status = 0x00000004
	StatusAlyBnd          Status = 0x00000005
	StatusInvPrtFlg       Status = 0x00000006
	StatusInvRegDlvFlg    Status = 0x00000007
	StatusSysErr          Status = 0x00000008
	StatusInvSrcAdr       Status = 0x0000000A
	StatusInvDstAdr       Status = 0x0000000B
	StatusInvMsgID        Status = 0x0000000C
	StatusBindFail        Status = 0x0000000D
	StatusInvPaswd        Status = 0x0000000E
	StatusInvSysID        Status = 0x0000000F
	StatusCancelFail      Status = 0x00000011
	StatusReplaceFail     Status = 0x00000013
	StatusMsgQFul         Status = 0x00000014
	StatusInvSerTyp       Status = 0x00000015
	StatusInvNumDe        Status = 0x00000033
	StatusInvDLName       Status = 0x00000034
	StatusInvDestFlag     Status = 0x00000040
	StatusInvSubRep       Status = 0x00000042
	StatusInvEsmClass     Status = 0x00000043
	StatusCntSubDL        Status = 0x00000044
	StatusSubmitFail      Status = 0x00000045
	StatusInvSrcTON       Status = 0x00000048
	StatusInvSrcNPI       Status = 0x00000049
	StatusInvDstTON       Status = 0x00000050
	StatusInvDstNPI       Status = 0x00000051
	StatusInvSysTyp       Status = 0x00000053
	StatusInvRepFlag      Status = 0x00000054
	StatusInvNumMsgs      Status = 0x00000055
	StatusThrottled       Status = 0x00000058
	StatusInvSched        Status = 0x00000061
	StatusInvExpiry       Status = 0x00000062
	StatusInvDftMsgID     Status = 0x00000063
	StatusTempAppErr      Status = 0x00000064
	StatusPermAppErr      Status = 0x00000065
	StatusRejeAppErr      Status = 0x00000066
	StatusQueryFail       Status = 0x00000067
	StatusInvOptParStream Status = 0x000000C0
	StatusOptParNotAllwd  Status = 0x000000C1
	StatusInvParLen       Status = 0x000000C2
	StatusMissingOptParam Status = 0x000000C3
	StatusInvOptParamVal  Status = 0x000000C4
	StatusDeliveryFailure Status = 0x000000FE
	StatusUnknownErr      Status = 0x000000FF
)

// CommandID is four byte PDU command identifier.
type CommandID uint32

// SMPP command set.
const (
	GenericNackID         CommandID = 0x80000000
	BindReceiverID        CommandID = 0x00000001
	BindReceiverRespID    CommandID = 0x80000001
	BindTransmitterID     CommandID = 0x00000002
	BindTransmitterRespID CommandID = 0x80000002
	QuerySmID             CommandID = 0x00000003
	QuerySmRespID         CommandID = 0x80000003
	SubmitSmID            CommandID = 0x00000004
	SubmitSmRespID        CommandID = 0x80000004
	DeliverSmID           CommandID = 0x00000005
	DeliverSmRespID       CommandID = 0x80000005
	UnbindID              CommandID = 0x00000006
	UnbindRespID          CommandID = 0x80000006
	ReplaceSmID           CommandID = 0x00000007
	ReplaceSmRespID       CommandID = 0x80000007
	CancelSmID            CommandID = 0x00000008
	CancelSmRespID        CommandID = 0x80000008
	BindTransceiverID     CommandID = 0x00000009
	BindTransceiverRespID CommandID = 0x80000009
	OutbindID             CommandID = 0x0000000B
	EnquireLinkID         CommandID = 0x00000015
	EnquireLinkRespID     CommandID = 0x80000015
	SubmitMultiID         CommandID = 0x00000021
	SubmitMultiRespID     CommandID = 0x80000021
	AlertNotificationID   CommandID = 0x00000102
	DataSmID              CommandID = 0x00000103
	DataSmRespID          CommandID = 0x80000103
)

// SMPP mandatory fields set.
const (
	SystemIDFld             string = "system_id"
	PasswordFld             string = "password"
	SystemTypeFld           string = "system_type"
	InterfaceVersionFld     string = "interface_version"
	AddrTonFld              string = "addr_ton"
	AddrNpiFld              string = "addr_npi"
	AddressRangeFld         string = "address_range"
	ServiceTypeFld          string = "service_type"
	SourceAddrTonFld        string = "source_addr_ton"
	SourceAddrNpiFld        string = "source_addr_npi"
	SourceAddrFld           string = "source_addr"
	DestAddrTonFld          string = "dest_addr_ton"
	DestAddrNpiFld          string = "dest_addr_npi"
	NumberOfDestsFld        string = "number_of_dests"
	DestFlagFld             string = "dest_flag"
	DlNameFld               string = "dl_name"
	DestinationAddrFld      string = "destination_addr"
	NoUnsuccessFld          string = "no_unsuccess"
	EsmClassFld             string = "esm_class"
	ProtocolIDFld           string = "protocol_id"
	PriorityFlagFld         string = "priority_flag"
	ScheduleDeliveryTimeFld string = "schedule_delivery_time"
	ValidityPeriodFld       string = "validity_period"
	RegisteredDeliveryFld   string = "registered_delivery"
	ReplaceIfPresentFlagFld string = "replace_if_present_flag"
	DataCodingFld           string = "data_coding"
	SmDefaultMsgIDFld       string = "sm_default_msg_id"
	SmLengthFld             string = "sm_length"
	ShortMessageFld         string = "short_message"
	MessageIDFld            string = "message_id"
	FinalDateFld            string = "final_date"
	MessageStateFld         string = "message_state"
	ErrorCodeFld            string = "error_code"
	EsmeAddrTonFld          string = "esme_addr_ton"
	EsmeAddrNpiFld          string = "esme_addr_npi"
	EsmeAddrFld             string = "esme_addr"
)

// TagID represents two byte optional tag identifier.
type TagID uint16

// PDU tags for optional fields.
const (
	TagDestAddrSubUnit        TagID = 0x0005
	TagDestNetworkType        TagID = 0x0006
	TagDestBearerType         TagID = 0x0007
	TagDestTelematicsID       TagID = 0x0008
	TagSourceAddrSubunit      TagID = 0x000D
	TagSourceNetworkType      TagID = 0x000E
	TagSourceBearerType       TagID = 0x000F
	TagSourceTelematicsID     TagID = 0x0010
	TagQosTimeToLive          TagID = 0x0017
	TagPayloadType            TagID = 0x0019
	TagAdditionalStatusInfoTe TagID = 0x001D
	TagReceiptedMessageID     TagID = 0x001E
	TagMsMsgWaitFacilities    TagID = 0x0030
	TagPrivacyIndicator       TagID = 0x0201
	TagSourceSubaddress       TagID = 0x0202
	TagDestSubaddress         TagID = 0x0203
	TagUserMessageReference   TagID = 0x0204
	TagUserResponseCode       TagID = 0x0205
	TagSourcePort             TagID = 0x020A
	TagDestinationPort        TagID = 0x020B
	TagSarMsgRefNum           TagID = 0x020C
	TagLanguageIndicator      TagID = 0x020D
	TagSarTotalSegments       TagID = 0x020E
	TagSarSegmentSeqnum       TagID = 0x020F
	TagScInterfaceVersion     TagID = 0x0210
	TagCallbackNumPresInd     TagID = 0x0302
	TagCallbackNumA           TagID = 0x0303
	TagNumberOfMessages       TagID = 0x0304
	TagCallbackNum            TagID = 0x0381
	TagDpfResult              TagID = 0x0420
	TagSetDPF                 TagID = 0x0421
	TagMsAvailabilityStatus   TagID = 0x0422
	TagNetworkErrorCode       TagID = 0x0423
	TagMessagePayload         TagID = 0x0424
	TagDeliveryFailureReason  TagID = 0x0425
	TagMoreMessagesToSend     TagID = 0x0426
	TagMessageState           TagID = 0x0427
	TagUssdServiceOp          TagID = 0x0501
	TagDisplayTime            TagID = 0x1201
	TagSmsSignal              TagID = 0x1203
	TagMsValidity             TagID = 0x1204
	TagAlertOnMessageDeliv    TagID = 0x130C
	TagItsReplyType           TagID = 0x1380
	TagItsSessionInfo         TagID = 0x1383
)
