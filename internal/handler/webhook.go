package handler

type WebHookHandler struct {

}

func NewWebHookHandler () WebHookHandler {
	return WebHookHandler{}
}

func (h *WebHookHandler) create()