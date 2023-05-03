package main

const (
	MSG_PAGE_SYSTEM_BASE = 0x00
	MSG_PAGE_MAIN_BASE   = 0x10
	MSG_PAGE_MAIN_GET    = 0x11
	MSG_PAGE_MAIN_SAVE   = 0x12
	MSG_PAGE_MAIN_IDLE   = 0x13
	MSG_PAGE_DETAIL_BASE = 0x20
	MSG_PAGE_DETAIL_GET  = 0x21
)

type Msg struct {
	Id    int
	Param string
	Data  interface{}
}
